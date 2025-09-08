package routes

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/controller"
	"github.com/donbarrigon/nuevo-proyecto/internal/middleware"
)

func Migration(r *app.Routes) {
	r.Prefix("db", func() {
		r.Get("/seed", controller.Seed)
		r.Get("/migrate", controller.Migrare)
		r.Get("/migrate/fresh", controller.Fresh)
		r.Get("/migrate/reset", controller.Reset)
		r.Get("/migrate/refresh", controller.Refresh)
		r.Get("/migrate/rollback", controller.Rollback)

	}, middleware.OnlyLocalhost)
}
