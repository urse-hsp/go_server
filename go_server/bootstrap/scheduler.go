// bootstrap/scheduler.go
package bootstrap

import "go-demo-server/pkg/httpclient/scheduler"

var Scheduler *scheduler.Scheduler

func InitScheduler() {
	Scheduler = scheduler.NewScheduler()
	Scheduler.Start()
}
