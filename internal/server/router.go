package server

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/pkg/user"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	user.Routes(router)

	return router
}
