package handlers

import (
	"net/http"

	"github.com/vlah-software-house/etsy-mock-api/internal/models"
)

// GET /v3/application/shops/{shop_id}/shipping-profiles
func (h *Handler) GetShopShippingProfiles(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "shops_r") {
		return
	}
	shopID, ok := extractPathID(r.URL.Path, "shops")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid shop_id")
		return
	}
	profiles := h.Store.GetShopShippingProfiles(shopID)
	if profiles == nil {
		profiles = []models.ShopShippingProfile{}
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   len(profiles),
		Results: profiles,
	})
}

// GET /v3/application/shops/{shop_id}/shipping-profiles/{shipping_profile_id}
func (h *Handler) GetShippingProfile(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "shops_r") {
		return
	}
	profileID, ok := extractPathID(r.URL.Path, "shipping-profiles")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid shipping_profile_id")
		return
	}
	profile, found := h.Store.GetShippingProfile(profileID)
	if !found {
		writeError(w, http.StatusNotFound, "Shipping profile not found")
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

// POST /v3/application/shops/{shop_id}/shipping-profiles
func (h *Handler) CreateShippingProfile(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "shops_w") {
		return
	}
	shopID, ok := extractPathID(r.URL.Path, "shops")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid shop_id")
		return
	}
	shop, found := h.Store.GetShop(shopID)
	if !found {
		writeError(w, http.StatusNotFound, "Shop not found")
		return
	}

	var body struct {
		Title            string `json:"title"`
		OriginCountryISO string `json:"origin_country_iso"`
		ProfileType      string `json:"profile_type"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if body.Title == "" || body.OriginCountryISO == "" {
		writeError(w, http.StatusBadRequest, "title and origin_country_iso are required")
		return
	}
	if body.ProfileType == "" {
		body.ProfileType = "manual"
	}

	profile := h.Store.CreateShippingProfile(shop.UserID, body.Title, body.OriginCountryISO, body.ProfileType)
	writeJSON(w, http.StatusCreated, profile)
}

// DELETE /v3/application/shops/{shop_id}/shipping-profiles/{shipping_profile_id}
func (h *Handler) DeleteShippingProfile(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "shops_w") {
		return
	}
	profileID, ok := extractPathID(r.URL.Path, "shipping-profiles")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid shipping_profile_id")
		return
	}
	if !h.Store.DeleteShippingProfile(profileID) {
		writeError(w, http.StatusNotFound, "Shipping profile not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
