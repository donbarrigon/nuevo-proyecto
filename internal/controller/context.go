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
	"github.com/donbarrigon/nuevo-proyecto/internal/app/request"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
	"github.com/donbarrigon/nuevo-proyecto/pkg/system"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MessageResource struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Context struct {
	Writer     http.ResponseWriter
	Request    *http.Request
	PathParams map[string]string
	User       *model.User
	Token      *model.Token
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

func (ctx *Context) GetBody(request any) system.Error {
	decoder := json.NewDecoder(ctx.Request.Body)
	if err := decoder.Decode(request); err != nil {
		return system.Errors.New(
			http.StatusBadRequest,
			"El cuerpo de la solicitud es incorrecto",
			"No se pudo decodificar el cuerpo de la solicitud: {error}",
			system.F{Key: "error", Value: err.Error()},
		)
	}
	defer ctx.Request.Body.Close()
	return nil
}

func (ctx *Context) GetMultiPartForm(req any) system.Error {
	err := ctx.Request.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		return &system.Err{
			Status:  http.StatusBadRequest,
			Message: system.Translate(ctx.Lang(), "El formulario no se pudo procesar"),
			Err:     system.Translate(ctx.Lang(), "Error al analizar el formulario: %v", err.Error()),
		}
	}

	form := ctx.Request.MultipartForm
	if form == nil {
		return &system.Err{
			Status:  http.StatusBadRequest,
			Message: system.Translate(ctx.Lang(), "Formulario no válido"),
			Err:     system.Translate(ctx.Lang(), "No se encontró un formulario multipart válido"),
		}
	}

	v := reflect.ValueOf(req)
	t := reflect.TypeOf(req)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		tag := field.Tag.Get("json")
		if tag == "-" {
			continue
		}

		var formKey string
		if tag != "" {
			formKey = strings.Split(tag, ",")[0]
		}
		if formKey == "" {
			formKey = field.Name
		}

		values := form.Value[formKey]
		if len(values) == 0 || !fieldValue.CanSet() {
			continue
		}

		switch field.Type.Kind() {
		case reflect.String:
			fieldValue.SetString(values[0])

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if intVal, err := strconv.ParseInt(values[0], 10, 64); err == nil {
				fieldValue.SetInt(intVal)
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if uintVal, err := strconv.ParseUint(values[0], 10, 64); err == nil {
				fieldValue.SetUint(uintVal)
			}

		case reflect.Float32, reflect.Float64:
			if floatVal, err := strconv.ParseFloat(values[0], 64); err == nil {
				fieldValue.SetFloat(floatVal)
			}

		case reflect.Bool:
			if boolVal, err := strconv.ParseBool(values[0]); err == nil {
				fieldValue.SetBool(boolVal)
			}

		case reflect.Slice:
			elemKind := field.Type.Elem().Kind()
			slice := reflect.MakeSlice(field.Type, 0, len(values))

			for _, val := range values {
				var converted reflect.Value
				switch elemKind {
				case reflect.String:
					converted = reflect.ValueOf(val)

				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if intVal, err := strconv.ParseInt(val, 10, 64); err == nil {
						converted = reflect.ValueOf(intVal).Convert(field.Type.Elem())
					}

				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					if uintVal, err := strconv.ParseUint(val, 10, 64); err == nil {
						converted = reflect.ValueOf(uintVal).Convert(field.Type.Elem())
					}

				case reflect.Float32, reflect.Float64:
					if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
						converted = reflect.ValueOf(floatVal).Convert(field.Type.Elem())
					}

				case reflect.Bool:
					if boolVal, err := strconv.ParseBool(val); err == nil {
						converted = reflect.ValueOf(boolVal).Convert(field.Type.Elem())
					}
				}

				if converted.IsValid() {
					slice = reflect.Append(slice, converted)
				}
			}

			fieldValue.Set(slice)
		}
	}

	return nil
}

