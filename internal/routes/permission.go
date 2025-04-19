package routes

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app/controller"
	"github.com/donbarrigon/nuevo-proyecto/internal/app/middleware"
)

// path: /Permission/
func permission() {
	Get("/dashboard/permissions", controller.PermissionIndex, middleware.Auth)
	Get("/dashboard/permissions/export", controller.PermissionExport, middleware.Auth)
	Get("/dashboard/permissions/{id}", controller.PermissionShow, middleware.Auth)
	Post("/dashboard/permissions", controller.PermissionStore, middleware.Auth)
	Patch("/dashboard/permissions/{id}", controller.PermissionUpdate, middleware.Auth)
	Delete("/dashboard/permissions/{id}", controller.PermissionDestroy, middleware.Auth)
	Patch("/dashboard/permissions/{id}/restore", controller.PermissionRestore, middleware.Auth)
	Delete("/dashboard/permissions/{id}/force-delete", controller.PermissionForceDelete, middleware.Auth)
}
