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
	Query *string `form:"query"`
}
type RequestPageQuery struct {
	v1.PageRequest
	RequestQuery
}

// ================= 响应 DTO =================

// 对外公开（别人能看到）
type PublicDTO struct {
}

// 私有（自己能看到）
type PrivateDTO struct {
}
