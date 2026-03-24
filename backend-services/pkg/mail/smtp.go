package mail

import (
	"fmt"
	"net/smtp"
	"strings"
)

// SendSMTP sends a plain-text email when host is non-empty. Uses LOGIN-capable SMTP (AWS SES, many providers).
func SendSMTP(host, port, user, pass, from string, to []string, subject, body string) error {
	if host == "" || len(to) == 0 {
		return nil
	}
	if port == "" {
		port = "587"
	}
	addr := host + ":" + port
	var auth smtp.Auth
	if user != "" {
		auth = smtp.PlainAuth("", user, pass, host)
	}
	hdr := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n",
		from, strings.Join(to, ","), subject)
	msg := []byte(hdr + body)
	return smtp.SendMail(addr, auth, from, to, msg)
}
