package crawler

import (
	"context"
	"dataPanel/serviceend/global"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
	"math"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultPoolSize    = 20
	browserKeepAlive   = 15 * time.Minute
	poolCleanInterval  = 1 * time.Minute
	defaultMinSize     = 2
	defaultScaleFactor = 1.5
	defaultTimeout     = 60 * time.Second
)

type BrowserContext struct {
	Ctx      context.Context
	Key      string
	Cancel   context.CancelFunc
	Valid    atomic.Bool
	LastUsed atomic.Value // 存储time.Time
	Cleaning atomic.Bool  // 清理状态标记
}

// BrowserPool 浏览器实例池管理结构
// 使用缓冲通道+sync.Map混合模式实现资源池：
// - 通道用于快速调度可用实例
// - sync.Map用于全生命周期跟踪
type BrowserPool struct {
	pool        chan *BrowserContext           // 缓冲池（快速调度）
	opts        []chromedp.ExecAllocatorOption // 浏览器启动配置
	minSize     atomic.Int32                   // 最小池容量
	maxSize     atomic.Int32                   // 最大池容量
	poolMu      sync.Mutex                     // 保护pool通道操作
	scaleFactor atomic.Uint64                  // 动态扩容因子
	pendingReqs atomic.Int32                   // 等待队列长度
	metrics     *PoolMetrics                   // 新增指标实例

}

var (
	instance            *BrowserPool // 单例实例
	once                sync.Once    // 单例控制
	antiDetectionScript = `
		Object.defineProperty(navigator, 'plugins', {
			get: () => [1, 2, 3],
			configurable: false
		});
		Object.defineProperty(navigator, 'languages', {
			get: () => ['zh-CN', 'zh'],
			configurable: false
		});
		window.chrome = undefined;
		delete window._cdc;
		delete window.__driver_evaluate;
	`
)

func GetBrowserPool() *BrowserPool {
	once.Do(func() {
		metrics := NewPoolMetrics(1 * time.Minute)
		pool := &BrowserPool{
			pool:    make(chan *BrowserContext, defaultMinSize),
			opts:    getBaseAllocatorOptions(),
			metrics: metrics,
		}
		pool.minSize.Store(defaultMinSize)
		pool.maxSize.Store(defaultPoolSize)
		pool.scaleFactor.Store(math.Float64bits(defaultScaleFactor))
		pool.metrics.ScaleFactor.Store(math.Float64bits(defaultScaleFactor))
		go pool.startPoolMaintenance()
		go metrics.Start()
		instance = pool
	})
	return instance
}

// Acquire 获取浏览器实例（基于key的智能调度）
func (bp *BrowserPool) Acquire(key string, timeout time.Duration, opts ...chromedp.ExecAllocatorOption) (*BrowserContext, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Nanoseconds()
		bp.metrics.AcquireTime.Add(duration)
		updateMax(&bp.metrics.MaxAcquireTime, duration)
	}()

	bp.pendingReqs.Add(1)
	defer bp.pendingReqs.Add(-1)
	var bc *BrowserContext
	// 优先查找同key的可用实例
	if bc = bp.findMatchingInstance(key); bc != nil {
		bp.metrics.Reused.Add(1)
		bc.LastUsed.Store(time.Now())
		bp.setTimeout(bc, timeout)
		return bc, nil
	}

	if opts != nil && len(opts) > 0 && bc == nil {
		bc, err := bp.createBrowserContext(key, opts...)
		if err == nil && bc != nil {
			bp.setTimeout(bc, timeout)
			return bc, nil
		} else {
			global.GvaLog.Error("获得实例失败", zap.Error(err))
		}
	}
	if bc == nil {
		bc, err := bp.createBrowserContext(key, opts...)
		if err == nil && bc != nil {
			bp.setTimeout(bc, timeout)
			return bc, nil
		} else {
			global.GvaLog.Error("获得实例失败", zap.Error(err))
		}
	}
	return nil, fmt.Errorf("获取实例失败")
}

