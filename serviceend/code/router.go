package code

import (
	"dataPanel/serviceend/common"
	"dataPanel/serviceend/common/middleware"
	"dataPanel/serviceend/global"
	"dataPanel/serviceend/router"

	"github.com/gin-gonic/gin"
)

func CreateGinServer() (engine *gin.Engine) {
	//创建gin 实例
	engine = gin.New()
	engine.Use(common.CatchError()) //全局异常处理
	engine.Use(middleware.Cors())   // 直接放行全部跨域请求
	g := engine.RouterGroup.Group(global.GvaConfig.System.ApplicationName)
	router.SetupRouter(g)
	global.GvaLog.Info("路由加载  GinServer register success")
	return engine
}
