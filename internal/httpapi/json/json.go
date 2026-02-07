package json

import (
	stdjson "encoding/json"
	"fmt"
	"net/http"
)

// Read decodes JSON from the request body into v.
func Read(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	dec := stdjson.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// Write sends a JSON response with the provided status code and payload.
func Write(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if payload == nil {
		return nil
	}
	enc := stdjson.NewEncoder(w)
	return enc.Encode(payload)
}

// Error payload shape: { "error": { "code": "<str>", "message": "<str>", "details": <optional> } }
type ErrorBody struct {
	Error struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Details interface{} `json:"details,omitempty"`
	} `json:"error"`
}

// ErrorCode writes a structured error with provided code and message.
func ErrorCode(w http.ResponseWriter, status int, code, message string, details interface{}) error {
	var b ErrorBody
	b.Error.Code = code
	b.Error.Message = message
	b.Error.Details = details
	return Write(w, status, b)
}

// Errorf formats and writes an error.
func Errorf(w http.ResponseWriter, status int, code, format string, a ...interface{}) error {
	return ErrorCode(w, status, code, fmt.Sprintf(format, a...), nil)
}
