package app

import (
	"net/http"
	"strings"
)

type ControllerFun func(ctx *HttpContext)
type MiddlewareFun func(func(ctx *HttpContext)) func(ctx *HttpContext)

type Router struct {
	path        string
	isVar       bool
	index       int
	controller  ControllerFun
	middlewares []MiddlewareFun
	routers     []*Router
}

type RouterData struct {
	Params      map[string]string
	Controller  ControllerFun
	Middlewares []MiddlewareFun
}

type Route struct {
	Path       []string
	IsVar      []bool
	Controller ControllerFun
	Middleware []MiddlewareFun
	Name       string
}

type Routes struct {
	routes      []*Route
	prefixes    []string
	middlewares []MiddlewareFun
}

func NewRoutes() *Routes {
	return &Routes{
		routes:      []*Route{},
		prefixes:    []string{},
		middlewares: []MiddlewareFun{},
	}
}

func (r *Routes) Get(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) *Routes {
	return r.SetRoute(http.MethodGet, path, ctrl, middlewares...)
}

func (r *Routes) Post(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) *Routes {
	return r.SetRoute(http.MethodPost, path, ctrl, middlewares...)
}

func (r *Routes) Patch(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) *Routes {
	return r.SetRoute(http.MethodPatch, path, ctrl, middlewares...)
}

func (r *Routes) Put(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) *Routes {
	return r.SetRoute(http.MethodPut, path, ctrl, middlewares...)
}

func (r *Routes) Delete(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) *Routes {
	return r.SetRoute(http.MethodDelete, path, ctrl, middlewares...)
}

func (r *Routes) Options(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) *Routes {
	return r.SetRoute(http.MethodOptions, path, ctrl, middlewares...)
}

func (r *Routes) Head(path string, ctrl ControllerFun, middlewares ...MiddlewareFun) *Routes {
	return r.SetRoute(http.MethodHead, path, ctrl, middlewares...)
}

func (r *Routes) Prefix(callback func(), prefix ...string) {
	for i, p := range prefix {
		prefix[i] = strings.Trim(p, "/")
	}
	r.prefixes = append(r.prefixes, prefix...)
	callback()
	r.prefixes = r.prefixes[:len(r.prefixes)-len(prefix)]
}

func (r *Routes) Use(callback func(), middlewares ...MiddlewareFun) {
	r.middlewares = append(r.middlewares, middlewares...)
	callback()
	r.middlewares = r.middlewares[:len(r.middlewares)-len(middlewares)]
}

func (r *Routes) Name(name string) {
	prefix := strings.Join(r.prefixes, ".")
	r.routes[len(r.routes)-1].Name = prefix + "." + name
}

func (r *Routes) SetRoute(method string, path string, ctrl ControllerFun, middlewares ...MiddlewareFun) *Routes {
	segments := r.prefixes
	segments = append(segments, strings.Split(strings.Trim(path, "/"), "/")...)

	var pathParts []string
	var isVars []bool

	pathParts = append(pathParts, method)
	isVars = append(isVars, false)

	for _, part := range segments {
		if strings.HasPrefix(part, ":") || (strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}")) {
			part := strings.Trim(part, ":{}")
			pathParts = append(pathParts, part)
			isVars = append(isVars, true)
		} else {
			pathParts = append(pathParts, part)
			isVars = append(isVars, false)
		}
	}

	pathParts = append(pathParts, method)

	newRoute := &Route{
		Path:       pathParts,
		IsVar:      isVars,
		Controller: ctrl,
		Middleware: append(r.middlewares, middlewares...),
	}

	r.routes = append(r.routes, newRoute)
	return r
}

// construlle las rutas optimizadas para luego buscar
// toma el array de rutas y las convierte en ramas
func (r *Router) Make(routes *Routes) {
	r.routers = []*Router{}
	// recorro las rutas
	for _, route := range routes.routes {
		// paso el router padre y la ruta, el indice es cero por que es la raiz
		r.add(0, route)
	}
}

// avansa recursivamente por el array del pat de la ruta y va creando ramas
func (r *Router) add(index int, route *Route) {
	// veo a ver si el path ya existe o es una variable existente
	for i, router := range r.routers {
		// if router.isVar == route.IsVar[index] || router.path == route.Path[index] {
		if router.isVar || router.path == route.Path[index] {
			// si existe, tiro palante con el router que ya existente
			r.routers[i].add(index+1, route)
			return
		}
	}

	// si no existe, creo un nuevo router
	newRouter := &Router{
		path:    route.Path[index],
		isVar:   route.IsVar[index],
		index:   index,
		routers: []*Router{}, // creo el array de rutas vacio para evitar errores en la recursividad
	}

	//agrego el router a la lista del router actual(el padre)
	r.routers = append(r.routers, newRouter)

	// si la ruta aun no termina, sigo adelante con el nuevo router
	if index < len(route.Path) {
		newRouter.add(index+1, route)
		return
	}

	// si la ruta termina, agrego el controlador
	newRouter.controller = route.Controller
	newRouter.middlewares = route.Middleware
}

// func (r *Router) Find(path string, rd *RouterData) {
// 	p := strings.Split(strings.Trim(path, "/"), "/")
// 	r.find(p, rd)
// }

// busca recursivamente por las ramas del router y si encuentra el controlador lo asigna en rd
// si el controlador de rd es nil, significa que no encontro la ruta y se debe manejar como 404
func (r *Router) Find(path []string, rd *RouterData) {
	for _, router := range r.routers {
		if router.isVar || path[router.index] == router.path {
			// si es variable se guarda
			if router.isVar {
				rd.Params[router.path] = path[router.index]
			}

			// si el path aun no termina, sigue adelante y si el siguiente no tiene rutas, se para
			if len(path) > router.index {
				router.Find(path, rd)
				return
			}

			// si es la ultima rama devuelve controlador
			if router.controller != nil {
				rd.Controller = router.controller
				rd.Middlewares = router.middlewares
				return
			}
		}
	}
}

func (router *Router) HandleFunction() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		rd := &RouterData{
			Params: map[string]string{},
		}
		pathSegments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		router.Find(pathSegments, rd)
		if rd.Controller != nil {
			ctx := NewHttpContext(w, r)
			ctx.Params = rd.Params
			router.Use(rd.Controller, rd.Middlewares...)(ctx)
		}
	}
}

func (r *Router) Use(function ControllerFun, middlewares ...MiddlewareFun) ControllerFun {
	for i := len(middlewares) - 1; i >= 0; i-- {
		function = middlewares[i](function)
	}
	return function
}
