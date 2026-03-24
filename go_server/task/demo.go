package task

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func demo() {
	// 创建定时任务
	c := cron.New(cron.WithSeconds())

	// 每5秒执行一次
	c.AddFunc("*/5 * * * * *", func() {
		fmt.Println("定时任务执行")
	})

	c.Start()

	select {}
}
