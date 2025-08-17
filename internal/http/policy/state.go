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
	return ctx.User.Can("create state")
}

func StateUpdate(ctx *app.HttpContext) app.Error {
	return ctx.User.Can("update state")
}

func StateDelete(ctx *app.HttpContext) app.Error {
	return ctx.User.Can("delete state")
}
