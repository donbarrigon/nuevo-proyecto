package routes

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	dbRoutes "github.com/donbarrigon/nuevo-proyecto/internal/database/routes"
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

	return r
}
