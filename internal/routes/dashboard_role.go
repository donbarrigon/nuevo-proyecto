package routes

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/app/controller"
	"github.com/donbarrigon/nuevo-proyecto/internal/app/middleware"
)

func DashboardRole(w http.ResponseWriter, r *http.Request) {
	ctx := controller.NewContext(w, r)
	if ctx.Request.Method == http.MethodGet {
		middleware.Auth(controller.IndexDashboardRole)(ctx)
	}
}
