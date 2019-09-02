package utils

import (
	mail "github.com/go-mail/mail"
	b "github.com/pickjunk/bgo"
)

// Mail send mail in HTML format
func Mail(to []string, title string, content string) error {
	if b.Config.Get("mail").Exists() {
		log.Panic().Str("field", "mail").Msg("config field not found")
	}
	host := b.Config.Get("mail.host").String()
	if host == "" {
		log.Panic().Str("field", "mail.host").Msg("config field not found")
	}
	port := int(b.Config.Get("mail.port").Int())
	if port == 0 {
		log.Panic().Str("field", "mail.port").Msg("config field not found")
	}
	user := b.Config.Get("mail.user").String()
	if user == "" {
		log.Panic().Str("field", "mail.user").Msg("config field not found")
	}
	passwd := b.Config.Get("mail.passwd").String()
	if passwd == "" {
		log.Panic().Str("field", "mail.passwd").Msg("config field not found")
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
	m.SetHeader("From", user)
	m.SetHeader("To", list...)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)

	d := mail.NewDialer(host, port, user, passwd)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	return d.DialAndSend(m)
}
