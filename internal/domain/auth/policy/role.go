package policy

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/shared/model"
)

func RoleViewAny(ctx *app.HttpContext, profile *model.Profile) app.Error {
	return ctx.Auth.Can("view role")
}

func RoleView(ctx *app.HttpContext, role *model.Role) app.Error {
	return ctx.Auth.Can("view role")
}

func RoleCreate(ctx *app.HttpContext, role *model.Role) app.Error {
	return ctx.Auth.Can("create role")
}

func RoleUpdate(ctx *app.HttpContext, role *model.Role) app.Error {
	return ctx.Auth.Can("update role")
}

func RoleDelete(ctx *app.HttpContext, role *model.Role) app.Error {
	return ctx.Auth.Can("delete role")
}
