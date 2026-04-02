package demo

import (
	"go-server/internal/controller"
	"go-server/internal/dao"
	"go-server/internal/middleware"
	"go-server/internal/router"
	"go-server/internal/service"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(deps router.RouterDeps, r *gin.RouterGroup) {
	// ================= 用户模块 =================
	// 初始化依赖
	userRepository := dao.NewUserRepository(deps.Repository)                  // dao
	userService := service.NewUserService(deps.Service, userRepository)       // service
	userController := controller.NewUserController(deps.Handler, userService) // controller

	user := r.Group("/user")
	{
		// ✅ 不需要登录
		user.POST("/login", userController.Login)
		user.POST("/register", userController.Create)

		// ✅ 需要登录
		auth := user.Group("/")
		auth.Use(middleware.StrictAuth(deps.JWT, deps.Logger)) // 依赖注入 JWT 和 Logger
		{
			auth.DELETE("/:id", userController.Delete) // 删除自己
			auth.PUT("/info", userController.Update)   // 修改自己
			auth.GET("/info", userController.Get)      // 当前toekn信息

			auth.GET("/:id", userController.GetDetail)  // 👈 detail
			auth.GET("/", userController.GetList)       // 用户列表
			auth.GET("/lists", userController.GetLists) // 分页用户列表
		}
	}
}
