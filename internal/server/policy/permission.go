package policy

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/server/model"
)

func PermissionViewAny(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("view permission")
}

func PermissionView(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("view permission")
}

func PermissionCreate(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("create permission")
}

func PermissionUpdate(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("update permission")
}

func PermissionDelete(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("delete permission")
}

func PermissionGrant(ctx *app.HttpContext, permission *model.Permission) app.Error {
	if err := ctx.Auth.Can("grant permission"); err != nil {
		return err
	}
	return ctx.Auth.Can(permission.Name)
}

func PermissionRevoke(ctx *app.HttpContext, permission *model.Permission) app.Error {
	if err := ctx.Auth.Can("revoke permission"); err != nil {
		return err
	}
	return ctx.Auth.Can(permission.Name)
}
