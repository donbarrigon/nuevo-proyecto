package controller

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/app/model"
	"github.com/donbarrigon/nuevo-proyecto/internal/app/resource"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
	"github.com/donbarrigon/nuevo-proyecto/pkg/errors"
	"github.com/donbarrigon/nuevo-proyecto/pkg/validate"
	"golang.org/x/crypto/bcrypt"
)

func UserShow(ctx *Context) {
	id := ctx.LastParam()

	user := &model.User{}
	if err := db.FindByHexID(user, id); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteJSON(http.StatusOK, resource.NewUserLoginResource(user, nil))
}

func UserStore(ctx *Context) {

	user := &model.User{}
	if err := ctx.GetBody(user); err != nil {
		ctx.WriteError(err)
		return
	}

	if err := user.Validate(ctx.Lang()); err != nil {
		ctx.WriteError(err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.WriteError(&errors.Err{
			Status:  http.StatusInternalServerError,
			Message: "No se logro encriptar la contraseña",
			Err:     err.Error(),
		})
		return
	}
	user.Password = string(hashedPassword)

	if err := db.Create(user); err != nil {
		ctx.WriteError(err)
		return
	}

	token := model.NewToken(user.ID)
	if err := db.Create(token); err != nil {
		ctx.WriteError(err)
		return
	}

	res := resource.NewUserLoginResource(user, token)
	ctx.WriteJSON(http.StatusOK, res)
}

func UserUpdate(ctx *Context) {

	id := ctx.LastParam()

	user := &model.User{}
	if err := db.FindByHexID(user, id); err != nil {
		ctx.WriteError(err)
		return
	}

	if err := ctx.GetBody(user); err != nil {
		ctx.WriteError(err)
		return
	}

	if err := user.Validate(ctx.Lang()); err != nil {
		ctx.WriteError(err)
		return
	}

	// esto es para que el usuario solo pueda modificarse asi mismo
	if user.ID != ctx.User.GetID() {
		ctx.WriteError(&errors.Err{
			Status:  http.StatusUnauthorized,
			Message: "No autorizado",
			Err:     ctx.TT("No esta autorizado para realizar esta accion"),
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.WriteError(&errors.Err{
			Status:  http.StatusInternalServerError,
			Message: "No se logro encriptar la contraseña",
			Err:     err.Error(),
		})
		return
	}
	user.Password = string(hashedPassword)

	if err := db.Update(user); err != nil {
		ctx.WriteError(err)
		return
	}

	res := resource.NewUserLoginResource(user, nil)
	ctx.WriteJSON(http.StatusOK, res)
}

func UserDestroy(ctx *Context) {

	id := ctx.LastParam()
	user := &model.User{}
	if err := db.FindByHexID(user, id); err != nil {
		ctx.WriteError(err)
		return
	}

	// esto espara que el usuario solo pueda eliminarse asi mismo
	if user.ID != ctx.User.GetID() {
		ctx.WriteError(&errors.Err{
			Status:  http.StatusUnauthorized,
			Message: "No autorizado",
			Err:     ctx.TT("No esta autorizado para realizar esta accion"),
		})
		return
	}

	if err := db.Delete(user); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
}

func Login(ctx *Context) {

	var req map[string]string
	if err := ctx.GetBody(req); err != nil {
		ctx.WriteError(err)
		return
	}

	user := &model.User{}
	if validate.Email(req["user"]) {
		if err := db.FindOneByField(user, "email", req["user"]); err != nil {
			ctx.WriteError(&errors.Err{
				Status:  http.StatusUnauthorized,
				Message: "No autorizado",
				Err:     ctx.TT("No autorizado"),
			})
			return
		}
	} else {
		if err := db.FindOneByField(user, "phone", req["user"]); err != nil {
			ctx.WriteError(&errors.Err{
				Status:  http.StatusUnauthorized,
				Message: "No autorizado",
				Err:     ctx.TT("No autorizado"),
			})
			return
		}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req["password"])); err != nil {
		ctx.WriteError(&errors.Err{
			Status:  http.StatusUnauthorized,
			Message: "No autorizado",
			Err:     ctx.TT("No autorizado"),
		})
		return
	}

	token := model.NewToken(user.ID)
	if err := db.Create(token); err != nil {
		ctx.WriteError(err)
		return
	}
	res := resource.NewUserLoginResource(user, token)
	ctx.WriteJSON(http.StatusOK, res)
}

func Logout(ctx *Context) {

	if ctx.Token == nil {
		ctx.WriteError(&errors.Err{
			Status:  http.StatusUnauthorized,
			Message: "No autorizado",
			Err:     ctx.TT("Token Invalido"),
		})
		return
	}
	if err := db.ForceDelete(ctx.Token); err != nil {
		ctx.WriteError(err)
		return
	}
	ctx.WriteNoContent()
}
