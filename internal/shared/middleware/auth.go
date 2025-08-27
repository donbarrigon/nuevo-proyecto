package middleware

import (
	"strings"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/shared/model"
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
		if err := accessToken.First("token", authToken); err != nil {
			ctx.ResponseError(app.Errors.Unauthorizedf("Token Expired."))
			return
		}

		if accessToken.ExpiresAt.Before(time.Now()) {
			ctx.ResponseError(app.Errors.Unauthorizedf("Token Expired. Please login again."))
			return
		}

		user := model.NewUser()
		if err := user.FindByID(accessToken.UserID); err != nil {
			ctx.ResponseError(app.Errors.Unauthorizedf("Token user not found."))
			return
		}

		accessToken.Refresh()
		if err := accessToken.Update(); err != nil {
			app.PrintError("Failed to update access token", app.Entry{Key: "error", Value: err})
		}

		ctx.Auth = app.NewAuth(user, accessToken)

		next(ctx)
	}
}
