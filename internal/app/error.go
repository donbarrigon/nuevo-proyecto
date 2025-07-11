package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver"
)

type Error interface {
	error
	GetStatus() int
	GetMessage() string
	GetErr() any
	SetStatus(code int)
	SetMessage(format string, ph ...F)
	SetErr(format any, ph ...F)
	Translate(lang string)                      // traduce el error
	Append(err *FieldError)                     // agrega error si es que hay con struct.
	Appendf(field string, text string, ph ...F) // agrega error si es que hay con valores por separado.
	Errors() Error                              // para el retorno de multiples errores usado en el request
}

type Err struct {
	Status    int                 `json:"-"`
	Message   string              `json:"message"`
	Err       any                 `json:"errors,omitempty"`
	ErrMap    map[string][]string `json:"-"`
	phMessage Fields              `json:"-"`
	phMap     map[string][]Fields `json:"-"`
}

type FieldError struct {
	FieldName    string
	Message      string
	Placeholders Fields
}

// variable para azucar sintactico
var Errors = Err{
	Status:    0,
	Message:   "",
	Err:       nil,
	ErrMap:    make(map[string][]string),
	phMessage: make(Fields, 0),
	phMap:     make(map[string][]Fields),
}

func (e *Err) New(status int, message string, err any, ph ...F) Error {
	return &Err{
		Status:    status,
		Message:   message,
		Err:       err,
		ErrMap:    make(map[string][]string),
		phMessage: ph,
	}
}

func (e *Err) NewEmpty() Error {
	return &Err{
		ErrMap: make(map[string][]string),
		phMap:  make(map[string][]Fields),
	}
}

func (e *Err) Error() string {
	msg := interpolatePlaceholders(e.Message, e.phMessage)
	return fmt.Sprintf("[%v] %v: \n%v", e.Status, msg, e.GetErr())
}

func (e *Err) GetStatus() int {
	return e.Status
}

func (e *Err) GetMessage() string {
	return e.Message
}

func (e *Err) GetErr() any {
	if e.ErrMap != nil {
		return e.ErrMap
	}
	if str, ok := e.Err.(string); ok {
		return interpolatePlaceholders(str, e.phMessage)
	}
	return e.Err
}

func (e *Err) SetStatus(code int) {
	e.Status = code
}

func (e *Err) SetMessage(format string, ph ...F) {
	e.Message = format
	if len(e.phMessage) == 0 {
		e.phMessage = ph
	} else {
		e.phMessage = append(e.phMessage, ph...)
	}
}

func (e *Err) SetErr(format any, ph ...F) {
	e.Err = format
	if len(e.phMessage) == 0 {
		e.phMessage = ph
	} else {
		e.phMessage = append(e.phMessage, ph...)
	}
}

func (e *Err) Translate(l string) {

	e.Message = Translate(l, e.Message, e.phMessage...)

	if e.ErrMap == nil {
		switch errVal := e.Err.(type) {
		case string:
			e.Err = Translate(l, errVal, e.phMessage...)
		case error:
			e.Err = Translate(l, errVal.Error())
		}
	} else {
		for key, messages := range e.ErrMap {
			for k, msg := range messages {
				e.ErrMap[key][k] = Translate(l, msg, e.phMap[key][k]...)
			}
		}
	}
}

func (e *Err) Appendf(field string, text string, ph ...F) {
	if text != "" {
		e.ErrMap[field] = append(e.ErrMap[field], text)
		e.phMap[field] = append(e.phMap[field], ph)
	}
}

func (e *Err) Append(err *FieldError) {
	if err != nil {
		e.ErrMap[err.FieldName] = append(e.ErrMap[err.FieldName], err.Message)
		e.phMap[err.FieldName] = append(e.phMap[err.FieldName], err.Placeholders)
	}
}

func (e *Err) Errors() Error {
	if e.ErrMap == nil {
		return nil
	}
	e.Status = http.StatusUnprocessableEntity
	e.Message = "Error de validación"
	e.Err = e.ErrMap
	return e
}

// ---------------------------------------------------------------- //
// funciones para crear erores estandarizados
func (e *Err) Mongo(err error) Error {

	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		return Errors.NoDocuments(err)
	case errors.Is(err, mongo.ErrClientDisconnected):
		return Errors.ClientDisconnected(err)
	case errors.Is(err, mongo.ErrNilDocument):
		return Errors.BadRequest(err)
	case errors.Is(err, context.DeadlineExceeded):
		return Errors.Timeout(err)
	case errors.As(err, &mongo.WriteException{}):
		return Errors.HandleWriteException(err)
	case errors.As(err, &mongo.CommandError{}):
		return Errors.Command(err)
	case errors.As(err, &mongo.BulkWriteException{}):
		return Errors.BulkWrite(err)
	case errors.As(err, &driver.Error{}):
		return Errors.Driver(err)
	default:
		return Errors.Unknown(err)
	}
}

func (e *Err) NotFound(err error) Error {
	return &Err{
		Status:  http.StatusNotFound,
		Message: "El recurso no existe",
		Err:     err.Error(),
	}
}

func (e *Err) NoDocuments(err error) Error {
	return &Err{
		Status:  http.StatusNotFound,
		Message: "No se encontraron registros",
		Err:     err.Error(),
	}
}

func (e *Err) ClientDisconnected(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Cliente desconectado",
		Err:     err.Error(),
	}
}

func (e *Err) BadRequest(err error) Error {
	return &Err{
		Status:  http.StatusBadRequest,
		Message: "Solicitud incorrecta",
		Err:     err.Error(),
	}
}

