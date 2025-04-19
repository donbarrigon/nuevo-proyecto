package server

import (
	"net/http"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/", HandleFunc)
	return router
}

func HandleFunc(w http.ResponseWriter, r *http.Request) {
	// ctx := controller.NewContext(w, r)
	// for _, route := range routes.Routes {

	// }
}
