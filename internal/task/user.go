package task

import (
	"context"
	"fmt"
	"go-server/internal/dao"
	userdto "go-server/internal/dto/user"
)

type UserTask interface {
	CheckUser(ctx context.Context) error
}

func NewUserTask(
	task *Task,
	userRepo dao.UserRepository,
) UserTask {
	return &userTask{
		userRepo: userRepo,
		Task:     task,
	}
}

type userTask struct {
	userRepo dao.UserRepository
	*Task
}

func (t userTask) CheckUser(ctx context.Context) error {
	// t.logger.Info("CheckUser")
	fmt.Println("定时任务执行")
	user, err := t.userRepo.GetList(ctx, userdto.RequestQuery{})
	if err != nil {
		return err
	}
	fmt.Println(user)

	return nil
}
