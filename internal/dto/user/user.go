package userdto

import v1 "go-server/api/v1"

// ================= 请求 DTO =================

// 登录 / 注册
type CreateRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 更新用户
type UpdateRequest struct {
	Username *string `json:"username"`
	Avatar   *string `json:"avatar" binding:"omitempty,url"`
	// Email    string `json:"email"`
	// Phone    string `json:"phone"`
}

type RequestQuery struct {
	Query *string `form:"query"`
}
type RequestPageQuery struct {
	v1.PageRequest
	RequestQuery
}

// ================= 响应 DTO =================

// 👉 对外公开（别人能看到）
type PublicDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// 👉 私有（自己能看到）
type PrivateDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

// ================= 登录返回 =================

type LoginResponse struct {
	Token string     `json:"token"`
	User  PrivateDTO `json:"user"`
}
