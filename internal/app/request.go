package app

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
)

func GetRequestBody(ctx *HandlerContext, request any) *ErrorJSON {
	decoder := json.NewDecoder(ctx.Request.Body)
	if err := decoder.Decode(request); err != nil {
		return &ErrorJSON{
			Status:  http.StatusBadRequest,
			Message: lang.M(ctx.Lang(), "app.bad-request"),
			Error:   "Failed to decode request body: " + err.Error(),
		}
	}
	defer ctx.Request.Body.Close()
	return nil
}

func AllowedMethods(ctx *HandlerContext, methods ...string) *ErrorJSON {

	if slices.Contains(methods, ctx.Request.Method) {
		return nil
	}

	return &ErrorJSON{
		Status:  http.StatusMethodNotAllowed,
		Message: lang.M(ctx.Lang(), "app.method-not-allowed"),
		Error:   "Method [" + ctx.Request.Method + "] Not Allowed",
	}
}
