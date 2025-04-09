package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/donbarrigon/nuevo-proyecto/internal/app/model"
	"github.com/donbarrigon/nuevo-proyecto/pkg/errors"
	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
)

type Filter struct {
	Key    string
	Filter string
	Value  string
}

type Sort struct {
	Field     string
	Direction int // 1 para asc, -1 para desc (Mongo style)
}

type QueryFilter struct {
	Filters         []*Filter
	Sort            []*Sort
	Page            int
	PerPage         int
	Cursor          string
	CursorDirection int // 1 para asc, -1 para desc (Mongo style)

}

var allFilters = []string{
	"eq",      // Igual a
	"ne",      // Distinto de (not equal)
	"gt",      // Mayor que (greater than)
	"gte",     // Mayor o igual que (greater than or equal)
	"lt",      // Menor que (less than)
	"lte",     // Menor o igual que (less than or equal)
	"lk",      // Contiene (similar a SQL LIKE, puede ser case insensitive)
	"ilk",     // LIKE sin distinción entre mayúsculas/minúsculas (PostgreSQL style)
	"in",      // Dentro de una lista de valores
	"nin",     // No dentro de una lista de valores (not in)
	"null",    // Es nulo
	"nnull",   // No es nulo
	"between", // Entre dos valores
}

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	User    *model.User
	Token   *model.Token
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

func (ctx *Context) GetBody(request any) errors.Error {
	decoder := json.NewDecoder(ctx.Request.Body)
	if err := decoder.Decode(request); err != nil {
		return &errors.Err{
			Status:  http.StatusBadRequest,
			Message: "El cuerpo de la solicitud es incorrecto",
			Err:     lang.TT(ctx.Lang(), "No se pudo decodificar el cuerpo de la solicitud") + ": " + err.Error(),
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

func (ctx *Context) LastParam() string {
	sections := strings.Split(strings.Trim(ctx.Request.URL.Path, "/"), "/")
	return sections[len(sections)-1]
}

func (ctx *Context) GetParams(params ...string) (map[string]string, errors.Error) {
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

func (ctx *Context) WriteJSON(status int, data any) {
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.WriteHeader(status)

	if err := json.NewEncoder(ctx.Writer).Encode(data); err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.Writer.WriteHeader(500)
		ctx.Writer.Write([]byte(fmt.Sprintf(`{"message":"Error","error":"%s"}`, lang.TT(ctx.Lang(), "No se pudo codificar la respuesta"))))
	}
}

func (ctx *Context) WriteError(err errors.Error) {
	err.Traslate(ctx.Lang())
	ctx.WriteJSON(err.GetStatus(), err)
}

func (ctx *Context) WriteNotFound() {
	ctx.WriteError(errors.NotFound(errors.New(lang.TT(ctx.Lang(), "El recurso [%v:%v] no existe", ctx.Request.Method, ctx.Request.URL.Path))))
}

func (ctx *Context) WriteNoContent() {
	ctx.Writer.WriteHeader(http.StatusNoContent)
}

func (ctx *Context) TT(s string, v ...any) string {
	return lang.TT(ctx.Lang(), s, v...)
}

func (ctx *Context) GetQueryFilter(allowFilters map[string][]string) *QueryFilter {
	query := ctx.Request.URL.Query()

	qf := &QueryFilter{
		Filters:         []*Filter{},
		Sort:            []*Sort{},
		Page:            1,
		PerPage:         15,
		Cursor:          "",
		CursorDirection: 1,
	}

	for rawKey, values := range query {
		if len(values) == 0 {
			continue
		}
		value := values[0]

		// --------------------
		// 1. Paginación
		// --------------------
		if rawKey == "page" {
			if page, err := strconv.Atoi(value); err == nil && page > 0 {
				qf.Page = page
				continue
			}
		}

		if rawKey == "per_page" {
			if perPage, err := strconv.Atoi(value); err == nil && perPage > 0 {
				qf.PerPage = perPage
			}
		}

		// --------------------
		// Cursor
		// --------------------
		if rawKey == "cursor[asc]" || rawKey == "cursor" {
			qf.Cursor = value
			qf.CursorDirection = 1
			continue
		}

		if rawKey == "cursor[desc]" {
			qf.Cursor = value
			qf.CursorDirection = -1
			continue
		}

		// --------------------
		// Sort
		// --------------------
		if rawKey == "sort" {
			fields := strings.Split(value, ",")
			for _, f := range fields {
				f = strings.TrimSpace(f)
				if f == "" {
					continue
				}
				direction := 1
				if strings.HasPrefix(f, "-") {
					direction = -1
					f = f[1:]
				} else if strings.HasPrefix(f, "+") {
					f = f[1:]
				} else if strings.HasPrefix(f, "[asc]") {
					f = f[5:]
				} else if strings.HasPrefix(f, "[desc]") {
					direction = -1
					f = f[6:]
				} else if strings.HasSuffix(f, ":asc") {
					f = f[:4]
				} else if strings.HasSuffix(f, ":desc") {
					direction = -1
					f = f[:5]
				}
				qf.Sort = append(qf.Sort, &Sort{
					Field:     f,
					Direction: direction,
				})
			}
			continue
		}

		// --------------------
		// Filtros
		// Filtro normal: ej. name[eq]=John
		// --------------------
		var key, filter string
		if strings.Contains(rawKey, "[") && strings.HasSuffix(rawKey, "]") {
			parts := strings.SplitN(rawKey, "[", 2)
			if len(parts) != 2 {
				continue
			}
			key = parts[0]
			filter = strings.TrimSuffix(parts[1], "]")

			if filter == "in" || filter == "nin" {
				value = strings.Join(values, ",")
			}
		} else {
			key = rawKey
			filter = "eq"
		}

		allowed, ok := allowFilters[key]
		if !ok {
			continue
		}

		isValid := false
		for _, af := range allowed {
			if af == filter {
				isValid = true
				break
			}
		}
		if !isValid {
			continue
		}

		qf.Filters = append(qf.Filters, &Filter{
			Key:    key,
			Filter: filter,
			Value:  value,
		})
	}

	return qf
}
