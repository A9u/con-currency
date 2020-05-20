package service

import (
	"con-currency/config"
	"net/smtp"
	"strings"
)

type Mailer interface {
	Send(to []string, from, subject, body string) (err error)
}

type emailProvider struct{}

func NewMailer() Mailer {
	return &emailProvider{}
}

func (em emailProvider) Send(to []string, from, subject, body string) (err error) {

	provider := config.GetString("mail_provider")
	auth := smtp.PlainAuth("", config.GetString("mail_username"), config.GetString("mail_password"), provider)

	recipients := strings.Join(to, ",")
	msg := []byte("To:" + recipients + "\r\nSubject: " + subject + "\r\n" + body)

	err = smtp.SendMail(config.GetString("mail_provider_port"), auth, from, to, msg)
	return
}
