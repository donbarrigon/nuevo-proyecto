package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app/controller"
	"github.com/donbarrigon/nuevo-proyecto/internal/app/model"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
	"github.com/donbarrigon/nuevo-proyecto/pkg/errors"
)

type ControllerFun func(ctx *controller.Context)

func Auth(next ControllerFun) ControllerFun {

	return func(ctx *controller.Context) {

		authHeader := ctx.Request.Header.Get("Authorization")

		if authHeader == "" {
			ctx.WriteError(&errors.Err{
				Status:  http.StatusUnauthorized,
				Message: ctx.TT("app.unauthorized"),
				Err:     ctx.TT("guard.auth.header"),
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.WriteError(&errors.Err{
				Status:  http.StatusUnauthorized,
				Message: ctx.TT("app.unauthorized"),
				Err:     ctx.TT("guard.auth.header-format"),
			})
			return
		}

		authToken := parts[1]
		tokenModel := &model.Token{}

		if err := db.Mongo.FindOneByField(tokenModel, "token", authToken); err != nil {
			ctx.WriteError(&errors.Err{
				Status:  http.StatusUnauthorized,
				Message: ctx.TT("app.unauthorized"),
				Err:     ctx.TT("guard.auth.expired") + " : " + err.Error(),
			})
			return
		}

		if tokenModel.ExpiresAt.Before(time.Now()) {
			ctx.WriteError(&errors.Err{
				Status:  http.StatusUnauthorized,
				Message: ctx.TT("app.unauthorized"),
				Err:     ctx.TT("guard.auth.expired"),
			})
			return
		}
		if _, err := db.Mongo.Update(tokenModel); err != nil {
			log.Println("[" + tokenModel.Token + "] " + err.Error())
		}

		userModel := &model.User{}
		if err := db.Mongo.FindByID(userModel, tokenModel.UserID); err != nil {
			ctx.WriteError(&errors.Err{
				Status:  http.StatusUnauthorized,
				Message: ctx.TT("app.unauthorized"),
				Err:     ctx.TT("guard.auth.invalid-token", err.Error()),
			})
			return
		}

		ctx.Token = tokenModel
		ctx.User = userModel

		next(ctx)
	}
}
