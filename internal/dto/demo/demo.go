package demodto

import (
	v1 "go-server/api/v1"
)

// ================= 请求 DTO =================
type CreateRequest struct {
	// gorm.Model
	Username string
}

type UpdateRequest struct {
}

type RequestQuery struct {
	Query  string `form:"query"`  // 查询
	Status *int   `form:"status"` // 状态
	Type   *int   `form:"type"`   // 类型
}

type RequestPageQuery struct {
	v1.PageRequest
	RequestQuery
}

// ================= 响应 DTO =================

type PageResponse struct {
	Data []PublicDTO `json:"data"` // 列表
	v1.PageSizeResponse
}

// 对外公开（别人能看到）
type PublicDTO struct {
}

// 私有（自己能看到）
type PrivateDTO struct {
	PublicDTO
}
