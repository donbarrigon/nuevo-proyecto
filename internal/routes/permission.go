package routes

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app/controller"
	"github.com/donbarrigon/nuevo-proyecto/internal/app/middleware"
)

// path: /Permission/
func permission() {
	Get("/dashboard/permissions", controller.PermissionIndex, middleware.Auth)
	Name("dashboard.permissions.index")

	Get("/dashboard/permissions/export", controller.PermissionExport, middleware.Auth)
	Name("dashboard.permissions.export")

	Get("/dashboard/permissions/{id}", controller.PermissionShow, middleware.Auth)
	Name("dashboard.permissions.show")

	Post("/dashboard/permissions", controller.PermissionStore, middleware.Auth)
	Name("dashboard.permissions.store")

	Patch("/dashboard/permissions/{id}", controller.PermissionUpdate, middleware.Auth)
	Name("dashboard.permissions.update")

	Delete("/dashboard/permissions/{id}", controller.PermissionDestroy, middleware.Auth)
	Name("dashboard.permissions.destroy")

	Patch("/dashboard/permissions/{id}/restore", controller.PermissionRestore, middleware.Auth)
	Name("dashboard.permissions.restore")

	Delete("/dashboard/permissions/{id}/force-delete", controller.PermissionForceDelete, middleware.Auth)
	Name("dashboard.permissions.force-delete")
}