func (e *Err) Timeout(err error) Error {
	return &Err{
		Status:  http.StatusRequestTimeout,
		Message: "Tiempo de espera agotado",
		Err:     err.Error(),
	}
}

func (e *Err) HandleWriteException(err error) Error {
	var writeEx mongo.WriteException
	if errors.As(err, &writeEx) {
		for _, we := range writeEx.WriteErrors {
			if we.Code == 11000 {
				return &Err{
					Status:  http.StatusConflict,
					Message: "Registro duplicado",
					Err:     err.Error(),
				}
			}
		}
	}
	return Errors.Write(err)
}

func (e *Err) Write(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Write error",
		Err:     err.Error(),
	}
}

func (e *Err) Update(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Update error",
		Err:     err.Error(),
	}
}

func (e *Err) Delete(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Delete error",
		Err:     err.Error(),
	}
}

func (e *Err) Restore(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Restore error",
		Err:     err.Error(),
	}
}

func (e *Err) ForceDelete(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Permanent delete error",
		Err:     err.Error(),
	}
}

func (e *Err) Command(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Command error",
		Err:     err.Error(),
	}
}

func (e *Err) BulkWrite(err error) Error {
	return &Err{
		Status:  http.StatusBadRequest,
		Message: "Bulk write error",
		Err:     err.Error(),
	}
}

func (e *Err) Driver(err error) Error {
	return &Err{
		Status:  http.StatusBadGateway,
		Message: "Driver error",
		Err:     err.Error(),
	}
}

func (e *Err) Unknown(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Unexpected error",
		Err:     err.Error(),
	}
}

func (e *Err) Unauthorized(err error) Error {
	return &Err{
		Status:  http.StatusUnauthorized,
		Message: "Unauthorized",
		Err:     err.Error(),
	}
}

func (e *Err) HexID(err error) Error {
	return &Err{
		Status:  http.StatusBadRequest,
		Message: "Invalid hexadecimal ID",
		Err:     err.Error(),
	}
}

func (e *Err) NotFoundf(format string, ph ...F) Error {
	return &Err{
		Status:    http.StatusNotFound,
		Message:   "Resource not found",
		Err:       format,
		phMessage: ph,
	}
}

func (e *Err) NoDocumentsf(format string, ph ...F) Error {
	return &Err{
		Status:    http.StatusNotFound,
		Message:   "No documents found",
		Err:       format,
		phMessage: ph,
	}
}

func (e *Err) ClientDisconnectedf(format string, ph ...F) Error {
	return &Err{
		Status:    http.StatusInternalServerError,
		Message:   "Client disconnected",
		Err:       format,
		phMessage: ph,
	}
}

func (e *Err) BadRequestf(format string, ph ...F) Error {
	return &Err{
		Status:    http.StatusBadRequest,
		Message:   "Bad request",
		Err:       format,
		phMessage: ph,
	}
}

func (e *Err) Timeoutf(format string, ph ...F) Error {
	return &Err{
		Status:    http.StatusRequestTimeout,
		Message:   "Request timeout",
		Err:       format,
		phMessage: ph,
	}
}

func (e *Err) Writef(format string, ph ...F) Error {
	return &Err{
		Status:    http.StatusInternalServerError,
		Message:   "Write error",
		Err:       format,
		phMessage: ph,
	}
}

func (e *Err) Updatef(format string, ph ...F) Error {
	return &Err{
		Status:    http.StatusInternalServerError,
		Message:   "Update error",
		Err:       format,
		phMessage: ph,
	}
}

func (e *Err) Deletef(format string, ph ...F) Error {
	return &Err{
		Status:    http.StatusInternalServerError,
		Message:   "Delete error",
		Err:       format,
		phMessage: ph,
	}
}

func (e *Err) Restoref(format string, ph ...F) Error {
	return &Err{
		Status:    http.StatusInternalServerError,
		Message:   "Restore error",
		Err:       format,
		phMessage: ph,
	}
}

func (e *Err) ForceDeletef(format string, ph ...F) Error {
	return &Err{
		Status:    http.StatusInternalServerError,
		Message:   "Permanent delete error",
		Err:       format,
		phMessage: ph,
	}
}

func (e *Err) Commandf(format string, ph ...F) Error {
	return &Err{
		Status:    http.StatusInternalServerError,
		Message:   "Command error",
		Err:       format,
		phMessage: ph,
	}
}

func (e *Err) BulkWritef(format string, ph ...F) Error {
	return &Err{
		Status:    http.StatusBadRequest,
		Message:   "Bulk write error",
		Err:       format,
		phMessage: ph,
	}
}

func (e *Err) Driverf(format string, ph ...F) Error {
	return &Err{
		Status:    http.StatusBadGateway,
		Message:   "Driver error",
		Err:       format,
		phMessage: ph,
	}
}

func (e *Err) Unknownf(format string, ph ...F) Error {
	return &Err{
		Status:    http.StatusInternalServerError,
		Message:   "Unexpected error",
		Err:       format,
		phMessage: ph,
	}
}

func (e *Err) Unauthorizedf(format string, ph ...F) Error {
	return &Err{
		Status:    http.StatusUnauthorized,
		Message:   "Unauthorized",
		Err:       format,
		phMessage: ph,
	}
}

func (e *Err) HexIDf(format string, ph ...F) Error {
	return &Err{
		Status:    http.StatusBadRequest,
		Message:   "Invalid hexadecimal ID",
		Err:       format,
		phMessage: ph,
	}
}
