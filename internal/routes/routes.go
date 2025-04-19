package routes

import (
	"net/http"
	"strings"

	"github.com/donbarrigon/nuevo-proyecto/internal/app/controller"
)

type ControllerFun func(ctx *controller.Context)
type MiddlewareFun func(func(ctx *controller.Context)) func(ctx *controller.Context)

func Use(function ControllerFun, middlewares ...MiddlewareFun) ControllerFun {
	for i := len(middlewares) - 1; i >= 0; i-- {
		function = middlewares[i](function)
	}
	return function
}

var Routes []Route

type Route struct {
	Method     string
	Path       []string
	IsVar      []bool
	Controller ControllerFun
	Middleware []MiddlewareFun
}

func init() {
	Routes = make([]Route, 0)
	// aca todas las funciones que crean rutas
	permission()

}

func Get(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	registerRoute(http.MethodGet, path, ctrl, middlewares...)
}

func Post(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	registerRoute(http.MethodPost, path, ctrl, middlewares...)
}

func Patch(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	registerRoute(http.MethodPatch, path, ctrl, middlewares...)
}

func Put(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	registerRoute(http.MethodPut, path, ctrl, middlewares...)
}

func Delete(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	registerRoute(http.MethodDelete, path, ctrl, middlewares...)
}

func Options(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	registerRoute(http.MethodOptions, path, ctrl, middlewares...)
}

func Head(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	registerRoute(http.MethodHead, path, ctrl, middlewares...)
}

func registerRoute(method, path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	var pathParts []string
	var isVars []bool

	for _, part := range segments {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			paramName := strings.Trim(part, "{}")
			pathParts = append(pathParts, paramName)
			isVars = append(isVars, true)
		} else {
			pathParts = append(pathParts, part)
			isVars = append(isVars, false)
		}
	}

	Routes = append(Routes, Route{
		Method:     method,
		Path:       pathParts,
		IsVar:      isVars,
		Controller: ctrl,
		Middleware: middlewares,
	})
}
