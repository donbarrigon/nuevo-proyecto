package policy

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
)

func CountryViewAny(ctx *app.HttpContext) app.Error {
	return nil
}

func CountryView(ctx *app.HttpContext) app.Error {
	return nil
}

func CountryCreate(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("create country")
}

func CountryUpdate(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("update country")
}

func CountryDelete(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("delete country")
}
