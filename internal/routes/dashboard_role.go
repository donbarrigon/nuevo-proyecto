package routes

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/app/controller"
)

func DashboardRole(w http.ResponseWriter, r *http.Request) {
	ctx := controller.NewContext(w, r)
	if ctx.Request.Method == http.MethodGet {
		//
	}
}
