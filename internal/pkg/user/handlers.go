package user

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/guard"
	"github.com/donbarrigon/nuevo-proyecto/internal/model"
)

func Show(w http.ResponseWriter, r *http.Request) {
	ctx := app.NewHandlerContext(w, r)
	if !guard.AuthToken(ctx) {
		return
	}
	userID := showUserRequest(ctx)
	if userID == "" {
		app.ResponseErrorJSON(ctx.Writer, "No existe", http.StatusNotFound, "No existe")
		return
	}

	userModel := &model.User{}
	if err := app.DB.FindByHexID(userModel, userID); err != nil {
		app.ResponseErrorJSON(ctx.Writer, "No se encontró el usuario: "+err.Error(), http.StatusNotFound, "No existe")
		return
	}
	app.ResponseJSON(ctx.Writer, NewUserStoreResource(userModel, nil), 200)
}

func Store(w http.ResponseWriter, r *http.Request) {
	ctx := app.NewHandlerContext(w, r)
	req := storeRequest(ctx)
	if req == nil {
		return
	}

	userModel, err := storeService(req)
	if err != nil {
		app.ResponseErrorJSON(ctx.Writer, err.Error(), 500, "No se guardo")
		return
	}

	tokenModel, err := tokenStoreService(userModel)
	if err != nil {
		app.ResponseErrorJSON(ctx.Writer, err.Error(), 500, "No se creo el token")
		return
	}

	ur := NewUserStoreResource(userModel, tokenModel)
	app.ResponseJSON(ctx.Writer, ur, 200)
}

func Update(w http.ResponseWriter, r *http.Request) {
	ctx := app.NewHandlerContext(w, r)
	if !guard.AuthToken(ctx) {
		return
	}

	req := updateRequest(ctx)
	if req == nil {
		return
	}

	userModel, err := updateService(req, ctx.User)
	if err != nil {
		app.ResponseErrorJSON(ctx.Writer, err.Error(), 500, "No se modifico")
		return
	}

	ur := NewUserStoreResource(userModel, nil)
	app.ResponseJSON(ctx.Writer, ur, 200)
}

func Delete(w http.ResponseWriter, r *http.Request) {

}

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := app.NewHandlerContext(w, r)
	req := loginRequest(ctx)
	if req == nil {
		message := "Las credenciales no son validas"
		app.ResponseErrorJSON(ctx.Writer, message, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userModel := loginService(req)
	if userModel == nil {
		message := "Las credenciales no son validas"
		app.ResponseErrorJSON(ctx.Writer, message, http.StatusUnauthorized, "Unauthorized")
		return
	}

	tokenModel, err := tokenStoreService(userModel)
	if err != nil {
		app.ResponseErrorJSON(ctx.Writer, err.Error(), 500, "Intentelo de nuevo mas tarde")
		return
	}

	ur := NewUserStoreResource(userModel, tokenModel)
	app.ResponseJSON(ctx.Writer, ur, 200)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	ctx := app.NewHandlerContext(w, r)
	if !guard.AuthToken(ctx) {
		return
	}

	if _, err := app.DB.Destroy(ctx.Token); err != nil {
		app.ResponseErrorJSON(ctx.Writer, err.Error(), 500, "Intentelo de nuevo mas tarde")
		return
	}
	w.WriteHeader(http.StatusOK)
}
