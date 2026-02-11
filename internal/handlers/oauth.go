package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gabriel/etsy-mock/internal/middleware"
)

// GET /v3/application/openapi-ping
func (h *Handler) OpenAPIPing(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"application_id": "mock-app-123",
	})
}

// POST /v3/application/scopes
func (h *Handler) CheckScopes(w http.ResponseWriter, r *http.Request) {
	scopes, ok := r.Context().Value(middleware.ContextScopes).([]string)
	if !ok {
		scopes = []string{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"scopes": scopes,
	})
}

// GET /v3/application/users/me
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusForbidden, "This endpoint requires OAuth2. Provide a Bearer token.")
		return
	}
	user, found := h.Store.GetUser(userID)
	if !found {
		writeError(w, http.StatusNotFound, "User not found")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// POST /v3/public/oauth/token — Token exchange and refresh
func (h *Handler) OAuthToken(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	grantType := r.FormValue("grant_type")
	clientID := r.FormValue("client_id")

	switch grantType {
	case "authorization_code":
		h.handleAuthCodeExchange(w, r, clientID)
	case "refresh_token":
		h.handleRefreshToken(w, r, clientID)
	default:
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Unsupported grant_type: %s", grantType))
	}
}

func (h *Handler) handleAuthCodeExchange(w http.ResponseWriter, r *http.Request, clientID string) {
	code := r.FormValue("code")
	codeVerifier := r.FormValue("code_verifier")

	if clientID == "" || code == "" || codeVerifier == "" {
		writeError(w, http.StatusBadRequest, "client_id, code, and code_verifier are required")
		return
	}

	// In mock mode, accept any code and verifier. Extract user from the code pattern.
	// Convention: code format is "mock_code_{user_id}_{scopes}" or just any string.
	// Default to user 1001 if we can't parse.
	userID := int64(1001)
	scopes := middleware.AllScopes()

	// Try to extract user ID from code like "mock_code_1001"
	if strings.HasPrefix(code, "mock_code_") {
		rest := strings.TrimPrefix(code, "mock_code_")
		parts := strings.SplitN(rest, "_", 2)
		if id, ok := parseID(parts[0]); ok {
			userID = id
		}
	}

	accessToken := generateToken(userID)
	refreshToken := "refresh_" + generateToken(userID)

	h.TokenStore.Store(&middleware.TokenEntry{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       userID,
		Scopes:       scopes,
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	})

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"access_token":  accessToken,
		"token_type":    "Bearer",
		"expires_in":    3600,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) handleRefreshToken(w http.ResponseWriter, r *http.Request, clientID string) {
	refreshToken := r.FormValue("refresh_token")
	if clientID == "" || refreshToken == "" {
		writeError(w, http.StatusBadRequest, "client_id and refresh_token are required")
		return
	}

	entry, ok := h.TokenStore.GetByRefresh(refreshToken)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	newAccessToken := generateToken(entry.UserID)
	newRefreshToken := "refresh_" + generateToken(entry.UserID)

	h.TokenStore.Store(&middleware.TokenEntry{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		UserID:       entry.UserID,
		Scopes:       entry.Scopes,
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	})

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"access_token":  newAccessToken,
		"token_type":    "Bearer",
		"expires_in":    3600,
		"refresh_token": newRefreshToken,
	})
}

func generateToken(userID int64) string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%d.%s", userID, hex.EncodeToString(b))
}
