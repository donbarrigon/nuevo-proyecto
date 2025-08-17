package routes

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/http/controller"
	"github.com/donbarrigon/nuevo-proyecto/internal/http/middleware"
)

func user(r *app.Routes) {

	r.Post("/users", controller.UserStore).
		Name("users.store")

	r.Post("/users/login", controller.Login).
		Name("users.login")

	r.Post("/users/logout", controller.Logout, middleware.Auth).
		Name("users.logout")

	r.Prefix(func() {
		r.Use(func() {

			r.Get("/users/{id}", controller.UserShow).
				Name("users.show")

			r.Patch("/users/{id}", controller.UserUpdate).
				Name("users.update")

			r.Delete("/users/{id}", controller.PermissionDestroy).
				Name("users.destroy")

		}, middleware.Auth)
	}, "dashboard")
}
