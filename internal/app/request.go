package app

import (
	"encoding/json"
	"net/http"
)

func GetRequestBody(ctx *HandlerContext, request any, method string) bool {
	if ctx.Request.Method != method {
		message := "Method " + ctx.Request.Method + " Not Allowed"
		ResponseErrorJSON(ctx.Writer, message, http.StatusMethodNotAllowed, "Method Not Allowed")
		return false
	}

	decoder := json.NewDecoder(ctx.Request.Body)
	if err := decoder.Decode(request); err != nil {
		message := "Failed to decode request body: " + err.Error()
		ResponseErrorJSON(ctx.Writer, message, http.StatusBadRequest, "Invalid Request Body")
		return false
	}
	defer ctx.Request.Body.Close()
	return true
}
