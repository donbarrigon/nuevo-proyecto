package user

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/guard"
	"github.com/donbarrigon/nuevo-proyecto/internal/model"
	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
)

func Show(w http.ResponseWriter, r *http.Request) {
	ctx := app.NewHandlerContext(w, r)
	if err := guard.AuthToken(ctx); err != nil {
		err.WriteResponse(w)
		return
	}

	userID, err := showUserRequest(ctx)
	if err != nil {
		err.WriteResponse(w)
		return
	}

	userModel := &model.User{}
	if err := app.Mongo.FindByHexID(userModel, userID); err != nil {
		app.ResponseErrorJSON(w, err.Error(), http.StatusNotFound, lang.M(ctx.Lang(), "app.not_found"))
		return
	}

	app.ResponseJSON(w, NewUserLoginResource(userModel, nil), http.StatusOK)
}

func Store(w http.ResponseWriter, r *http.Request) {
	ctx := app.NewHandlerContext(w, r)
	req, err := storeRequest(ctx)
	if err != nil {
		err.WriteResponse(w)
		return
	}

	userModel, err := storeService(ctx, req)
	if err != nil {
		err.WriteResponse(w)
		return
	}

	tokenModel, err := tokenStoreService(ctx, userModel)
	if err != nil {
		err.WriteResponse(w)
		return
	}

	ulr := NewUserLoginResource(userModel, tokenModel)
	app.ResponseJSON(ctx.Writer, ulr, 200)
}

func Update(w http.ResponseWriter, r *http.Request) {
	ctx := app.NewHandlerContext(w, r)
	if err := guard.AuthToken(ctx); err != nil {
		err.WriteResponse(w)
		return
	}

	req, err := updateRequest(ctx)
	if err != nil {
		err.WriteResponse(w)
		return
	}

	userModel, err := updateService(ctx, req)
	if err != nil {
		err.WriteResponse(w)
		return
	}

	ur := NewUserLoginResource(userModel, nil)
	app.ResponseJSON(ctx.Writer, ur, 200)
}

func Delete(w http.ResponseWriter, r *http.Request) {

}

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := app.NewHandlerContext(w, r)
	req, err := loginRequest(ctx)
	if err != nil {
		err.WriteResponse(w)
		return
	}

	userModel, err := loginService(ctx, req)
	if err != nil {
		err.WriteResponse(w)
		return
	}

	tokenModel, err := tokenStoreService(ctx, userModel)
	if err != nil {
		err.WriteResponse(w)
		return
	}

	ur := NewUserLoginResource(userModel, tokenModel)
	app.ResponseJSON(ctx.Writer, ur, 200)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	ctx := app.NewHandlerContext(w, r)
	if err := guard.AuthToken(ctx); err != nil {
		err.WriteResponse(w)
		return
	}

	if _, err := app.Mongo.Destroy(ctx.Token); err != nil {
		app.ResponseErrorJSON(ctx.Writer, err.Error(), 500, "Intentelo de nuevo mas tarde")
		return
	}
	w.WriteHeader(http.StatusOK)
}
