.PHONY: server migration task build dev tidy swag install

# ====================
# 基础运行
# ====================

server:
	go run cmd/server/main.go

migration:
	go run cmd/migration/main.go

task:
	go run cmd/task/main.go

# ====================
# 构建
# ====================

build:
	go build -o app cmd/server/main.go

# ====================
# 开发流程
# ====================

dev: migration server

# ====================
# 依赖管理
# ====================

tidy:
	go mod tidy

install:
	go mod download

# ====================
# 文档
# ====================

swag:
	swag init -g main.go -d cmd/server,internal,pkg,api
# 	swag init -g cmd/server/main.go



# make server      # 启动服务
# make migration   # 执行迁移
# make task        # 启动任务
# make dev         # 一键开发（先迁移再启动）
# make tidy        # 清理依赖
# make swag        # 生成接口文档
# make build       # 打包