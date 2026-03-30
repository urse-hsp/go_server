# scripts
	go mod tidy 
	更新依赖
	自动补全缺失依赖
	自动写入 go.sum
	清理没用的依赖

 swag init -g cmd/server/main.go // 更新文文档
 go run cmd/server/main.go // 启动服务

 


# 安装依赖
go get github.com/gin-gonic/gin // web框架
go get gorm.io/gorm // 核心库 操作数据库表
go get gorm.io/driver/mysql // 数据库驱动
go get github.com/golang-jwt/jwt/v5
go get github.com/spf13/viper  // 配置文件
go get go.uber.org/zap // 日志


# 文件结构
project
 ├ config // 配置文件
 └ router // 路由
 ├ controller // 控制器
 ├ service // 服务
 ├ dao // 数据访问
 ├ model // 数据模型
 ├ middleware // 中间件
 ├ pkg // 公共包
 ├ utils // 工具包


# 第一个接口 gin
r := gin.Default()
r.GET("/ping", func(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
})
r.Run(":8080")



# gin 取值

// userID := c.Query("user_id") // Query get
// userID := c.PostForm("user_id") // PostForm post
// userID := c.Param("user_id") //  path 路径
// 如果 age 参数不存在，则默认为 "18"
	age := c.DefaultQuery("age", "18") 



# 📏 GORM 的默认规则
GORM 不需要你显式指定表名，它有一套默认的自动映射规则。
GORM 会根据你的 Model 结构体名称，自动推导出数据库的表名。规则非常简单：
结构体名称（复数形式） + 蛇形命名

| 结构体名称 (Model) | 默认推断的表名 (Table) | 解释 |
| :--- | :--- | :--- |
| `User` | `users` | 单数变复数 |
| `BlogPost` | `blog_posts` | 驼峰变蛇形，且变复数 |
| `Employee` | `employees` | 单数变复数 |
| `MyUser` | `my_users` | 保持前缀，变复数 |




# middleware 本质是：服务级逻辑（HTTP 层）
比如：
	JWT 鉴权
	日志
	跨域
	权限控制
这些都是：项目强绑定的