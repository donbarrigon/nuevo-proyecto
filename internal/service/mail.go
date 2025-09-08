package service

import (
	"net/smtp"
	"strings"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
)

func SendMail(subject string, body string, to ...string) {

	if app.Env.MAIL_USERNAME == "tuemail@gmail.com" {
		app.PrintWarning("Failed to send email: no email configured")
		return
	}

	for i := range 3 {

		auth := smtp.PlainAuth(app.Env.MAIL_IDENTITY, app.Env.MAIL_USERNAME, app.Env.MAIL_PASSWORD, app.Env.MAIL_HOST)

		msg := []byte(
			"From: " + app.Env.MAIL_FROM_NAME + " <" + app.Env.MAIL_USERNAME + ">\r\n" +
				"To: " + strings.Join(to, ",") + "\r\n" +
				"Subject: " + subject + "\r\n" +
				"MIME-Version: 1.0\r\n" +
				"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
				"\r\n" +
				body + "\r\n",
		)

		err := smtp.SendMail(app.Env.MAIL_HOST+":"+app.Env.MAIL_PORT, auth, app.Env.MAIL_USERNAME, to, msg)
		if err != nil {
			app.PrintError("Failed to send email: try :try to: :to error: :error", app.E("error", err), app.E("to", to), app.E("try", i))
			time.Sleep(15 * time.Second)
			continue
		}

		return
	}
}
