package httpapi

import (
	"net/http"
	"strings"
)

// simple router using stdlib. expects paths:
// POST /v1/units
// POST /v1/units/{unitId}/versions
// GET  /healthz
func NewRouter(api API) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		switch {
		case r.Method == http.MethodGet && p == "/healthz":
			api.handleHealth(w, r)
			return
		case r.Method == http.MethodPost && p == "/v1/units":
			api.handleCreateUnit(w, r)
			return
		case r.Method == http.MethodPost && strings.HasPrefix(p, "/v1/units/") && strings.HasSuffix(p, "/versions"):
			// extract unit key between
			parts := strings.Split(p, "/")
			if len(parts) >= 5 {
				unitKey := parts[3]
				api.handleCreateVersion(w, r, unitKey)
				return
			}
		}
		http.NotFound(w, r)
	})
}
