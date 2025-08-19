package routes

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/http/controller"
	"github.com/donbarrigon/nuevo-proyecto/internal/http/middleware"
)

func user(r *app.Routes) {

	r.Get("users/:id", controller.UserShow, middleware.Auth).
		Name("users.show")

	r.Get("confirm/:code", controller.UserConfirmEmail).
		Name("users.confirm-email")

	r.Get("revert-email-change/:code", controller.UserRevertEmail).
		Name("users.revert-email-change")

	r.Post("users", controller.UserStore).
		Name("users.store")

	r.Post("users/login", controller.Login).
		Name("users.login")

	r.Post("users/logout", controller.Logout, middleware.Auth).
		Name("users.logout")

	r.Prefix("dashboard", func() {

		r.Get("users", controller.UserIndex).
			Name("users.index")

		r.Get("users/trashed", controller.UserTrashed).
			Name("users.trashed")

		r.Patch("users/:id/profile", controller.UserUpdateProfile).
			Name("users.update-profile")

		r.Patch("users/:id/email", controller.UserUpdateEmail).
			Name("users.update-email")

		r.Patch("users/:id/password", controller.UserUpdatePassword).
			Name("users.update-password")

		r.Delete("users/:id", controller.UserDestroy).
			Name("users.destroy")

		r.Patch("users/:id/restore", controller.UserDestroy).
			Name("users.restore")

	}, middleware.Auth)
}
