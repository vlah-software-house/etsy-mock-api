package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gabriel/etsy-mock/internal/middleware"
	"github.com/gabriel/etsy-mock/internal/models"
	"github.com/gabriel/etsy-mock/internal/store"
)

type Handler struct {
	Store      *store.Store
	TokenStore *middleware.TokenStore
}

func New(s *store.Store, ts *middleware.TokenStore) *Handler {
	return &Handler{Store: s, TokenStore: ts}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, models.ErrorResponse{Error: msg})
}

func pathParam(r *http.Request, prefix string) string {
	path := r.URL.Path
	idx := strings.Index(path, prefix)
	if idx == -1 {
		return ""
	}
	rest := path[idx+len(prefix):]
	if slashIdx := strings.Index(rest, "/"); slashIdx != -1 {
		return rest[:slashIdx]
	}
	return rest
}

func parseID(s string) (int64, bool) {
	id, err := strconv.ParseInt(s, 10, 64)
	return id, err == nil
}

func queryInt(r *http.Request, key string, def int) int {
	s := r.URL.Query().Get(key)
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

func queryString(r *http.Request, key string, def string) string {
	s := r.URL.Query().Get(key)
	if s == "" {
		return def
	}
	return s
}

func decodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// requireScope checks if the request has the required OAuth2 scope.
// Returns true if scope is present (continue handling), false if error was written.
func requireScope(w http.ResponseWriter, r *http.Request, scope string) bool {
	if !middleware.HasScope(r, scope) {
		writeJSON(w, http.StatusForbidden, models.ErrorResponse{
			Error: "This endpoint requires OAuth2 scope: " + scope + ". Provide a Bearer token with this scope.",
		})
		return false
	}
	return true
}

// extractSegment extracts a path segment at the given index from a split path.
// Example: /v3/application/shops/123/listings -> segments = ["", "v3", "application", "shops", "123", "listings"]
func extractPathID(path string, after string) (int64, bool) {
	parts := strings.Split(path, "/")
	for i, p := range parts {
		if p == after && i+1 < len(parts) {
			return parseID(parts[i+1])
		}
	}
	return 0, false
}
