package controller

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
	. "github.com/donbarrigon/nuevo-proyecto/internal/database/db/qb"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/model"
	"github.com/donbarrigon/nuevo-proyecto/internal/http/policy"
	"github.com/donbarrigon/nuevo-proyecto/internal/http/resource"
	"github.com/donbarrigon/nuevo-proyecto/internal/http/validator"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

func UserIndex(ctx *app.HttpContext) {
	if err := policy.UserViewAny(ctx); err != nil {
		ctx.WriteError(err)
		return
	}

	users := []*model.User{}
	if err := db.Find(&model.User{}, &users, Document(WithOutTrashed())); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteJSON(http.StatusOK, users)
}

func UserShow(ctx *app.HttpContext) {
	if err := policy.UserViewAny(ctx); err != nil {
		ctx.WriteError(err)
		return
	}

	id := ctx.Params["id"]
	user := &model.User{}
	withProfile, unwindProfile := user.WithProfile()

	err := db.AggregateOne(user, mongo.Pipeline{
		Match(Where("id", Eq(id))),
		user.WithRoles(),
		user.WhithPermissions(),
		withProfile,
		unwindProfile,
	})
	if err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteJSON(http.StatusOK, resource.NewUserLoginResource(user, nil))
}

func UserStore(ctx *app.HttpContext) {

	if err := policy.UserCreate(ctx); err != nil {
		ctx.WriteError(err)
		return
	}

	validator := &validator.StoreUser{}
	if err := ctx.ValidateBody(validator); err != nil {
		ctx.WriteError(err)
		return
	}

	user := &model.User{}
	db.Fill(user, validator)

	hashedPassword, er := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if er != nil {
		ctx.WriteError(&app.Err{
			Status:  http.StatusInternalServerError,
			Message: "Password encryption failed",
			Err:     er.Error(),
		})
		return
	}
	user.Password = string(hashedPassword)

	if err := db.Create(user); err != nil {
		ctx.WriteError(err)
		return
	}

	token, err := model.NewAccessToken(user.ID)
	if err != nil {
		ctx.WriteError(err)
		return
	}

	res := resource.NewUserLoginResource(user, token)
	ctx.WriteJSON(http.StatusOK, res)
}

func UserUpdate(ctx *app.HttpContext) {

	id := ctx.Get("id")

	user := &model.User{}
	if err := db.FindByHexID(user, id); err != nil {
		ctx.WriteError(err)
		return
	}

	if err := ctx.GetBody(user); err != nil {
		ctx.WriteError(err)
		return
	}

	// if err := user.Validate(ctx.Lang()); err != nil {
	// 	ctx.WriteError(err)
	// 	return
	// }

	// esto es para que el usuario solo pueda modificarse asi mismo
	if user.GetID() != ctx.User.GetID() {
		ctx.WriteError(&app.Err{
			Status:  http.StatusUnauthorized,
			Message: "No autorizado",
			Err:     "No esta autorizado para realizar esta accion",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.WriteError(&app.Err{
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

func UserDestroy(ctx *app.Context) {

	id := ctx.Get("id")
	user := &model.User{}
	if err := db.FindByHexID(user, id); err != nil {
		ctx.WriteError(err)
		return
	}

	// esto espara que el usuario solo pueda eliminarse asi mismo
	if user.GetID() != ctx.User.GetID() {
		ctx.WriteError(&app.Err{
			Status:  http.StatusUnauthorized,
			Message: "No autorizado",
			Err:     "No esta autorizado para realizar esta accion",
		})
		return
	}

	if err := db.Delete(user); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
}

func Login(ctx *app.Context) {

	var req map[string]string
	if err := ctx.GetBody(req); err != nil {
		ctx.WriteError(err)
		return
	}

	user := &model.User{}
	if err := app.Email("email", req["user"]); err == nil {
		if err := db.FindOneByField(user, "email", req["user"]); err != nil {
			ctx.WriteError(&app.Err{
				Status:  http.StatusUnauthorized,
				Message: "No autorizado",
				Err:     "No autorizado",
			})
			return
		}
	} else {
		if err := db.FindOneByField(user, "phone", req["user"]); err != nil {
			ctx.WriteError(&app.Err{
				Status:  http.StatusUnauthorized,
				Message: "No autorizado",
				Err:     "No autorizado",
			})
			return
		}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req["password"])); err != nil {
		ctx.WriteError(&app.Err{
			Status:  http.StatusUnauthorized,
			Message: "No autorizado",
			Err:     "No autorizado",
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

func Logout(ctx *app.Context) {

	if ctx.Token == nil {
		ctx.WriteError(&app.Err{
			Status:  http.StatusUnauthorized,
			Message: "No autorizado",
			Err:     "Token Invalido",
		})
		return
	}
	// if err := db.ForceDelete(ctx.Token); err != nil {
	// 	ctx.WriteError(err)
	// 	return
	// }
	ctx.WriteNoContent()
}
