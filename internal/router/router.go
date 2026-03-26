package router

import (
	"go-demo-server/internal/controller"
	"go-demo-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "ok")
	})

	api := r.Group("/api")

	// ================= 用户模块 =================
	user := api.Group("/user")
	{
		// ✅ 不需要登录
		user.POST("/login", controller.Login)
		user.POST("/register", controller.Create)

		// ✅ 需要登录（只针对“自己”）
		auth := user.Group("/")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.DELETE("/:id", controller.DeleteUser)
			auth.GET("/info", controller.GetUserInfo) // 自己
			auth.PUT("/info", controller.UpdateUser)  // 修改自己
		}
	}

	// ================= 用户资源（别人） =================
	users := api.Group("/users")
	{
		users.GET("", controller.GetUserList)       // 用户列表
		users.GET("/list", controller.GetUserLists) // 分页用户列表

		// 查看别人（是否需要登录看你业务）
		users.GET("/:id", controller.GetUserDetail) // 👈 detail
	}

	return r
}
