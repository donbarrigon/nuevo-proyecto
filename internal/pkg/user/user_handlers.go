package user

import (
	"net/http"
	"strings"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/model"
	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
)

func Show(ctx *app.ControllerContext) {

	params := strings.Split(ctx.Request.URL.Path, "/")
	userID := params[len(params)-1]

	userModel := &model.User{}
	if err := app.Mongo.FindByHexID(userModel, userID); err != nil {
		app.NewErrorJSON(lang.M(ctx.Lang(), "app.not_found"), err.Error(), http.StatusNotFound).WriteResponse(ctx.Writer)
		return
	}
	app.ResponseJSON(ctx.Writer, NewUserLoginResource(userModel, nil), http.StatusOK)
}

func Store(ctx *app.ControllerContext) {

	req := &UserRequest{}
	if err := app.GetRequest(ctx, req, http.MethodPost); err != nil {
		err.WriteResponse(ctx.Writer)
		return
	}

	userModel, tokenModel, err := storeService(ctx, req)
	if err != nil {
		err.WriteResponse(ctx.Writer)
		return
	}

	ulr := NewUserLoginResource(userModel, tokenModel)
	app.ResponseJSON(ctx.Writer, ulr, 200)
}

func Update(ctx *app.ControllerContext) {

	req := &UserRequest{}
	if err := app.GetRequest(ctx, req, http.MethodPost); err != nil {
		err.WriteResponse(ctx.Writer)
		return
	}

	userModel, err := updateService(ctx, req)
	if err != nil {
		err.WriteResponse(ctx.Writer)
		return
	}

	ur := NewUserLoginResource(userModel, nil)
	app.ResponseJSON(ctx.Writer, ur, 200)
}

func Delete(ctx *app.ControllerContext) {

	if err := deleteService(ctx); err != nil {
		err.WriteResponse(ctx.Writer)
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
}

func Login(ctx *app.ControllerContext) {

	req, err := loginRequest(ctx)
	if err != nil {
		err.WriteResponse(ctx.Writer)
		return
	}

	userModel, tokenModel, err := loginService(ctx, req)
	if err != nil {
		err.WriteResponse(ctx.Writer)
		return
	}

	ur := NewUserLoginResource(userModel, tokenModel)
	app.ResponseJSON(ctx.Writer, ur, 200)
}

func Logout(ctx *app.ControllerContext) {

	if _, err := app.Mongo.Destroy(ctx.Token); err != nil {
		app.NewErrorJSON(lang.M(ctx.Lang(), "user.logout.destroy"), err.Error(), 500).WriteResponse(ctx.Writer)
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
}
