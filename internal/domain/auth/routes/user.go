package routes

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/domain/auth/controller"
	"github.com/donbarrigon/nuevo-proyecto/internal/shared/middleware"
)

func User(r *app.Routes) {

	r.Get("users/:id", controller.UserShow, middleware.Auth).
		Name("users.show")

	r.Post("users", controller.UserStore).
		Name("users.store")

	r.Post("users/login", controller.Login).
		Name("users.login")

	r.Post("users/logout", controller.Logout, middleware.Auth).
		Name("users.logout")

	r.Post("users/forgot-password", controller.UserForgotPassword).
		Name("users.forgot-password")

	r.Patch("users/confirm/:id/:code", controller.UserConfirmEmail).
		Name("users.confirm-email")

	r.Patch("users/revert-email-change/:id/:code", controller.UserRevertEmail).
		Name("users.revert-email-change")

	r.Patch("users/reset-password/:id/:code", controller.UserResetPassword).
		Name("users.reset-password")

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

		r.Put("users/:id/restore", controller.UserDestroy).
			Name("users.restore")

			//roles
		r.Get("roles", controller.RoleIndex).
			Name("roles.index")

		r.Get("roles/trashed", controller.RoleTrashed).
			Name("roles.trashed")

		r.Get("roles/:id", controller.RoleShow).
			Name("roles.show")

		r.Post("roles", controller.RoleStore).
			Name("roles.store")

		r.Patch("roles/:id", controller.RoleUpdate).
			Name("roles.update")

		r.Delete("roles/:id", controller.RoleDestroy).
			Name("roles.destroy")

		r.Put("roles/:id/restore", controller.RoleRestore).
			Name("roles.restore")

		r.Patch("roles/:id/grant", controller.RoleGrant).
			Name("roles.grant")

		r.Patch("roles/:id/revoke", controller.RoleRevoke).
			Name("roles.revoke")

		//permissions
		r.Get("permissions", controller.PermissionIndex).
			Name("permissions.index")

		r.Get("permissions/trashed", controller.PermissionTrashed).
			Name("permissions.trashed")

		r.Get("permissions/:id", controller.PermissionShow).
			Name("permissions.show")

		r.Post("permissions", controller.PermissionStore).
			Name("permissions.store")

		r.Patch("permissions/:id", controller.PermissionUpdate).
			Name("permissions.update")

		r.Delete("permissions/:id", controller.PermissionDestroy).
			Name("permissions.destroy")

		r.Put("permissions/:id/restore", controller.PermissionRestore).
			Name("permissions.restore")

		r.Patch("permissions/:id/grant", controller.PermissionGrant).
			Name("permissions.grant")

		r.Patch("permissions/:id/revoke", controller.PermissionRevoke).
			Name("permissions.revoke")
	}, middleware.Auth)
}
