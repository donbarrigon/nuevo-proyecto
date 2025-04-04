package core

import "fmt"

type Error interface {
	error
	GetStatus() int
	GetMessage() string
	GetErr() any
}

type Err struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
	Err     any    `json:"error"`
}

func NewError(status int, message string, err any) Error {
	return &Err{
		Status:  status,
		Message: message,
		Err:     err,
	}
}

func (e *Err) Error() string {
	return fmt.Sprintf("%v: %v", e.Message, e.Err)
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
