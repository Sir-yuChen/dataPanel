package crawler

import (
	"dataPanel/serviceend/global"
	"fmt"
	"go.uber.org/zap"
	"math"
	"sync/atomic"
	"time"
)

// PoolMetrics 浏览器池监控指标
type PoolMetrics struct {
	Created         atomic.Int64  // 总创建次数
	Reused          atomic.Int64  // 总复用次数
	Active          atomic.Int32  // 活跃实例数
	Pending         atomic.Int32  // 等待中请求数
	ScaleFactor     atomic.Uint64 // 当前扩容因子
	AcquireTime     atomic.Int64  // 总获取耗时(ns)
	ReleaseTime     atomic.Int64  // 总释放耗时(ns)
	CreateTime      atomic.Int64  // 总创建耗时(ns)
	ScaleTime       atomic.Int64  // 总扩容耗时(ns)
	CleanTime       atomic.Int64  // 总清理耗时(ns)
	MemoryUsage     atomic.Uint64 // 累计内存使用量
	MaxMemoryUsage  atomic.Uint64 // 峰值内存使用量
	MaxAcquireTime  atomic.Int64  // 最大获取耗时(ns)
	MaxReleaseTime  atomic.Int64  // 最大释放耗时(ns)
	collectInterval time.Duration
	stopChan        chan struct{}
}

// 创建新监控实例
func NewPoolMetrics(interval time.Duration) *PoolMetrics {
	return &PoolMetrics{
		collectInterval: interval,
		stopChan:        make(chan struct{}),
	}
}

// Start 启动指标采集协程
func (m *PoolMetrics) Start() {
	go m.monitor()
}

// Stop 停止指标采集
func (m *PoolMetrics) Stop() {
	close(m.stopChan)
}

func (m *PoolMetrics) monitor() {
	ticker := time.NewTicker(m.collectInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			report := m.Report()
			// 添加字段翻译层
			zhReport := map[string]interface{}{
				"活跃实例数":  report["active"],
				"总创建次数":  report["created"],
				"等待请求数":  report["pending"],
				"总复用次数":  report["reused"],
				"最大获取耗时": report["max_acquire_ms"],
				"平均获取耗时": report["avg_acquire_ms"],
				"平均释放耗时": report["avg_release_ms"],
				"平均创建耗时": report["avg_create_ns"],
				"平均扩容耗时": report["avg_scale_ns"],
				"扩容因子":   report["scale_factor"],
				"总内存占用":  report["total_mem_mb"],
				"峰值内存":   report["max_mem_mb"],
			}
			global.GvaLog.Info("浏览器池监控指标报告",
				zap.Any("指标详情", zhReport))
		case <-m.stopChan:
			return
		}
	}
}

// Report 生成监控报告（线程安全）
func (m *PoolMetrics) Report() map[string]interface{} {
	total := m.Created.Load() + m.Reused.Load()
	avgDivisor := func(total int64) int64 {
		if total == 0 {
			return 1 // 防止除零
		}
		return total
	}

	return map[string]interface{}{
		"created":        m.Created.Load(),
		"reused":         m.Reused.Load(),
		"active":         m.Active.Load(),
		"pending":        m.Pending.Load(),
		"avg_acquire_ms": formatNsToMs(m.AcquireTime.Load(), avgDivisor(total)),
		"avg_release_ms": formatNsToMs(m.ReleaseTime.Load(), avgDivisor(total)),
		"total_mem_mb":   formatMB(m.MemoryUsage.Load()),
		"max_mem_mb":     formatMB(m.MaxMemoryUsage.Load()),
		"max_acquire_ms": time.Duration(m.MaxAcquireTime.Load()).Milliseconds(),
		"scale_factor":   math.Float64frombits(m.ScaleFactor.Load()),
		"avg_create_ns":  m.CreateTime.Load() / avgDivisor(m.Created.Load()),
		"avg_scale_ns":   m.ScaleTime.Load() / avgDivisor(m.Created.Load()),
	}
}
func formatNsToMs(totalNs, count int64) float64 {
	if count == 0 {
		return 0
	}
	f := float64(totalNs/count) / 1e6
	return math.Round(f*100) / 100
}

func formatMB(bytes uint64) float64 {
	mb := float64(bytes) / 1024 / 1024
	return math.Round(mb*100) / 100
}

// 格式化耗时
func (m *PoolMetrics) formatDuration(totalNs, count int64) string {
	if count == 0 {
		return "0ms"
	}
	return fmt.Sprintf("%.2fms", float64(totalNs/count)/1e6)
}

// 获取内存指标快照
func (m *PoolMetrics) GetMemoryStats() (uint64, uint64) {
	return m.MemoryUsage.Load(), m.MaxMemoryUsage.Load()
}

// 重置内存指标
func (m *PoolMetrics) ResetMemoryStats() {
	m.MemoryUsage.Store(0)
	m.MaxMemoryUsage.Store(0)
}

// 通用最大值更新函数
func updateMax(target *atomic.Int64, newVal int64) {
	for {
		old := target.Load()
		if newVal <= old {
			break
		}
		if target.CompareAndSwap(old, newVal) {
			break
		}
	}
}
