package userdto

import (
	v1 "go-server/api/v1"
)

// ================= 请求 DTO =================

// 登录 / 注册
type CreateRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

// 更新用户
type UpdateRequest struct {
	Username *string `json:"username" binding:"omitempty"`   // 用户名
	Avatar   *string `json:"avatar" binding:"omitempty,url"` // 头像
	// Email    string `json:"email"`
	// Phone    string `json:"phone"`
}

type RequestQuery struct {
	Query string `form:"query"`
}

type RequestPageQuery struct {
	v1.PageRequest
	RequestQuery
}

// ================= 响应 DTO =================

type UserPageResponse struct {
	Data []PublicDTO `json:"data"` // 列表
	v1.PageSizeResponse
}

// 👉 对外公开（别人能看到）
type PublicDTO struct {
	ID       uint   `json:"id"`       // ID
	Username string `json:"username"` // 用户名
	Avatar   string `json:"avatar"`   // 头像
}

// 👉 私有（自己能看到）
type PrivateDTO struct {
	PublicDTO
	Email string `json:"email"` // 邮箱
	Phone string `json:"phone"` // 手机号
}

// ================= 登录返回 =================

type LoginResponse struct {
	Token string     `json:"token"` // token
	User  PrivateDTO `json:"user"`  // 用户信息
}
