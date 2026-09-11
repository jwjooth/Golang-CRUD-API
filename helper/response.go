package helper

import (
	"encoding/json"
	"net/http"
)

// SuccessEnvelope is the standard success body.
type SuccessEnvelope struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

// ErrorEnvelope is the standard error body.
type ErrorEnvelope struct {
	Error string `json:"error"`
}

// WriteJSON writes v as JSON with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteSuccess writes a 2xx envelope: {"message":..., "data":..., "meta"?}.
func WriteSuccess(w http.ResponseWriter, status int, message string, data any) {
	WriteJSON(w, status, SuccessEnvelope{Message: message, Data: data})
}

// WriteSuccessWithMeta writes a 2xx envelope including pagination meta.
func WriteSuccessWithMeta(w http.ResponseWriter, status int, message string, data any, meta any) {
	WriteJSON(w, status, SuccessEnvelope{Message: message, Data: data, Meta: meta})
}

// WriteError writes a 4xx/5xx envelope: {"error": message}.
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, ErrorEnvelope{Error: message})
}

// WriteAppError maps an AppError (or generic error) to a JSON error body.
func WriteAppError(w http.ResponseWriter, err error) {
	WriteError(w, CodeOf(err), MessageOf(err))
}
