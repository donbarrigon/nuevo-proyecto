package routes

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	dbController "github.com/donbarrigon/nuevo-proyecto/internal/database/controller"
	auth "github.com/donbarrigon/nuevo-proyecto/internal/domain/auth/routes"
)

func GetApi() *app.Routes {
	r := &app.Routes{}
	r.Prefix("api", func() {
		// aca todas las funciones que crean rutas
		auth.User(r)

	})
	// rutas para las migraciones y seed
	r.Get("seed/run", dbController.SeedRun)
	r.Get("seed/list", dbController.SeedList)
	r.Get("seed/run/:seed", dbController.SeedForce)

	return r
}
