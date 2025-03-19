package guard

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func AuthToken(ctx *app.HandlerContext) bool {
	authHeader := ctx.Request.Header.Get("Authorization")

	if authHeader == "" {
		app.ResponseErrorJSON(ctx.Writer, "Invalid authorization header", http.StatusUnauthorized, "Unauthorized")
		return false
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		app.ResponseErrorJSON(ctx.Writer, "Invalid Authorization header format", http.StatusUnauthorized, "Unauthorized")
		return false
	}

	authToken := parts[1]
	tokenModel := &model.Token{}
	filter := bson.D{{Key: "token", Value: authToken}}
	if err := app.DB.FindOne(tokenModel, filter); err != nil {
		app.ResponseErrorJSON(ctx.Writer, "Invalid token "+err.Error(), http.StatusUnauthorized, "Unauthorized")
		return false
	}

	if tokenModel.ExpiresAt.Before(time.Now()) {
		app.ResponseErrorJSON(ctx.Writer, "Token has expired", http.StatusUnauthorized, "Unauthorized")
		return false
	}

	tokenModel.Refresh()
	if _, err := app.DB.Update(tokenModel, tokenModel.ID); err != nil {
		log.Println("Error updating token [" + tokenModel.Token + "]")
		log.Println(tokenModel)
	}

	userModel := &model.User{}
	if err := app.DB.FindByID(userModel, tokenModel.UserID); err != nil {
		app.ResponseErrorJSON(ctx.Writer, "Invalid user "+err.Error(), http.StatusUnauthorized, "Unauthorized")
		return false
	}

	ctx.User = userModel
	ctx.Token = tokenModel
	return true
}
