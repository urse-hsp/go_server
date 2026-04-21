.PHONY: help dev server task migration \
        build clean \
        docker-build up down restart logs ps \
        migrate deploy \
        tidy install swag

# ====================
# 变量
# ====================
APP_NAME=demo
DOCKER_COMPOSE_FILE=deploy/docker-compose/docker-compose.yml

# ====================
# 帮助
# ====================
help:
	@echo "make dev        - 本地开发（server + migration）"
	@echo "make server     - 本地运行 server"
	@echo "make task       - 本地运行 task"
	@echo "make migration  - 本地执行 migration"
	@echo ""
	@echo "make build      - 构建所有二进制"
	@echo "make clean      - 清理构建产物"
	@echo ""
	@echo "make docker-build - 构建 Docker 镜像"
	@echo "make up         - 启动 server + task"
	@echo "make down       - 停止服务"
	@echo "make logs       - 查看日志"
	@echo ""
	@echo "make migrate    - 执行数据库迁移"
	@echo "make deploy     - 一键部署"

# ====================
# 本地开发
# ====================
server:
	APP_CONF=config/local.yml go run cmd/server/main.go

task:
	APP_CONF=config/local.yml go run cmd/task/main.go

migration:
	APP_CONF=config/local.yml go run cmd/migration/main.go

dev: migration server

# ====================
# 构建
# ====================
build:
	@echo ">>> building binaries..."
	go build -o bin/server cmd/server/main.go
	go build -o bin/task cmd/task/main.go
	go build -o bin/migration cmd/migration/main.go

clean:
	rm -rf bin/

# ====================
# Docker
# ====================

# 构建镜像
docker-build:
	docker compose -p $(APP_NAME) -f $(DOCKER_COMPOSE_FILE) build
	
# 启动服务
up:
	docker compose -p $(APP_NAME) -f $(DOCKER_COMPOSE_FILE) up -d server task

# 停止服务
down:
	docker compose -p $(APP_NAME) -f $(DOCKER_COMPOSE_FILE) down

# 重启服务
restart:
	docker compose -p $(APP_NAME) -f $(DOCKER_COMPOSE_FILE) restart

logs:
	docker compose -p $(APP_NAME) -f $(DOCKER_COMPOSE_FILE) logs -f

ps:
	docker compose -p $(APP_NAME) -f $(DOCKER_COMPOSE_FILE) ps

# ====================
# Migration（重点）
# ====================
migrate:
	docker compose -p $(APP_NAME) -f $(DOCKER_COMPOSE_FILE) run --rm migration

# ====================
# 一键部署（推荐）
# ====================
deploy: docker-build migrate up
	@echo ">>> deploy success 🚀"
# 构建+启动 docker compose up -d --build 

# ====================
# 依赖管理
# ====================
tidy:
	go mod tidy

install:
	go mod download

# ====================
# Swagger
# ====================
swag:
	swag init -g cmd/server/main.go -d cmd/server,internal,pkg,api