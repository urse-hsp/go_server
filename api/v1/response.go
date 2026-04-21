package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RESTful + 统一错误结构（混合模式）
// 成功（RESTful）
// 失败（统一 response）

// ===== 统一出口 =====
func writeJSON(c *gin.Context, httpStatus int, res any) {
	c.JSON(httpStatus, res)
}

// ===== 核心响应 =====
func response(c *gin.Context, httpStatus int, code int, msg string, data any) {
	writeJSON(c, httpStatus, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// =================  成功 =================

// 200
func Success(c *gin.Context, data any) {
	writeJSON(c, http.StatusOK, data)
}

// 201
func Created(c *gin.Context, data any) {
	writeJSON(c, http.StatusCreated, data)
}

// 204（无返回）
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// =================  错误 =================

func Fail(c *gin.Context, httpStatus int, code int, msg string) {
	response(c, httpStatus, code, msg, nil)
}

// 常用错误快捷方法
func BadRequest(c *gin.Context, msg string) {
	Fail(c, http.StatusBadRequest, 400, msg)
}

func Unauthorized(c *gin.Context, msg ...string) {
	defaultMsg := "未授权"
	if len(msg) > 0 {
		defaultMsg = msg[0]
	}
	Fail(c, http.StatusUnauthorized, 401, defaultMsg)
}

func Forbidden(c *gin.Context, msg ...string) {
	defaultMsg := "禁止访问"
	if len(msg) > 0 {
		defaultMsg = msg[0]
	}
	Fail(c, http.StatusForbidden, 403, defaultMsg)
}

func NotFound(c *gin.Context) {
	Fail(c, http.StatusNotFound, 404, "资源不存在")
}

func ServerError(c *gin.Context, msg ...string) {
	defaultMsg := "服务器错误"
	if len(msg) > 0 {
		defaultMsg = msg[0]
	}
	Fail(c, http.StatusInternalServerError, 500, defaultMsg)
}

// =================  分页 =================

// List 成功返回（分页）
func List[T any](c *gin.Context, data []T, total int, page int, pageSize int) {
	writeJSON(c, http.StatusOK, PageResponse[T]{
		Data: data,
		PageSizeResponse: PageSizeResponse{
			Total: total,
			PageRequest: PageRequest{
				Page:     page,
				PageSize: pageSize,
			},
		},
	})
}
