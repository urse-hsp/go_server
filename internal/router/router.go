package router

import (
	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

// type RouterDeps struct {
// 	Logger      *log.Logger
// 	Config      *viper.Viper
// 	JWT         *jwt.JWT
// }

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "ok")
	})

	// 全局中间件
	r.Use(
		middleware.CORSMiddleware(),
		// middleware.RequestLogMiddleware(bootstrap.Logger), // 依赖注入日志组件
	)

	api := r.Group("/api")

	// ================= 用户模块 =================
	InitUserRouter(api)

	return r
}
