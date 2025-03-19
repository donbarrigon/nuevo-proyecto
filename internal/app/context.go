package app

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/model"
)

type HandlerContext struct {
	Request *http.Request
	Writer  http.ResponseWriter
	User    *model.User
	Token   *model.Token
}

func NewHandlerContext(w http.ResponseWriter, r *http.Request) *HandlerContext {
	return &HandlerContext{
		Request: r,
		Writer:  w,
	}
}