// GetMultiPartForm analiza multipart/form-data y asigna los valores al struct `req`
// func (ctx *Context) GetMultiPartForm(req any) system.Error {
// 	err := ctx.Request.ParseMultipartForm(32 << 20) // 32 MB por defecto
// 	if err != nil {
// 		return &system.Err{
// 			Status:  http.StatusBadRequest,
// 			Message: "El formulario no se pudo procesar",
// 			Err:     system.Translate(ctx.Lang(), "Error al analizar el formulario") + ": " + err.Error(),
// 		}
// 	}

// 	form := ctx.Request.MultipartForm
// 	if form == nil {
// 		return &system.Err{
// 			Status:  http.StatusBadRequest,
// 			Message: "Formulario no válido",
// 			Err:     system.Translate(ctx.Lang(), "No se encontró un formulario multipart válido"),
// 		}
// 	}

// 	v := reflect.ValueOf(req)
// 	t := reflect.TypeOf(req)

// 	// Si es puntero, desreferenciar
// 	if t.Kind() == reflect.Ptr {
// 		t = t.Elem()
// 		v = v.Elem()
// 	}

// 	for i := 0; i < t.NumField(); i++ {
// 		field := t.Field(i)
// 		formKey := field.Tag.Get("form")
// 		if formKey == "" {
// 			formKey = field.Name
// 		}

// 		values := form.Value[formKey]
// 		if len(values) == 0 {
// 			continue // No hay valor para este campo
// 		}

// 		fieldValue := v.Field(i)
// 		if !fieldValue.CanSet() {
// 			continue
// 		}

// 		// Solo manejamos string por simplicidad aquí
// 		if field.Type.Kind() == reflect.String {
// 			fieldValue.SetString(values[0])
// 		}
// 		// Puedes extender esto para otros tipos si lo necesitas
// 	}

// 	return nil
// }

func (ctx *Context) ValidateBody(req any) system.Error {
	if err := ctx.GetBody(req); err != nil {
		return err
	}

	v := reflect.ValueOf(req)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			if err := request.Validate(ctx.Lang(), item); err != nil {
				return err
			}
		}
	default:
		if err := request.Validate(ctx.Lang(), req); err != nil {
			return err
		}
	}

	return nil
}

func (ctx *Context) ValidateMultiPartForm(req any) system.Error {
	if err := ctx.GetMultiPartForm(req); err != nil {
		return err
	}
	if err := request.Validate(ctx.Lang(), req); err != nil {
		return err
	}
	return nil
}

func (ctx *Context) ValidateRequest(req any) system.Error {
	contentType := ctx.Request.Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(contentType, "multipart/form-data"):
		return ctx.ValidateMultiPartForm(req)

	case strings.HasPrefix(contentType, "application/json"),
		strings.HasPrefix(contentType, "application/*+json"):
		return ctx.ValidateBody(req)

	default:
		return &system.Err{
			Status:  http.StatusUnsupportedMediaType,
			Message: system.Translate(ctx.Lang(), "Tipo de contenido no soportado"),
			Err:     system.Translate(ctx.Lang(), "Tipo de contenido no soportado: %v", contentType),
		}
	}
}

func (ctx *Context) Get(param string, defaultValue ...string) string {
	if value := ctx.PathParams["id"]; value != "" {
		return value
	}
	value := ctx.Request.URL.Query().Get(param)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
	}
	return value
}

func (ctx *Context) WriteJSON(status int, data any) {
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.WriteHeader(status)

	if err := json.NewEncoder(ctx.Writer).Encode(data); err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.Writer.WriteHeader(500)
		ctx.Writer.Write([]byte(fmt.Sprintf(`{"message":"Error","error":"%s"}`, system.Translate(ctx.Lang(), "No se pudo codificar la respuesta"))))
	}
}

func (ctx *Context) WriteError(err system.Error) {
	err.Translate(ctx.Lang())
	ctx.WriteJSON(err.GetStatus(), err)
}

func (ctx *Context) WriteNotFound() {
	ctx.WriteError(system.Errors.NotFoundf("El recurso [%v:%v] no existe", ctx.Request.Method, ctx.Request.URL.Path))
}

