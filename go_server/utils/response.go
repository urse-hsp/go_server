package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RESTful + 统一错误结构（混合模式）
// 成功（RESTful）
// 失败（统一 response）

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	// Data any    `json:"data,omitempty"`
}

type PageData struct {
	List     any `json:"list"`
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

// ===== 核心出口（唯一 JSON 出口）=====
func writeJSON(c *gin.Context, httpStatus int, res any) {
	c.JSON(httpStatus, res)
}

// 核心方法
func response(c *gin.Context, httpStatus int, code int, msg string, data any) {
	// response(c, http.StatusOK, 0, "", PageData{
	writeJSON(c, httpStatus, Response{
		Code: code,
		Msg:  msg,
		// Data: data,
	})
}

// ===== 成功 =====

// 200
func Success(c *gin.Context, data ...any) {
	// response(c, http.StatusOK, 0, "", data)
	writeJSON(c, http.StatusOK, data)
}

// // 200 带 msg
// func SuccessWithMsg(c *gin.Context, data any, msg string) {
// 	response(c, http.StatusOK, 0, msg, data)
// }

// 201
func Created(c *gin.Context, data any) {
	// response(c, http.StatusCreated, 0, "created", data)
	writeJSON(c, http.StatusCreated, data)
}

// 204（无返回）
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// ===== 错误 =====

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

func Forbidden(c *gin.Context) {
	Fail(c, http.StatusForbidden, 403, "禁止访问")
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

// List 成功返回（分页）
func List(c *gin.Context, list any, total int, page int, pageSize int) {
	writeJSON(c, http.StatusOK, PageData{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
