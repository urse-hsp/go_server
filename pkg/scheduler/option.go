// option.go（可选：扩展用）

package scheduler

import "time"

type Option func(*BaseJob)

// 设置重试
func WithRetry(retry int, interval time.Duration) Option {
	return func(j *BaseJob) {
		j.retry = retry
		j.interval = interval
	}
}

// 设置超时
func WithTimeout(timeout time.Duration) Option {
	return func(j *BaseJob) {
		j.timeout = timeout
	}
}
