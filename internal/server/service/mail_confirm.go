package service

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/server/model"
)

func SendEmailConfirm(user *model.User) {
	switch app.Env.APP_LOCALE {
	case "es":
		sendEmailConfirmEs(user)
	default:
		sendEmailConfirmEn(user)
	}
}

func sendEmailConfirmEs(user *model.User) {

	verificationCode := model.NewVerificationCode()
	if err := verificationCode.Generate(user.ID, "email-verification"); err != nil {
		app.PrintError("Failed to generate verification code", app.E("error", err))
		return
	}
	subject := "Confirma tu cuenta en " + app.Env.APP_NAME

	body := `
    <h1>Bienvenido a ` + app.Env.APP_NAME + `</h1>
    <p>Gracias por registrarte. Para completar tu registro, haz clic en el siguiente enlace:</p>
    <p>
        <a href="` + app.Env.APP_URL + `/users/confirm/` + user.ID.Hex() + `/` + verificationCode.Code + `" 
           style="display:inline-block;padding:10px 20px;background:#0069d9;color:#fff;
                  text-decoration:none;border-radius:5px;">
           Confirmar mi correo
        </a>
    </p>
    <p>Si no fuiste tú quien se registró, puedes ignorar este mensaje.</p>
    `

	SendMail(subject, body, user.Email)
}

func sendEmailConfirmEn(user *model.User) {
	verificationCode := model.NewVerificationCode()
	if err := verificationCode.Generate(user.ID, "email-verification"); err != nil {
		app.PrintError("Failed to generate verification code", app.E("error", err))
		return
	}

	subject := "Confirm your account on " + app.Env.APP_NAME

	body := `
    <h1>Welcome to ` + app.Env.APP_NAME + `</h1>
    <p>Thank you for signing up. To complete your registration, please click the link below:</p>
    <p>
        <a href="` + app.Env.APP_URL + `/users/confirm/` + user.ID.Hex() + `/` + verificationCode.Code + `" 
           style="display:inline-block;padding:10px 20px;background:#0069d9;color:#fff;
                  text-decoration:none;border-radius:5px;">
           Confirm my email
        </a>
    </p>
    <p>If you did not create this account, you can safely ignore this message.</p>
    `

	SendMail(subject, body, user.Email)
}
