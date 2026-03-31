package middleware

// 认证中间件
// JWT 登录校验 + 用户信息注入 + 日志增强

import (
	v1 "go-server/api/v1"
	"go-server/pkg/jwt"
	"go-server/pkg/log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 严格认证（必须登录才能访问）
func StrictAuth(j *jwt.JWT, logger *log.Logger) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		if j == nil {
			logger.Error("jwt not initialized")
			v1.ServerError(ctx, "jwt 未初始化")
			ctx.Abort()
			return
		}

		tokenString := ctx.Request.Header.Get("Authorization") // token 从 header 获取
		if tokenString == "" {
			logger.WithContext(ctx).Warn("No token", zap.Any("data", map[string]interface{}{
				"url":    ctx.Request.URL,
				"params": ctx.Params,
			}))
			v1.Unauthorized(ctx, "未登录，请携带 Token")
			ctx.Abort()
			return
		}

		if !strings.HasPrefix(tokenString, "Bearer ") {
			v1.Unauthorized(ctx, "Token 格式错误")
			ctx.Abort()
			return
		}
		claims, err := j.ParseToken(tokenString)
		if err != nil {
			logger.WithContext(ctx).Error("token error", zap.Any("data", map[string]interface{}{
				"url":    ctx.Request.URL,
				"params": ctx.Params,
			}), zap.Error(err))
			v1.Unauthorized(ctx, "Token 无效或已过期")
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)       // 注入用户信息
		recoveryLoggerFunc(ctx, logger) // 加入日志上下文
		ctx.Next()
	}
}

// 非严格认（登录了就识别你，没登录也能访问）
func NoStrictAuth(j *jwt.JWT, logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")
		if tokenString == "" {
			tokenString, _ = ctx.Cookie("accessToken")
		}
		if tokenString == "" {
			tokenString = ctx.Query("accessToken")
		}
		if tokenString == "" {
			ctx.Next()
			return
		}

		claims, err := j.ParseToken(tokenString)
		if err != nil {
			ctx.Next()
			return
		}

		ctx.Set("claims", claims)
		recoveryLoggerFunc(ctx, logger)
		ctx.Next()
	}
}

func recoveryLoggerFunc(ctx *gin.Context, logger *log.Logger) {
	if userInfo, ok := ctx.MustGet("claims").(*jwt.MyCustomClaims); ok {
		logger.WithValue(ctx, zap.String("UserId", strconv.Itoa(int(userInfo.UserId))))
	}
}
