package router

import (
	"dataPanel/serviceend/controller"

	"github.com/gin-gonic/gin"
)

// RouterInf 所有controller都要实现该接口设置路由 方便统一添加路由组前缀 多服务器上线使用 此处进行统一中间件等初始化操作
type RouterInf interface {
	SetupRouter(g *gin.RouterGroup)
}

func SetupRouter(g *gin.RouterGroup) {
	controller.NewHelloController().SetupRouter(g)
	controller.NewCommonController().SetupRouter(g)
}
