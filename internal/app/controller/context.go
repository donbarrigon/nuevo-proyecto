package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
	"github.com/donbarrigon/nuevo-proyecto/pkg/errors"
	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
)

type Context interface {
	W() http.ResponseWriter
	R() *http.Request
}
type Controller struct {
	Writer  http.ResponseWriter
	Request *http.Request
	User    db.MongoModel
	Token   db.MongoModel
}

// func NewContext(w http.ResponseWriter, r *http.Request) *Context {
// 	return &Controller{
// 		Writer:  w,
// 		Request: r,
// 	}
// }

func (ctx *Controller) W() http.ResponseWriter {
	return ctx.Writer
}

func (ctx *Controller) R() *http.Request {
	return ctx.Request
}

func (ctx *Controller) Lang() string {
	return ctx.Request.Header.Get("Accept-Language")
}

func (ctx *Controller) GetBody(request any) errors.Error {
	decoder := json.NewDecoder(ctx.Request.Body)
	if err := decoder.Decode(request); err != nil {
		return &errors.Err{
			Status:  http.StatusBadRequest,
			Message: lang.TT(ctx.Lang(), "Solicitud incorrecta"),
			Err:     lang.TT(ctx.Lang(), "No se pudo decodificar el cuerpo de la solicitud") + ": " + err.Error(),
		}
	}
	defer ctx.Request.Body.Close()
	return nil
}

func (ctx *Controller) Get(param string, defaultValue string) string {
	value := ctx.Request.URL.Query().Get(param)
	if value == "" {
		return defaultValue
	}
	return value
}

func (ctx *Controller) GetParam() string {
	sections := strings.Split(strings.Trim(ctx.Request.URL.Path, "/"), "/")
	return sections[len(sections)-1]
}

func (ctx *Controller) GetParams(params ...string) (map[string]string, errors.Error) {
	sections := strings.Split(strings.Trim(ctx.Request.URL.Path, "/"), "/")
	numberOfSections := len(sections)
	numberOfParams := len(params)
	result := make(map[string]string, 0)
	if numberOfSections < numberOfParams {
		return nil, &errors.Err{
			Status:  http.StatusBadRequest,
			Message: lang.TT(ctx.Lang(), "Solicitud incorrecta"),
			Err:     lang.TT(ctx.Lang(), "Faltan parámetros en la solicitud"),
		}
	}
	for i := 0; i < numberOfParams; i++ {
		result[params[i]] = sections[numberOfSections-(numberOfParams-i)]
	}
	return result, nil
}

func (ctx *Controller) WriteJSON(status int, data any) {
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.WriteHeader(status)

	if err := json.NewEncoder(ctx.Writer).Encode(data); err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.Writer.WriteHeader(500)
		ctx.Writer.Write([]byte(fmt.Sprintf(`{"message":"Error","error":"%s"}`, lang.TT(ctx.Lang(), "No se pudo codificar la respuesta"))))
	}
}

func (ctx *Controller) WriteError(err errors.Error) {
	ctx.WriteJSON(err.GetStatus(), err)
}

func (ctx *Controller) TT(s string, v ...any) string {
	return lang.TT(ctx.Lang(), s, v...)
}
