package mail

import (
	gomail "gopkg.in/mail.v2"
	"os"
)

func SendMail(from, to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	pass := os.Getenv("EMAIL_API_KEY")

	d := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", pass)
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
