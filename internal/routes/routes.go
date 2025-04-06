package routes

import "github.com/donbarrigon/nuevo-proyecto/internal/app/controller"

type ControllerFun func(ctx *controller.Context)
type MiddlewareFun func(ControllerFun) ControllerFun

func Use(function ControllerFun, middlewares ...MiddlewareFun) ControllerFun {
	for i := len(middlewares) - 1; i >= 0; i-- {
		function = middlewares[i](function)
	}
	return function
}
