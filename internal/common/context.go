package com

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	User    Model
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
	}
}

func (ctx *Context) Lang() string {
	return ctx.Request.Header.Get("Accept-Language")
}

func (ctx *Context) GetBody(request any) Error {
	decoder := json.NewDecoder(ctx.Request.Body)
	if err := decoder.Decode(request); err != nil {
		return &Err{
			Status:  http.StatusBadRequest,
			Message: TT(ctx.Lang(), "Solicitud incorrecta"),
			Err:     TT(ctx.Lang(), "No se pudo decodificar el cuerpo de la solicitud") + ": " + err.Error(),
		}
	}
	defer ctx.Request.Body.Close()
	return nil
}

func (ctx *Context) Get(param string, defaultValue string) string {
	value := ctx.Request.URL.Query().Get(param)
	if value == "" {
		return defaultValue
	}
	return value
}

func (ctx *Context) WriteJSON(status int, data any) {
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.WriteHeader(status)

	if err := json.NewEncoder(ctx.Writer).Encode(data); err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.Writer.WriteHeader(500)
		ctx.Writer.Write([]byte(fmt.Sprintf(`{"message":"Error","error":"%s"}`, TT(ctx.Lang(), "No se pudo codificar la respuesta"))))
	}
}

func (ctx *Context) WriteError(err Error) {
	ctx.WriteJSON(err.GetStatus(), err)
}
