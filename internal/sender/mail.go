package sender

import (
	"accountbalance/pkg/email"
	"accountbalance/pkg/setting"
)

func setupMail(setting *setting.Config) *email.Mail {
	m := email.NewMail(setting.MailConfig.Host, setting.MailConfig.Port)
	m.Login(setting.MailConfig.UserName, setting.MailConfig.Password)
	return m
}

func SendMail(setting *setting.Config, msg string) error {
	m := setupMail(setting)
	return m.Send(setting.MailConfig.Subject, msg, setting.MailConfig.To)
}
