package middleware

import (
	"bytes"
	"dataPanel/serviceend/global"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io/ioutil"
	"time"
)

func RequestApiLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		requestID := uuid.New().String()

		// 统一参数容器
		params := map[string]interface{}{
			"query": c.Request.URL.Query(),
		}

		// 根据 Content-Type 处理不同参数类型
		contentType := c.ContentType()
		switch {
		case contentType == "application/json":
			// 读取并重置请求体
			body, _ := ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

			// 解析 JSON
			var jsonData interface{}
			if err := json.Unmarshal(body, &jsonData); err == nil {
				params["body"] = jsonData
			} else {
				params["raw_body"] = string(body) // 解析失败时记录原始内容
			}

		case c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH":
			if err := c.Request.ParseForm(); err == nil {
				params["form"] = c.Request.PostForm
			} else {
				global.GvaLog.Error("Parse form error",
					zap.String("uuid", requestID),
					zap.Error(err))
			}
		}

		// 公共日志字段
		commonFields := []zap.Field{
			zap.String("uuid", requestID),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("content_type", contentType),
			zap.Any("params", params), // 统一参数字段
		}
		// 请求开始日志（简化字段）
		global.GvaLog.Info("Request start", commonFields...)

		// 处理耗时测量
		processingStart := time.Now()
		c.Next()
		processingTime := time.Since(processingStart)

		endFields := []zap.Field{
			zap.String("uuid", requestID),
			zap.String("path", c.Request.URL.Path),
			zap.Duration("total_latency", time.Since(startTime)),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("processing_time", processingTime),
		}
		global.GvaLog.Info("Request End", endFields...)
	}
}
