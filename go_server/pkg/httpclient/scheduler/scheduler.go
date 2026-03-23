// scheduler.go（核心调度器）
package scheduler

import (
	"log"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
}

// 初始化构造
func NewScheduler() *Scheduler {
	c := cron.New(
		cron.WithSeconds(), // 支持秒级
		cron.WithChain(
			cron.Recover(cron.DefaultLogger), // 防 panic
		),
	)

	return &Scheduler{cron: c}
}

// 添加任务
func (s *Scheduler) AddJob(spec string, job Job) error {
	_, err := s.cron.AddFunc(spec, func() {
		job.Run()
	})
	return err
}

// 启动
func (s *Scheduler) Start() {
	log.Println("[Scheduler] 启动")
	s.cron.Start()
}

// 停止
func (s *Scheduler) Stop() {
	log.Println("[Scheduler] 停止")
	ctx := s.cron.Stop()
	<-ctx.Done()
}
