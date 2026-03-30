package router

import (
	"go-server/internal/controller"
	"go-server/internal/dao"
	"go-server/internal/middleware"
	"go-server/internal/service"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(r *gin.RouterGroup) {
	// ================= 用户模块 =================
	user := r.Group("/user")

	userRepository := dao.NewUserRepository()                   // dao
	userService := service.NewUserService(userRepository)       // service
	userController := controller.NewUserController(userService) // controller
	{
		// ✅ 不需要登录
		user.POST("/login", userController.Login)
		user.POST("/register", userController.Create)

		// ✅ 需要登录（只针对“自己”）
		auth := user.Group("/")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.DELETE("/:id", userController.DeleteUser)
			auth.GET("/info", userController.GetUserInfo) // 自己
			auth.PUT("/info", userController.UpdateUser)  // 修改自己
		}
	}

	// ================= 用户资源（别人） =================
	users := r.Group("/users")
	{
		users.GET("", userController.GetUserList)       // 用户列表
		users.GET("/list", userController.GetUserLists) // 分页用户列表

		// 查看别人（是否需要登录看你业务）
		users.GET("/:id", userController.GetUserDetail) // 👈 detail
	}
}
