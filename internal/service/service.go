package service

import (
	"go-server/internal/bootstrap"
	"go-server/internal/dao"
	"go-server/pkg/jwt"
	"go-server/pkg/log"
	"go-server/pkg/sid"
)

// Service 服务
// 小写命名[私有]只能在 service 包内部用
// 负责业务逻辑处理，(调用 Repository/dao 进行数据访问)，调用其他工具包进行辅助功能（如日志、JWT、Sid 等）
type Service struct {
	logger   *log.Logger
	sid      *sid.Sid
	jwt      *jwt.JWT
	tm       dao.Transaction
	rdbCache *bootstrap.RDBCache
}

func NewService(
	tm dao.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
	jwt *jwt.JWT,
	rdbCache *bootstrap.RDBCache,
) *Service {
	return &Service{
		logger:   logger,
		sid:      sid,
		jwt:      jwt,
		tm:       tm,
		rdbCache: rdbCache,
	}
}

// 手动判空赋值
func AssignIfNotNil[T any](dst *T, src *T) {
	if src != nil {
		*dst = *src
	}
}
