package core

import (
	"net/http"
)

type HandlerFunc func(ctx *Context)

func HandlerGet(handler HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler(&Context{Writer: w, Request: r})
		} else {
			http.NotFound(w, r)
		}
	}
}

func HandlerPost(handler HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler(&Context{Writer: w, Request: r})
		} else {
			http.NotFound(w, r)
		}
	}
}

func HandlerPut(handler HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			handler(&Context{Writer: w, Request: r})
		} else {
			http.NotFound(w, r)
		}
	}
}

func HandlerDelete(handler HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			handler(&Context{Writer: w, Request: r})
		} else {
			http.NotFound(w, r)
		}
	}
}

func HandlerPatch(handler HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			handler(&Context{Writer: w, Request: r})
		} else {
			http.NotFound(w, r)
		}
	}
}

func HandlerOptions(handler HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			handler(&Context{Writer: w, Request: r})
		} else {
			http.NotFound(w, r)
		}
	}
}

func HandlerHead(handler HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			handler(&Context{Writer: w, Request: r})
		} else {
			http.NotFound(w, r)
		}
	}
}

func HandlerConnect(handler HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodConnect {
			handler(&Context{Writer: w, Request: r})
		} else {
			http.NotFound(w, r)
		}
	}
}

func HandlerTrace(handler HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodTrace {
			handler(&Context{Writer: w, Request: r})
		} else {
			http.NotFound(w, r)
		}
	}
}
