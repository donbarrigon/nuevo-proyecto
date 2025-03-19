package user

import "net/http"

func Routes(router *http.ServeMux) {

	router.HandleFunc("/user/show", Show)
	router.HandleFunc("/user/store", Store)
	router.HandleFunc("/user/update", Update)
	router.HandleFunc("/user/delete", Delete)
	router.HandleFunc("/user/login", Login)
	router.HandleFunc("/user/logout", Logout)

}
