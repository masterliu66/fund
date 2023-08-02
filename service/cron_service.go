package service

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

func StartCron(spec string, cmd func(), channel chan int) {

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
		return
	}

	Cron := cron.New(cron.WithLocation(loc))
	// 执行定时任务
	_, err = Cron.AddFunc(spec, cmd)
	if err != nil {
		fmt.Println(err)
	}

	//启动
	Cron.Start()
	defer Cron.Stop()

	for message := range channel {
		fmt.Println(message)
	}
}
