package mail

import (
	"bytes"
	"html/template"
	"os"
)

func VerifyEmail(email string, callback string) error {
	t, err := template.ParseFiles("templates/verify.html")
	if err != nil {
		return err
	}

	data := struct {
		Email string
		Link  string
	}{
		Email: email,
		Link:  callback,
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return err
	}

	from := os.Getenv("EMAIL_FROM")
	err = SendMail(from, email, "Verify your email", body.String())
	if err != nil {
		return err
	}

	return nil
}
