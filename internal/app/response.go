package app

import (
	"encoding/json"
	"net/http"
)

type ErrorJSON struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
	Error   any    `json:"error"`
}

func (e *ErrorJSON) WriteResponse(w http.ResponseWriter) {
	ResponseJSON(w, e, e.Status)
}

func NewErrorJSON(err any) *ErrorJSON {
	return &ErrorJSON{
		Status:  http.StatusInternalServerError,
		Message: "Error",
		Error:   err,
	}
}

type OkJSON struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewOkJSON(data any) *OkJSON {
	return &OkJSON{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    data,
	}
}

func (o *OkJSON) WriteResponse(w http.ResponseWriter) {
	ResponseJSON(w, o, o.Status)
}

func ResponseJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"Error","error":"failed to encode response"}`))
	}
}

func ResponseOkJSON(w http.ResponseWriter, data any, status int, message string) {
	okJSON := &OkJSON{
		Message: message,
		Status:  status,
		Data:    data,
	}
	okJSON.WriteResponse(w)
}

func ResponseErrorJSON(w http.ResponseWriter, err any, status int, message string) {
	errorJSON := &ErrorJSON{
		Message: message,
		Status:  status,
		Error:   err,
	}
	errorJSON.WriteResponse(w)
}
