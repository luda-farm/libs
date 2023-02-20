package mailer

import (
	"fmt"
	"net/smtp"
)

type (
	Mailer struct {
		host     string
		port     int
		username string
		password string
	}

	MailerConfig struct {
		Host     string
		Port     int
		Username string
		Password string
	}

	loginAuth struct {
		username string
		password string
	}
)

func NewMailer(config MailerConfig) Mailer {
	return Mailer{
		host:     config.Host,
		port:     config.Port,
		username: config.Username,
		password: config.Password,
	}
}

func (conf Mailer) Send(subject, body, recipient string, additionalRecipients ...string) error {
	addr := fmt.Sprintf("%s:%d", conf.host, conf.port)
	content := "MIME-version: 1.0\n" + `Content-Type: text/html; charset="UTF-8";` + "\r\n"
	content += fmt.Sprintf("From: %s\r\n", conf.username)
	content += fmt.Sprintf("To: %s\r\n", conf.username)
	content += fmt.Sprintf("Subject: %s\r\n", subject)
	content += fmt.Sprintf("\r\n%s\r\n", body)
	return smtp.SendMail(
		addr, loginAuth{username: conf.username, password: conf.password}, conf.username,
		append(additionalRecipients, recipient), []byte(content),
	)
}

func (auth loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (auth loginAuth) Next(cmd []byte, more bool) ([]byte, error) {
	if !more {
		return nil, nil
	}

	switch string(cmd) {
	case "Username:":
		return []byte(auth.username), nil
	case "Password:":
		return []byte(auth.password), nil
	default:
		return nil, fmt.Errorf("unknown cmd '%s' in login auth", string(cmd))
	}
}
