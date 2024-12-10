package mail

import (
	"bytes"
	"html/template"
	"os"
)

func ResetPassword(email string, callback string) error {
	t, err := template.ParseFiles("templates/reset.html")
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
	err = SendMail(from, email, "Reset you password", body.String())
	if err != nil {
		return err
	}

	return nil
}
