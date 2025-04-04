package com

import "fmt"

type Error interface {
	error
	GetStatus() int
	GetMessage() string
	GetErr() any
	TraslateError(lang string) string
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
	return e.TraslateError(LANG)
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

func (e *Err) TraslateError(lang string) string {
	return fmt.Sprintf(`{"message":"%s","error":"%v"}`, TT(lang, e.Message), e.Err)
}
