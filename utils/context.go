// Context 数据读取工具
package utils

import "github.com/gin-gonic/gin"

// 获取用户ID
func GetUserID(c *gin.Context) uint {
	if v, ok := c.Get("user_id"); ok {
		if userID, ok := v.(uint); ok {
			return userID
		}
	}
	return 0
}

// 获取用户名
func GetUserName(c *gin.Context) string {
	if v, ok := c.Get("user_name"); ok {
		if name, ok := v.(string); ok {
			return name
		}
	}
	return ""
}
