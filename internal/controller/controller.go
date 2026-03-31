package controller

import (
	"go-server/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// type Handler struct {
// 	logger *log.Logger
// }

// func NewHandler(
// 	logger *log.Logger,
// ) *Handler {
// 	return &Handler{
// 		logger: logger,
// 	}
// }

func GetUserIdFromCtx(ctx *gin.Context) uint {
	v, exists := ctx.Get("claims")
	if !exists {
		return 0
	}
	return v.(*jwt.MyCustomClaims).UserId
}

// // 获取用户ID
// func GetUserID(c *gin.Context) uint {
// 	if v, ok := c.Get("user_id"); ok {
// 		if userID, ok := v.(uint); ok {
// 			return userID
// 		}
// 	}
// 	return 0
// }

// // 获取用户名
// func GetUserName(c *gin.Context) string {
// 	if v, ok := c.Get("user_name"); ok {
// 		if name, ok := v.(string); ok {
// 			return name
// 		}
// 	}
// 	return ""
// }
