package routes

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/app/controller"
	"github.com/donbarrigon/nuevo-proyecto/internal/app/middleware"
)

func User(w http.ResponseWriter, r *http.Request) {

	ctx := controller.NewContext(w, r)

	if r.Method == http.MethodGet {

		if r.URL.Path == "/user/loguot" {
			middleware.Auth(controller.Logout)(ctx)
			return
		}

		middleware.Auth(controller.UserShow)(ctx)
		return
	}

	if r.Method == http.MethodPost {

		if r.URL.Path == "/user/login" {
			controller.Login(ctx)
			return
		}

		controller.UserStore(ctx)
		return
	}

	if r.Method == http.MethodPatch {
		middleware.Auth(controller.UserUpdate)(ctx)
		return
	}

	if r.Method == http.MethodDelete {
		middleware.Auth(controller.UserDestroy)(ctx)
		return
	}

	ctx.WriteNotFound()
}
