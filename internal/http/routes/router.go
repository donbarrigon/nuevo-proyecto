package routes

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
)

func GetApi() *app.Routes {
	r := &app.Routes{}
	r.Prefix("api", func() {
		// aca todas las funciones que crean rutas
		user(r)
		//permission(r)
	})

	return r
}
