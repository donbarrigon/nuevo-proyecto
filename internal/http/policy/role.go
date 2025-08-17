package policy

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/model"
)

func RoleViewAny(ctx *app.HttpContext, profile *model.Profile) app.Error {
	return ctx.User.Can("view role")
}

func RoleView(ctx *app.HttpContext, role *model.Role) app.Error {
	return ctx.User.Can("view role")
}

func RoleCreate(ctx *app.HttpContext, role *model.Role) app.Error {
	return ctx.User.Can("create role")
}

func RoleUpdate(ctx *app.HttpContext, role *model.Role) app.Error {
	return ctx.User.Can("update role")
}

func RoleDelete(ctx *app.HttpContext, role *model.Role) app.Error {
	return ctx.User.Can("delete role")
}
