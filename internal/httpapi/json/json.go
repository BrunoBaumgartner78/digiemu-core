package json

import (
	stdjson "encoding/json"
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

// Error writes a JSON error message: {"error": "message"}
func Error(w http.ResponseWriter, status int, message string) error {
	return Write(w, status, map[string]string{"error": message})
}
