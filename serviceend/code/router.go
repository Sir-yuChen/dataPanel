package code

import (
	_ "dataPanel/docs"
	"dataPanel/serviceend/common/middleware"
	"dataPanel/serviceend/global"
	"dataPanel/serviceend/model"
	"dataPanel/serviceend/router"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"reflect"
)

func CreateGinServer() (engine *gin.Engine) {
	//创建gin 实例
	engine = gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册 model.LocalTime 类型的自定义校验规则
		v.RegisterCustomTypeFunc(ValidateJSONDateType, model.LocalTime{})
	}
	engine.Use(middleware.CatchError()) //全局异常处理
	engine.Use(middleware.Cors())       // 直接放行全部跨域请求
	engine.Use(middleware.RequestApiLog())
	g := engine.RouterGroup.Group(global.GvaConfig.System.ApplicationName)
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.SetupRouter(g)
	global.GvaLog.Info("路由加载成功  GinServer register success")
	return engine
}

// binding:"required" 并不能正常工作问题解决
func ValidateJSONDateType(field reflect.Value) interface{} {
	if field.Type() == reflect.TypeOf(model.LocalTime{}) {
		timeStr := field.Interface().(model.LocalTime).String()
		// 0001-01-01 00:00:00 是 go 中 time.Time 类型的空值
		// 这里返回 Nil 则会被 validator 判定为空值，而无法通过 `binding:"required"` 规则
		if timeStr == "0001-01-01 00:00:00" {
			return nil
		}
		return timeStr
	}
	return nil
}
