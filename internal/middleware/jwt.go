// middleware/jwt.go
package middleware

import (
	v1 "go-demo-server/api/v1"
	"go-demo-server/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// // 白名单
		// whiteList := []string{
		// 	"/api/user/login",
		// 	"/api/user/register",
		// }
		// // 检查当前请求路径是否在白名单中
		// for _, path := range whiteList {
		// 	if c.Request.URL.Path == path {
		// 		c.Next() // 是白名单，直接放行，不验证 Token
		// 		return
		// 	}
		// }

		// 1. 获取 Token
		tokenStr := c.GetHeader("Authorization")

		// 2. 简单的格式检查 (通常格式是 "Bearer <token>")
		if tokenStr == "" {
			v1.Unauthorized(c, "未登录，请携带 Token")
			c.Abort() // 终止后续流程
			return
		}

		// 去除 "Bearer " 前缀 (可选，看前端怎么传)
		parts := strings.SplitN(tokenStr, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenStr = parts[1]
		}

		// 3. 核心拦截步骤：验证 Token 有效性
		token, err := jwt.ParseToken(tokenStr)
		if err != nil || !token.Valid {
			v1.Unauthorized(c, "Token 无效或已过期")
			c.Abort() // 验证失败，终止后续流程
			return
		}

		// 4. 验证通过：将用户信息存入 Context
		// 这样后续的 Handler 就可以通过 c.Get("user_id") 拿到当前登录用户
		claims := token.Claims.(*jwt.Claims)
		c.Set("user_id", claims.UserID)
		// c.Set("user_name", claims.UserName)

		// 5. 放行：继续执行后续的 Handler
		c.Next()
	}
}
