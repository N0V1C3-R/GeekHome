package utils

import (
	"github.com/go-gomail/gomail"
	"os"
	"strconv"
)

var (
	mailHost string
	email    string
	pwd      string
)

func AdminSendEmail(toAddr, topic, body string) error {
	mailHost = os.Getenv("MAIL_HOST")
	port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		return err
	}
	email = os.Getenv("ADMIN_MAIL_ADDR")
	pwd = os.Getenv("ADMIN_MAIL_PWD")

	d := gomail.NewDialer(mailHost, port, email, pwd)
	m := genMsgObj(toAddr, topic, body)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SysSendEmail(toAddr, topic, body string) error {
	mailHost = os.Getenv("MAIL_HOST")
	port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		return err
	}
	email = os.Getenv("SYS_MAIL_ADDR")
	pwd = os.Getenv("SYS_MAIL_PWD")

	d := gomail.NewDialer(mailHost, port, email, pwd)
	m := genMsgObj(toAddr, topic, body)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func genMsgObj(toAddr, topic, body string) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", email)
	m.SetHeader("To", toAddr)
	m.SetHeader("Subject", topic)
	m.SetBody("text/plain", body)

	return m
}
