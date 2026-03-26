package dto

// ================= 请求 DTO =================

// 登录 / 注册
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 更新用户
type UserUpdateRequest struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	// Email    string `json:"email"`
	// Phone    string `json:"phone"`
}

// ================= 响应 DTO =================

// 👉 对外公开（别人能看到）
type UserPublicDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// 👉 私有（自己能看到）
type UserPrivateDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

// ================= 登录返回 =================

type LoginResponse struct {
	Token string         `json:"token"`
	User  UserPrivateDTO `json:"user"`
}
