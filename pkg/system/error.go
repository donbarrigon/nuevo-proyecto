package system

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
	SetMessage(format string, a ...any)
	SetErr(format any, a ...any)
	Translate(lang string)           // traduce el error
	Append(field string, err string) // agrega error si es que hay
	Errors() Error                   // para el retorno de multiples errores usado en el request
}

type Err struct {
	Status  int                 `json:"-"`
	Message string              `json:"message"`
	Err     any                 `json:"error,omitempty"`
	ErrMap  map[string][]string `json:"-"`
	am      []any               `json:"-"`
	ae      []any               `json:"-"`
}

// variable para azucar sintactico
var Errors = Err{
	Status:  0,
	Message: "",
	Err:     nil,
	ErrMap:  make(map[string][]string),
}

func (e *Err) New() Error {
	return &Err{
		ErrMap: make(map[string][]string),
	}
}

func (e *Err) Error() string {
	return fmt.Sprintf("[%v] %v: \n%v", e.Status, e.Message, e.Err)
}

func (e *Err) GetStatus() int {
	return e.Status
}

func (e *Err) GetMessage() string {
	return fmt.Sprintf(e.Message, e.am...)
}

func (e *Err) GetErr() any {
	if e.ErrMap != nil {
		return e.ErrMap
	}
	if str, ok := e.Err.(string); ok {
		return fmt.Sprintf(str, e.ae...)
	}
	return e.Err
}

func (e *Err) SetStatus(code int) {
	e.Status = code
}

func (e *Err) SetMessage(format string, a ...any) {
	e.Message = format
	e.am = a
}

func (e *Err) SetErr(format any, a ...any) {
	e.Err = format
	e.ae = a
}

func (e *Err) Translate(l string) {
	e.Message = Translate(l, e.Message, e.am...)

	switch errVal := e.Err.(type) {
	case string:
		e.Err = Translate(l, errVal, e.ae...)
	case error:
		e.Err = Translate(l, errVal.Error())
	}
}

func (e *Err) Append(field string, text string) {
	// if e.ErrMap == nil {
	// 	e.ErrMap = make(map[string][]string)
	// }
	if text != "" {
		e.ErrMap[field] = append(e.ErrMap[field], text)
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

// -------------------------------- //
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
		Message: "Error al escribir el registro",
		Err:     err.Error(),
	}
}

func (e *Err) Update(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Error al modificar el registro",
		Err:     err.Error(),
	}
}

func (e *Err) Delete(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Error al eliminar el registro",
		Err:     err.Error(),
	}
}

func (e *Err) Restore(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Error al restaurar el registro",
		Err:     err.Error(),
	}
}

func (e *Err) ForceDelete(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Error al eliminar el registro permanentemente",
		Err:     err.Error(),
	}
}

func (e *Err) Command(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Error de comando",
		Err:     err.Error(),
	}
}

func (e *Err) BulkWrite(err error) Error {
	return &Err{
		Status:  http.StatusBadRequest,
		Message: "Error en escritura masiva",
		Err:     err.Error(),
	}
}

func (e *Err) Driver(err error) Error {
	return &Err{
		Status:  http.StatusBadGateway,
		Message: "Error del driver",
		Err:     err.Error(),
	}
}

func (e *Err) Unknown(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Error inesperado",
		Err:     err.Error(),
	}
}

func (e *Err) Unauthorized(err error) Error {
	return &Err{
		Status:  http.StatusUnauthorized,
		Message: "No autorizado",
		Err:     err.Error(),
	}
}

func (e *Err) HexID(err error) Error {
	return &Err{
		Status:  http.StatusBadRequest,
		Message: "El id no es un hexadecimal válido",
		Err:     err.Error(),
	}
}

func (e *Err) NotFoundf(format string, a ...any) Error {
	return &Err{
		Status:  http.StatusNotFound,
		Message: "El recurso no existe",
		Err:     fmt.Sprintf(format, a...),
	}
}

func (e *Err) NoDocumentsf(format string, a ...any) Error {
	return &Err{
		Status:  http.StatusNotFound,
		Message: "No se encontraron registros",
		Err:     fmt.Sprintf(format, a...),
	}
}

func (e *Err) ClientDisconnectedf(format string, a ...any) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Cliente desconectado",
		Err:     fmt.Sprintf(format, a...),
	}
}

func (e *Err) BadRequestf(format string, a ...any) Error {
	return &Err{
		Status:  http.StatusBadRequest,
		Message: "Solicitud incorrecta",
		Err:     fmt.Sprintf(format, a...),
	}
}

func (e *Err) Timeoutf(format string, a ...any) Error {
	return &Err{
		Status:  http.StatusRequestTimeout,
		Message: "Tiempo de espera agotado",
		Err:     fmt.Sprintf(format, a...),
	}
}

func (e *Err) Writef(format string, a ...any) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Error al escribir el registro",
		Err:     fmt.Sprintf(format, a...),
	}
}

func (e *Err) Updatef(format string, a ...any) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Error al modificar el registro",
		Err:     fmt.Sprintf(format, a...),
	}
}

func (e *Err) Deletef(format string, a ...any) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Error al eliminar el registro",
		Err:     fmt.Sprintf(format, a...),
	}
}

func (e *Err) Restoref(format string, a ...any) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Error al restaurar el registro",
		Err:     fmt.Sprintf(format, a...),
	}
}

func (e *Err) ForceDeletef(format string, a ...any) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Error al eliminar el registro permanentemente",
		Err:     fmt.Sprintf(format, a...),
	}
}

func (e *Err) Commandf(format string, a ...any) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Error de comando",
		Err:     fmt.Sprintf(format, a...),
	}
}

func (e *Err) BulkWritef(format string, a ...any) Error {
	return &Err{
		Status:  http.StatusBadRequest,
		Message: "Error en escritura masiva",
		Err:     fmt.Sprintf(format, a...),
	}
}

func (e *Err) Driverf(format string, a ...any) Error {
	return &Err{
		Status:  http.StatusBadGateway,
		Message: "Error del driver",
		Err:     fmt.Sprintf(format, a...),
	}
}

func (e *Err) Unknownf(format string, a ...any) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Error inesperado",
		Err:     fmt.Sprintf(format, a...),
	}
}

func (e *Err) Unauthorizedf(format string, a ...any) Error {
	return &Err{
		Status:  http.StatusUnauthorized,
		Message: "No autorizado",
		Err:     fmt.Sprintf(format, a...),
	}
}

func (e *Err) HexIDf(format string, a ...any) Error {
	return &Err{
		Status:  http.StatusBadRequest,
		Message: "El id no es un hexadecimal válido",
		Err:     fmt.Sprintf(format, a...),
	}
}
