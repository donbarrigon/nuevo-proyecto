package service

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/model"
)

func SendEmailForgotPassword(user *model.User) {
	if app.Env.APP_LOCALE == "es" {
		sendEmailForgotPasswordEs(user)
	} else {
		sendEmailForgotPasswordEn(user)
	}
}

func sendEmailForgotPasswordEs(user *model.User) {
	resetCode := model.NewVerificationCode()
	if err := resetCode.Generate(user.ID, "reset-password", nil); err != nil {
		app.PrintError("Failed to generate reset code", app.E("error", err))
		return
	}

	subject := "Restablece tu contraseña en " + app.Env.APP_NAME

	body := `
    <h1>Hola ` + user.Profile.Nickname + `</h1>
    <p>Recibimos una solicitud para restablecer tu contraseña en ` + app.Env.APP_NAME + `.</p>
    <p>Se creara una nueva contraseña haciendo clic en el siguiente enlace:</p>
    <p>
        <a href="` + app.Env.APP_URL + `/users/reset-password/` + user.ID.Hex() + `/` + resetCode.Code + `" 
           style="display:inline-block;padding:10px 20px;background:#007bff;color:#fff;
                  text-decoration:none;border-radius:5px;">
           Restablecer contraseña
        </a>
    </p>
    <p>Si no solicitaste este cambio, simplemente ignora este mensaje. Tu contraseña seguirá siendo la misma.</p>
    <p>Si no puedes hacer clic, copia y pega este enlace en tu navegador:</p>
    <p>` + app.Env.APP_URL + `/users/reset-password/` + user.ID.Hex() + `/` + resetCode.Code + `</p>
    <br>
    <p>Equipo de ` + app.Env.APP_NAME + `</p>
    `

	SendMail(subject, body, user.Email)
}

func sendEmailForgotPasswordEn(user *model.User) {
	resetCode := model.NewVerificationCode()
	if err := resetCode.Generate(user.ID, "reset-password", nil); err != nil {
		app.PrintError("Failed to generate reset code", app.E("error", err))
		return
	}

	subject := "Reset your password at " + app.Env.APP_NAME

	body := `
    <h1>Hello ` + user.Profile.Nickname + `</h1>
    <p>We received a request to reset your password at ` + app.Env.APP_NAME + `.</p>
    <p>A new password will be created by clicking the link below:</p>
    <p>
        <a href="` + app.Env.APP_URL + `/users/reset-password/` + user.ID.Hex() + `/` + resetCode.Code + `" 
           style="display:inline-block;padding:10px 20px;background:#007bff;color:#fff;
                  text-decoration:none;border-radius:5px;">
           Reset Password
        </a>
    </p>
    <p>If you did not request this change, simply ignore this email. Your current password will remain valid.</p>
    <p>If you cannot click the button, copy and paste this link into your browser:</p>
    <p>` + app.Env.APP_URL + `/users/reset-password/` + user.ID.Hex() + `/` + resetCode.Code + `</p>
    <br>
    <p>The ` + app.Env.APP_NAME + ` Team</p>
    `

	SendMail(subject, body, user.Email)
}
