package user

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/middleware"
)

func Routes(router *http.ServeMux) {

	router.HandleFunc("/user/show", app.HandlerGet(middleware.AuthToken(Show)))
	router.HandleFunc("/user/create", app.HandlerGet(Store))
	router.HandleFunc("/user/store", app.HandlerPost(Store))
	router.HandleFunc("/user/edit", app.HandlerGet(middleware.AuthToken(Update)))
	router.HandleFunc("/user/update", app.HandlerPut(middleware.AuthToken(Update)))
	router.HandleFunc("/user/delete", app.HandlerDelete(middleware.AuthToken(Delete)))
	router.HandleFunc("/user/login", app.HandlerGet(Login))
	router.HandleFunc("/user/login/start", app.HandlerPost(Login))
	router.HandleFunc("/user/logout", app.HandlerGet(middleware.AuthToken(Logout)))
	router.HandleFunc("/user/logout/exit", app.HandlerPost(middleware.AuthToken(Logout)))
	router.HandleFunc("/role/store", app.HandlerGet(middleware.AuthToken(StoreRole)))

}
