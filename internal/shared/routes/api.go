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
	r.Prefix("db", func() {
		r.Get("seed", dbController.SeedRun)
		r.Get("seed/tracker", dbController.SeedTracker)

		r.Get("migrate", dbController.Migrare)
		r.Get("migrate/rollback", dbController.Rollback)
	})
	return r
}
