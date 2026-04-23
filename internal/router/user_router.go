package router

import (
	"go-server/internal/client"
	"go-server/internal/controller"
	"go-server/internal/dao"
	"go-server/internal/middleware"
	"go-server/internal/service"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(deps RouterDeps, r *gin.RouterGroup) {
	// ================= 用户模块 =================
	// 初始化依赖
	userRepository := dao.NewUserRepository(deps.Repository)
	userClient := client.NewUserClient()                                            // dao
	userService := service.NewUserService(deps.Service, userRepository, userClient) // service
	userController := controller.NewUserController(deps.Handler, userService)       // controller

	user := r.Group("/user")
	// ✅ 不需要登录
	{
		user.POST("/login", userController.Login)
		user.POST("/register", userController.Create)
		// redis缓存http接口 用户列表 https://randomuser.me/api/
		user.GET("/http", userController.HttpUserList)

	}
	// ✅ 需要登录
	strictAuthRouter := user.Group("").Use(middleware.StrictAuth(deps.JWT, deps.Logger))
	{
		strictAuthRouter.DELETE("/:id", userController.Delete) // 删除
		strictAuthRouter.PUT("/info", userController.Update)   // 修改
		strictAuthRouter.GET("/info", userController.Get)      // 当前toekn信息
	}
	// ✅ 不强制登录
	noStrictAuth := user.Group("").Use(middleware.NoStrictAuth(deps.JWT, deps.Logger))
	{
		noStrictAuth.GET("/lists", userController.GetPageList) // 分页列表
		noStrictAuth.GET("", userController.GetList)           // 列表
		noStrictAuth.GET("/:id", userController.GetDetail)     // detail
	}
}
