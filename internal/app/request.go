package app

import (
	"encoding/json"
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
)

type Request interface {
	Validate(language string) map[string][]string
}

func AllowedMethod(ctx *HandlerContext, method string) *ErrorJSON {

	if method == ctx.Request.Method {
		return nil
	}

	return &ErrorJSON{
		Status:  http.StatusMethodNotAllowed,
		Message: lang.M(ctx.Lang(), "app.method-not-allowed"),
		Error:   "Method [" + ctx.Request.Method + "] Not Allowed",
	}
}

func AllowedMethods(ctx *HandlerContext, methods ...string) *ErrorJSON {

	// if slices.Contains(methods, ctx.Request.Method) {
	// 	return nil
	// }
	for _, method := range methods {
		if method == ctx.Request.Method {
			return nil
		}
	}

	return &ErrorJSON{
		Status:  http.StatusMethodNotAllowed,
		Message: lang.M(ctx.Lang(), "app.method-not-allowed"),
		Error:   "Method [" + ctx.Request.Method + "] Not Allowed",
	}
}

func GetBodyRequest(ctx *HandlerContext, request any) *ErrorJSON {
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

func GetRequest(ctx *HandlerContext, request Request, methods ...string) *ErrorJSON {

	if err := AllowedMethods(ctx, methods...); err != nil {
		return err
	}

	if err := GetBodyRequest(ctx, request); err != nil {
		return err
	}

	if errMap := request.Validate(ctx.Lang()); len(errMap) > 0 {
		return &ErrorJSON{
			Status:  http.StatusUnprocessableEntity,
			Message: lang.M(ctx.Lang(), "app.unprocessable-entity"),
			Error:   errMap,
		}
	}

	return nil
}
