package service

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"strconv"
)

func SendMessageToMail(msg string) {

	// 定义收件人
	mailTo := []string{
		"xxx@xxx.com",
	}
	// 邮件主题为"Hello"
	subject := "基金估值走势分析"

	// 邮件正文
	body := msg

	err := sendMail(mailTo, subject, body)
	if err != nil {
		log.Println(err)
		fmt.Println("send fail")
		return
	}

	fmt.Println("send successfully")
}

func sendMail(mailTo []string, subject string, body string) error {

	//定义邮箱服务器连接信息，如果是网易邮箱 pass填密码，qq邮箱填授权码
	mailConn := map[string]string{
		"user": "xxx@xxx.com",
		"pass": "123456",
		"host": "xxx.xxx.com",
		"port": "25",
	}

	// 转换端口类型为int
	port, _ := strconv.Atoi(mailConn["port"])

	m := gomail.NewMessage()

	// 这种方式可以添加别名，即“XX官方”
	// 说明：如果是用网易邮箱账号发送，以下方法别名可以是中文，如果是qq企业邮箱，以下方法用中文别名，会报错，需要用上面此方法转码
	m.SetHeader("From", m.FormatAddress(mailConn["user"], "XX官方"))
	// 这种方式可以添加别名，即“FB Sample”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code>
	//m.SetHeader("From", "FB Sample"+"<"+mailConn["user"]+">")
	//m.SetHeader("From", mailConn["user"])
	// 发送给多个用户
	m.SetHeader("To", mailTo...)
	// 设置邮件主题
	m.SetHeader("Subject", subject)
	// 设置邮件正文
	m.SetBody("text/html", body)

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)

	return err
}
