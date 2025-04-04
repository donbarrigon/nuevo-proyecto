package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	User    Model
	Token   Model
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

func (ctx *Context) GetParam() string {
	sections := strings.Split(strings.Trim(ctx.Request.URL.Path, "/"), "/")
	return sections[len(sections)-1]
}

func (ctx *Context) GetParams(params ...string) (map[string]string, Error) {
	sections := strings.Split(strings.Trim(ctx.Request.URL.Path, "/"), "/")
	numberOfSections := len(sections)
	numberOfParams := len(params)
	result := make(map[string]string, 0)
	if numberOfSections < numberOfParams {
		return nil, &Err{
			Status:  http.StatusBadRequest,
			Message: TT(ctx.Lang(), "Solicitud incorrecta"),
			Err:     TT(ctx.Lang(), "Faltan parámetros en la solicitud"),
		}
	}
	for i := 0; i < numberOfParams; i++ {
		result[params[i]] = sections[numberOfSections-(numberOfParams-i)]
	}
	return result, nil
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

func (ctx *Context) TT(s string, v ...any) string {
	return TT(ctx.Lang(), s, v...)
}
