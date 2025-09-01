package app

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuthInterface interface {
	GetID() bson.ObjectID
	GetUserID() bson.ObjectID
	Can(permissionName ...string) Error
	HasRole(roleName ...string) Error
}

type Validator interface {
	PrepareForValidation(ctx *HttpContext) Error
}

type MessageResource struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type HttpContext struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Params  map[string]string
	Auth    AuthInterface
}

func NewHttpContext(w http.ResponseWriter, r *http.Request) *HttpContext {
	return &HttpContext{
		Writer:  w,
		Request: r,
	}
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
			phMessage: List{{"error", err.Error()}},
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
	if errPFV != nil {
		errPFV = errPFV.Errors()
	}
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
	return errPFV
}

func (ctx *HttpContext) GetParam(param string, defaultValue string) string {
	if value := ctx.Params[param]; value != "" {
		return value
	}
	return defaultValue
}

func (ctx *HttpContext) GetInput(param string) string {
	return ctx.Request.URL.Query().Get(param)
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
		Entry{"method", ctx.Request.Method},
		Entry{"path", ctx.Request.URL.Path},
	))
}

func (ctx *HttpContext) ResponseMessage(code int, data any, message string, ph ...Entry) {
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
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

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
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

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
