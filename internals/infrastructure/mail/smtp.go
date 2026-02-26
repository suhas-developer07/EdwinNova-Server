package mail

import (
	"fmt"
	"net/smtp"
	"os"
)

type SMTPClient struct {
	host string
	port string
	auth smtp.Auth
	from string
}

func NewSMTPClient() (*SMTPClient, error) {

	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")

	if host == "" || port == "" || user == "" || pass == "" {
		return nil, fmt.Errorf("missing smtp environment variables")
	}

	auth := smtp.PlainAuth("", user, pass, host)

	return &SMTPClient{
		host: host,
		port: port,
		auth: auth,
		from: user,
	}, nil
}

func (s *SMTPClient) Send(to, subject, body string) error {

	msg := []byte(
		"Subject: " + subject + "\r\n" +
			"MIME-version: 1.0;\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
			body,
	)

	return smtp.SendMail(
		s.host+":"+s.port,
		s.auth,
		s.from,
		[]string{to},
		msg,
	)
}