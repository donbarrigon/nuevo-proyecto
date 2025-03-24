package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/model"
	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
)

func AuthToken(next func(ctx *app.ControllerContext)) func(ctx *app.ControllerContext) {

	return func(ctx *app.ControllerContext) {

		authHeader := ctx.Request.Header.Get("Authorization")

		if authHeader == "" {
			err := app.ErrorJSON{
				Message: lang.M(ctx.Lang(), "app.unauthorized"),
				Error:   lang.M(ctx.Lang(), "guard.auth.header"),
				Status:  http.StatusUnauthorized,
			}
			err.WriteResponse(ctx.Writer)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			err := app.ErrorJSON{
				Message: lang.M(ctx.Lang(), "app.unauthorized"),
				Error:   lang.M(ctx.Lang(), "guard.auth.header-format"),
				Status:  http.StatusUnauthorized,
			}
			err.WriteResponse(ctx.Writer)
			return
		}

		authToken := parts[1]
		tokenModel := &model.Token{}

		// collection := app.Mongo.Database.Collection(tokenModel.CollectionName())
		// filter := bson.D{
		// 	{Key: "token", Value: authToken},
		// 	{Key: "expiresAt", Value: bson.D{{Key: "$gt", Value: time.Now()}}}, // Solo si no ha expirado
		// }
		// update := bson.D{{Key: "$set", Value: bson.D{{Key: "expiresAt", Value: tokenModel.RefreshTime()}}}}
		// result, err2 := collection.UpdateOne(context.TODO(), filter, update)
		// if err2 != nil {
		// 	return &app.ErrorJSON{
		// 		Status:  http.StatusUnauthorized,
		// 		Message: lang.M(ctx.Lang(), "app.unauthorized"),
		// 		Error:   lang.M(ctx.Lang(), "guard.auth.expired") + " : " + err2.Error(),
		// 	}
		// }
		// if result.MatchedCount == 0 {
		// 	return &app.ErrorJSON{
		// 		Status:  http.StatusUnauthorized,
		// 		Message: lang.M(ctx.Lang(), "app.unauthorized"),
		// 		Error:   lang.M(ctx.Lang(), "guard.auth.expired"),
		// 	}
		// }

		if err := app.Mongo.FindOneByField(tokenModel, "token", authToken); err != nil {
			err := app.ErrorJSON{
				Message: lang.M(ctx.Lang(), "app.unauthorized"),
				Error:   lang.M(ctx.Lang(), "guard.auth.expired") + " : " + err.Error(),
				Status:  http.StatusUnauthorized,
			}
			err.WriteResponse(ctx.Writer)
			return
		}

		if tokenModel.ExpiresAt.Before(time.Now()) {
			err := app.ErrorJSON{
				Message: lang.M(ctx.Lang(), "app.unauthorized"),
				Error:   lang.M(ctx.Lang(), "guard.auth.expired"),
				Status:  http.StatusUnauthorized,
			}
			err.WriteResponse(ctx.Writer)
		}
		tokenModel.Refresh()
		if _, err := app.Mongo.Update(tokenModel); err != nil {
			log.Println("[" + tokenModel.Token + "] " + err.Error())
		}

		userModel := &model.User{}
		if err := app.Mongo.FindByID(userModel, tokenModel.UserID); err != nil {
			err := app.ErrorJSON{
				Message: lang.M(ctx.Lang(), "app.unauthorized"),
				Error:   lang.M(ctx.Lang(), "guard.auth.invalid-token", err.Error()),
				Status:  http.StatusUnauthorized,
			}
			err.WriteResponse(ctx.Writer)
		}

		ctx.Token = tokenModel
		ctx.User = userModel

		next(ctx)
	}
}
