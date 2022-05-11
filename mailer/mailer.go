package mailer

import (
	"fmt"
	"net/smtp"
)

type (
	loginAuth struct{}
)

var (
	Host     string
	Port     int
	Username string
	Password string
)

func Send(subject, body, recipient string, additionalRecipients ...string) error {
	addr := fmt.Sprintf("%s:%d", Host, Port)
	content := "MIME-version: 1.0\n" + `Content-Type: text/html; charset="UTF-8";` + "\r\n"
	content += fmt.Sprintf("From: %s\r\n", Username)
	content += fmt.Sprintf("To: %s\r\n", Username)
	content += fmt.Sprintf("Subject: %s\r\n", subject)
	content += fmt.Sprintf("\r\n%s\r\n", body)
	return smtp.SendMail(
		addr, loginAuth{}, Username,
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
		return []byte(Username), nil
	case "Password:":
		return []byte(Password), nil
	default:
		return nil, fmt.Errorf("unknown cmd '%s' in login auth", string(cmd))
	}
}
