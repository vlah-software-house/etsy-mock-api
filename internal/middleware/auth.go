package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type contextKey string

const (
	ContextKeystring contextKey = "keystring"
	ContextScopes    contextKey = "scopes"
	ContextUserID    contextKey = "user_id"
)

// TokenEntry represents a stored mock OAuth token.
type TokenEntry struct {
	AccessToken  string
	RefreshToken string
	UserID       int64
	Scopes       []string
	ExpiresAt    time.Time
}

// TokenStore holds mock OAuth tokens. Thread-safe.
type TokenStore struct {
	mu     sync.RWMutex
	tokens map[string]*TokenEntry // keyed by access_token
}

func NewTokenStore() *TokenStore {
	ts := &TokenStore{tokens: make(map[string]*TokenEntry)}
	// Pre-seed tokens for testing convenience
	ts.tokens["test-token-alice"] = &TokenEntry{
		AccessToken: "test-token-alice", RefreshToken: "refresh-alice",
		UserID: 1001, Scopes: AllScopes(), ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	ts.tokens["test-token-bob"] = &TokenEntry{
		AccessToken: "test-token-bob", RefreshToken: "refresh-bob",
		UserID: 1002, Scopes: AllScopes(), ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	return ts
}

func (ts *TokenStore) Get(accessToken string) (*TokenEntry, bool) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	t, ok := ts.tokens[accessToken]
	if ok && time.Now().After(t.ExpiresAt) {
		return nil, false
	}
	return t, ok
}

func (ts *TokenStore) GetByRefresh(refreshToken string) (*TokenEntry, bool) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	for _, t := range ts.tokens {
		if t.RefreshToken == refreshToken {
			return t, true
		}
	}
	return nil, false
}

func (ts *TokenStore) Store(entry *TokenEntry) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.tokens[entry.AccessToken] = entry
}

func AllScopes() []string {
	return []string{
		"address_r", "address_w", "billing_r", "cart_r", "cart_w",
		"email_r", "favorites_r", "favorites_w", "feedback_r",
		"listings_d", "listings_r", "listings_w",
		"profile_r", "profile_w", "recommend_r", "recommend_w",
		"shops_r", "shops_w", "transactions_r", "transactions_w",
	}
}

// MockAuth validates API key (keystring:shared_secret format) and optional OAuth bearer token.
func MockAuth(tokenStore *TokenStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for public OAuth endpoints and ping
			if strings.HasPrefix(r.URL.Path, "/v3/public/") ||
				r.URL.Path == "/ping" ||
				r.URL.Path == "/v3/application/openapi-ping" {
				next.ServeHTTP(w, r)
				return
			}

			apiKey := r.Header.Get("x-api-key")
			authHeader := r.Header.Get("Authorization")

			// Validate x-api-key format: keystring:shared_secret
			if apiKey == "" {
				writeAuthError(w, http.StatusUnauthorized, "Missing x-api-key header. Format: keystring:shared_secret")
				return
			}

			parts := strings.SplitN(apiKey, ":", 2)
			if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
				writeAuthError(w, http.StatusUnauthorized, "Invalid x-api-key format. Expected keystring:shared_secret")
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeystring, parts[0])

			// Check for OAuth bearer token
			if strings.HasPrefix(authHeader, "Bearer ") {
				token := authHeader[7:]
				entry, ok := tokenStore.Get(token)
				if !ok {
					writeAuthError(w, http.StatusUnauthorized, "Invalid or expired OAuth access token")
					return
				}
				ctx = context.WithValue(ctx, ContextScopes, entry.Scopes)
				ctx = context.WithValue(ctx, ContextUserID, entry.UserID)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// HasScope checks if the current request has a given scope.
func HasScope(r *http.Request, scope string) bool {
	scopes, ok := r.Context().Value(ContextScopes).([]string)
	if !ok {
		return false
	}
	for _, s := range scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// GetUserID extracts the authenticated user ID from the context.
func GetUserID(r *http.Request) (int64, bool) {
	uid, ok := r.Context().Value(ContextUserID).(int64)
	return uid, ok
}

func writeAuthError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"error":"%s"}`, msg)
}

// RateLimitHeaders adds Etsy-style rate limit headers to every response.
func RateLimitHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-limit-per-second", "10")
		w.Header().Set("x-remaining-this-second", "9")
		w.Header().Set("x-limit-per-day", "10000")
		w.Header().Set("x-remaining-today", "9999")
		next.ServeHTTP(w, r)
	})
}

// RequestLogger logs each request with method, path, and duration.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// CORS adds permissive CORS headers for local development.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-api-key")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// JSONContent sets content-type to application/json.
func JSONContent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
