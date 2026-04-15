package router

import (
	"go-server/internal/controller"
	"go-server/internal/dao"
	"go-server/internal/middleware"
	"go-server/internal/service"

	"github.com/gin-gonic/gin"
)

func InitDemoRouter(deps RouterDeps, r *gin.RouterGroup) {
	// ================= **模块 =================
	// 初始化依赖
	Repository := dao.NewDemoRepository(deps.Repository)              // dao
	Service := service.NewDemoService(deps.Service, Repository)       // service
	Controller := controller.NewDemoController(deps.Handler, Service) // controller

	noAuthRouter := r.Group("/demo")
	// ✅ 不需要登录
	{
	}
	// ✅ 需要登录
	strictAuthRouter := noAuthRouter.Group("").Use(middleware.StrictAuth(deps.JWT, deps.Logger))
	{
		strictAuthRouter.POST("", Controller.Create)       // create
		strictAuthRouter.GET("", Controller.GetPageList)   // get
		strictAuthRouter.PUT("/:id", Controller.Update)    // edit
		strictAuthRouter.DELETE("/:id", Controller.Delete) // delete
		strictAuthRouter.GET("/:id", Controller.GetDetail) // detail
	}
	// // ✅ 不强制登录
	// noStrictAuth := noAuthRouter.Group("").Use(middleware.NoStrictAuth(deps.JWT, deps.Logger))
	// {
	// 	noStrictAuth.GET("/lists", Controller.GetPageList) // 分页列表
	// }
}
