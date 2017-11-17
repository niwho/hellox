package common

import (
	"gopkg.in/gomail.v2"
)

func SendMail() error {
	m := gomail.NewMessage()
	m.SetHeader("From", "pusic_push@126.com")
	m.SetHeader("To", "wutao@imusics.net", "niwhocn@gmail.com")
	// m.SetAddressHeader("Cc", "pusic_push@imusics.net", "pusic_push@imusics.net")
	m.SetHeader("Subject", "推送每日统计")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
	//	m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.126.com", 465, "pusic_push", "ps1503")

	// Send the email to Bob, Cora and Dan.
	return d.DialAndSend(m)
}