func (ctx *Context) WriteMessage(code int, data any, message string, v ...any) {
	ctx.WriteJSON(code, &MessageResource{
		Message: system.Translate(ctx.Lang(), message, v...),
		Data:    data,
	})
}

func (ctx *Context) WriteSuccess(data any) {
	ctx.WriteJSON(http.StatusOK, &MessageResource{
		Message: system.Translate(ctx.Lang(), "Solicitud procesada con éxito"),
		Data:    data,
	})
}

func (ctx *Context) WriteCreated(data any) {
	ctx.WriteJSON(http.StatusCreated, &MessageResource{
		Message: system.Translate(ctx.Lang(), "Recurso creado exitosamente"),
		Data:    data,
	})
}

func (ctx *Context) WriteUpdated(data any) {
	ctx.WriteJSON(http.StatusOK, &MessageResource{
		Message: system.Translate(ctx.Lang(), "Recurso actualizado exitosamente"),
		Data:    data,
	})
}

func (ctx *Context) WriteDeleted(data any) {
	ctx.WriteJSON(http.StatusOK, &MessageResource{
		Message: system.Translate(ctx.Lang(), "Recurso eliminado exitosamente"),
		Data:    data,
	})
}

func (ctx *Context) WriteRestored(data any) {
	ctx.WriteJSON(http.StatusOK, &MessageResource{
		Message: system.Translate(ctx.Lang(), "Recurso restaurado exitosamente"),
		Data:    data,
	})
}

func (ctx *Context) WriteForceDeleted(data any) {
	ctx.WriteJSON(http.StatusOK, &MessageResource{
		Message: system.Translate(ctx.Lang(), "Recurso eliminado permanentemente"),
		Data:    data,
	})
}

func (ctx *Context) WriteNoContent() {
	ctx.Writer.WriteHeader(http.StatusNoContent)
}

func (ctx *Context) TT(s string, v ...any) string {
	return system.Translate(ctx.Lang(), s, v...)
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
		err := &system.Err{
			Status:  http.StatusInternalServerError,
			Message: "Error al escribir el csv",
			Err:     system.Translate(ctx.Lang(), "Los datos no son un slice de structs"),
		}
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
		err := system.Errors.NoDocumentsf(system.Translate(ctx.Lang(), "No hay datos"))
		ctx.WriteError(err)
		return
	}

	first := val.Index(0)
	elemType := first.Type()

	var headers []string
	var fields []int

	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		tag = strings.Split(tag, ",")[0]
		headers = append(headers, tag)
		fields = append(fields, i)
	}
	writer.Write(headers)

	for i := 0; i < val.Len(); i++ {
		var record []string
		elem := val.Index(i)

		for _, j := range fields {
			fieldVal := elem.Field(j)

			if fieldVal.Type() == reflect.TypeOf(bson.ObjectID{}) {
				objID := fieldVal.Interface().(bson.ObjectID)
				record = append(record, objID.Hex()) // sin comillas manuales
				continue
			}

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
			}
		}
		writer.Write(record)
	}
	writer.Flush()

	ctx.Writer.Header().Set("Content-Type", "text/csv")
	ctx.Writer.Header().Set("Content-Disposition", "attachment;filename="+fileName+".csv")
	ctx.Writer.Write(buffer.Bytes())
}

// func (ctx *Context) WriteCSV(fileName string, data any, comma ...rune) {
// 	val := reflect.ValueOf(data)

// 	if val.Kind() != reflect.Slice {
// 		err := errors.NewError(
// 			http.StatusInternalServerError,
// 			"Error al escribir el csv",
// 			errors.New(system.Translate(ctx.Lang(), "Los datos no son un slice de structs")),
// 		)
// 		ctx.WriteError(err)
// 		return
// 	}

// 	var buffer bytes.Buffer
// 	writer := csv.NewWriter(&buffer)

// 	del := ';'
// 	if len(comma) > 0 {
// 		del = comma[0]
// 	}
// 	writer.Comma = del

