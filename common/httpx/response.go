package httpx

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func OK(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusOK, Response{Code: 0, Message: "ok", Data: data})
}

func Fail(w http.ResponseWriter, status, code int, msg string) {
	writeJSON(w, status, Response{Code: code, Message: msg})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
