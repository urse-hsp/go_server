package dto

import "go-demo-server/model"

// dto 返回给前端的结构体

// 登录请求 DTO
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 更新请求 DTO
type UserUpdateRequest struct {
	// ID       uint   `json:"id" binding:"required"` // 更新时必须要有 ID
	// Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3"`
	Avatar   string `json:"avatar" binding:"omitempty,url"`
}

// User 响应结构体
type UserDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// User 登录响应
type LoginResponse struct {
	User  UserDTO `json:"user"`
	Token string  `json:"token"`
}

// 转换函数
func ToUserDTO(u *model.User) UserDTO {
	// func ToUserDTO[T model.User](u UserDTO) UserDTO {
	return UserDTO{
		ID:       uint(u.ID),
		Username: u.Username,
	}
}
