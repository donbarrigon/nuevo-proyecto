package routes

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/app/controller"
	"github.com/donbarrigon/nuevo-proyecto/internal/app/middleware"
)

func UserRoutes(w http.ResponseWriter, r *http.Request) {
	ctx := controller.NewUserController(w, r)
	middleware.Auth(ctx, ctx.Store)
	// if r.Method == http.MethodGet {

	// 	if r.URL.Path == "/user/loguot" {
	// 		middleware.Auth(handler.Logout)(ctx)
	// 		return
	// 	}

	// 	middleware.Auth(Show)(ctx)
	// 	return
	// }

	// if r.Method == http.MethodPost {

	// 	if r.URL.Path == "/user/login" {
	// 		Login(ctx)
	// 		return
	// 	}

	// 	Store(ctx)
	// 	return
	// }

	// if r.Method == http.MethodPatch {
	// 	middleware.Auth(Update)(ctx)
	// 	return
	// }

	// if r.Method == http.MethodDelete {
	// 	middleware.Auth(Destroy)(ctx)
	// 	return
	// }

	// http.NotFound(w, r)
}
