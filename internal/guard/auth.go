package guard

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/model"
)

func AuthToken(ctx *app.HandlerContext) *app.ErrorJSON {
	authHeader := ctx.Request.Header.Get("Authorization")

	if authHeader == "" {
		return &app.ErrorJSON{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Invalid authorization header",
		}
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return &app.ErrorJSON{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Invalid Authorization header format",
		}
	}

	authToken := parts[1]
	tokenModel := &model.Token{}
	if err := app.Mongo.FindByHexID(tokenModel, authToken); err != nil {
		return &app.ErrorJSON{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Invalid token",
		}
	}

	if tokenModel.ExpiresAt.Before(time.Now()) {
		return &app.ErrorJSON{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Token has expired",
		}
	}

	tokenModel.Refresh()
	if _, err := app.Mongo.Update(tokenModel); err != nil {
		log.Println("Error updating token [" + tokenModel.Token + "]")
		log.Println(err)
		log.Println(tokenModel)
	}

	userModel := &model.User{}
	if err := app.Mongo.FindByID(userModel, tokenModel.UserID); err != nil {
		return &app.ErrorJSON{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Invalid user " + err.Error(),
		}
	}

	ctx.User = userModel
	ctx.Token = tokenModel
	return nil
}
