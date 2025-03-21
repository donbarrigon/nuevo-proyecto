package app

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/model"
)

type HandlerContext struct {
	Writer  http.ResponseWriter
	Request *http.Request
	User    *model.User
	Token   *model.Token
}

func NewHandlerContext(w http.ResponseWriter, r *http.Request) *HandlerContext {
	return &HandlerContext{
		Writer:  w,
		Request: r,
	}
}

func (h *HandlerContext) Lang() string {
	return h.Request.Header.Get("Accept-Language")
}