// 	if val.Len() == 0 {
// 		err := errors.NoDocuments(errors.New(system.Translate(ctx.Lang(), "No hay datos")))
// 		ctx.WriteError(err)
// 		return
// 	}

// 	first := val.Index(0)
// 	elemType := first.Type()

// 	var headers []string
// 	var fields []int // Índices de campos válidos

// 	// Encabezados filtrando por tag json
// 	for i := 0; i < elemType.NumField(); i++ {
// 		field := elemType.Field(i)
// 		tag := field.Tag.Get("json")
// 		if tag == "" || tag == "-" {
// 			continue
// 		}
// 		// Cortar por coma por si hay `json:"name,omitempty"`
// 		tag = strings.Split(tag, ",")[0]
// 		headers = append(headers, tag)
// 		fields = append(fields, i)
// 	}
// 	writer.Write(headers)

// 	// Datos
// 	for i := 0; i < val.Len(); i++ {
// 		var record []string
// 		elem := val.Index(i)

// 		for _, j := range fields {
// 			fieldVal := elem.Field(j)
// 			switch fieldVal.Kind() {
// 			case reflect.String:
// 				record = append(record, fieldVal.String())
// 			case reflect.Int, reflect.Int64:
// 				record = append(record, fmt.Sprintf("%d", fieldVal.Int()))
// 			case reflect.Float64:
// 				record = append(record, fmt.Sprintf("%f", fieldVal.Float()))
// 			case reflect.Bool:
// 				record = append(record, fmt.Sprintf("%t", fieldVal.Bool()))
// 			case reflect.Struct:
// 				if t, ok := fieldVal.Interface().(time.Time); ok {
// 					record = append(record, t.Format(time.RFC3339))
// 				} else {
// 					jsonVal, _ := json.Marshal(fieldVal.Interface())
// 					record = append(record, string(jsonVal))
// 				}
// 			case reflect.Slice, reflect.Map, reflect.Array:
// 				jsonVal, _ := json.Marshal(fieldVal.Interface())
// 				record = append(record, string(jsonVal))
// 			default:
// 				record = append(record, fmt.Sprintf("%v", fieldVal.Interface()))
// 				// Para cualquier otro tipo (interface, pointer, etc.)
// 				// jsonVal, _ := json.Marshal(fieldVal.Interface())
// 				// record = append(record, string(jsonVal))
// 			}
// 		}
// 		writer.Write(record)
// 	}
// 	writer.Flush()

// 	ctx.Writer.Header().Set("Content-Type", "text/csv")
// 	ctx.Writer.Header().Set("Content-Disposition", "attachment;filename="+fileName+".csv")
// 	ctx.Writer.Write(buffer.Bytes())
// }

// Fill llena los campos del modelo con los valores del request,
// pero solo si el campo del modelo tiene la etiqueta fillable
func Fill(model any, request any) system.Error {
	modelValue := reflect.ValueOf(model)
	requestValue := reflect.ValueOf(request)

	if modelValue.Kind() != reflect.Ptr || requestValue.Kind() != reflect.Ptr {
		return system.Errors.Unknownf("Los parámetros model y request deben ser punteros")
	}

	modelValue = modelValue.Elem()
	requestValue = requestValue.Elem()

	if modelValue.Kind() != reflect.Struct || requestValue.Kind() != reflect.Struct {
		return system.Errors.Unknownf("Los parámetros model y request deben ser estructuras")
	}

	modelType := modelValue.Type()

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		// if fillable, ok := field.Tag.Lookup("fillable"); ok && fillable == "true" {
		if _, ok := field.Tag.Lookup("fillable"); ok {
			fieldName := field.Name

			requestField := requestValue.FieldByName(fieldName)

			if requestField.IsValid() && requestField.Type().AssignableTo(field.Type) {
				modelField := modelValue.Field(i)
				if modelField.CanSet() {
					modelField.Set(requestField)
				}
			}
		}
	}
	return nil
}
