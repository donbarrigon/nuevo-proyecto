package routes

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
)

func GetApi() *app.Routes {
	r := &app.Routes{}
	r.Prefix(func() {
		// aca todas las funciones que crean rutas
		user(r)
		//permission(r)
	}, "api")

	return r
}
