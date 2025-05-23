package utils

import (
	"dataPanel/serviceend/global"
	"dataPanel/serviceend/model"
	"sync"

	"go.uber.org/zap"
)

// sync.Once 实现（线程安全推荐）
var (
	messageBusInstance *MessageBus
	once               sync.Once
)

// 全局前端消息提示Bus 发布/订阅模式
type MessageBus struct {
	Subscriptions map[string]map[chan model.MessageDialogModel]bool
	mu            sync.RWMutex
}

func NewMessageBus() *MessageBus {
	once.Do(func() {
		messageBusInstance = &MessageBus{
			Subscriptions: make(map[string]map[chan model.MessageDialogModel]bool),
		}
	})
	return messageBusInstance
}

// Subscribe 添加订阅通道（避免重复添加）
func (eb *MessageBus) Subscribe(eventType string, ch chan model.MessageDialogModel) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	if _, exists := eb.Subscriptions[eventType]; !exists {
		eb.Subscriptions[eventType] = make(map[chan model.MessageDialogModel]bool)
	}

	if eb.Subscriptions[eventType][ch] {
		global.GvaLog.Warn("已存在重复订阅事件类型", zap.Any("eventType", eventType))
		return // 避免重复订阅
	}
	eb.Subscriptions[eventType][ch] = true
	global.GvaLog.Info("添加订阅通道成功", zap.Any("eventType", eventType))
}

// Unsubscribe 取消订阅（不再关闭通道）
func (eb *MessageBus) Unsubscribe(eventType string, ch chan model.MessageDialogModel) bool {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	subs, ok := eb.Subscriptions[eventType]
	if !ok {
		return false
	}
	if !subs[ch] {
		return false
	}

	delete(subs, ch)
	return true
}

// Publish 发布事件（优化锁范围）
func (eb *MessageBus) Publish(eventType string, data model.MessageDialogModel) {
	//异常发布消息无法获取锁
	if !eb.mu.TryRLock() {
		global.GvaLog.Error("无法获取锁", zap.Any("eventType", eventType), zap.Any("data", data))
		return
	}
	subsMap, ok := eb.Subscriptions[eventType]
	if !ok {
		eb.mu.RUnlock()
		global.GvaLog.Error("订阅事件不存在", zap.Any("eventType", eventType), zap.Any("data", data))
		return
	}
	subs := make([]chan model.MessageDialogModel, 0, len(subsMap))
	for ch := range subsMap {
		subs = append(subs, ch)
	}
	eb.mu.RUnlock()

	var closedChannels []chan model.MessageDialogModel
	// 发送消息并检测关闭的通道
	for _, ch := range subs {
		func(c chan model.MessageDialogModel) {
			defer func() {
				if r := recover(); r != nil {
					closedChannels = append(closedChannels, c)
					global.GvaLog.Warn("检测到消息通道已关闭",
						zap.String("eventType", eventType),
						zap.Any("channel", c))
				}
			}()
			select {
			case c <- data:
				global.GvaLog.Info("发布消息成功",
					zap.String("eventType", eventType),
					zap.Any("data", data))
			default:
				global.GvaLog.Warn("消息队列已满，丢弃消息",
					zap.String("eventType", eventType),
					zap.Any("data", data))
			}
		}(ch)
	}
	// 3. 自动清理已关闭的通道（需要时加锁）
	if len(closedChannels) > 0 {
		eb.mu.Lock()
		defer eb.mu.Unlock()
		subsMap, ok := eb.Subscriptions[eventType]
		if !ok {
			return
		}
		for _, closedCh := range closedChannels {
			delete(subsMap, closedCh)
			global.GvaLog.Warn("从订阅列表移除已关闭通道",
				zap.String("eventType", eventType))
		}
	}
}
