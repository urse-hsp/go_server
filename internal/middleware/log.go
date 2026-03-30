// 请求日志 + 响应日志中间件
// 把每一次接口的请求和返回完整记录下来（用于排查问题 / 监控 / 追踪）
package middleware

import (
	"bytes"
	"go-server/pkg/log"
	"io"
	"time"

	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/duke-git/lancet/v2/random"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 请求日志（在请求进入接口前执行）
func RequestLogMiddleware(logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// The configuration is initialized once per request
		uuid, err := random.UUIdV4() // 生成一个随机的 UUID 作为 trace ID （每个请求都有一个唯一 ID）
		if err != nil {
			return
		}
		trace := cryptor.Md5String(uuid)
		// 记录请求信息
		logger.WithValue(ctx, zap.String("trace", trace))
		logger.WithValue(ctx, zap.String("request_method", ctx.Request.Method))
		logger.WithValue(ctx, zap.Any("request_headers", ctx.Request.Header))
		logger.WithValue(ctx, zap.String("request_url", ctx.Request.URL.String()))
		if ctx.Request.Body != nil {
			bodyBytes, _ := ctx.GetRawData()                            // 读取请求体
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 关键点
			logger.WithValue(ctx, zap.String("request_params", string(bodyBytes)))
		}
		logger.WithContext(ctx).Info("Request") // 打印日志
		ctx.Next()
	}
}

// 响应日志（在请求返回接口后执行）
func ResponseLogMiddleware(logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 拦截响应
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw
		startTime := time.Now()
		ctx.Next()
		duration := time.Since(startTime).String()
		logger.WithContext(ctx).Info("Response", zap.Any("response_body", blw.body.String()), zap.Any("time", duration))
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)                  // 记录响应
	return w.ResponseWriter.Write(b) // 正常返回
}
