package policy

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
)

func StateViewAny(ctx *app.HttpContext) app.Error {
	return nil
}

func StateView(ctx *app.HttpContext) app.Error {
	return nil
}

func StateCreate(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("create state")
}

func StateUpdate(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("update state")
}

func StateDelete(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("delete state")
}
