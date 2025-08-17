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
	return ctx.User.Can("create city")
}

func CityUpdate(ctx *app.HttpContext) app.Error {
	return ctx.User.Can("update city")
}

func CityDelete(ctx *app.HttpContext) app.Error {
	return ctx.User.Can("delete city")
}
