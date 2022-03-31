package mailer

import (
	"fmt"
	"net/smtp"
)

type (
	Message struct {
		Sender     string
		Recipients []string
		Subject    string
		Body       string
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

func (msg Message) Send() error {
	addr := fmt.Sprintf("%s:%d", Host, Port)
	auth := loginAuth{
		username: Username,
		password: Password,
	}
	body := "MIME-version: 1.0\n" + `Content-Type: text/html; charset="UTF-8";` + "\r\n"
	body += fmt.Sprintf("From: %s\r\n", msg.Sender)
	body += fmt.Sprintf("To: %s\r\n", msg.Sender)
	body += fmt.Sprintf("Subject: %s\r\n", msg.Subject)
	body += fmt.Sprintf("\r\n%s\r\n", msg.Body)
	fmt.Println(addr)
	return smtp.SendMail(addr, auth, msg.Sender, msg.Recipients, []byte(body))
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
