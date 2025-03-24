package user

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
)

// func Routes(router *http.ServeMux) {

// 	router.HandleFunc("/user/show", app.HandlerGet(middleware.AuthToken(Show)))
// 	router.HandleFunc("/user/create", app.HandlerGet(Store))
// 	router.HandleFunc("/user/store", app.HandlerPost(Store))
// 	router.HandleFunc("/user/edit", app.HandlerGet(middleware.AuthToken(Update)))
// 	router.HandleFunc("/user/update", app.HandlerPut(middleware.AuthToken(Update)))
// 	router.HandleFunc("/user/delete", app.HandlerDelete(middleware.AuthToken(Delete)))
// 	router.HandleFunc("/user/login", app.HandlerGet(Login))
// 	router.HandleFunc("/user/login/start", app.HandlerPost(Login))
// 	router.HandleFunc("/user/logout", app.HandlerGet(middleware.AuthToken(Logout)))
// 	router.HandleFunc("/user/logout/exit", app.HandlerPost(middleware.AuthToken(Logout)))
// 	router.HandleFunc("/role/store", app.HandlerGet(middleware.AuthToken(StoreRole)))

// }

func Routes(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		if r.URL.Path == "/user/login" {
			Logout(app.NewHandlerContext(w, r))
			return
		}
		Show(app.NewHandlerContext(w, r))
		return
	}
	if r.Method == http.MethodPost {
		if r.URL.Path == "/user/loguot" {
			Logout(app.NewHandlerContext(w, r))
			return
		}
		Store(app.NewHandlerContext(w, r))
		return
	}
	if r.Method == http.MethodPut {
		Update(app.NewHandlerContext(w, r))
		return
	}
	if r.Method == http.MethodDelete {
		Delete(app.NewHandlerContext(w, r))
		return
	}
	http.NotFound(w, r)
}
