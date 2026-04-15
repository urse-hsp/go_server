package controller

import (
	v1 "go-server/api/v1"
	demodto "go-server/internal/dto/demo"
	"go-server/internal/service"

	"github.com/gin-gonic/gin"
)

func NewDemoController(handler *Handler, s service.DemoService) *demoController {
	return &demoController{
		Handler: handler,
		Service: s,
	}
}

type demoController struct {
	*Handler
	Service service.DemoService // 依赖注入
}

// ================= 创建 =================

// @Summary **创建
// @Tags DEMO
// @Accept json
// @Produce json
// @Param data body demodto.CreateRequest true "注册参数"
// @Success 201 {object} demodto.UserPrivateDTO
// @Router /demo [post]

func (u *demoController) Create(c *gin.Context) {
	var req demodto.CreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		v1.BadRequest(c, "请求参数错误")
		return
	}

	user, err := u.Service.Create(c, req)
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	v1.Created(c, demodto.ToPrivateDTO(user))
}

// ================= 删除id信息 =================

// @Summary **删除
// @Tags DEMO
// @Produce json
// @Param id path int true "ID"
// @Success 204 {string} string "No Content"
// @Router /demo/{id} [delete]

func (u *demoController) Delete(c *gin.Context) {
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

// @Summary **更新
// @Tags DEMO
// @Accept json
// @Produce json
// @Param data body demodto.UpdateRequest true "更新参数"
// @Success 200 {object} demodto.UserPrivateDTO
// @Router /demo [put]

func (u *demoController) Update(c *gin.Context) {
	id, ok := GetUintID(c, "id")
	if !ok {
		return
	}

	var req demodto.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.BadRequest(c, "请求参数错误"+err.Error())
		return
	}

	user, err := u.Service.Update(c, id, req)
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	v1.Success(c, demodto.ToPrivateDTO(user))
}

// ================= 获取id详情 =================

// @Summary 获取详情
// @Tags DEMO
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} demodto.UserPublicDTO
// @Router /demo/{id} [get]

func (u *demoController) GetDetail(c *gin.Context) {
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
		v1.Success(c, demodto.ToPrivateDTO(user))
	} else {
		v1.Success(c, demodto.ToPublicDTO(user))
	}
}

// ================= 列表 =================

// @Summary **列表
// @Tags DEMO
// @Produce json
// @Param data query demodto.RequestQuery false "查询参数"
// @Success 200 {object} []demodto.UserPublicDTO
// @Router /user [get]

func (u *demoController) GetList(c *gin.Context) {
	var q demodto.RequestQuery

	if err := c.ShouldBindQuery(&q); err != nil {
		v1.BadRequest(c, "参数错误")
		return
	}

	users, err := u.Service.GetList(c, q)
	if err != nil {
		v1.BadRequest(c, err.Error())
		return
	}

	list := demodto.ListToPublic(users)

	v1.Success(c, list)
}

// ================= 分页列表 =================

// @Summary **列表-分页
// @Tags DEMO
// @Produce json
// @Param data query demodto.RequestPageQuery false "查询参数"
// @Success 200 {object} v1.PageResponse
// @Router /demo/lists [get]

func (u *demoController) GetPageList(c *gin.Context) {
	var q demodto.RequestPageQuery

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

	list := demodto.ListToPublic(users)

	v1.List(c, list, int(total), q.Page, q.PageSize)
}
