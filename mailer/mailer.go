package mailer

import (
	"fmt"
	"net/smtp"
)

type (
	message struct {
		subject    string
		body       string
		recipients []string
	}

	loginAuth struct {
		username, password string
	}
)

var (
	Host     string
	Port     int
	Username string
	Password string
)

func NewMessage(subject, body, recipient string, additionalRecipients ...string) message {
	return message{
		subject:    subject,
		body:       body,
		recipients: append(additionalRecipients, recipient),
	}
}

func (msg message) Send() error {
	addr := fmt.Sprintf("%s:%d", Host, Port)
	auth := loginAuth{
		username: Username,
		password: Password,
	}
	body := "MIME-version: 1.0\n" + `Content-Type: text/html; charset="UTF-8";` + "\r\n"
	body += fmt.Sprintf("From: %s\r\n", Username)
	body += fmt.Sprintf("To: %s\r\n", Username)
	body += fmt.Sprintf("Subject: %s\r\n", msg.subject)
	body += fmt.Sprintf("\r\n%s\r\n", msg.body)
	return smtp.SendMail(addr, auth, Username, msg.recipients, []byte(body))
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
