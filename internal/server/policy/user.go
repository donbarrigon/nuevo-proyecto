package policy

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/server/model"
)

func UserViewAny(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("view user")

}

func UserView(ctx *app.HttpContext, user *model.User) app.Error {
	if ctx.Auth.GetUserID() == user.ID {
		return nil
	}
	return ctx.Auth.Can("view user")
}

func UserCreate(ctx *app.HttpContext) app.Error {
	return nil
}

func UserUpdate(ctx *app.HttpContext, user *model.User) app.Error {
	if ctx.Auth.GetUserID() == user.ID {
		return nil
	}
	return ctx.Auth.Can("update user")
}

func UserDelete(ctx *app.HttpContext) app.Error {
	return ctx.Auth.Can("delete user")
}
