package bootstrap

import (
	"context"
	"go-server/internal/dao"
	"go-server/internal/task"
	"go-server/pkg/server/scheduler"
)

func RegisterTasks(s *scheduler.Scheduler, t *task.Task, r *dao.Repository) {
	// task.Demo()

	userRepository := dao.NewUserRepository(r)
	userTask := task.NewUserTask(t, userRepository)

	job := scheduler.NewJob(
		"user",
		func(ctx context.Context) error {
			userTask.CheckUser(ctx)
			return nil
		},
	)

	s.AddJob("0 04 21 * * *", job)
}
