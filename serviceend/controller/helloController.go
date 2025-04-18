package controller

import (
	"dataPanel/serviceend/common/response"
	"dataPanel/serviceend/global"
	"dataPanel/serviceend/service"

	"github.com/gin-gonic/gin"
)

type HelloController struct{}

var (
	helloService = service.ServiceGroupApp.HelloService
)

func NewHelloController() *HelloController {
	return &HelloController{}
}
func (h *HelloController) SetupRouter(g *gin.RouterGroup) {
	helloRouter := g.Group("/")
	{
		helloRouter.GET("/hello", h.GetHello) // 健康监测
	}
}

func (h *HelloController) GetHello(ctx *gin.Context) {
	global.GvaLog.Info("测试接口请求Hello")
	if str := helloService.Hello(); len(str) > 0 {
		response.OkWithDetailed(nil, str, ctx)
	} else {
		response.FailWithMessage("系统异常", ctx)
	}
}
