package policy

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/model"
)

func PermissionViewAny(ctx *app.HttpContext, profile *model.Profile) app.Error {
	return ctx.User.Can("view permission")
}

func PermissionView(ctx *app.HttpContext, permission *model.Permission) app.Error {
	return ctx.User.Can("view permission")
}

func PermissionCreate(ctx *app.HttpContext, permission *model.Permission) app.Error {
	return ctx.User.Can("create permission")
}

func PermissionUpdate(ctx *app.HttpContext, permission *model.Permission) app.Error {
	return ctx.User.Can("update permission")
}

func PermissionDelete(ctx *app.HttpContext, permission *model.Permission) app.Error {
	return ctx.User.Can("delete permission")
}
