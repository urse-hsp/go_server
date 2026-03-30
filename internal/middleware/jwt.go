// JWT 中间件，验证 JWT Token 的有效性
package middleware

import (
	v1 "go-server/api/v1"
	"go-server/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")

		if tokenStr == "" {
			v1.Unauthorized(c, "未登录，请携带 Token")
			c.Abort()
			return
		}

		if !strings.HasPrefix(tokenStr, "Bearer ") {
			v1.Unauthorized(c, "Token 格式错误")
			c.Abort()
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		token, err := jwt.ParseToken(tokenStr)
		if err != nil || !token.Valid {
			v1.Unauthorized(c, "Token 无效或已过期")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*jwt.Claims)
		if !ok {
			v1.Unauthorized(c, "Token 解析失败")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)

		c.Next()
	}
}
