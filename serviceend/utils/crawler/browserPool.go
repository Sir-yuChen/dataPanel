package crawler

import (
	"context"
	"dataPanel/serviceend/global"
	"fmt"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultPoolSize     = 10
	browserKeepAlive    = 3 * time.Minute
	poolCleanInterval   = 30 * time.Second
	defaultMinSize      = 2
	defaultScaleFactor  = 1.5
	healthCheckInterval = 15 * time.Second
)

type BrowserContext struct {
	Ctx      context.Context
	Cancel   context.CancelFunc
	Valid    bool
	LastUsed time.Time
}

type BrowserPool struct {
	pool        chan *BrowserContext
	opts        []chromedp.ExecAllocatorOption
	mu          sync.RWMutex
	activeCnt   int
	minSize     int
	maxSize     int
	createdCnt  int64
	reusedCnt   int64
	scaleFactor float64
	pendingReqs int32
	scaleMutex  sync.Mutex
}

var (
	pools               = make(map[string]*BrowserPool)
	poolMutex           sync.RWMutex
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

// 获取浏览器池
func GetBrowserPool(poolKey string, opts []chromedp.ExecAllocatorOption) (*BrowserPool, error) {
	poolMutex.Lock()
	defer poolMutex.Unlock()

	if pool, exists := pools[poolKey]; exists {
		return pool, nil
	}

	baseOpts := getBaseAllocatorOptions()
	fullOpts := append(baseOpts, opts...)

	pool := &BrowserPool{
		pool:        make(chan *BrowserContext, defaultMinSize),
		opts:        fullOpts,
		maxSize:     defaultPoolSize,
		minSize:     defaultMinSize,
		scaleFactor: defaultScaleFactor,
	}

	go pool.startPoolMaintenance()
	// 新增监控日志协程
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			metrics := pool.Metrics()
			global.GvaLog.Info("浏览器池监控指标",
				zap.String("浏览器池KEY", poolKey),
				zap.Any("指标", metrics),
			)
		}
	}()

	pools[poolKey] = pool
	return pool, nil
}

// 动态扩容判断
func (bp *BrowserPool) shouldScaleUp() bool {
	bp.scaleMutex.Lock()
	defer bp.scaleMutex.Unlock()

	return len(bp.pool) == 0 &&
		bp.activeCnt >= bp.maxSize &&
		atomic.LoadInt32(&bp.pendingReqs) > int32(bp.maxSize/2)
}

// 扩容执行逻辑
func (bp *BrowserPool) scaleUp() {
	bp.scaleMutex.Lock()
	defer bp.scaleMutex.Unlock()

	if bp.scaleFactor > 2.0 {
		global.GvaLog.Warn("缩放因子过大已自动修正",
			zap.Float64("原值", bp.scaleFactor))
		bp.scaleFactor = 2.0
	}

	newMax := int(float64(bp.maxSize) * bp.scaleFactor)
	if newMax > bp.maxSize*2 {
		newMax = bp.maxSize * 2
	}

	for i := bp.maxSize; i < newMax; i++ {
		bc, err := bp.createBrowserContext()
		if err != nil {
			global.GvaLog.Error("实例创建失败", zap.Error(err))
			continue
		}
		select {
		case bp.pool <- bc:
			bp.activeCnt++
		default:
			bc.Cancel()
		}
	}
	bp.maxSize = newMax
}

