package routes

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
)

func GetApi() {
	r := &app.Routes{}
	// aca todas las funciones que crean rutas
	r.Prefix(func() {
		user(r)
		permission(r)
	}, "api")

}
