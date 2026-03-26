// 任务封装
package scheduler

import (
	"context"
	"log"
	"sync"
	"time"
)

type Job interface {
	Run()
	Name() string
}

// 基础任务
type BaseJob struct {
	name     string
	handler  func(ctx context.Context) error
	retry    int
	interval time.Duration
	timeout  time.Duration

	mu      sync.Mutex
	running bool
}

// 创建任务信息
func NewJob(name string, handler func(ctx context.Context) error, opts ...Option) *BaseJob {
	job := &BaseJob{
		name:     name,
		handler:  handler,
		retry:    0,
		interval: time.Second,
		timeout:  30 * time.Second,
	}

	for _, opt := range opts {
		opt(job)
	}

	return job
}

func (j *BaseJob) Name() string {
	return j.name
}

func (j *BaseJob) Run() {
	j.mu.Lock()
	if j.running {
		log.Printf("[Job:%s] 跳过（上次未完成）", j.name)
		j.mu.Unlock()
		return
	}
	j.running = true
	j.mu.Unlock()

	defer func() {
		j.mu.Lock()
		j.running = false
		j.mu.Unlock()
	}()

	start := time.Now()

	defer func() {
		if err := recover(); err != nil {
			log.Printf("[Job:%s] panic: %v", j.name, err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), j.timeout)
	defer cancel()

	log.Printf("[Job:%s] 开始执行", j.name)

	err := j.execute(ctx)

	if err != nil {
		log.Printf("[Job:%s] 失败: %v", j.name, err)
	} else {
		log.Printf("[Job:%s] 成功,耗时:%v", j.name, time.Since(start))
	}
}

func (j *BaseJob) execute(ctx context.Context) error {
	var err error

	for i := 0; i <= j.retry; i++ {
		err = j.handler(ctx)
		if err == nil {
			return nil
		}

		log.Printf("[Job:%s] 重试 %d/%d", j.name, i+1, j.retry)

		time.Sleep(j.interval)
	}

	return err
}
