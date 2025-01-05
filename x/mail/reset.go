package mail

import (
	"bytes"
	"html/template"
	"os"
)

// ResetPassword sends an email to the user to reset their password.
// The callback contains a code that the user can use to reset their password.
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
