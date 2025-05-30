package crawler

import (
	"context"
	"dataPanel/serviceend/global"
	"errors"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
	"math"
	"strings"
	"sync"
	"time"
)

type Crawler struct {
	opts         []chromedp.ExecAllocatorOption
	poolKey      string
	timeout      time.Duration
	maxRetries   int           // 最大重试次数
	failureCount int           // 失败计数器
	circuitOpen  bool          // 熔断状态
	lastFailure  time.Time     // 最后失败时间
	cooldown     time.Duration // 熔断冷却时间
	circuitMutex sync.Mutex    // 熔断锁
}

const (
	defaultCooldown = 60 * time.Second // 默认熔断冷却时间
	maxRetries      = 3                // 默认最大重试
)

func NewCrawler(poolKey string, opts []chromedp.ExecAllocatorOption) *Crawler {
	return &Crawler{
		poolKey:    poolKey,
		opts:       opts,
		timeout:    30 * time.Second,
		maxRetries: maxRetries,
		cooldown:   defaultCooldown,
	}
}

// 新增熔断配置方法
func (c *Crawler) WithCircuitBreaker(maxRetries int, cooldown time.Duration) *Crawler {
	c.maxRetries = maxRetries
	c.cooldown = cooldown
	return c
}

// 检查熔断状态
func (c *Crawler) isCircuitOpen() bool {
	c.circuitMutex.Lock()
	defer c.circuitMutex.Unlock()

	if c.circuitOpen {
		if time.Since(c.lastFailure) > c.cooldown {
			c.circuitOpen = false // 自动恢复
			c.failureCount = 0
			return false
		}
		return true
	}
	return false
}

// 更新熔断状态
func (c *Crawler) recordFailure() {
	c.circuitMutex.Lock()
	defer c.circuitMutex.Unlock()

	c.failureCount++
	if c.failureCount >= c.maxRetries {
		c.circuitOpen = true
		c.lastFailure = time.Now()
	}
}

// 重置熔断器
func (c *Crawler) resetCircuit() {
	c.circuitMutex.Lock()
	defer c.circuitMutex.Unlock()

	c.failureCount = 0
	c.circuitOpen = false
}
func (c *Crawler) WithTimeout(timeout time.Duration) *Crawler {
	c.timeout = timeout
	return c
}

func (c *Crawler) GetHTML(url, waitSelector string) (string, error) {
	actions := []chromedp.Action{
		chromedp.Navigate(url),
		chromedp.WaitVisible(waitSelector),
		chromedp.WaitReady(waitSelector),
		chromedp.InnerHTML("body", new(string)),
	}
	return c.executeActions(actions...)
}

func (c *Crawler) ExecuteActions(actions ...chromedp.Action) (string, error) {
	return c.executeActions(actions...)
}

func (c *Crawler) executeActions(actions ...chromedp.Action) (string, error) {
	if c.isCircuitOpen() {
		return "", errors.New("circuit breaker is open")
	}

	browserPool, err := GetBrowserPool(c.poolKey, c.opts)
	if err != nil {
		c.recordFailure()
		global.GvaLog.Error("获取浏览器池失败", zap.Error(err))
		return "", err
	}

	ctx, err := browserPool.Acquire()
	if err != nil || ctx == nil {
		global.GvaLog.Error("获取浏览上下文失败", zap.Error(err))
		return "", err
	}
	var htmlContent string
	finalActions := append(actions, chromedp.InnerHTML("body", &htmlContent))

	// 错误分类处理函数
	shouldRetry := func(err error) bool {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return true
		case strings.Contains(err.Error(), "net::ERR_CONNECTION_RESET"):
			return true
		case strings.Contains(err.Error(), "Session terminated"):
			return true
		case strings.Contains(err.Error(), "handshake failed"):
			return true
		default:
			return false
		}
	}

	var (
		lastErr      error
		totalElapsed time.Duration
		baseTimeout  = c.timeout
	)
	for attempt := 0; attempt < c.maxRetries; attempt++ {
		currentTimeout := baseTimeout + time.Duration(math.Pow(2, float64(attempt)))*time.Second
		timeoutCtx, timeoutCancel := context.WithTimeout(ctx.Ctx, currentTimeout)

		startTime := time.Now()
		err = chromedp.Run(timeoutCtx, finalActions...)
		elapsed := time.Since(startTime)
		timeoutCancel()

		if err == nil {
			c.resetCircuit()
			return htmlContent, nil
		}

		if errors.Is(err, context.DeadlineExceeded) {
			totalElapsed += elapsed
			if totalElapsed > c.timeout*3 {
				lastErr = errors.New("maximum timeout renewal exceeded")
				break
			}
			continue
		}

		if !shouldRetry(err) {
			lastErr = err
			break
		}

		global.GvaLog.Warn("可恢复错误触发重试",
			zap.Int("attempt", attempt+1),
			zap.Error(err))
	}

	c.recordFailure()
	global.GvaLog.Error("操作最终失败",
		zap.String("pool", c.poolKey),
		zap.Duration("timeout", c.timeout),
		zap.Int("retries", c.maxRetries),
		zap.Error(lastErr))

	return "", lastErr
}
