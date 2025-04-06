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

type User struct{}

func (u *User) Show(ctx *Context) {
	id := ctx.GetParam()

	user := &model.User{}
	if err := db.Mongo.FindByHexID(user, id); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteJSON(http.StatusOK, resource.NewUserLoginResource(user, nil))
}

func (u *User) Store(ctx *Context) {

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
			Message: ctx.TT("No se logro encriptar la contraseña"),
			Err:     err.Error(),
		})
		return
	}
	user.Password = string(hashedPassword)

	if _, err := db.Mongo.Create(user); err != nil {
		ctx.WriteError(err)
		return
	}

	token := model.NewToken(user.ID)
	if _, err := db.Mongo.Create(token); err != nil {
		ctx.WriteError(err)
		return
	}

	res := resource.NewUserLoginResource(user, token)
	ctx.WriteJSON(http.StatusOK, res)
}

func (u *User) Update(ctx *Context) {

	user := &model.User{}
	id := ctx.GetParam()

	if err := db.Mongo.FindByHexID(user, id); err != nil {
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
			Message: ctx.TT("No autorizado"),
			Err:     ctx.TT("No esta autorizado para realizar esta accion"),
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.WriteError(&errors.Err{
			Status:  http.StatusInternalServerError,
			Message: ctx.TT("No se logro encriptar la contraseña"),
			Err:     err.Error(),
		})
		return
	}
	user.Password = string(hashedPassword)

	if _, err := db.Mongo.Update(user); err != nil {
		ctx.WriteError(err)
		return
	}

	res := resource.NewUserLoginResource(user, nil)
	ctx.WriteJSON(http.StatusOK, res)
}

func (u *User) Destroy(ctx *Context) {

	id := ctx.GetParam()
	user := &model.User{}
	if err := db.Mongo.FindByHexID(user, id); err != nil {
		ctx.WriteError(err)
		return
	}

	// esto espara que el usuario solo pueda eliminarse asi mismo
	if user.ID != ctx.User.GetID() {
		ctx.WriteError(&errors.Err{
			Status:  http.StatusUnauthorized,
			Message: ctx.TT("No autorizado"),
			Err:     ctx.TT("No esta autorizado para realizar esta accion"),
		})
		return
	}

	if _, err := db.Mongo.Delete(user); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
}

func (u *User) Login(ctx *Context) {

	var req map[string]string
	if err := ctx.GetBody(req); err != nil {
		ctx.WriteError(err)
		return
	}

	user := &model.User{}
	if validate.Email(req["user"]) {
		if err := db.Mongo.FindOneByField(user, "email", req["user"]); err != nil {
			ctx.WriteError(&errors.Err{
				Status:  http.StatusUnauthorized,
				Message: ctx.TT("No autorizado"),
				Err:     ctx.TT("No autorizado"),
			})
			return
		}
	} else {
		if err := db.Mongo.FindOneByField(user, "phone", req["user"]); err != nil {
			ctx.WriteError(&errors.Err{
				Status:  http.StatusUnauthorized,
				Message: ctx.TT("No autorizado"),
				Err:     ctx.TT("No autorizado"),
			})
			return
		}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req["password"])); err != nil {
		ctx.WriteError(&errors.Err{
			Status:  http.StatusUnauthorized,
			Message: ctx.TT("No autorizado"),
			Err:     ctx.TT("No autorizado"),
		})
		return
	}

	token := model.NewToken(user.ID)
	if _, err := db.Mongo.Create(token); err != nil {
		ctx.WriteError(err)
		return
	}
	res := resource.NewUserLoginResource(user, token)
	ctx.WriteJSON(http.StatusOK, res)
}

func (u *User) Logout(ctx *Context) {

	if ctx.Token == nil {
		ctx.WriteError(&errors.Err{
			Status:  http.StatusUnauthorized,
			Message: ctx.TT("No autorizado"),
			Err:     ctx.TT("Token Invalido"),
		})
		return
	}
	if _, err := db.Mongo.Delete(ctx.Token); err != nil {
		ctx.WriteError(err)
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
}
