package middleware

import (
	"log"
	"strings"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app/controller"
	"github.com/donbarrigon/nuevo-proyecto/internal/app/model"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
	"github.com/donbarrigon/nuevo-proyecto/pkg/errors"
	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
)

func Auth(next func(ctx *controller.Context)) func(ctx *controller.Context) {

	return func(ctx *controller.Context) {

		authHeader := ctx.Request.Header.Get("Authorization")

		if authHeader == "" {
			ctx.WriteError(errors.SUnauthorized(lang.TT(ctx.Lang(), "El encabezado de autorización está vacío. Se requiere un token.")))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.WriteError(errors.SUnauthorized(lang.TT(ctx.Lang(), "Formato de token inválido. Se esperaba un 'Bearer token'.")))
			return
		}

		authToken := parts[1]
		tokenModel := &model.Token{}

		if err := db.FindOneByField(tokenModel, "token", authToken); err != nil {
			ctx.WriteError(errors.SUnauthorized(lang.TT(ctx.Lang(), "El token no existe o no es válido. Verifique su autenticación.")))
			return
		}

		if tokenModel.ExpiresAt.Before(time.Now()) {
			ctx.WriteError(errors.SUnauthorized(lang.TT(ctx.Lang(), "El token ha expirado. Por favor, vuelva a autenticar.")))
			return
		}
		if err := db.Update(tokenModel); err != nil {
			log.Println("[" + tokenModel.Token + "] " + err.Error())
		}

		userModel := &model.User{}
		if err := db.FindByID(userModel, tokenModel.UserID); err != nil {
			ctx.WriteError(errors.SUnauthorized(lang.TT(ctx.Lang(), "El token no existe o no es válido. Verifique su autenticación.")))
			return
		}

		ctx.Token = tokenModel
		ctx.User = userModel

		next(ctx)
	}
}

func Token(next func(ctx *controller.Context)) func(ctx *controller.Context) {

	return func(ctx *controller.Context) {

		next(ctx)
	}
}
