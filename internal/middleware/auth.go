package middleware

import (
	"strings"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	. "github.com/donbarrigon/nuevo-proyecto/internal/app/qb"
	"github.com/donbarrigon/nuevo-proyecto/internal/model"
)

func Auth(next func(ctx *app.HttpContext)) func(ctx *app.HttpContext) {

	return func(ctx *app.HttpContext) {

		authHeader := ctx.Request.Header.Get("Authorization")

		if authHeader == "" {
			ctx.ResponseError(app.Errors.Unauthorizedf("The 'Authorization' header is required."))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.ResponseError(app.Errors.Unauthorizedf("The 'Authorization' header is invalid."))
			return
		}
		authToken := parts[1]

		accessToken := model.NewAccessToken()
		if err := accessToken.AggregateOne(Pipeline(
			Match(Where("token", Eq(authToken))),
			With("users", "user_id", "_id", "user"),
			Unwind("$user"),
		)); err != nil {
			ctx.ResponseError(app.Errors.Unauthorizedf("Token not found. :error", app.Entry{Key: "error", Value: err.Error()}))
			return
		}

		if accessToken.ExpiresAt.Before(time.Now()) {
			ctx.ResponseError(app.Errors.Unauthorizedf("Token Expired."))
			return
		}

		if accessToken.User.DeletedAt != nil {
			ctx.ResponseError(app.Errors.Unauthorizedf("User Inactive or Deleted."))
			return
		}

		if err := accessToken.Refresh(); err != nil {
			app.PrintError("Failed to update access token", app.Entry{Key: "error", Value: err.Error()})
		}

		// ctx.Writer.Header().Set("Authorization", "Bearer "+accessToken.Token)
		ctx.Auth = accessToken

		next(ctx)
	}
}
