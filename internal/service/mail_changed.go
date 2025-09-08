package service

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/model"
)

func SendEmailChanged(user *model.User, oldEmail string) {
	switch app.Env.APP_LOCALE {
	case "es":
		sendEmailChangedEs(user, oldEmail)
	default:
		sendEmailChangedEn(user, oldEmail)
	}
}

func sendEmailChangedEs(user *model.User, oldEmail string) {
	revertCode := model.NewVerificationCode()
	if err := revertCode.Generate(user.ID, "email-change-revert", map[string]string{"old_email": oldEmail}); err != nil {
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
        <a href="` + app.Env.APP_URL + `/users/revert-email-change/` + user.ID.Hex() + `/` + revertCode.Code + `" 
           style="display:inline-block;padding:10px 20px;background:#dc3545;color:#fff;
                  text-decoration:none;border-radius:5px;">
           Revertir cambio de correo
        </a>
    </p>
    <p>Si no puedes hacer clic, copia y pega este enlace en tu navegador:</p>
    <p>` + app.Env.APP_URL + `/users/revert-email-change/` + user.ID.Hex() + `/` + revertCode.Code + `</p>
    `

	// Se envía al email ANTIGUO, no al nuevo
	SendMail(subject, body, oldEmail)
}

func sendEmailChangedEn(user *model.User, oldEmail string) {
	// Generate a code to revert the change
	revertCode := model.NewVerificationCode()
	if err := revertCode.Generate(user.ID, "email-change-revert", map[string]string{"old_email": oldEmail}); err != nil {
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
        <a href="` + app.Env.APP_URL + `/users/revert-email-change/` + user.ID.Hex() + `/` + revertCode.Code + `" 
           style="display:inline-block;padding:10px 20px;background:#dc3545;color:#fff;
                  text-decoration:none;border-radius:5px;">
           Revert email change
        </a>
    </p>
    <p>If the button doesn’t work, copy and paste this link into your browser:</p>
    <p>` + app.Env.APP_URL + `/users/revert-email-change/` + user.ID.Hex() + `/` + revertCode.Code + `</p>
    `

	// Send to the OLD email, not the new one
	SendMail(subject, body, oldEmail)
}