// 获取浏览器上下文
func (bp *BrowserPool) Acquire() (*BrowserContext, error) {
	atomic.AddInt32(&bp.pendingReqs, 1)
	defer atomic.AddInt32(&bp.pendingReqs, -1)

	if bp.shouldScaleUp() {
		go bp.scaleUp()
	}

	select {
	case bc := <-bp.pool:
		bp.mu.Lock()
		defer bp.mu.Unlock()
		if bc.isValid() {
			atomic.AddInt64(&bp.reusedCnt, 1)
			bc.LastUsed = time.Now()
			return bc, nil
		}
		bc.Cancel()
		bp.activeCnt--
	default:
		if bp.activeCnt < bp.maxSize {
			const maxRetries = 2
			for i := 0; i < maxRetries; i++ {
				browserContext, err := bp.createBrowserContext()
				if err == nil {
					bp.mu.Lock()
					defer bp.mu.Unlock() //Possible resource leak, 'defer' is called in the 'for' loop
					atomic.AddInt64(&bp.createdCnt, 1)
					bp.activeCnt++
					browserContext.LastUsed = time.Now()
					select {
					case bp.pool <- browserContext:
						return browserContext, nil
					default:
						return browserContext, nil
					}
				}
				if i == maxRetries-1 {
					global.GvaLog.Error("浏览器实例创建失败",
						zap.Error(err),
						zap.Int("retries", maxRetries))
					return nil, err
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
	// Fallback 逻辑
	bp.mu.Lock()
	defer bp.mu.Unlock()
	return bp.createBrowserContext()
}

// 释放浏览器上下文
func (bp *BrowserPool) Release(bc *BrowserContext) {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	if !bc.Valid {
		bc.Cancel()
		bp.activeCnt--
		return
	}

	bc.LastUsed = time.Now()

	if err := chromedp.Run(bc.Ctx, chromedp.Navigate("about:blank")); err != nil {
		bc.markInvalid()
		bc.Cancel()
		bp.activeCnt--
		return
	}

	select {
	case bp.pool <- bc:
	default:
		bc.Cancel()
		bp.activeCnt--
	}
}

// 创建浏览器上下文
func (bp *BrowserPool) createBrowserContext() (*BrowserContext, error) {
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(
		context.Background(),
		append(
			getBaseAllocatorOptions(),
			chromedp.Flag("disable-notifications", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-3d-apis", false), // 已确认无效的flag
		)...,
	)

	ctx, cancelCtx := chromedp.NewContext(allocCtx)

	// 执行反检测脚本
	if err := chromedp.Run(ctx,
		chromedp.Evaluate(antiDetectionScript, nil),
	); err != nil {
		cancelCtx()
		cancelAlloc()
		return nil, fmt.Errorf("反检测初始化失败: %w", err)
	}

	bc := &BrowserContext{
		Ctx:      ctx,
		Cancel:   func() { cancelCtx(); cancelAlloc() },
		Valid:    true,
		LastUsed: time.Now(),
	}

	go bp.monitorBrowserContext(bc)
	return bc, nil
}

// 上下文健康监控
func (bp *BrowserPool) monitorBrowserContext(bc *BrowserContext) {
	defer func() {
		if r := recover(); r != nil {
			global.GvaLog.Error("监控协程异常",
				zap.Any("panic", r),
				zap.Any("context", bc))
		}
	}()

	ticker := time.NewTicker(healthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !bc.isValid() {
				return
			}
			if err := chromedp.Run(bc.Ctx,
				chromedp.Navigate("about:blank"),
				chromedp.WaitVisible("body"),
			); err != nil {
				global.GvaLog.Warn("健康检查失败,当前实例已放弃",
					zap.Error(err),
					zap.Any("context", bc))
				bc.markInvalid()
				bp.activeCnt--
				return
			}
			bc.LastUsed = time.Now()
		case <-bc.Ctx.Done():
			return
		}
	}
}

// 池维护协程
func (bp *BrowserPool) startPoolMaintenance() {
	ticker := time.NewTicker(poolCleanInterval)
	defer ticker.Stop()

	for range ticker.C {
		bp.cleanExpiredInstances()
	}
}

// 清理过期实例
func (bp *BrowserPool) cleanExpiredInstances() {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	keepCount := bp.minSize
	if keepCount < 1 {
		keepCount = 1
	}

	for len(bp.pool) > keepCount {
		select {
		case bc := <-bp.pool:
			if time.Since(bc.LastUsed) > browserKeepAlive || !bc.Valid {
				bc.Cancel()
				bp.activeCnt--
				global.GvaLog.Warn("浏览器实例超时/无效,已被清除")
			} else {
				bp.pool <- bc
			}
		default:
			return
		}
	}
}

// 上下文有效性检查
func (bc *BrowserContext) isValid() bool {
	if !bc.Valid {
		return false
	}
	select {
	case <-bc.Ctx.Done():
		bc.Valid = false
		return false
	default:
		return time.Since(bc.LastUsed) < browserKeepAlive
	}
}

// 标记上下文失效
func (bc *BrowserContext) markInvalid() {
	bc.Valid = false
	bc.Cancel()

}

// 池监控指标
func (bp *BrowserPool) Metrics() map[string]interface{} {
	validCount := 0
	poolSize := len(bp.pool)

	// 计算有效实例数
	for i := 0; i < poolSize; i++ {
		select {
		case bc := <-bp.pool:
			if bc.isValid() {
				validCount++
			}
			bp.pool <- bc
		default:
			break
		}
	}

	return map[string]interface{}{
		"总创建数":  atomic.LoadInt64(&bp.createdCnt),
		"总复用数":  atomic.LoadInt64(&bp.reusedCnt),
		"活跃实例数": bp.activeCnt,
		"池容量":   cap(bp.pool),
		"有效实例数": len(bp.pool),
		"缩放因子":  bp.scaleFactor,
		"健康率": float64(atomic.LoadInt64(&bp.reusedCnt)) /
			float64(atomic.LoadInt64(&bp.createdCnt)+atomic.LoadInt64(&bp.reusedCnt)),
		"等待请求": atomic.LoadInt32(&bp.pendingReqs),
	}
}

// 动态配置池参数
func (bp *BrowserPool) Configure(minSize int, maxSize int, scaleFactor float64) {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	bp.minSize = minSize
	bp.maxSize = maxSize
	bp.scaleFactor = scaleFactor

	if cap(bp.pool) != maxSize {
		newPool := make(chan *BrowserContext, maxSize)
		close(bp.pool)
		for bc := range bp.pool {
			if bc.isValid() {
				newPool <- bc
			} else {
				bc.Cancel()
				bp.activeCnt--
			}
		}
		bp.pool = newPool
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
		chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
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
		// 禁用/dev/shm共享内存（解决Docker环境问题）
		chromedp.Flag("disable-dev-shm-usage", true),
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
	}
	return options
}
