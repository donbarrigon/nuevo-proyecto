package handler

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/app/user/model"
	"github.com/donbarrigon/nuevo-proyecto/internal/app/user/resource"
	"github.com/donbarrigon/nuevo-proyecto/internal/core"
	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
	"github.com/donbarrigon/nuevo-proyecto/pkg/validate"
	"golang.org/x/crypto/bcrypt"
)

func show(ctx *core.Context) {

	id := ctx.GetParam()

	user := &model.User{}
	if err := core.Mongo.FindByHexID(user, id); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteJSON(http.StatusOK, resource.NewUserLoginResource(user, nil))
}

func store(ctx *core.Context) {

	user := &model.User{}
	if err := ctx.GetBody(user); err != nil {
		ctx.WriteError(err)
		return
	}

	if err := user.Validate(ctx); err != nil {
		ctx.WriteError(err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.WriteError(&core.Err{
			Status:  http.StatusInternalServerError,
			Message: ctx.TT("No se logro encriptar la contraseña"),
			Err:     err.Error(),
		})
		return
	}
	user.Password = string(hashedPassword)

	if _, err := core.Mongo.Create(user); err != nil {
		ctx.WriteError(err)
		return
	}

	token := model.NewToken()
	token.UserID = user.ID
	if _, err := core.Mongo.Create(token); err != nil {
		ctx.WriteError(err)
		return
	}

	res := resource.NewUserLoginResource(user, token)
	ctx.WriteJSON(http.StatusOK, res)
}

func update(ctx *core.Context) {

	user := &model.User{}
	id := ctx.GetParam()

	if err := core.Mongo.FindByHexID(user, id); err != nil {
		ctx.WriteError(err)
		return
	}

	if err := ctx.GetBody(user); err != nil {
		ctx.WriteError(err)
		return
	}

	if err := user.Validate(ctx); err != nil {
		ctx.WriteError(err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.WriteError(&core.Err{
			Status:  http.StatusInternalServerError,
			Message: ctx.TT("No se logro encriptar la contraseña"),
			Err:     err.Error(),
		})
		return
	}
	user.Password = string(hashedPassword)

	if _, err := core.Mongo.Update(user); err != nil {
		ctx.WriteError(err)
		return
	}

	res := resource.NewUserLoginResource(user, nil)
	ctx.WriteJSON(http.StatusOK, res)
}

func destroy(ctx *core.Context) {

	id := ctx.GetParam()
	user := &model.User{}
	if err := core.Mongo.FindByHexID(user, id); err != nil {
		ctx.WriteError(err)
		return
	}

	if _, err := core.Mongo.Delete(user); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
}

func login(ctx *core.Context) {

	var req map[string]string
	if err := ctx.GetBody(req); err != nil {
		ctx.WriteError(err)
		return
	}

	user := &model.User{}
	if validate.Email(req["user"]) {
		if err := core.Mongo.FindOneByField(user, "email", req["user"]); err != nil {
			ctx.WriteError(&core.Err{
				Status:  http.StatusUnauthorized,
				Message: lang.M(ctx.Lang(), "user.service.unautorized"),
				Err:     lang.M(ctx.Lang(), "user.service.unautorized"),
			})
			return
		}
	} else {
		if err := core.Mongo.FindOneByField(user, "phone", req["user"]); err != nil {
			ctx.WriteError(&core.Err{
				Status:  http.StatusUnauthorized,
				Message: lang.M(ctx.Lang(), "user.service.unautorized"),
				Err:     lang.M(ctx.Lang(), "user.service.unautorized"),
			})
			return
		}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req["password"])); err != nil {
		ctx.WriteError(&core.Err{
			Status:  http.StatusUnauthorized,
			Message: lang.M(ctx.Lang(), "user.service.unautorized"),
			Err:     lang.M(ctx.Lang(), "user.service.unautorized"),
		})
		return
	}

	token := model.NewToken()
	token.UserID = user.ID
	if _, err := core.Mongo.Create(token); err != nil {
		ctx.WriteError(err)
		return
	}
	res := resource.NewUserLoginResource(user, token)
	ctx.WriteJSON(http.StatusOK, res)
}

func logout(ctx *core.Context) {

	if _, err := core.Mongo.Delete(ctx.Token); err != nil {
		ctx.WriteError(err)
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
}
