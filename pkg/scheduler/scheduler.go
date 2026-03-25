// scheduler.go（调度器）
package scheduler

import (
	"log"
	"sync"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
	mu   sync.Mutex
	jobs map[string]cron.EntryID
}

// 初始化 scheduler[定时任务]
func NewScheduler() *Scheduler {
	c := cron.New(
		cron.WithSeconds(),
		cron.WithChain(
			cron.Recover(cron.DefaultLogger),
		),
	)

	return &Scheduler{
		cron: c,
		jobs: make(map[string]cron.EntryID),
	}
}

// 添加任务
func (s *Scheduler) AddJob(spec string, job Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, err := s.cron.AddFunc(spec, func() {
		job.Run()
	})
	if err != nil {
		return err
	}

	s.jobs[job.Name()] = id
	log.Printf("[Scheduler] 添加任务: %s (%s)", job.Name(), spec)
	return nil
}

// 删除任务
func (s *Scheduler) RemoveJob(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id, ok := s.jobs[name]; ok {
		s.cron.Remove(id)
		delete(s.jobs, name)
		log.Printf("[Scheduler] 删除任务: %s", name)
	}
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
