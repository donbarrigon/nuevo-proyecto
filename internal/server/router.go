package server

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/routes"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/user/", routes.UserRoutes)

	return router
}