func (bp *BrowserPool) createBrowserContext(key string, opts ...chromedp.ExecAllocatorOption) (*BrowserContext, error) {
	var allocCtx context.Context
	var cancelAlloc func()
	if opts != nil && len(opts) > 0 {
		allocOpts := append(bp.opts, opts...)
		allocCtx, cancelAlloc = chromedp.NewExecAllocator(context.Background(), allocOpts...)
	} else {
		allocCtx, cancelAlloc = chromedp.NewExecAllocator(context.Background(), bp.opts...)
	}
	ctx, cancelCtx := chromedp.NewContext(allocCtx, chromedp.WithLogf(global.GvaLog.Sugar().Infof))
	if err := chromedp.Run(ctx,
		chromedp.Evaluate(antiDetectionScript, nil),
	); err != nil {
		cancelCtx()
		cancelAlloc()
		return nil, fmt.Errorf("浏览器实例初始化失败: %w", err)
	}

	bc := &BrowserContext{
		Ctx:    ctx,
		Key:    key,
		Cancel: func() { cancelCtx(); cancelAlloc() },
	}
	bc.Valid.Store(true)
	bc.LastUsed.Store(time.Now())
	bp.metrics.Created.Add(1)
	bp.metrics.Active.Add(1)
	//创建完毕则执行一个空action列表，目的chromedp当前的API设计逻辑是只会在第一次调用Run()的时候创建headless-chrome进程
	if err := chromedp.Run(bc.Ctx, make([]chromedp.Action, 0, 1)...); err != nil {
		global.GvaLog.Error("浏览器实例预执行失败", zap.String("key", bc.Key), zap.Error(err))
		return bc, nil
	} else {
		global.GvaLog.Info("浏览器实例预执行成功", zap.String("key", bc.Key))
		return bc, nil
	}
	global.GvaLog.Info("创建浏览器实例成功", zap.String("key", bc.Key))
	return bc, nil
}

func (bp *BrowserPool) Release(bc *BrowserContext) {
	go bp.asyncRelease(bc) // 入口改为异步
}

func (bp *BrowserPool) asyncRelease(bc *BrowserContext) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Nanoseconds()
		bp.metrics.ReleaseTime.Add(duration)
		updateMax(&bp.metrics.MaxReleaseTime, duration)
	}()
	// 前置检查移出锁外
	if !bc.Valid.Load() || errors.Is(bc.Ctx.Err(), context.Canceled) {
		bp.destroyContext(bc)
		global.GvaLog.Warn("释放浏览器前置检查,执行销毁",
			zap.String("key", bc.Key), zap.Any("valid", bc.Valid.Load()),
			zap.Error(bc.Ctx.Err()))
		return
	}
	bp.poolMu.Lock()
	defer bp.poolMu.Unlock()
	if bc.Cleaning.Load() {
		return // 避免重复清理
	}
	// 使用CAS保证状态标记原子性
	if !bc.Cleaning.CompareAndSwap(false, true) {
		return
	}
	go func() {
		// 使用浏览器实例的上下文创建子上下文
		cleanupCtx, cancel := context.WithTimeout(bc.Ctx, defaultTimeout)
		defer cancel()
		// 新增有效性检查
		if errors.Is(cleanupCtx.Err(), context.Canceled) || !bc.Valid.Load() {
			global.GvaLog.Warn("上下文已失效，跳过清理",
				zap.String("key", bc.Key))
			return
		}
		if err := chromedp.Run(cleanupCtx,
			chromedp.Navigate("about:blank"),
			network.ClearBrowserCache(),
			network.ClearBrowserCookies(),
		); err != nil && !errors.Is(err, context.Canceled) {
			global.GvaLog.Warn("清理操作失败",
				zap.String("key", bc.Key),
				zap.Error(err))
		}
		bc.Cleaning.Store(false)
	}()
	// 增强入队逻辑
	select {
	case bp.pool <- bc:
		bp.metrics.Reused.Add(1)
		bc.LastUsed.Store(time.Now())
	default:
		currentCap := cap(bp.pool)
		scaleFactor := int(math.Float64frombits(bp.scaleFactor.Load()))
		newCap := currentCap * scaleFactor

		// 容量保护逻辑
		if newCap <= currentCap {
			newCap = currentCap + 1
		}
		if newCap > int(bp.maxSize.Load()) {
			newCap = int(bp.maxSize.Load())
		}

		if newCap > currentCap {
			newPool := make(chan *BrowserContext, newCap)
			close(bp.pool)
			for len(bp.pool) > 0 {
				select {
				case ctx := <-bp.pool:
					newPool <- ctx
				default:
					break
				}
			}
			newPool <- bc
			bp.pool = newPool
		} else {
			go bp.destroyContext(bc)
			global.GvaLog.Info("通道已达到最大容量,放弃实例,执行销毁", zap.String("key", bc.Key))
		}
	}
}

