package policy

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
)

func CityViewAny(ctx *app.HttpContext) app.Error {
	return nil
}

func CityView(ctx *app.HttpContext) app.Error {
	return nil
}

func CityCreate(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("create city")
}

func CityUpdate(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("update city")
}

func CityDelete(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("delete city")
}
