package controller

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app/model"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
	"github.com/donbarrigon/nuevo-proyecto/pkg/errors"
	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
)

type MessageResource struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
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

func (ctx *Context) WriteMessage(code int, data any, message string, v ...any) {
	ctx.WriteJSON(code, &MessageResource{
		Message: lang.TT(ctx.Lang(), message, v...),
		Data:    data,
	})
}

func (ctx *Context) WriteSuccess(data any) {
	ctx.WriteJSON(http.StatusOK, &MessageResource{
		Message: lang.TT(ctx.Lang(), "Solicitud procesada con éxito"),
		Data:    data,
	})
}

func (ctx *Context) WriteCreated(data any) {
	ctx.WriteJSON(http.StatusCreated, &MessageResource{
		Message: lang.TT(ctx.Lang(), "Recurso creado exitosamente"),
		Data:    data,
	})
}

func (ctx *Context) WriteUpdated(data any) {
	ctx.WriteJSON(http.StatusOK, &MessageResource{
		Message: lang.TT(ctx.Lang(), "Recurso actualizado exitosamente"),
		Data:    data,
	})
}

func (ctx *Context) WriteDeleted(data any) {
	ctx.WriteJSON(http.StatusOK, &MessageResource{
		Message: lang.TT(ctx.Lang(), "Recurso eliminado exitosamente"),
		Data:    data,
	})
}

func (ctx *Context) WriteRestored(data any) {
	ctx.WriteJSON(http.StatusOK, &MessageResource{
		Message: lang.TT(ctx.Lang(), "Recurso restaurado exitosamente"),
		Data:    data,
	})
}

func (ctx *Context) WriteForceDeleted(data any) {
	ctx.WriteJSON(http.StatusOK, &MessageResource{
		Message: lang.TT(ctx.Lang(), "Recurso eliminado permanentemente"),
		Data:    data,
	})
}

func (ctx *Context) WriteNoContent() {
	ctx.Writer.WriteHeader(http.StatusNoContent)
}

func (ctx *Context) TT(s string, v ...any) string {
	return lang.TT(ctx.Lang(), s, v...)
}

func (ctx *Context) GetQueryFilter(allowFilters map[string][]string) *db.QueryFilter {
	query := ctx.Request.URL.Query()

	qf := db.NewQueryFilter()

	qf.Path = ctx.Request.URL.Path

	for rawKey, values := range query {
		if len(values) == 0 {
			continue
		}

		// --------------------
		// Paginación normal: ej. GET /api/ventas?page=123&per_page=50
		// --------------------
		if rawKey == "page" {
			if page, err := strconv.Atoi(values[0]); err == nil && page > 0 {
				qf.Page = page
			}
			continue
		}

		if rawKey == "per_page" {
			if perPage, err := strconv.Atoi(values[0]); err == nil && perPage > 0 {
				qf.PerPage = perPage
			}
			continue
		}

		// --------------------
		// Cursor normal:    ej. GET /api/ventas?cursor=643a1f2e5c8e4d3a7b2f1a9b
		// Cursor normal:    ej. GET /api/ventas?cursor[asc]=643a1f2e5c8e4d3a7b2f1a9b
		// Cursor invertido: ej. GET /api/ventas?cursor[desc]=643a1f2e5c8e4d3a7b2f1a9b
		// --------------------
		if rawKey == "cursor[asc]" || rawKey == "cursor" {
			qf.Cursor = values[0]
			qf.CursorDirection = 1
			continue
		}

		if rawKey == "cursor[desc]" {
			qf.Cursor = values[0]
			qf.CursorDirection = -1
			continue
		}

		// --------------------
		// Sort normal: ej. GET /api/ventas?sort=cliente_id,producto_id
		// Sort multi:  ej. GET /api/ventas?sort=cliente_id&sort=producto_id
		// Sort bad but it works:  ej. GET /api/ventas?sort=cliente&sort=producto,category
		// --------------------
		if rawKey == "sort" {
			for _, value := range values {
				fields := strings.Split(value, ",")
				for _, field := range fields {
					field = strings.TrimSpace(field)
					if field == "" {
						continue
					}
					direction := 1
					if strings.HasPrefix(field, "-") {
						direction = -1
						field = field[1:]
					} else if strings.HasPrefix(field, "+") {
						field = field[1:]
					} else if strings.HasPrefix(field, "[asc]") {
						field = field[5:]
					} else if strings.HasPrefix(field, "[desc]") {
						direction = -1
						field = field[6:]
					} else if strings.HasSuffix(field, ":asc") {
						field = field[:4]
					} else if strings.HasSuffix(field, ":desc") {
						direction = -1
						field = field[:5]
					}
					allowed := allowFilters[field]
					// if slices.Contains(allow, "sortable") {
					// 	qf.AppendGrouBy(field)
					// }
					for _, allow := range allowed {
						if allow == "sortable" {
							qf.AppendSort(field, direction)
							break
						}
					}
				}
			}
			continue
		}

		// --------------------
		// GroupBy normal: ej. GET /api/ventas?groupby=cliente_id,producto_id
		// GroupBy multi:  ej. GET /api/ventas?groupby=cliente_id&groupby=producto_id
		// GroupBy bad but it works:  ej. GET /api/ventas?groupby=cliente_id&groupby=producto_id,category_id
		// --------------------
		if rawKey == "groupby" {
			for _, value := range values {
				fields := strings.Split(value, ",")
				for _, field := range fields {
					allowed := allowFilters[field]
					// if slices.Contains(allow, "groupable") {
					// 	qf.AppendGrouBy(field)
					// }
					for _, allow := range allowed {
						if allow == "groupable" {
							qf.AppendGrouBy(field)
							break
						}
					}
				}
			}
			continue
		}
		// --------------------
		// trash normal:  ej. GET /api/ventas?trash=only
		// trash normal:  ej. GET /api/ventas?trash=with
		// trash default: ej. GET /api/ventas?trash=without
		// --------------------
		if rawKey == "trash" {
			if values[0] == "without" || values[0] == "0" {
				qf.Trash = 0
			} else if values[0] == "with" || values[0] == "1" {
				qf.Trash = 1
			} else if values[0] == "only" || values[0] == "2" {
				qf.Trash = 2
			}
			continue
		}
		// --------------------
		// Filtro normal:  ej. GET /api/ventas?name[lk]=Andres
		// Filtro default: ej. GET /api/ventas?name=Andres -> GET /api/ventas?name[eq]=Andres
		// --------------------
		var key, filter, value string
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

		for _, af := range allowed {
			if af == filter {
				qf.AppendFilter(key, filter, value)
				break
			}
		}
	}

	return qf
}