// 根据key查找匹配实例
func (bp *BrowserPool) findMatchingInstance(key string) *BrowserContext {
	bp.poolMu.Lock()
	defer bp.poolMu.Unlock()

	// 快速遍历通道中的实例
	tempPool := make([]*BrowserContext, 0, len(bp.pool))
	var found *BrowserContext
	for i := 0; i < len(bp.pool); i++ {
		bc := <-bp.pool
		if bc != nil && bc.Key == key && bc.isValid() {
			found = bc
			break
		}
		tempPool = append(tempPool, bc)
	}
	// 重建通道
	for _, bc := range tempPool {
		if bc != nil {
			bp.pool <- bc
		}
	}
	return found
}

// 池维护协程,清理过期实例
func (bp *BrowserPool) startPoolMaintenance() {
	ticker := time.NewTicker(poolCleanInterval)
	defer ticker.Stop()

	for range ticker.C {
		bp.cleanExpiredInstances()
	}
}

// 清理过期实例
func (bp *BrowserPool) cleanExpiredInstances() {
	var validContexts []*BrowserContext
	bp.poolMu.Lock()
	defer bp.poolMu.Unlock()

	// 快速清理通道中的过期实例
	for {
		select {
		case bc := <-bp.pool:
			if bc.isValid() {
				validContexts = append(validContexts, bc)
			} else {
				go bp.destroyContext(bc)
				global.GvaLog.Warn("实例已过期,执行销毁",
					zap.String("key", bc.Key))
			}
		default:
			break
		}
		break
	}

	// 重新填充有效实例
	for _, bc := range validContexts {
		select {
		case bp.pool <- bc:
		default:
			go bp.destroyContext(bc)
			global.GvaLog.Warn("通道已满,放弃实例,并销毁",
				zap.String("key", string(bc.Key)))
		}
	}
}

// 关闭浏览器池
func (bp *BrowserPool) Close() {
	global.GvaLog.Info("开始关闭浏览器池")
	//遍历通道
	for bc := range bp.pool {
		if bc != nil {
			go bp.destroyContext(bc)
		}
	}
	global.GvaLog.Info("开始关闭浏览器池成功")
}

// 上下文有效性检查
func (bc *BrowserContext) isValid() bool {
	if !bc.Valid.Load() {
		return false
	}
	// 双重检查上下文状态
	select {
	case <-bc.Ctx.Done():
		bc.Valid.Store(false)
		return false
	default:
	}
	// 检查是否超时：如果设置超时时间，则正常检查，如果未设置超时时间则最后操作时间记录+浏览器最大存活时间小于等于当前时间，则标记为失效
	if deadline, ok := bc.Ctx.Deadline(); ok {
		remaining := time.Until(deadline)
		if remaining < 100*time.Millisecond { // 提前失效保护
			bc.Valid.Store(false)
			global.GvaLog.Warn("实例已超时,已标记为失效",
				zap.String("key", bc.Key))
			return false
		}
	} else {
		lastUsed := bc.LastUsed.Load().(time.Time)
		if time.Since(lastUsed) > browserKeepAlive { // 1秒缓冲
			bc.Valid.Store(false)
			global.GvaLog.Warn("实例未设置超时时间,已超过最大存活时间,已标记为失效",
				zap.String("key", bc.Key),
				zap.Time("last_used", lastUsed),
				zap.Duration("keep_alive", browserKeepAlive))
			return false
		}
	}
	return true
}

