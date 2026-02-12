package handlers

import (
	"net/http"

	"github.com/vlah-software-house/etsy-mock-api/internal/models"
)

// GET /v3/application/users/{user_id}
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "email_r") {
		return
	}
	userID, ok := extractPathID(r.URL.Path, "users")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid user_id")
		return
	}
	user, found := h.Store.GetUser(userID)
	if !found {
		writeError(w, http.StatusNotFound, "User not found")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// GET /v3/application/users/{user_id}/addresses
func (h *Handler) GetUserAddresses(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "address_r") {
		return
	}
	userID, ok := extractPathID(r.URL.Path, "users")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid user_id")
		return
	}
	addrs := h.Store.GetUserAddresses(userID)
	if addrs == nil {
		addrs = []models.UserAddress{}
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   len(addrs),
		Results: addrs,
	})
}

// GET /v3/application/users/{user_id}/shops
func (h *Handler) GetUserShops(w http.ResponseWriter, r *http.Request) {
	userID, ok := extractPathID(r.URL.Path, "users")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid user_id")
		return
	}

	var shops []models.Shop
	for _, shop := range h.Store.Shops {
		if shop.UserID == userID {
			shops = append(shops, *shop)
		}
	}
	if shops == nil {
		shops = []models.Shop{}
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   len(shops),
		Results: shops,
	})
}
