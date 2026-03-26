package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPage(c *gin.Context) (page int, pageSize int) {
	page = 1
	pageSize = 10

	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}

	if ps := c.Query("pageSize"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	return
}
