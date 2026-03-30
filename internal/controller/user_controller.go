package controller

import (
	"fmt"
	v1 "go-server/api/v1"
	userdto "go-server/internal/dto/user"
	"go-server/internal/model"
	"go-server/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type UserController interface {
	Login(c *gin.Context)
	Create(c *gin.Context)
	DeleteUser(c *gin.Context)
	GetUserInfo(c *gin.Context)
	UpdateUser(c *gin.Context)
	GetUserList(c *gin.Context)
	GetUserDetail(c *gin.Context)
	GetUserLists(c *gin.Context)
}

type userController struct {
	userService service.UserService // 依赖注入
}

func NewUserController(s service.UserService) UserController {
	return &userController{
		userService: s,
	}
}

//
// var userService = service.NewUserService() // 能用但不推荐，代码写死了，

// ================= 登录 =================

// @Summary 用户登录
// @Description 输入账号密码获取 token
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body userdto.LoginRequest true "登录参数"
// @Success 200 {object} userdto.LoginResponse
// @Router /user/login [post]
func (u *userController) Login(c *gin.Context) {
	var req userdto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		v1.BadRequest(c, "请求参数错误") // // 如果 JSON 里没传 username 或 password，就报错
		return
	}

	user, token, err := u.userService.Login(req.Username, req.Password)
	if err != nil {
		v1.Unauthorized(c, "用户名或密码不对") // 401 用户名或密码不对
		return
	}

	v1.Success(c, userdto.LoginResponse{
		Token: token,
		User:  userdto.ToUserPrivateDTO(user),
	})
}

// ================= 注册 =================

// @Summary 用户注册
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body userdto.LoginRequest true "注册参数"
// @Success 200 {object} userdto.UserPrivateDTO
// @Router /user/register [post]
func (u *userController) Create(c *gin.Context) {
	var req userdto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		v1.BadRequest(c, "请求参数错误")
		return
	}

	user, err := u.userService.Create(req.Username, req.Password)
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	v1.Created(c, userdto.ToUserPrivateDTO(user))
}

// ================= 获取当前用户 =================

// @Summary 获取当前用户信息
// @Tags 用户
// @Produce json
// @Success 200 {object} userdto.UserPrivateDTO
// @Router /user/info [get]
func (u *userController) GetUserInfo(c *gin.Context) {
	userID := v1.GetUserID(c)

	user, err := u.userService.GetUserDetail(userID)
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	v1.Success(c, userdto.ToUserPrivateDTO(user))
}

// ================= 更新当前用户 =================

// @Summary 更新当前用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body userdto.UserUpdateRequest true "更新参数"
// @Success 200 {object} userdto.UserPrivateDTO
// @Router /user/info [put]
func (u *userController) UpdateUser(c *gin.Context) {
	var req userdto.UserUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		v1.BadRequest(c, "请求参数错误"+err.Error())
		return
	}

	var userModel model.User

	if err := copier.Copy(&userModel, &req); err != nil {
		v1.BadRequest(c, "数据转换错误")
		return
	}

	userID := v1.GetUserID(c)

	user, err := u.userService.UpdateUser(userModel, userID)
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	v1.Success(c, userdto.ToUserPrivateDTO(user))
}

// ================= 获取他人用户 =================

// @Summary 获取用户详情
// @Tags 用户
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} userdto.UserPublicDTO
// @Router /users/{id} [get]
func (u *userController) GetUserDetail(c *gin.Context) {
	idStr := c.Param("id")

	// ID 通常是正整数 → 建议用 ParseUint 并转 uint
	// 普通整数字符串 → Atoi 更简单[支付负数]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		v1.BadRequest(c, "无效的用户ID")
		return
	}

	user, err := u.userService.GetUserDetail(uint(id))
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	currentUserID := v1.GetUserID(c)

	// 权限控制：自己 vs 他人
	if currentUserID == uint(id) {
		v1.Success(c, userdto.ToUserPrivateDTO(user))
	} else {
		v1.Success(c, userdto.ToUserPublicDTO(user))
	}
}

// ================= 删除用户 =================

// @Summary 删除用户
// @Tags 用户
// @Produce json
// @Param id path int true "用户ID"
// @Success 204 {string} string "No Content"
// @Router /user/{id} [delete]
func (u *userController) DeleteUser(c *gin.Context) {
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
	if err := u.userService.DeleteUser(uint(id)); err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	v1.NoContent(c)
}

// ================= 用户列表 =================
// @Summary 用户列表
// @Tags 用户
// @Produce json
// @Success 200 {object} []userdto.UserPublicDTO
// @Router /users [get]
func (u *userController) GetUserList(c *gin.Context) {
	users, err := u.userService.GetUserList()
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	list := userdto.UserListToPublic(users)

	v1.Success(c, list)
}

// ================= 用户列表 分页 =================

// @Summary 用户列表 分页
// @Tags 用户
// @Produce json
// @Success 200 {object} v1.PageResponse
// @Router /users/list [get]
func (u *userController) GetUserLists(c *gin.Context) {
	page, pageSize := v1.GetPage(c)

	users, total, err := u.userService.GetUserLists(page, pageSize)
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	list := userdto.UserListToPublic(users)

	v1.List(c, list, int(total), page, pageSize)
}
