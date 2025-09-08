package service

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/server/model"
)

func SendEmailPasswordChanged(user *model.User) {
	switch app.Env.APP_LOCALE {
	case "es":
		sendEmailPasswordChangedEs(user)
	default:
		sendEmailPasswordChangedEn(user)
	}
}

func sendEmailPasswordChangedEs(user *model.User) {
	verificationCode := model.NewVerificationCode()
	if err := verificationCode.Generate(user.ID, "reset-password"); err != nil {
		app.PrintError("Failed to generate verification code", app.E("error", err))
		return
	}

	subject := "Notificación de cambio de contraseña en " + app.Env.APP_NAME

	body := `
    <h1>Hola ` + user.Profile.Nickname + `</h1>
    <p>Queremos informarte que tu contraseña en ` + app.Env.APP_NAME + ` ha sido cambiada exitosamente.</p>
    <p>Si fuiste tú quien realizó este cambio, no necesitas hacer nada más.</p>
    <p>Si <strong>no fuiste tú</strong>, por favor restablece tu contraseña inmediatamente o contacta a nuestro soporte.</p>
    <p>
        <a href="` + app.Env.APP_URL + `/users/reset-password/` + user.ID.Hex() + `/` + verificationCode.Code + `" 
           style="display:inline-block;padding:10px 20px;background:#dc3545;color:#fff;
                  text-decoration:none;border-radius:5px;">
           Restablecer contraseña
        </a>
    </p>
    `

	SendMail(subject, body, user.Email)
}

func sendEmailPasswordChangedEn(user *model.User) {
	verificationCode := model.NewVerificationCode()
	if err := verificationCode.Generate(user.ID, "reset-password"); err != nil {
		app.PrintError("Failed to generate verification code", app.E("error", err))
		return
	}

	subject := "Password Change Notification at " + app.Env.APP_NAME

	body := `
    <h1>Hello ` + user.Profile.Nickname + `</h1>
    <p>We want to let you know that your password at ` + app.Env.APP_NAME + ` has been successfully changed.</p>
    <p>If you made this change, no further action is required.</p>
    <p>If <strong>you did not</strong> make this change, please reset your password immediately or contact our support team.</p>
    <p>
        <a href="` + app.Env.APP_URL + `/users/reset-password/` + user.ID.Hex() + `/` + verificationCode.Code + `" 
           style="display:inline-block;padding:10px 20px;background:#dc3545;color:#fff;
                  text-decoration:none;border-radius:5px;">
           Reset Password
        </a>
    </p>
    `

	SendMail(subject, body, user.Email)
}
