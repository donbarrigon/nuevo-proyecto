package routes

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/http/controller"
)

// path: /Permission/
func permission(r *app.Routes) {
	r.Prefix(func() {
		r.Get("/permissions", controller.PermissionIndex).
			Name("permissions.index")

		r.Get("/permissions/export", controller.PermissionExport).
			Name("permissions.export")

		r.Get("/permissions/:id", controller.PermissionShow).
			Name("permissions.show")

		r.Post("/permissions", controller.PermissionStore).
			Name("permissions.store")

		r.Patch("/permissions/:id", controller.PermissionUpdate).
			Name("permissions.update")

		r.Delete("/permissions/:id", controller.PermissionDestroy).
			Name("permissions.destroy")

		r.Patch("/permissions/:id/restore", controller.PermissionRestore).
			Name("permissions.restore")

		r.Delete("/permissions/:id/force-delete", controller.PermissionForceDelete).
			Name("permissions.force-delete")
	}, "dashboard")
}
