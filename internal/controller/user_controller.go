package controller

import (
	"fmt"
	v1 "go-server/api/v1"
	userdto "go-server/internal/dto/user"
	"go-server/internal/service"

	"github.com/gin-gonic/gin"
)

func NewUserController(handler *Handler, s service.UserService) *userController {
	return &userController{
		Handler: handler,
		Service: s,
	}
}

type userController struct {
	*Handler
	Service service.UserService // 依赖注入
}

// ================= 登录 =================

// @Summary 用户登录
// @Description 输入账号密码获取 token
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body userdto.CreateRequest true "登录参数"
// @Success 200 {object} userdto.LoginResponse
// @Router /user/login [post]
func (u *userController) Login(c *gin.Context) {
	var req userdto.CreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		v1.BadRequest(c, "请求参数错误") // // 如果 JSON 里没传 username 或 password，就报错
		return
	}

	user, token, err := u.Service.Login(c, req.Username, req.Password)
	if err != nil {
		v1.Unauthorized(c, "用户名或密码不对") // 401 用户名或密码不对
		return
	}

	v1.Success(c, userdto.LoginResponse{
		Token: token,
		User:  userdto.ToPrivateDTO(user),
	})
}

// ================= 创建 =================

// @Summary 用户注册
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body userdto.CreateRequest true "注册参数"
// @Success 201 {object} userdto.UserPrivateDTO
// @Router /user/register [post]
func (u *userController) Create(c *gin.Context) {
	var req userdto.CreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		v1.BadRequest(c, "请求参数错误")
		return
	}

	user, err := u.Service.Create(c, req)
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	v1.Created(c, userdto.ToPrivateDTO(user))
}

// ================= 删除id信息 =================

// @Summary 删除用户
// @Tags 用户
// @Produce json
// @Param id path int true "用户ID"
// @Success 204 {string} string "No Content"
// @Router /user/{id} [delete]
func (u *userController) Delete(c *gin.Context) {
	id, ok := GetUintID(c, "id")
	if !ok {
		return
	}

	currentUserID := GetUserIdFromCtx(c)

	// 只允许删除自己（可扩展管理员）
	if currentUserID != uint(id) {
		v1.Forbidden(c, "无权限删除他人")
		return
	}

	if err := u.Service.Delete(c, uint(id)); err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	v1.NoContent(c)
}

// ================= 更新当前id信息 =================

// @Summary 更新当前用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body userdto.UpdateRequest true "更新参数"
// @Success 200 {object} userdto.UserPrivateDTO
// @Router /user/info [put]
func (u *userController) Update(c *gin.Context) {
	userID := GetUserIdFromCtx(c)

	var req userdto.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.BadRequest(c, "请求参数错误"+err.Error())
		return
	}

	user, err := u.Service.Update(c, userID, req)
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	v1.Success(c, userdto.ToPrivateDTO(user))
}

// ================= 获取token信息 =================

// @Summary 获取当前用户信息
// @Tags 用户
// @Produce json
// @Success 200 {object} userdto.UserPrivateDTO
// @Router /user/info [get]
func (u *userController) Get(c *gin.Context) {
	// userID := v1.GetUserID(c)
	userID := GetUserIdFromCtx(c)
	fmt.Print(userID, "userID")

	user, err := u.Service.GetDetail(c, userID)
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	v1.Success(c, userdto.ToPrivateDTO(user))
}

// ================= 获取id详情 =================

// @Summary 获取用户详情
// @Tags 用户
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} userdto.UserPublicDTO
// @Router /user/{id} [get]
func (u *userController) GetDetail(c *gin.Context) {
	id, ok := GetUintID(c, "id")
	if !ok {
		return
	}

	user, err := u.Service.GetDetail(c, uint(id))
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	currentUserID := GetUserIdFromCtx(c)

	// 权限控制：自己 vs 他人
	if currentUserID == uint(id) {
		v1.Success(c, userdto.ToPrivateDTO(user))
	} else {
		v1.Success(c, userdto.ToPublicDTO(user))
	}
}

// ================= 列表 =================

// @Summary 用户列表
// @Tags 用户
// @Produce json
// @Param data query userdto.RequestQuery false "查询参数"
// @Success 200 {object} []userdto.UserPublicDTO
// @Router /user [get]
func (u *userController) GetList(c *gin.Context) {
	var q userdto.RequestQuery

	if err := c.ShouldBindQuery(&q); err != nil {
		v1.BadRequest(c, "参数错误")
		return
	}

	users, err := u.Service.GetList(c, q)
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	list := userdto.ListToPublic(users)

	v1.Success(c, list)
}

// ================= 分页列表 =================

// @Summary 用户列表-分页
// @Tags 用户
// @Produce json
// @Param data query userdto.RequestPageQuery false "查询参数"
// @Success 200 {object} v1.PageResponse
// @Router /user/lists [get]
func (u *userController) GetPageList(c *gin.Context) {
	var q userdto.RequestPageQuery

	if err := c.ShouldBindQuery(&q); err != nil {
		v1.BadRequest(c, "参数错误")
		return
	}

	q.Normalize()
	users, total, err := u.Service.GetPageList(c, q)
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	list := userdto.ListToPublic(users)

	v1.List(c, list, int(total), q.Page, q.PageSize)
}
