package httpapi

// This file defines the common HTTP/JSON boundary. Keeping JSON decoding and
// error formatting here means every handler returns the same response shape.

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// errorBody is serialized as {"error":{"code":"...","message":"..."}}.
type errorBody struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// writeJSON sets the response metadata before encoding a Go value as JSON.
// Headers must be set before WriteHeader because that sends them to the client.
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if value != nil {
		_ = json.NewEncoder(w).Encode(value)
	}
}

// writeError gives machines a stable code and people a readable message.
func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorBody{Error: apiError{Code: code, Message: message}})
}

// decodeJSON converts the request body into a handler input struct. Returning a
// bool lets handlers use `if !decodeJSON(...) { return }` without repeating
// the same HTTP error code. Unknown fields and a second JSON value are rejected
// so frontend/backend contract mistakes fail visibly instead of being ignored.
func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "Проверьте формат отправленных данных.")
		return false
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "invalid_json", "Запрос должен содержать один JSON-объект.")
		return false
	}
	return true
}
