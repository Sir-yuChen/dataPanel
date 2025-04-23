package exposed

import (
	"context"
	"dataPanel/serviceend/global"
)

// HelloWails 暴露给wails得 struct
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

// GetHello 自定义暴露的方法
func (h *HelloWails) GetHello() string {
	global.GvaLog.Info("Hello wailis3")
	return "Hello Wails3"
}
