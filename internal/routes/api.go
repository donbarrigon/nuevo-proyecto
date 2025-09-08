package routes

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	dbroutes "github.com/donbarrigon/nuevo-proyecto/internal/database/routes"
	"github.com/donbarrigon/nuevo-proyecto/internal/server/controller"
)

func GetAll() *app.Routes {
	r := &app.Routes{}
	r.Prefix("api", func() {
		// aca todas las funciones que crean rutas de la api
		user(r)

	})
	// rutas para las migraciones y seed
	dbroutes.Migration(r)

	// rutas de testeo estas deben estar es en el frontal pero las pongo aca para hacer pruebas
	// cuando crees el frontal eliminalas
	r.Get("users/confirm/:id/:code", controller.UserConfirmEmail).
		Name("users.confirm-email")

	r.Get("users/revert-email-change/:id/:code", controller.UserRevertEmail).
		Name("users.revert-email-change")

	r.Get("users/reset-password/:id/:code", controller.UserResetPassword).
		Name("users.reset-password")

	return r
}
