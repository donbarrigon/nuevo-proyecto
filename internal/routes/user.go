package routes

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/app/controller"
	"github.com/donbarrigon/nuevo-proyecto/internal/app/middleware"
)

func UserRoutes(w http.ResponseWriter, r *http.Request) {

	ctx := controller.NewContext(w, r)
	if r.Method == http.MethodGet {

		if r.URL.Path == "/user/loguot" {
			middleware.Auth((&controller.User{}).Logout)(ctx)
			return
		}

		middleware.Auth((&controller.User{}).Show)(ctx)
		return
	}

	if r.Method == http.MethodPost {

		if r.URL.Path == "/user/login" {
			(&controller.User{}).Login(ctx)
			return
		}

		(&controller.User{}).Store(ctx)
		return
	}

	if r.Method == http.MethodPatch {
		middleware.Auth((&controller.User{}).Update)(ctx)
		return
	}

	if r.Method == http.MethodDelete {
		middleware.Auth((&controller.User{}).Destroy)(ctx)
		return
	}

	ctx.NotFound()
}
