package services

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
)

func SendAccountActivationEmail(mailTo []string, activationToken string) {

	// Sender data.
	from := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")

	// Receiver email address.
	to := mailTo

	// smtp server configuration.
	smtpHost := os.Getenv("MAIL_HOST")
	smtpPort := os.Getenv("MAIL_PORT")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("./store/template/account_activation_email.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Account Activation Email \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Address string
	}{
		Address: os.Getenv("APP_URL") + "/user/activation/" + activationToken,
	})

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, "semih@me.com", to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
}
