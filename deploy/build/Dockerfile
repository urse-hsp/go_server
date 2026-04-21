# 构建阶段
FROM golang:1.26 AS builder

WORKDIR /app
COPY . .

ARG APP_RELATIVE_PATH

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ${APP_RELATIVE_PATH}

# 运行阶段
FROM alpine:latest

WORKDIR /app

ARG APP_CONF
ENV APP_CONF=${APP_CONF}

COPY --from=builder /app/app .
COPY --from=builder /app/config ./config
COPY --from=builder /app/cmd/migration/sql ./cmd/migration/sql


CMD ["./app"]


# docker build -t shop-server:v1 -f deploy/build/Dockerfile --build-arg APP_CONF=config/prod.yml  --build-arg APP_RELATIVE_PATH=./cmd/server .
# docker run -d -p 8080:8080 shop-server:v1

