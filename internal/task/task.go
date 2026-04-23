package task

import (
	"go-server/internal/dao"
	"go-server/pkg/log"
	"go-server/pkg/sid"
)

type Task struct {
	logger *log.Logger
	sid    *sid.Sid
	// jwt    *jwt.JWT
	tm dao.Transaction
}

func NewTask(
	tm dao.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
) *Task {
	return &Task{
		logger: logger,
		sid:    sid,
		tm:     tm,
	}
}
