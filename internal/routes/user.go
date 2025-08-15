package routes

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/http/controller"
	"github.com/donbarrigon/nuevo-proyecto/internal/http/middleware"
)

func user(r *app.Routes) {
	r.Get("/users/{id}", controller.UserShow).
		Name("users.show")

	r.Post("/users/login", controller.Login).
		Name("users.login")

	r.Use(func() {
		r.Post("/users", controller.UserStore, middleware.Auth).
			Name("users.store")

		r.Patch("/users/{id}", controller.UserUpdate, middleware.Auth).
			Name("users.update")

		r.Delete("/users/{id}", controller.PermissionDestroy, middleware.Auth).
			Name("users.destroy")

		r.Post("/users/logout", controller.Logout, middleware.Auth).
			Name("users.logout")
	}, middleware.Auth)

}