func (ctx *Context) WriteCSV(fileName string, data any, comma ...rune) {
	val := reflect.ValueOf(data)

	if val.Kind() != reflect.Slice {
		err := errors.NewError(
			http.StatusInternalServerError,
			"Error al escribir el csv",
			errors.New(lang.TT(ctx.Lang(), "Los datos no son un slice de structs")),
		)
		ctx.WriteError(err)
		return
	}

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	del := ';'
	if len(comma) > 0 {
		del = comma[0]
	}
	writer.Comma = del

	if val.Len() == 0 {
		err := errors.NoDocuments(errors.New(lang.TT(ctx.Lang(), "No hay datos")))
		ctx.WriteError(err)
		return
	}

	first := val.Index(0)
	elemType := first.Type()

	var headers []string
	var fields []int // Índices de campos válidos

	// Encabezados filtrando por tag json
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		// Cortar por coma por si hay `json:"name,omitempty"`
		tag = strings.Split(tag, ",")[0]
		headers = append(headers, tag)
		fields = append(fields, i)
	}
	writer.Write(headers)

	// Datos
	for i := 0; i < val.Len(); i++ {
		var record []string
		elem := val.Index(i)

		for _, j := range fields {
			fieldVal := elem.Field(j)
			switch fieldVal.Kind() {
			case reflect.String:
				record = append(record, fieldVal.String())
			case reflect.Int, reflect.Int64:
				record = append(record, fmt.Sprintf("%d", fieldVal.Int()))
			case reflect.Float64:
				record = append(record, fmt.Sprintf("%f", fieldVal.Float()))
			case reflect.Bool:
				record = append(record, fmt.Sprintf("%t", fieldVal.Bool()))
			case reflect.Struct:
				if t, ok := fieldVal.Interface().(time.Time); ok {
					record = append(record, t.Format(time.RFC3339))
				} else {
					jsonVal, _ := json.Marshal(fieldVal.Interface())
					record = append(record, string(jsonVal))
				}
			case reflect.Slice, reflect.Map, reflect.Array:
				jsonVal, _ := json.Marshal(fieldVal.Interface())
				record = append(record, string(jsonVal))
			default:
				record = append(record, fmt.Sprintf("%v", fieldVal.Interface()))
				// Para cualquier otro tipo (interface, pointer, etc.)
				// jsonVal, _ := json.Marshal(fieldVal.Interface())
				// record = append(record, string(jsonVal))
			}
		}
		writer.Write(record)
	}
	writer.Flush()

	ctx.Writer.Header().Set("Content-Type", "text/csv")
	ctx.Writer.Header().Set("Content-Disposition", "attachment;filename="+fileName+".csv")
	ctx.Writer.Write(buffer.Bytes())
}

func Fill(model any, request any) {
	modelVal := reflect.ValueOf(model)
	requestVal := reflect.ValueOf(request)

	// Deben ser punteros a structs
	if modelVal.Kind() != reflect.Ptr || requestVal.Kind() != reflect.Ptr {
		return
	}

	modelVal = modelVal.Elem()
	requestVal = requestVal.Elem()

	if modelVal.Kind() != reflect.Struct || requestVal.Kind() != reflect.Struct {
		return
	}

	modelType := modelVal.Type()

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		jsonTag := field.Tag.Get("json")

		// Ignorar si no tiene tag json o si es '-'
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		jsonName := strings.Split(jsonTag, ",")[0]
		if jsonName == "" || jsonName == "-" {
			continue
		}

		// Buscar el campo en request por el mismo nombre
		reqField := requestVal.FieldByName(field.Name)
		if !reqField.IsValid() {
			continue
		}

		// Solo si el campo es asignable y compatible
		modelField := modelVal.Field(i)
		if modelField.CanSet() && reqField.Type().AssignableTo(modelField.Type()) {
			modelField.Set(reqField)
		}
	}
}
