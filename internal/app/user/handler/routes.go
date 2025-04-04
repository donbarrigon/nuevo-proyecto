package handler

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/core"
)

func Routes(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		if r.URL.Path == "/user/login" {
			logout(core.NewContext(w, r))
			return
		}
		show(core.NewContext(w, r))
		return
	}
	if r.Method == http.MethodPost {
		if r.URL.Path == "/user/loguot" {
			logout(core.NewContext(w, r))
			return
		}
		store(core.NewContext(w, r))
		return
	}
	if r.Method == http.MethodPut {
		update(core.NewContext(w, r))
		return
	}
	if r.Method == http.MethodDelete {
		delete(core.NewContext(w, r))
		return
	}
	http.NotFound(w, r)
}
