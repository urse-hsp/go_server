package controller

import (
	"fmt"
	v1 "go-demo-server/api/v1"
	"go-demo-server/internal/dto"
	"go-demo-server/internal/model"
	"go-demo-server/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// ================= 登录 =================

// @Summary 用户登录
// @Description 输入账号密码获取 token
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body dto.LoginRequest true "登录参数"
// @Success 200 {object} dto.LoginResponse
// @Router /api/user/login [post]
func Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		v1.BadRequest(c, "请求参数错误") // // 如果 JSON 里没传 username 或 password，就报错
		return
	}

	user, token, err := service.Login(req.Username, req.Password)
	if err != nil {
		v1.Unauthorized(c) // 401 用户名或密码不对
		return
	}

	v1.Success(c, dto.LoginResponse{
		Token: token,
		User:  toPrivateDTO(user),
	})
}

// ================= 注册 =================

// @Summary 用户注册
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body dto.LoginRequest true "注册参数"
// @Success 200 {object} dto.UserPrivateDTO
// @Router /api/user/register [post]
func Create(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		v1.BadRequest(c, "请求参数错误")
		return
	}

	user, err := service.Create(req.Username, req.Password)
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	v1.Created(c, toPrivateDTO(user))
}

// ================= 获取当前用户 =================

// @Summary 获取当前用户信息
// @Tags 用户
// @Produce json
// @Success 200 {object} dto.UserPrivateDTO
// @Router /api/user/info [get]
func GetUserInfo(c *gin.Context) {
	userID := v1.GetUserID(c)

	user, err := service.GetUserDetail(userID)
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	v1.Success(c, toPrivateDTO(user))
}

// ================= 更新当前用户 =================

// @Summary 更新当前用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body dto.UserUpdateRequest true "更新参数"
// @Success 200 {object} dto.UserPrivateDTO
// @Router /api/user/info [put]
func UpdateUser(c *gin.Context) {
	var req dto.UserUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		v1.BadRequest(c, "请求参数错误")
		return
	}

	var userModel model.User

	if err := copier.Copy(&userModel, &req); err != nil {
		v1.BadRequest(c, "数据转换错误")
		return
	}

	userID := v1.GetUserID(c)

	user, err := service.UpdateUser(userModel, userID)
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	v1.Success(c, toPrivateDTO(user))
}

// ================= 获取他人用户 =================

// @Summary 获取用户详情
// @Tags 用户
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} dto.UserPublicDTO
// @Router /api/users/{id} [get]
func GetUserDetail(c *gin.Context) {
	idStr := c.Param("id")

	// ID 通常是正整数 → 建议用 ParseUint 并转 uint
	// 普通整数字符串 → Atoi 更简单[支付负数]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		v1.BadRequest(c, "无效的用户ID")
		return
	}

	user, err := service.GetUserDetail(uint(id))
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	currentUserID := v1.GetUserID(c)

	// 权限控制：自己 vs 他人
	if currentUserID == uint(id) {
		v1.Success(c, toPrivateDTO(user))
	} else {
		v1.Success(c, toPublicDTO(user))
	}
}

// ================= 删除用户 =================

// @Summary 删除用户
// @Tags 用户
// @Produce json
// @Param id path int true "用户ID"
// @Success 204 {string} string "No Content"
// @Router /api/user/{id} [delete]
func DeleteUser(c *gin.Context) {
	fmt.Print("删除用户\n")

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		v1.BadRequest(c, "无效的用户ID")
		return
	}

	currentUserID := v1.GetUserID(c)

	// 只允许删除自己（可扩展管理员）
	if currentUserID != uint(id) {
		v1.Forbidden(c, "无权限删除他人")
		return
	}
	fmt.Print(idStr, "9999\n")
	if err := service.DeleteUser(uint(id)); err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	v1.NoContent(c)
}

// ================= 用户列表 =================
// @Summary 用户列表
// @Tags 用户
// @Produce json
// @Success 200 {object} []dto.UserPublicDTO
// @Router /api/users [get]
func GetUserList(c *gin.Context) {
	users, err := service.GetUserList()
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	list := make([]dto.UserPublicDTO, 0)

	for _, u := range users {
		list = append(list, toPublicDTO(&u))
	}

	v1.Success(c, list)
}

// ================= 用户列表 分页 =================

// @Summary 用户列表 分页
// @Tags 用户
// @Produce json
// @Success 200 {object} dto.PageResponse
// @Router /api/users/list [get]
func GetUserLists(c *gin.Context) {
	page, pageSize := v1.GetPage(c)

	users, total, err := service.GetUserLists(page, pageSize)
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	v1.List(c, users, int(total), page, pageSize)
}

// ================= DTO 转换 =================

// 👉 他人可见
func toPublicDTO(u *model.User) dto.UserPublicDTO {
	return dto.UserPublicDTO{
		ID:       u.ID,
		Username: u.Username,
		Avatar:   u.Avatar,
	}
}

// 👉 自己可见
func toPrivateDTO(u *model.User) dto.UserPrivateDTO {
	return dto.UserPrivateDTO{
		ID:       u.ID,
		Username: u.Username,
		Avatar:   u.Avatar,
		// Email:    u.Email,
		// Phone:    u.Phone,
	}
}
