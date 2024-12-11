package mail

import (
	"github.com/google/uuid"
	gomail "gopkg.in/mail.v2"
	"os"
)

type MailerProvider interface {
	Provide(orgID uuid.UUID) Mailer
}

type Mailer interface {
	SendMail(from, to, subject, body string) error
}

type Manager struct {
	host string
	port int
	user string
	pass string
}

func NewMailerProvider(host string, port int, user, pass string) *Manager {
	return &Manager{
		host: host,
		port: port,
		user: user,
		pass: pass,
	}
}

func (m *Manager) Provide(orgID uuid.UUID) Mailer {
	return &MailerImpl{
		host: m.host,
		port: m.port,
		user: m.user,
		pass: m.pass,
	}
}

type MailerImpl struct {
	host string
	port int
	user string
	pass string
}

func (m *MailerImpl) SendMail(from, to, subject, body string) error {
	return SendMail(from, to, subject, body)
}

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
