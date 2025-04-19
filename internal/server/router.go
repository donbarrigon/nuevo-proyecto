package server

import (
	"net/http"
	"strings"

	"github.com/donbarrigon/nuevo-proyecto/internal/app/controller"
	"github.com/donbarrigon/nuevo-proyecto/internal/routes"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", HandleFunction)
	return router
}

func HandleFunction(w http.ResponseWriter, r *http.Request) {
	ctx := controller.NewContext(w, r)
	params := map[string]string{}
	pathSegments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	pathLength := len(pathSegments)
	for _, route := range routes.Routes {

		if route.Method != r.Method {
			continue
		}

		if pathLength != len(route.Path) {
			continue
		}

		isThePath := true
		for i := 0; i < pathLength; i++ {
			if !route.IsVar[i] && route.Path[i] != pathSegments[i] {
				isThePath = false
				break
			}
		}

		if !isThePath {
			continue
		}

		for i := 0; i < pathLength; i++ {
			if route.IsVar[i] {
				params[route.Path[i]] = pathSegments[i]
			}
		}

		ctx.PathParams = params
		Use(route.Controller, route.Middleware...)(ctx)
		return
	}
	ctx.WriteNotFound()
}

func Use(function routes.ControllerFun, middlewares ...routes.MiddlewareFun) routes.ControllerFun {
	for i := len(middlewares) - 1; i >= 0; i-- {
		function = middlewares[i](function)
	}
	return function
}