// 上下文销毁方法
func (bp *BrowserPool) destroyContext(bc *BrowserContext) {
	if bc == nil {
		global.GvaLog.Error("实例销毁失败,实例为空")
		return
	}
	// 记录销毁前的状态
	global.GvaLog.Info("开始销毁浏览器实例",
		zap.String("key", bc.Key),
		zap.Bool("valid", bc.Valid.Load()),
		zap.Error(bc.Ctx.Err()))

	// 执行取消操作
	if bc.Cancel != nil {
		bc.Cancel()
	}

	// 确保资源释放
	bp.metrics.Active.Add(-1)
	bp.metrics.Created.Add(-1)

	// 状态标记
	bc.Valid.Store(false)
	//销毁进程
	chromedp.Cancel(bc.Ctx)

	global.GvaLog.Info("完成实例销毁",
		zap.String("key", bc.Key),
		zap.Time("last_used", bc.LastUsed.Load().(time.Time)))
}

// Configure 动态配置池参数
func (bp *BrowserPool) Configure(minSize int, maxSize int, scaleFactor float64) error {
	if minSize < 1 || maxSize < minSize || scaleFactor < 1.0 {
		return fmt.Errorf("invalid parameters")
	}
	bp.minSize.Store(int32(minSize))
	bp.maxSize.Store(int32(maxSize))
	bp.scaleFactor.Store(math.Float64bits(scaleFactor))
	if cap(bp.pool) != maxSize {
		newPool := make(chan *BrowserContext, maxSize)
		close(bp.pool)
		for bc := range bp.pool {
			if bc.isValid() {
				newPool <- bc
			} else {
				bp.destroyContext(bc)
				global.GvaLog.Warn("实例已无效,执行销毁",
					zap.String("key", string(bc.Key)))
			}
		}
		bp.pool = newPool
	}
	return nil
}

func (bp *BrowserPool) setTimeout(bc *BrowserContext, timeout time.Duration) *BrowserContext {
	if bc == nil {
		return bc
	}
	// 处理零值情况，使用默认超时时间
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	lastUsed, ok := bc.LastUsed.Load().(time.Time)
	if !ok {
		global.GvaLog.Warn("无效的 LastUsed 类型",
			zap.String("key", bc.Key),
			zap.Any("actual_type", fmt.Sprintf("%T", bc.LastUsed.Load())))
		return bc
	}

	// 计算新超时时间并确保不超过最大存活时间
	newTimeout := min(timeout, browserKeepAlive)
	updatedLastUsed := lastUsed.Add(newTimeout)
	bc.LastUsed.Store(updatedLastUsed)

	// 获取原始取消函数链
	originalCancel := bc.Cancel

	// 上下文超时设置（保持上下文链完整）
	if deadline, ok := bc.Ctx.Deadline(); ok {
		newDeadline := updatedLastUsed
		if newDeadline.After(deadline) {
			// 创建继承原有上下文的新实例
			newCtx, cancel := context.WithDeadline(bc.Ctx, newDeadline)
			// 新增剩余时间计算
			remaining := time.Until(newDeadline)
			// 构建新的取消函数链
			bc.Cancel = func() {
				cancel()
				originalCancel() // 保持原有取消逻辑
			}
			bc.Ctx = newCtx

			global.GvaLog.Info("浏览器实例超时时间已更新",
				zap.String("key", bc.Key),
				zap.Duration("added_timeout", newTimeout),
				zap.Duration("remaining_time", remaining),
				zap.Time("new_deadline", newDeadline))
		}
	} else {
		// 继承浏览器分配器的原始上下文
		allocCtx := context.Background()
		if parentCtx := contextGetParent(bc.Ctx); parentCtx != nil {
			allocCtx = parentCtx
		}
		newCtx, cancel := context.WithDeadline(allocCtx, updatedLastUsed)

		// 保持完整的取消链
		bc.Cancel = func() {
			cancel()
			originalCancel()
		}
		bc.Ctx = newCtx

		global.GvaLog.Debug("新增上下文超时设置",
			zap.String("key", bc.Key),
			zap.Duration("timeout", newTimeout))
	}

	return bc
}

