package controller

import (
	"fmt"
	v1 "go-server/api/v1"
	"go-server/pkg/jwt"
	"go-server/pkg/log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger *log.Logger
}

func NewHandler(
	logger *log.Logger,
) *Handler {
	return &Handler{
		logger: logger,
	}
}

// 从请求中解析 uint 参数
func ParseUintParam(c *gin.Context, key string) (uint, error) {
	idStr := c.Param(key)
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("无效的ID")
	}
	return uint(id64), nil
}

// 从上下文中获取用户信息
func GetUserIdFromCtx(ctx *gin.Context) uint {
	v, exists := ctx.Get("claims")
	if !exists {
		return 0
	}
	return v.(*jwt.MyCustomClaims).UserId
}

// path 参数解析工具函数
func GetUintID(c *gin.Context, key string) (uint, bool) {
	idStr := c.Param(key)

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		v1.BadRequest(c, "无效的"+key)
		return 0, false
	}

	return uint(id), true
}
