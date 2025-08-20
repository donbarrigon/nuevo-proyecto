package routes

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	auth "github.com/donbarrigon/nuevo-proyecto/internal/domain/auth/routes"
)

func GetApi() *app.Routes {
	r := &app.Routes{}
	r.Prefix("api", func() {
		// aca todas las funciones que crean rutas
		auth.User(r)

	})

	return r
}
