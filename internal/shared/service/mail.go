package service

import (
	"net/smtp"
	"strings"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/shared/model"
)

func SendMail(subject string, body string, to ...string) {

	if app.Env.MAIL_USERNAME == "tuemail@gmail.com" {
		app.PrintWarning("Failed to send email: no email configured")
		return
	}
	auth := smtp.PlainAuth("", app.Env.MAIL_USERNAME, app.Env.MAIL_PASSWORD, app.Env.MAIL_HOST)

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
		app.PrintError("Failed to send email: to: :to :error", app.E("error", err), app.E("to", to))
		return
	}
}

func SendVerificationEmail(user *model.User) {
	switch app.Env.APP_LOCALE {
	case "es":
		sendVerificationEmailEs(user)
	default:
		sendVerificationEmailEn(user)
	}
}

func SendEmailChangeNotification(user *model.User, oldEmail string) {
	switch app.Env.APP_LOCALE {
	case "es":
		sendEmailChangeNotificationEs(user, oldEmail)
	default:
		sendEmailChangeNotificationEn(user, oldEmail)
	}
}

func sendVerificationEmailEs(user *model.User) {

	vericationCode := model.NewVerificationCode()
	if err := vericationCode.Generate(user.ID, "email_verification"); err != nil {
		app.PrintError("Failed to generate verification code", app.E("error", err))
		return
	}
	subject := "Confirma tu cuenta en " + app.Env.APP_NAME

	body := `
    <h1>Bienvenido a ` + app.Env.APP_NAME + `</h1>
    <p>Gracias por registrarte. Para completar tu registro, haz clic en el siguiente enlace:</p>
    <p>
        <a href="` + app.Env.APP_URL + `/confirm/` + vericationCode.Code + `" 
           style="display:inline-block;padding:10px 20px;background:#0069d9;color:#fff;
                  text-decoration:none;border-radius:5px;">
           Confirmar mi correo
        </a>
    </p>
    <p>Si no fuiste tú quien se registró, puedes ignorar este mensaje.</p>
    `

	SendMail(subject, body, user.Email)
}

func sendEmailChangeNotificationEs(user *model.User, oldEmail string) {
	// Generamos un código para revertir el cambio
	revertCode := model.NewVerificationCode()
	if err := revertCode.Generate(user.ID, "email_change_revert", map[string]string{"old_email": oldEmail}); err != nil {
		app.PrintError("Failed to generate revert code", app.E("error", err))
		return
	}

	subject := "Tu correo en " + app.Env.APP_NAME + " ha sido actualizado"

	body := `
    <h1>Hola de nuevo en ` + app.Env.APP_NAME + `</h1>
    <p>Queremos informarte que tu dirección de correo fue actualizada recientemente.</p>
    <p>Nuevo correo: <strong>` + user.Email + `</strong></p>
    <p>Si realizaste este cambio, no necesitas hacer nada.</p>
    <p>Pero si <strong>NO fuiste tú</strong>, puedes revertir el cambio haciendo clic en el siguiente enlace:</p>
    <p>
        <a href="` + app.Env.APP_URL + `/revert-email-change/` + revertCode.Code + `" 
           style="display:inline-block;padding:10px 20px;background:#dc3545;color:#fff;
                  text-decoration:none;border-radius:5px;">
           Revertir cambio de correo
        </a>
    </p>
    <p>Si no puedes hacer clic, copia y pega este enlace en tu navegador:</p>
    <p>` + app.Env.APP_URL + `/revert-email-change?token=` + revertCode.Code + `</p>
    `

	// Se envía al email ANTIGUO, no al nuevo
	SendMail(subject, body, oldEmail)
}

func sendVerificationEmailEn(user *model.User) {
	verificationCode := model.NewVerificationCode()
	if err := verificationCode.Generate(user.ID, "email_verification"); err != nil {
		app.PrintError("Failed to generate verification code", app.E("error", err))
		return
	}

	subject := "Confirm your account on " + app.Env.APP_NAME

	body := `
    <h1>Welcome to ` + app.Env.APP_NAME + `</h1>
    <p>Thank you for signing up. To complete your registration, please click the link below:</p>
    <p>
        <a href="` + app.Env.APP_URL + `/confirm/` + verificationCode.Code + `" 
           style="display:inline-block;padding:10px 20px;background:#0069d9;color:#fff;
                  text-decoration:none;border-radius:5px;">
           Confirm my email
        </a>
    </p>
    <p>If you did not create this account, you can safely ignore this message.</p>
    `

	SendMail(subject, body, user.Email)
}

func sendEmailChangeNotificationEn(user *model.User, oldEmail string) {
	// Generate a code to revert the change
	revertCode := model.NewVerificationCode()
	if err := revertCode.Generate(user.ID, "email_change_revert", map[string]string{"old_email": oldEmail}); err != nil {
		app.PrintError("Failed to generate revert code", app.E("error", err))
		return
	}

	subject := "Your email address on " + app.Env.APP_NAME + " has been updated"

	body := `
    <h1>Hello again from ` + app.Env.APP_NAME + `</h1>
    <p>We want to let you know that your email address was recently updated.</p>
    <p>New email: <strong>` + user.Email + `</strong></p>
    <p>If you made this change, no further action is required.</p>
    <p>But if <strong>you did NOT make this change</strong>, you can revert it by clicking the link below:</p>
    <p>
        <a href="` + app.Env.APP_URL + `/revert-email-change/` + revertCode.Code + `" 
           style="display:inline-block;padding:10px 20px;background:#dc3545;color:#fff;
                  text-decoration:none;border-radius:5px;">
           Revert email change
        </a>
    </p>
    <p>If the button doesn’t work, copy and paste this link into your browser:</p>
    <p>` + app.Env.APP_URL + `/revert-email-change?token=` + revertCode.Code + `</p>
    `

	// Send to the OLD email, not the new one
	SendMail(subject, body, oldEmail)
}
