package router

import (
	"go-server/internal/controller"
	"go-server/internal/dao"
	"go-server/internal/middleware"
	"go-server/internal/service"
	"go-server/pkg/jwt"
	"go-server/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type RouterDeps struct {
	Logger     *log.Logger
	Config     *viper.Viper
	Repository *dao.Repository     // dao层工具包
	Service    *service.Service    // 业务层工具包
	Handler    *controller.Handler // 控制层工具包
	JWT        *jwt.JWT
}

func SetupRouter(deps RouterDeps) *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "ok")
	})

	// 全局中间件
	r.Use(
		middleware.CORSMiddleware(),
		middleware.RequestLogMiddleware(deps.Logger),  // 依赖注入日志组件
		middleware.ResponseLogMiddleware(deps.Logger), // 依赖注入日志组件
	)

	api := r.Group("/api")

	// ================= 用户模块 =================
	InitUserRouter(deps, api)

	return r
}
