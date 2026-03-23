// job.go（任务封装）
package scheduler

import (
	"log"
	"time"
)

type Job interface {
	Run()
}

// 通用任务结构
type BaseJob struct {
	Name     string
	Handler  func() error
	Retry    int
	Interval time.Duration
}

func (j *BaseJob) Run() {
	start := time.Now()

	defer func() {
		if err := recover(); err != nil {
			log.Printf("[Job:%s] panic: %v", j.Name, err)
		}
	}()

	log.Printf("[Job:%s] 开始执行", j.Name)

	err := j.execute()

	if err != nil {
		log.Printf("[Job:%s] 执行失败: %v", j.Name, err)
	} else {
		log.Printf("[Job:%s] 执行成功,耗时:%v", j.Name, time.Since(start))
	}
}

func (j *BaseJob) execute() error {
	var err error

	for i := 0; i <= j.Retry; i++ {
		err = j.Handler()
		if err == nil {
			return nil
		}

		log.Printf("[Job:%s] 重试 %d/%d", j.Name, i+1, j.Retry)

		time.Sleep(j.Interval)
	}

	return err
}
