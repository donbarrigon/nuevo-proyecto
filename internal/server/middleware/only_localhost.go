package middleware

import (
	"net"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
)

func OnlyLocalhost(next func(ctx *app.HttpContext)) func(ctx *app.HttpContext) {

	return func(ctx *app.HttpContext) {

		if !app.Env.DB_MIGRATION_ENABLE {
			app.PrintWarning("DB_MIGRATION_ENABLE is false")
			ctx.ResponseNoContent()
			return
		}

		host, _, er := net.SplitHostPort(ctx.Request.RemoteAddr)
		if er != nil {
			app.PrintWarning("Fail to get remote address: " + er.Error())
			ctx.ResponseNoContent()
			return
		}
		if host == "127.0.0.1" || host == "::1" {
			next(ctx)
			return
		}
		app.PrintWarning("Remote address is not localhost")
		ctx.ResponseNoContent()

	}
}
