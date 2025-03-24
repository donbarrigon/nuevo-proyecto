package app

import (
	"net/http"
)

type HandlerFunc func(ctx *ControllerContext)

func HandlerGet(handler HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler(&ControllerContext{Writer: w, Request: r})
		} else {
			http.NotFound(w, r)
		}
	}
}

func HandlerPost(handler HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler(&ControllerContext{Writer: w, Request: r})
		} else {
			http.NotFound(w, r)
		}
	}
}

func HandlerPut(handler HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			handler(&ControllerContext{Writer: w, Request: r})
		} else {
			http.NotFound(w, r)
		}
	}
}

func HandlerDelete(handler HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			handler(&ControllerContext{Writer: w, Request: r})
		} else {
			http.NotFound(w, r)
		}
	}
}

func HandlerPatch(handler HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			handler(&ControllerContext{Writer: w, Request: r})
		} else {
			http.NotFound(w, r)
		}
	}
}

func HandlerOptions(handler HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			handler(&ControllerContext{Writer: w, Request: r})
		} else {
			http.NotFound(w, r)
		}
	}
}

func HandlerHead(handler HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			handler(&ControllerContext{Writer: w, Request: r})
		} else {
			http.NotFound(w, r)
		}
	}
}

func HandlerConnect(handler HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodConnect {
			handler(&ControllerContext{Writer: w, Request: r})
		} else {
			http.NotFound(w, r)
		}
	}
}

func HandlerTrace(handler HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodTrace {
			handler(&ControllerContext{Writer: w, Request: r})
		} else {
			http.NotFound(w, r)
		}
	}
}
