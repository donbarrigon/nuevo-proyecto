package app

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

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserInterface interface {
	GetID() bson.ObjectID
	Can(permissionName string) Error
}

type TokenInterface interface {
	GetID() bson.ObjectID
	Can(permissionName string) Error
}

type Validator interface {
	PrepareForValidation(ctx *HttpContext) Error
}

type MessageResource struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Auth struct {
	User  UserInterface
	Token TokenInterface
}

type HttpContext struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Params  map[string]string
	Auth    *Auth
}

func NewHttpContext(w http.ResponseWriter, r *http.Request) *HttpContext {
	return &HttpContext{
		Writer:  w,
		Request: r,
	}
}

func (a *Auth) Can(permissionName string) Error {
	return a.Token.Can(permissionName)
}

func (a *Auth) UserID() string {
	return a.User.GetID().Hex()
}

func (ctx *HttpContext) Lang() string {
	return ctx.Request.Header.Get("Accept-Language")
}

func (ctx *HttpContext) GetBody(request any) Error {
	decoder := json.NewDecoder(ctx.Request.Body)
	if err := decoder.Decode(request); err != nil {
		return &Err{
			Status:    http.StatusBadRequest,
			Message:   "The request body is invalid",
			Err:       "Could not decode the request body: {error}",
			phMessage: Fields{{Key: "error", Value: err.Error()}},
		}
	}
	defer ctx.Request.Body.Close()
	return nil
}

func (ctx *HttpContext) ValidateBody(req Validator) Error {

	if err := ctx.GetBody(req); err != nil {
		return err
	}

	errPFV := req.PrepareForValidation(ctx)

	// v := reflect.ValueOf(req)
	// if v.Kind() == reflect.Ptr {
	// 	v = v.Elem()
	// }

	// switch v.Kind() {
	// case reflect.Slice, reflect.Array:
	// 	for i := 0; i < v.Len(); i++ {
	// 		item := v.Index(i).Interface()
	// 		if err := Validate(item); err != nil {
	// 			return err
	// 		}
	// 	}
	// default:
	// 	if err := Validate(req); err != nil {
	// 		return err
	// 	}
	// }

	if err := Validate(ctx, req); err != nil {
		if errPFV != nil {
			errMap, phMap := errPFV.GetMap()
			for key, valor := range errMap {
				for i, msg := range valor {
					err.Append(&FieldError{
						FieldName:    key,
						Message:      msg,
						Placeholders: phMap[key][i],
					})
				}
			}
		}
		return err.Errors()
	}
	return errPFV.Errors()
}

// func (ctx *HttpContext) GetMultiPartForm(req any) Error {
// 	err := ctx.Request.ParseMultipartForm(32 << 20) // 32 MB
// 	if err != nil {
// 		return &Err{
// 			Status:    http.StatusBadRequest,
// 			Message:   "The form could not be processed",
// 			Err:       "Failed to parse the form: {error}",
// 			phMessage: Fields{{Key: "error", Value: err.Error()}},
// 		}
// 	}

// 	form := ctx.Request.MultipartForm
// 	if form == nil {
// 		return &Err{
// 			Status:  http.StatusBadRequest,
// 			Message: "Invalid form",
// 			Err:     "No valid multipart form found",
// 		}
// 	}

// 	v := reflect.ValueOf(req)
// 	t := reflect.TypeOf(req)

// 	if t.Kind() == reflect.Ptr {
// 		t = t.Elem()
// 		v = v.Elem()
// 	}

// 	for i := 0; i < t.NumField(); i++ {
// 		field := t.Field(i)
// 		fieldValue := v.Field(i)

// 		tag := field.Tag.Get("json")
// 		if tag == "-" {
// 			continue
// 		}

// 		var formKey string
// 		if tag != "" {
// 			formKey = strings.Split(tag, ",")[0]
// 		}
// 		if formKey == "" {
// 			formKey = field.Name
// 		}

// 		values := form.Value[formKey]
// 		if len(values) == 0 || !fieldValue.CanSet() {
// 			continue
// 		}

// 		switch field.Type.Kind() {
// 		case reflect.String:
// 			fieldValue.SetString(values[0])

// 		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 			if intVal, err := strconv.ParseInt(values[0], 10, 64); err == nil {
// 				fieldValue.SetInt(intVal)
// 			}

// 		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
// 			if uintVal, err := strconv.ParseUint(values[0], 10, 64); err == nil {
// 				fieldValue.SetUint(uintVal)
// 			}

// 		case reflect.Float32, reflect.Float64:
// 			if floatVal, err := strconv.ParseFloat(values[0], 64); err == nil {
// 				fieldValue.SetFloat(floatVal)
// 			}

// 		case reflect.Bool:
// 			if boolVal, err := strconv.ParseBool(values[0]); err == nil {
// 				fieldValue.SetBool(boolVal)
// 			}

