package main

import (
	"fmt"
	"fund/dao"
	"fund/model"
	"fund/service"
	"fund/template"
	"fund/web"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"
)

func main() {

	fmt.Println("start……")

	run()
}

func run() {
	channel1 := make(chan int)
	channel2 := make(chan int)

	defer func(Db *sqlx.DB) {
		err := Db.Close()
		if err != nil {
			fmt.Println(err)
		}
		close(channel1)
		close(channel2)
	}(dao.Db)

	go service.StartCron("30 14 * * *", sendFundsStrategyToMail, channel1)
	go service.StartCron("30 23 * * *", insertFunds, channel2)

	router := web.NewRouter()

	s := &http.Server{
		Addr:           ":8000",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	if nil != err {
		fmt.Println("server error", err)
	}
}

func sendFundsStrategyToMail() {

	reports, err := service.CalFundsStrategy(model.Funds)
	if err != nil {
		fmt.Println(err)
		return
	}

	service.SendMessageToMail(template.GetFundReportTemplate(reports))
}

func insertFunds() {
	service.InsertFunds(model.Funds)
}

func insertHistoryFunds() {

	newFunds := []string{
		"501009", // 生物科技
		"501090", // 消费龙头
		"519671", // 300价值
		"162412", // 中证医疗
		"501050", // 50AH
		"161017", // 中证500
		"090010", // 中证红利
		"100032", // 中证红利
		"110003", // 上证50
		"012348", // 恒生科技
		"006327", // 中国互联
		"164906", // 海外互联c
	}

	service.InsertFundHistory(newFunds)
}
