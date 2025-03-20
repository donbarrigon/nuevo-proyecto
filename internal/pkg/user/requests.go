package user

import (
	"net/http"
	"strings"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
	"github.com/donbarrigon/nuevo-proyecto/pkg/validate"
)

type UserRequest struct {
	ID       string              `json:"-"`
	Name     string              `json:"name"`
	Email    string              `json:"email"`
	Phone    string              `json:"phone"`
	Password string              `json:"password"`
	ctx      *app.HandlerContext `json:"-"`
}

type LoginRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func validateRequest(ctx *app.HandlerContext, request *UserRequest) map[string][]string {
	errMap := make(map[string][]string, 0)
	errFields := make([]string, 0)

	if strings.TrimSpace(request.Name) != "" {
		if len(request.Name) < 3 {
			errFields = append(errFields, lang.M(request.ctx.Lang(), "app.request.min", 3))
		}
		if len(request.Name) > 255 {
			errFields = append(errFields, lang.M(request.ctx.Lang(), "app.request.max", 255))
		}
		if len(errFields) > 0 {
			errMap["name"] = errFields
			errFields = make([]string, 0)
		}
	} else {
		errMap["name"] = []string{lang.M(request.ctx.Lang(), "app.request.required")}
	}

	if strings.TrimSpace(request.Email) != "" {
		if len(request.Email) > 255 {
			errFields = append(errFields, lang.M(request.ctx.Lang(), "app.request.max", 255))
		}
		if !validate.Email(request.Email) {
			errFields = append(errFields, lang.M(request.ctx.Lang(), "app.request.email"))
		}
		if len(errFields) > 0 {
			errMap["email"] = errFields
			errFields = make([]string, 0)
		}
	}

	if strings.TrimSpace(request.Phone) != "" {
		if len(request.Phone) < 5 {
			errFields = append(errFields, lang.M(request.ctx.Lang(), "app.request.min", 5))
		}
		if len(request.Phone) > 255 {
			errFields = append(errFields, lang.M(request.ctx.Lang(), "app.request.max", 255))
		}
		if len(errFields) > 0 {
			errMap["phone"] = errFields
			errFields = make([]string, 0)
		}
	}

	if strings.TrimSpace(request.Password) != "" {
		if len(request.Password) < 8 {
			errFields = append(errFields, lang.M(request.ctx.Lang(), "app.request.min", 8))
		}
		if len(request.Password) > 32 {
			errFields = append(errFields, lang.M(request.ctx.Lang(), "app.request.max", 32))
		}
		if len(errFields) > 0 {
			errMap["password"] = errFields
		}
	} else {
		errMap["password"] = []string{lang.M(request.ctx.Lang(), "app.request.required")}
	}

	if strings.TrimSpace(request.Email) == "" && strings.TrimSpace(request.Phone) == "" {
		errMap["email"] = []string{"El email o telefono son requeridos"}
		errMap["phone"] = []string{"El email o telefono son requeridos"}
	}

	return errMap
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

func storeRequest(ctx *app.HandlerContext) (*UserRequest, *app.ErrorJSON) {
	if err := app.AllowedMethods(ctx, http.MethodPost); err != nil {
		return nil, err
	}

	request := &UserRequest{ctx: ctx}
	if err := app.GetRequestBody(ctx, request); err != nil {
		return nil, err
	}

	if errMap := validateRequest(ctx, request); len(errMap) > 0 {
		return nil, &app.ErrorJSON{
			Status:  http.StatusUnprocessableEntity,
			Message: lang.M(ctx.Lang(), "app.unprocessable-entity"),
			Error:   errMap,
		}
	}

	return request, nil
}

func updateRequest(ctx *app.HandlerContext) (*UserRequest, *app.ErrorJSON) {
	if err := app.AllowedMethods(ctx, http.MethodPut); err != nil {
		return nil, err
	}

	request := &UserRequest{}
	if err := app.GetRequestBody(ctx, request); err != nil {
		return nil, err
	}

	userId := ctx.Request.URL.Query().Get("u")
	if userId == "" {
		return nil, &app.ErrorJSON{
			Status:  http.StatusBadRequest,
			Message: "Bad request",
			Error:   "Query params is invalid",
		}
	}
	request.ID = userId

	if errMap := validateUserRequest(request); len(errMap) > 0 {
		return nil, &app.ErrorJSON{
			Status:  http.StatusUnprocessableEntity,
			Message: "Invalid request",
			Error:   errMap,
		}
	}

	return request, nil
}

func showUserRequest(ctx *app.HandlerContext) (string, *app.ErrorJSON) {
	if err := app.AllowedMethods(ctx, http.MethodGet); err != nil {
		return "", err
	}
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

func loginRequest(ctx *app.HandlerContext) (*LoginRequest, *app.ErrorJSON) {
	if err := app.AllowedMethods(ctx, http.MethodGet); err != nil {
		return nil, err
	}
	request := &LoginRequest{}
	if err := app.GetRequestBody(ctx, request); err != nil {
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
