package routes

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app/controller"
)

// path: /Permission/
func permission() {
	Get("/permissions", controller.PermissionIndex)
	Name("permissions.index")

	Get("/permissions/export", controller.PermissionExport)
	Name("permissions.export")

	Get("/permissions/{id}", controller.PermissionShow)
	Name("permissions.show")

	Post("/permissions", controller.PermissionStore)
	Name("permissions.store")

	Patch("/permissions/{id}", controller.PermissionUpdate)
	Name("permissions.update")

	Delete("/permissions/{id}", controller.PermissionDestroy)
	Name("permissions.destroy")

	Patch("/permissions/{id}/restore", controller.PermissionRestore)
	Name("permissions.restore")

	Delete("/permissions/{id}/force-delete", controller.PermissionForceDelete)
	Name("permissions.force-delete")
}
