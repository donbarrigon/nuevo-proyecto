package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
)

type Error interface {
	error
	GetStatus() int
	GetMessage() string
	GetErr() any
	Traslate(lang string, v ...any)
	Append(field string, err string)
	Errors() Error
}

type Err struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
	Err     any    `json:"error"`
	errors  map[string][]string
}

func New(err string) error {
	return errors.New(err)
}

func NewError(status int, message string, err any) Error {
	return &Err{
		Status:  status,
		Message: message,
		Err:     err,
	}
}

func (e *Err) Error() string {
	return fmt.Sprintf("[%v] %v: %v", e.Status, e.Message, e.Err)
}

func (e *Err) GetStatus() int {
	return e.Status
}

func (e *Err) GetMessage() string {
	return e.Message
}

func (e *Err) GetErr() any {
	return e.Err
}

func (e *Err) Traslate(l string, v ...any) {
	e.Message = lang.TT(l, e.Message, v...)
}

func (e *Err) Append(field string, err string) {
	if e.errors == nil {
		e.errors = make(map[string][]string)
	}
	e.errors[field] = append(e.errors[field], err)
}

func (e *Err) Errors() Error {
	if e.errors == nil {
		return nil
	}
	e.Status = http.StatusUnprocessableEntity
	e.Message = "Error de validación"
	e.Err = e.errors
	return e
}
