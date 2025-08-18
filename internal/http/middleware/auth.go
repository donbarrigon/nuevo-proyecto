package middleware

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
)

func Auth(next func(ctx *app.HttpContext)) func(ctx *app.HttpContext) {

	return func(ctx *app.HttpContext) {

		app.Log.Print("Auth middleware")
		next(ctx)
		return

		// authHeader := ctx.Request.Header.Get("Authorization")

		// if authHeader == "" {
		// 	ctx.WriteError(app.Errors.Unauthorizedf("El encabezado de autorización está vacío. Se requiere un token."))
		// 	return
		// }

		// parts := strings.Split(authHeader, " ")
		// if len(parts) != 2 || parts[0] != "Bearer" {
		// 	ctx.WriteError(app.Errors.Unauthorizedf("Formato de token inválido. Se esperaba un 'Bearer token'."))
		// 	return
		// }

		// authToken := parts[1]
		// tokenModel := &model.Token{}
		// if err := db.FindOneByField(tokenModel, "token", authToken); err != nil {
		// 	ctx.WriteError(app.Errors.Unauthorizedf("El token no existe o no es válido. Verifique su autenticación."))
		// 	return
		// }

		// if tokenModel.ExpiresAt.Before(time.Now()) {
		// 	ctx.WriteError(app.Errors.Unauthorizedf("El token ha expirado. Por favor, vuelva a autenticar."))
		// 	return
		// }
		// if err := db.Update(tokenModel); err != nil {
		// 	log.Println("[" + tokenModel.Token + "] " + err.Error())
		// }

		// userModel := &model.User{}
		// if err := db.FindByID(userModel, tokenModel.UserID); err != nil {
		// 	ctx.WriteError(app.Errors.Unauthorizedf("El token no existe o no es válido. Verifique su autenticación."))
		// 	return
		// }

		// ctx.Token = tokenModel
		// ctx.User = userModel

		// next(ctx)
	}
}
