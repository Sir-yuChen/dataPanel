package crawler

import (
	"context"
	"dataPanel/serviceend/global"
	"errors"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
	"strings"
	"time"
)

type Crawler struct {
	opts    []chromedp.ExecAllocatorOption
	poolKey string
	timeout time.Duration
}

func NewCrawler() *Crawler {
	return &Crawler{
		timeout: 5 * time.Minute,
	}
}
func (c *Crawler) WithTimeout(timeout time.Duration) *Crawler {
	c.timeout = timeout
	return c
}
func (c *Crawler) WithPoolKey(poolKey string) *Crawler {
	c.poolKey = poolKey
	return c
}
func (c *Crawler) WithOpts(opts []chromedp.ExecAllocatorOption) *Crawler {
	c.opts = opts
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

func (c *Crawler) executeActions(actions ...chromedp.Action) (string, error) {
	browserPool := GetBrowserPool()
	ctx, err := browserPool.Acquire(c.poolKey, c.timeout, c.opts...)
	if err != nil || ctx == nil {
		global.GvaLog.Error("获取浏览上下文失败", zap.Error(err))
		return "", err
	}
	defer func() {
		browserPool.Release(ctx)
	}()

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
	)
	for attempt := 0; attempt < 3; attempt++ {
		startTime := time.Now()
		err = chromedp.Run(ctx.Ctx, finalActions...)
		elapsed := time.Since(startTime)
		if err == nil {
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
	global.GvaLog.Error("操作最终失败",
		zap.String("pool", c.poolKey),
		zap.Duration("timeout", c.timeout),
		zap.Error(lastErr))

	return "", lastErr
}
