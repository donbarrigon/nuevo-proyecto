package routes

import (
	"net/http"
	"strings"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/http/middleware"
)

type ControllerFun func(ctx *app.Context)
type MiddlewareFun func(func(ctx *app.Context)) func(ctx *app.Context)

type Route struct {
	Method     string
	Path       []string
	IsVar      []bool
	Controller ControllerFun
	Middleware []MiddlewareFun
	Name       string
}

var Routes []Route
var routesPrefix string
var namePrefix string
var middlewaresPrefix []MiddlewareFun

func init() {
	Routes = make([]Route, 0)
	// aca todas las funciones que crean rutas
	user()

	Prefix("/dashboard", func() {
		permission()
	}, middleware.Auth)

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

func Prefix(prefix string, callback func(), middlewares ...MiddlewareFun) {

	routesPrefix = prefix
	namePrefix = strings.Replace(prefix, "/", ".", -1)

	if !strings.HasPrefix(routesPrefix, "/") {
		routesPrefix = "/" + routesPrefix
	}
	routesPrefix = strings.TrimSuffix(routesPrefix, "/")

	namePrefix = strings.TrimPrefix(namePrefix, ".")
	if !strings.HasSuffix(namePrefix, ".") {
		namePrefix = namePrefix + "."
	}

	middlewaresPrefix = middlewares

	callback()

	routesPrefix = ""
	namePrefix = ""
}

func registerRoute(method, path string, ctrl ControllerFun, middlewares ...MiddlewareFun) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	path = routesPrefix + path
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
		Middleware: append(middlewaresPrefix, middlewares...),
	})
}

func Name(name string) {
	Routes[len(Routes)-1].Name = namePrefix + name
}
