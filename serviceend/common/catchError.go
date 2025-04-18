package common

import (
	"dataPanel/serviceend/common/ApiReturn"
	"dataPanel/serviceend/common/response"
	"dataPanel/serviceend/global"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CatchError() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				url := c.Request.URL
				method := c.Request.Method
				//断言失败 ok false/true
				e, ok := err.(ApiReturn.ApiReturnCode)
				if ok {
					global.GvaLog.Error("全局异常Handler", zap.Any("url", url), zap.Any("method", method), zap.Any("Error", err))
					response.WithApiReturn(e, c)
					c.Abort()
					return
				}
				// 没有定义 错误
				global.GvaLog.Error("未知错误类型", zap.Any("url", url), zap.Any("method", method), zap.Any("Error", err))
				unknownErr := ApiReturn.UnknownErr
				unknownErr.Msg = err.(string)
				response.WithApiReturn(unknownErr, c)
				c.Abort()
			}
		}()
		c.Next()
	}
}
