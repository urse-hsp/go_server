package scheduler

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()

	// * * * * * *
	// | | | | | |
	// 秒 分 时 日 月 周
	// 每 5 秒执行一次
	c.AddFunc("*/5 * * * * *", func() {
		fmt.Println("定时任务执行")
	})

	c.Start()

	select {} // 阻塞主线程
}
