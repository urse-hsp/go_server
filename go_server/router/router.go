package router

import (
	"go-demo-server/controller"
	"go-demo-server/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 注册路由
	r.GET("/", func(c *gin.Context) {
		c.String(200, "ok")
		// c.JSON(200, "ok")
		// c.JSON(200, []string{"a", "b", "c"})
		// c.JSON(200, gin.H{
		// 	"msg": "ok",
		// })
	})

	// // 第一版
	api := r.Group("/api")
	// api.Use(middleware.AuthMiddleware())
	// {
	// 	api.POST("/user/login", controller.Login)
	// 	api.POST("/user/create", controller.Create)
	// 	// api.GET("/user/info", controller.GetUserInfo)
	// 	api.GET("/users/info", controller.GetUserInfo)
	// }

	{
		// 缓存接口测试 用户列表 https://randomuser.me/api/
		api.GET("/user/list", controller.GetUserList)
	}

	// ✅ 公开接口（不需要登录）
	auth := api.Group("/user")
	{
		auth.POST("/login", controller.Login)
		auth.POST("/register", controller.Create)
	}

	// ✅ 需要登录
	user := api.Group("/user")
	user.Use(middleware.AuthMiddleware())
	{
		user.DELETE("/:id", controller.DeleteUser)
		user.GET("/detail", controller.GetUserDetail)
		user.GET("/info", controller.GetUserInfo)
		user.POST("/info", controller.UpdateUser)
	}

	return r
}
