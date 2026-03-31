package service

import (
	"go-server/internal/bootstrap"
	"go-server/pkg/jwt"
	"go-server/pkg/log"
	"go-server/pkg/sid"
)

// Service 服务
// 小写命名[私有]只能在 service 包内部用
type Service struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     bootstrap.Transaction
}

func NewService(
	tm bootstrap.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
	jwt *jwt.JWT,
) *Service {
	return &Service{
		logger: logger,
		sid:    sid,
		jwt:    jwt,
		tm:     tm,
	}
}
