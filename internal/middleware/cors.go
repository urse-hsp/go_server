// CORS 中间件 解决“前端跨域请求被浏览器拦截”的问题
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin")) // 允许来源
		c.Header("Access-Control-Allow-Credentials", "true")           // 允许携带 Cookie

		// 处理预检请求
		if method == "OPTIONS" {
			c.Header("Access-Control-Allow-Methods", c.GetHeader("Access-Control-Request-Method")) // 返回允许信息
			c.Header("Access-Control-Allow-Headers", c.GetHeader("Access-Control-Request-Headers"))
			c.Header("Access-Control-Max-Age", "7200") // 缓存 2 小时
			c.AbortWithStatus(http.StatusNoContent)    // 结束请求
			return
		}
		c.Next()
	}
}
