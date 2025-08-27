package routes

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	dbRoutes "github.com/donbarrigon/nuevo-proyecto/internal/database/routes"
	authController "github.com/donbarrigon/nuevo-proyecto/internal/domain/auth/controller"
	authRoutes "github.com/donbarrigon/nuevo-proyecto/internal/domain/auth/routes"
)

func GetApi() *app.Routes {
	r := &app.Routes{}
	r.Prefix("api", func() {
		// aca todas las funciones que crean rutas de la api
		authRoutes.User(r)

	})
	// rutas para las migraciones y seed
	dbRoutes.Migration(r)

	// rutas de testeo estas deben estar es en el frontal pero demomento aca las pongo para testearlas
	r.Get("users/confirm/:id/:code", authController.UserConfirmEmail).
		Name("users.confirm-email")

	r.Get("users/revert-email-change/:id/:code", authController.UserRevertEmail).
		Name("users.revert-email-change")
	return r
}
