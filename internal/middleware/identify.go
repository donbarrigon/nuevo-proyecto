package middleware

import (
	"log"
	"strings"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
	"github.com/donbarrigon/nuevo-proyecto/internal/model"
)

func Identify(next func(ctx *app.Context)) func(ctx *app.Context) {

	return func(ctx *app.Context) {
		user := &model.User{}
		token := &model.Token{}

		authHeader := ctx.Request.Header.Get("Authorization")

		if authHeader == "" {
			ctx.Token = token.Anonymous()
			ctx.User = user.Anonymous()
			next(ctx)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.Token = token.Anonymous()
			ctx.User = user.Anonymous()
			next(ctx)
			return
		}

		authToken := parts[1]
		if err := db.FindOneByField(token, "token", authToken); err != nil {
			ctx.Token = token.Anonymous()
			ctx.User = user.Anonymous()
			next(ctx)
			return
		}

		if token.ExpiresAt.Before(time.Now()) {
			ctx.Token = token.Anonymous()
			ctx.User = user.Anonymous()
			next(ctx)
			return
		}

		if token.ID != token.Anonymous().ID {
			if err := db.Update(token); err != nil {
				log.Println("[" + token.Token + "] " + err.Error())
			}
		}

		if err := db.FindByID(user, token.UserID); err != nil {
			ctx.Token = token.Anonymous()
			ctx.User = user.Anonymous()
			next(ctx)
			return
		}

		ctx.Token = token
		ctx.User = user
		next(ctx)
	}
}