// 		case reflect.Slice:
// 			elemKind := field.Type.Elem().Kind()
// 			slice := reflect.MakeSlice(field.Type, 0, len(values))

// 			for _, val := range values {
// 				var converted reflect.Value
// 				switch elemKind {
// 				case reflect.String:
// 					converted = reflect.ValueOf(val)

// 				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 					if intVal, err := strconv.ParseInt(val, 10, 64); err == nil {
// 						converted = reflect.ValueOf(intVal).Convert(field.Type.Elem())
// 					}

// 				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
// 					if uintVal, err := strconv.ParseUint(val, 10, 64); err == nil {
// 						converted = reflect.ValueOf(uintVal).Convert(field.Type.Elem())
// 					}

// 				case reflect.Float32, reflect.Float64:
// 					if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
// 						converted = reflect.ValueOf(floatVal).Convert(field.Type.Elem())
// 					}

// 				case reflect.Bool:
// 					if boolVal, err := strconv.ParseBool(val); err == nil {
// 						converted = reflect.ValueOf(boolVal).Convert(field.Type.Elem())
// 					}
// 				}

// 				if converted.IsValid() {
// 					slice = reflect.Append(slice, converted)
// 				}
// 			}

// 			fieldValue.Set(slice)
// 		}
// 	}

// 	return nil
// }

// func (ctx *HttpContext) ValidateMultiPartForm(req any) Error {
// 	if err := ctx.GetMultiPartForm(req); err != nil {
// 		return err
// 	}
// 	if err := Validate(req); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (ctx *HttpContext) ValidateRequest(req any) Error {
// 	contentType := ctx.Request.Header.Get("Content-Type")

// 	switch {
// 	case strings.HasPrefix(contentType, "multipart/form-data"):
// 		return ctx.ValidateMultiPartForm(req)

// 	case strings.HasPrefix(contentType, "application/json"),
// 		strings.HasPrefix(contentType, "application/*+json"):
// 		return ctx.ValidateBody(req)

// 	default:
// 		return &Err{
// 			Status:    http.StatusUnsupportedMediaType,
// 			Message:   "Unsupported content type",
// 			Err:       "Unsupported content type: {contentType}",
// 			phMessage: Fields{{Key: "contentType", Value: contentType}},
// 		}
// 	}
// }

func (ctx *HttpContext) GetParam(param string, defaultValue string) string {
	if value := ctx.Params[param]; value != "" {
		return value
	}
	return defaultValue
}

func (ctx *HttpContext) GetInput(param string, defaultValue string) string {
	if value := ctx.Request.URL.Query().Get(param); value != "" {
		return value
	}
	return defaultValue
}

func (ctx *HttpContext) ResponseJSON(status int, data any) {
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.WriteHeader(status)

	if err := json.NewEncoder(ctx.Writer).Encode(data); err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.Writer.WriteHeader(500)
		ctx.Writer.Write([]byte(Translate(ctx.Lang(), `{"message": "Error", "error": "Could not encode the response"}`)))
	}
}

func (ctx *HttpContext) ResponseError(err Error) {
	err.Translate(ctx.Lang())
	ctx.ResponseJSON(err.GetStatus(), err)
}

func (ctx *HttpContext) ResponseNotFound() {
	ctx.ResponseError(Errors.NotFoundf("The resource [{method}:{path}] does not exist",
		F{Key: "method", Value: ctx.Request.Method},
		F{Key: "path", Value: ctx.Request.URL.Path},
	))
}

func (ctx *HttpContext) ResponseMessage(code int, data any, message string, ph ...F) {
	ctx.ResponseJSON(code, &MessageResource{
		Message: Translate(ctx.Lang(), message, ph...),
		Data:    data,
	})
}

func (ctx *HttpContext) ResponseOk(data any) {
	ctx.ResponseJSON(http.StatusOK, data)
}

func (ctx *HttpContext) ResponseCreated(data any) {
	ctx.ResponseJSON(http.StatusCreated, data)
}

func (ctx *HttpContext) ResponseNoContent() {
	ctx.Writer.WriteHeader(http.StatusNoContent)
}

func (ctx *HttpContext) GetQueryFilter(allowFilters map[string][]string) *QueryFilter {
	query := ctx.Request.URL.Query()

	qf := NewQueryFilter()

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

func (ctx *HttpContext) ResponseCSV(fileName string, data any, comma ...rune) {
	val := reflect.ValueOf(data)

	if val.Kind() != reflect.Slice {
		err := &Err{
			Status:  http.StatusInternalServerError,
			Message: "Error writing CSV",
			Err:     "Data is not a slice of structs",
		}
		ctx.ResponseError(err)
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
		err := Errors.NoDocumentsf("No data available")
		ctx.ResponseError(err)
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
