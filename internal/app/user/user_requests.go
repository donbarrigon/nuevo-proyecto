package user

import (
	"net/http"
	"strings"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
)

type UserRequest struct {
	ID       string `json:"-"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func validateLoginRequest(request *LoginRequest) bool {

	if strings.TrimSpace(request.User) == "" {
		return false
	}
	if len(request.User) < 5 {
		return false
	}
	if len(request.User) > 255 {
		return false
	}

	if strings.TrimSpace(request.Password) == "" {
		return false
	}
	if len(request.Password) < 8 {
		return false
	}
	if len(request.Password) > 32 {
		return false
	}

	return true
}

func showUserRequest(ctx *app.ControllerContext) (string, *app.ErrorJSON) {
	userId := ctx.Request.URL.Query().Get("u")
	if userId == "" {
		return "", &app.ErrorJSON{
			Status:  http.StatusBadRequest,
			Message: "Bad request",
			Error:   "Query params is invalid",
		}
	}
	return userId, nil
}

func loginRequest(ctx *app.ControllerContext) (*LoginRequest, *app.ErrorJSON) {
	if err := app.AllowedMethods(ctx, http.MethodGet); err != nil {
		return nil, err
	}
	request := &LoginRequest{}
	if err := app.GetBodyRequest(ctx, request); err != nil {
		return nil, err
	}

	if !validateLoginRequest(request) {
		return nil, &app.ErrorJSON{
			Status:  http.StatusUnauthorized,
			Message: lang.M(ctx.Lang(), "user.service.unautorized"),
			Error:   lang.M(ctx.Lang(), "user.service.unautorized"),
		}
	}

	return request, nil
}
