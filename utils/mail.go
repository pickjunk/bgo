package utils

import (
	mail "github.com/go-mail/mail"
	b "github.com/pickjunk/bgo"
)

// Mail send mail in HTML format
func Mail(to []string, title string, content string) error {
	cfg, ok := b.Config["mail"].(map[string]interface{})
	if !ok {
		b.Log.Panic("config [mail] not found")
	}
	mailHost, ok := cfg["host"].(string)
	if !ok {
		b.Log.Panic("config [mail.host] not found")
	}
	mailPort, ok := cfg["port"].(int)
	if !ok {
		b.Log.Panic("config [mail.port] not found")
	}
	mailUser, ok := cfg["user"].(string)
	if !ok {
		b.Log.Panic("config [mail.user] not found")
	}
	mailPasswd, ok := cfg["passwd"].(string)
	if !ok {
		b.Log.Panic("config [mail.passwd] not found")
	}

	// filter empty
	var list []string
	for _, t := range to {
		if t != "" {
			list = append(list, t)
		}
	}

	if len(list) == 0 {
		return nil
	}

	m := mail.NewMessage()
	m.SetHeader("From", mailUser)
	m.SetHeader("To", list...)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)

	d := mail.NewDialer(mailHost, mailPort, mailUser, mailPasswd)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	return d.DialAndSend(m)
}
