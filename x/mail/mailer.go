package mail

import (
	"github.com/google/uuid"
	gomail "gopkg.in/mail.v2"
	"os"
)

// MailerProvider provides a mailer for a given project.
type MailerProvider interface {
	Provide(projectID uuid.UUID) Mailer
}

// Mailer is an interface for sending emails.
type Mailer interface {
	SendMail(from, to, subject, body string) error
}

// SimpleMailer is a mailer provider.
type SimpleMailer struct {
	host string
	port int
	user string
	pass string
}

func NewMailerProvider(host string, port int, user, pass string) *SimpleMailer {
	return &SimpleMailer{
		host: host,
		port: port,
		user: user,
		pass: pass,
	}
}

func (m *SimpleMailer) Provide(orgID uuid.UUID) Mailer {
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
