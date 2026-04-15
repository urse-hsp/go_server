package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// 分页默认值兜底
func (p *PageRequest) Normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 || p.PageSize > 100 {
		p.PageSize = 10
	}
}

// 获取分页参数
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
