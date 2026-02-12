package handlers

import (
	"net/http"

	"github.com/vlah-software-house/etsy-mock-api/internal/models"
)

// GET /v3/application/shops/{shop_id}/reviews
func (h *Handler) GetShopReviews(w http.ResponseWriter, r *http.Request) {
	shopID, ok := extractPathID(r.URL.Path, "shops")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid shop_id")
		return
	}
	limit := queryInt(r, "limit", 25)
	offset := queryInt(r, "offset", 0)

	reviews, total := h.Store.GetShopReviews(shopID, limit, offset)
	if reviews == nil {
		reviews = []models.ListingReview{}
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   total,
		Results: reviews,
	})
}

// GET /v3/application/listings/{listing_id}/reviews
func (h *Handler) GetListingReviews(w http.ResponseWriter, r *http.Request) {
	listingID, ok := extractPathID(r.URL.Path, "listings")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid listing_id")
		return
	}
	limit := queryInt(r, "limit", 25)
	offset := queryInt(r, "offset", 0)

	reviews, total := h.Store.GetListingReviews(listingID, limit, offset)
	if reviews == nil {
		reviews = []models.ListingReview{}
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   total,
		Results: reviews,
	})
}
