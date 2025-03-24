package app

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/model"
)

type ControllerContext struct {
	Writer  http.ResponseWriter
	Request *http.Request
	User    *model.User
	Token   *model.Token
}

func NewHandlerContext(w http.ResponseWriter, r *http.Request) *ControllerContext {
	return &ControllerContext{
		Writer:  w,
		Request: r,
	}
}

func (h *ControllerContext) Lang() string {
	return h.Request.Header.Get("Accept-Language")
}
