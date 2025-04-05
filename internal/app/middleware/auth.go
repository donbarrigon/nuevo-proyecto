package middleware

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app/controller"
)

func Auth(ctx controller.Context, next func()) {
	next()
}

// func Auth2(next func(ctx *core.Context)) func(ctx *core.Context) {

// 	return func(ctx *core.Context) {

// 		authHeader := ctx.Request.Header.Get("Authorization")

// 		if authHeader == "" {
// 			ctx.WriteError(&core.Err{
// 				Status:  http.StatusUnauthorized,
// 				Message: ctx.TT("app.unauthorized"),
// 				Err:     ctx.TT("guard.auth.header"),
// 			})
// 			return
// 		}

// 		parts := strings.Split(authHeader, " ")
// 		if len(parts) != 2 || parts[0] != "Bearer" {
// 			ctx.WriteError(&core.Err{
// 				Status:  http.StatusUnauthorized,
// 				Message: ctx.TT("app.unauthorized"),
// 				Err:     ctx.TT("guard.auth.header-format"),
// 			})
// 			return
// 		}

// 		authToken := parts[1]
// 		tokenModel := &model.Token{}

// 		if err := core.Mongo.FindOneByField(tokenModel, "token", authToken); err != nil {
// 			ctx.WriteError(&core.Err{
// 				Status:  http.StatusUnauthorized,
// 				Message: ctx.TT("app.unauthorized"),
// 				Err:     ctx.TT("guard.auth.expired") + " : " + err.Error(),
// 			})
// 			return
// 		}

// 		if tokenModel.ExpiresAt.Before(time.Now()) {
// 			ctx.WriteError(&core.Err{
// 				Status:  http.StatusUnauthorized,
// 				Message: ctx.TT("app.unauthorized"),
// 				Err:     ctx.TT("guard.auth.expired"),
// 			})
// 			return
// 		}
// 		if _, err := core.Mongo.Update(tokenModel); err != nil {
// 			log.Println("[" + tokenModel.Token + "] " + err.Error())
// 		}

// 		userModel := &model.User{}
// 		if err := core.Mongo.FindByID(userModel, tokenModel.UserID); err != nil {
// 			ctx.WriteError(&core.Err{
// 				Status:  http.StatusUnauthorized,
// 				Message: ctx.TT("app.unauthorized"),
// 				Err:     ctx.TT("guard.auth.invalid-token", err.Error()),
// 			})
// 			return
// 		}

// 		ctx.Token = tokenModel
// 		ctx.User = userModel

// 		next(ctx)
// 	}
// }
//
