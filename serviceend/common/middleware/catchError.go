package middleware

import (
	"dataPanel/serviceend/common/ApiReturn"
	"dataPanel/serviceend/common/response"
	"dataPanel/serviceend/global"
	"net"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CatchError() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取用户的请求信息
				httpRequest, _ := httputil.DumpRequest(c.Request, true)
				// 链接中断，客户端中断连接为正常行为，不需要记录堆栈信息
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						errStr := strings.ToLower(se.Error())
						if strings.Contains(errStr, "broken pipe") || strings.Contains(errStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				// 链接中断的情况
				if brokenPipe {
					global.GvaLog.Error("连接中断异常",
						zap.Time("time", time.Now()),
						zap.Any("url", c.Request.URL.Path),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error))
					c.Abort()
					// 链接已断开，无法写状态码
					return
				}
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
				global.GvaLog.Error("未知错误类型", zap.Time("时间", time.Now()), zap.Any("url", url), zap.Any("method", method), zap.Any("Error", err))
				unknownErr := ApiReturn.Err
				unknownErr.Msg = err.(string)
				response.WithApiReturn(unknownErr, c)
				c.Abort()
			}
		}()
		c.Next()
	}
}
