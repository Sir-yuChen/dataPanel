package exposed

import (
	"context"
	"dataPanel/serviceend/global"
)

// 暴露给wails得 struct
type HelloWails struct {
	ctx context.Context
}

func NewHelloWails() *HelloWails {
	return &HelloWails{}
}
func (a *HelloWails) SetCtx(ctx context.Context) *HelloWails {
	a.ctx = ctx
	return a
}

// 自定义暴露的方法
func (h *HelloWails) GetHello() string {
	global.GvaLog.Info("Hello wailis")
	return "Hello Wails3"
}
