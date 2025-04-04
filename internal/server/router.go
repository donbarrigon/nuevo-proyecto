package server

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/app/user"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	//user.Routes(router)
	router.HandleFunc("/user/", user.Routes)

	return router
}
