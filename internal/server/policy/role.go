package policy

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/server/model"
)

func RoleViewAny(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("view role")
}

func RoleView(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("view role")
}

func RoleCreate(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("create role")
}

func RoleUpdate(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("update role")
}

func RoleDelete(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("delete role")
}

func RoleGrant(ctx *app.HttpContext, role *model.Role) app.Error {
	if err := ctx.Auth.Can("grant role"); err != nil {
		return err
	}
	return ctx.Auth.HasRole(role.Name)
}

func RoleRevoke(ctx *app.HttpContext, role *model.Role) app.Error {
	if err := ctx.Auth.Can("revoke role"); err != nil {
		return err
	}
	return ctx.Auth.HasRole(role.Name)
}
