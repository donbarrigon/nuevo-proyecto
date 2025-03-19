package user

import (
	"net/http"
	"strings"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/pkg/validate"
)

type UserRequest struct {
	ID       string
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func validateUserRequest(request *UserRequest) map[string][]string {
	errMap := make(map[string][]string, 0)
	errFields := make([]string, 0)

	if strings.TrimSpace(request.Name) != "" {
		if len(request.Name) < 3 {
			errFields = append(errFields, "Minimo 3 caracteres")
		}
		if len(request.Name) > 255 {
			errFields = append(errFields, "Maximo 255 caracteres")
		}
		if len(errFields) > 0 {
			errMap["name"] = errFields
			errFields = make([]string, 0)
		}
	} else {
		errMap["name"] = []string{"Es requerido"}
	}

	if strings.TrimSpace(request.Email) != "" {
		if len(request.Email) > 255 {
			errFields = append(errFields, "Maximo 255 caracteres")
		}
		if !validate.Email(request.Email) {
			errFields = append(errFields, "El email no es valido")
		}
		if len(errFields) > 0 {
			errMap["email"] = errFields
			errFields = make([]string, 0)
		}
	}

	if strings.TrimSpace(request.Phone) != "" {
		if len(request.Phone) < 5 {
			errFields = append(errFields, "Minimo 5 caracteres")
		}
		if len(request.Phone) > 255 {
			errFields = append(errFields, "Maximo 255 caracteres")
		}
		if len(errFields) > 0 {
			errMap["phone"] = errFields
			errFields = make([]string, 0)
		}
	}

	if strings.TrimSpace(request.Password) != "" {
		if len(request.Password) < 8 {
			errFields = append(errFields, "Minimo 8 caracteres")
		}
		if len(request.Password) > 32 {
			errFields = append(errFields, "Maximo 32 caracteres")
		}
		if len(errFields) > 0 {
			errMap["password"] = errFields
		}
	} else {
		errMap["password"] = []string{"Es requerido"}
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

func storeRequest(ctx *app.HandlerContext) *UserRequest {
	request := &UserRequest{}
	if !app.GetRequestBody(ctx, request, http.MethodPost) {
		return nil
	}

	if errMap := validateUserRequest(request); len(errMap) > 0 {
		app.ResponseErrorJSON(ctx.Writer, errMap, http.StatusUnprocessableEntity, "Unprocessable Entity")
		return nil
	}

	return request
}

func updateRequest(ctx *app.HandlerContext) *UserRequest {
	request := &UserRequest{}
	if !app.GetRequestBody(ctx, request, http.MethodPut) {
		return nil
	}

	if errMap := validateUserRequest(request); len(errMap) > 0 {
		app.ResponseErrorJSON(ctx.Writer, errMap, http.StatusUnprocessableEntity, "Unprocessable Entity")
		return nil
	}

	userId := ctx.Request.URL.Query().Get("u")
	if userId == "" {
		app.ResponseErrorJSON(ctx.Writer, "Query params is invalid", http.StatusBadRequest, "Bad request")
		return nil
	}
	request.ID = userId

	return request
}

func showUserRequest(ctx *app.HandlerContext) string {
	if ctx.Request.Method != http.MethodGet {
		message := "Method " + ctx.Request.Method + " Not Allowed"
		app.ResponseErrorJSON(ctx.Writer, message, http.StatusMethodNotAllowed, "Method Not Allowed")
		return ""
	}
	return ctx.Request.URL.Query().Get("u")
}

func loginRequest(ctx *app.HandlerContext) *LoginRequest {
	request := &LoginRequest{}
	if !app.GetRequestBody(ctx, request, http.MethodPost) {
		return nil
	}

	if !validateLoginRequest(request) {
		return nil
	}

	return request
}
