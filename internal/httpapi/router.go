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
			// expecting: /v1/units/{key}/versions
			parts := strings.Split(p, "/")
			// parts: ["", "v1", "units", "{key}", "versions"]
			if len(parts) == 5 && parts[1] == "v1" && parts[2] == "units" && parts[4] == "versions" {
				unitKey := parts[3]
				if unitKey == "" {
					http.NotFound(w, r)
					return
				}
				api.handleCreateVersion(w, r, unitKey)
				return
			}

		case (r.Method == http.MethodPut || r.Method == http.MethodGet) && strings.HasPrefix(p, "/v1/units/") && strings.HasSuffix(p, "/meaning"):
			parts := strings.Split(p, "/")
			if len(parts) == 5 && parts[1] == "v1" && parts[2] == "units" && parts[4] == "meaning" {
				unitKey := parts[3]
				if unitKey == "" {
					http.NotFound(w, r)
					return
				}
				if r.Method == http.MethodPut {
					api.handleSetMeaning(w, r, unitKey)
					return
				}
				api.handleGetMeaning(w, r, unitKey)
				return
			}
		}
		http.NotFound(w, r)
	})
}
