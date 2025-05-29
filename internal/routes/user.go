package routes

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app/controller"
	"github.com/donbarrigon/nuevo-proyecto/internal/app/middleware"
)

func user() {
	Get("/users/{id}", controller.UserShow)
	Name("users.show")

	Post("/users", controller.UserStore, middleware.Auth)
	Name("users.store")

	Patch("/users/{id}", controller.UserUpdate, middleware.Auth)
	Name("users.update")

	Delete("/users/{id}", controller.PermissionDestroy, middleware.Auth)
	Name("users.destroy")

	Post("/users/login", controller.Login)
	Name("users.login")

	Post("/users/logout", controller.Logout, middleware.Auth)
	Name("users.logout")
}
