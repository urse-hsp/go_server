package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取分页参数
func GetPageParams(c *gin.Context) (page int, pageSize int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	return
}