func contextGetParent(ctx context.Context) context.Context {
	if ctx == nil {
		return nil
	}
	switch ctx.(type) { // 移除未使用的变量c
	case interface{ Deadline() (time.Time, bool) }:
		return ctx
	case interface{ Value(key any) any }:
		return ctx
	case interface{ Done() <-chan struct{} }:
		return ctx
	case interface{ Err() error }:
		return ctx
	default:
		if reflect.TypeOf(ctx).String() == "*context.timerCtx" {
			if parent := reflect.ValueOf(ctx).Elem().FieldByName("parent"); parent.IsValid() {
				return parent.Interface().(context.Context)
			}
		}
		return nil
	}
}

func getBaseAllocatorOptions() []chromedp.ExecAllocatorOption {
	options := []chromedp.ExecAllocatorOption{
		// 无头模式运行浏览器（无GUI界面）
		chromedp.Flag("headless", true),
		// 禁用图片加载（提升爬取速度）
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		// 禁用GPU加速（避免无头模式下的渲染问题）
		chromedp.Flag("disable-gpu", true),
		// 禁用后台网络请求（提升性能）
		chromedp.Flag("disable-background-networking", true),
		// 启用网络服务相关特性（优化资源管理）
		chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess,PreciseMemoryInfo"),
		// 禁用后台定时器节流（提升定时任务准确性）
		chromedp.Flag("disable-background-timer-throttling", true),
		// 禁止后台遮挡窗口优化（保持进程活跃）
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		// 禁用崩溃报告组件（避免干扰）
		chromedp.Flag("disable-breakpad", true),
		// 关闭客户端反钓鱼检测（提升性能）
		chromedp.Flag("disable-client-side-phishing-detection", true),
		// 禁用默认应用程序（减少资源占用）
		chromedp.Flag("disable-default-apps", true),
		// 禁用所有扩展程序（保证环境纯净）
		chromedp.Flag("disable-extensions", true),
		// 关闭特定浏览器特性（优化兼容性）
		chromedp.Flag("disable-features", "site-per-process,Translate,BlinkGenPropertyTrees"),
		// 禁用挂起监控（提升稳定性）
		chromedp.Flag("disable-hang-monitor", true),
		// 禁用IPC洪水攻击防护（提升通信效率）
		chromedp.Flag("disable-ipc-flooding-protection", true),
		// 禁用弹窗拦截（避免漏抓弹窗内容）
		chromedp.Flag("disable-popup-blocking", true),
		// 禁用重新提交表单提示（保持流程连贯）
		chromedp.Flag("disable-prompt-on-repost", true),
		// 禁用渲染进程后台化（保持渲染优先级）
		chromedp.Flag("disable-renderer-backgrounding", true),
		// 关闭浏览器同步功能（避免账号关联）
		chromedp.Flag("disable-sync", true),
		// 强制使用sRGB色彩配置（统一渲染效果）
		chromedp.Flag("force-color-profile", "srgb"),
		// 仅记录基础指标（减少数据收集）
		chromedp.Flag("metrics-recording-only", true),
		// 禁用安全浏览自动更新（提升启动速度）
		chromedp.Flag("safebrowsing-disable-auto-update", true),
		// 启用自动化控制标志（隐藏自动化特征）
		chromedp.Flag("enable-automation", true),
		// 使用基础密码存储（避免系统密钥环依赖）
		chromedp.Flag("password-store", "basic"),
		// 启用模拟密钥链（兼容无密钥系统环境）
		chromedp.Flag("use-mock-keychain", true),
		chromedp.Flag("disable-notifications", true),
		chromedp.Flag("enable-precise-memory-info", true),
		chromedp.Flag("enable-memory-info", true),
		chromedp.Flag("disable-features", "site-per-process,Translate,BlinkGenPropertyTrees,OutOfBlinkCors"), // 调整禁用特性
		// 允许使用共享内存
		chromedp.Flag("disable-dev-shm-usage", true),
	}
	return options
}
