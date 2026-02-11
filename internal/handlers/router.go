package handlers

import (
	"net/http"
	"strings"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/v3/application/", h.route)

	// OAuth2 public endpoints (no auth required)
	mux.HandleFunc("/v3/public/oauth/token", h.OAuthToken)

	// Health check / ping (no auth needed)
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Admin endpoint to reset data
	mux.HandleFunc("/admin/reset", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "POST only")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "reset"})
	})
}

func (h *Handler) route(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// OpenAPI ping
	if path == "/v3/application/openapi-ping" {
		h.OpenAPIPing(w, r)
		return
	}

	// Scopes check
	if path == "/v3/application/scopes" && r.Method == http.MethodPost {
		h.CheckScopes(w, r)
		return
	}

	// Shipping carriers (public)
	if path == "/v3/application/shipping-carriers" {
		h.GetShippingCarriers(w, r)
		return
	}

	// Buyer Taxonomy
	if strings.HasPrefix(path, "/v3/application/buyer-taxonomy/nodes") {
		if strings.Contains(path, "/properties") {
			h.GetBuyerTaxonomyProperties(w, r)
			return
		}
		h.GetBuyerTaxonomyNodes(w, r)
		return
	}

	// Seller Taxonomy
	if strings.HasPrefix(path, "/v3/application/seller-taxonomy/nodes") {
		if strings.Contains(path, "/properties") {
			h.GetSellerTaxonomyProperties(w, r)
			return
		}
		h.GetSellerTaxonomyNodes(w, r)
		return
	}

	// Batch listings
	if path == "/v3/application/listings/batch" {
		h.GetListingsBatch(w, r)
		return
	}

	// Active listings (global search) - must be before individual listing routes
	if path == "/v3/application/listings/active" {
		h.FindAllActiveListings(w, r)
		return
	}

	// Listing by ID routes (not under shops)
	if strings.HasPrefix(path, "/v3/application/listings/") {
		h.routeListings(w, r)
		return
	}

	// User routes
	if path == "/v3/application/users/me" {
		h.GetMe(w, r)
		return
	}
	if strings.HasPrefix(path, "/v3/application/users/") {
		h.routeUsers(w, r)
		return
	}

	// Shop routes
	if strings.HasPrefix(path, "/v3/application/shops") {
		h.routeShops(w, r)
		return
	}

	writeError(w, http.StatusNotFound, "Endpoint not found")
}

func (h *Handler) routeListings(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	listingPath := strings.TrimPrefix(path, "/v3/application/listings/")
	parts := strings.Split(listingPath, "/")

	if len(parts) >= 2 {
		switch parts[1] {
		case "inventory":
			h.GetListingInventory(w, r)
			return
		case "reviews":
			h.GetListingReviews(w, r)
			return
		case "personalization":
			h.GetListingPersonalization(w, r)
			return
		case "videos":
			if len(parts) == 2 {
				h.GetListingVideos(w, r)
			} else {
				h.GetListingVideo(w, r)
			}
			return
		case "images":
			h.GetListingImages(w, r)
			return
		case "products":
			// /listings/{id}/products/{pid}/offerings/{oid}
			if len(parts) >= 4 && parts[3] == "offerings" {
				h.GetListingOffering(w, r)
				return
			}
		}
	}

	// Single listing by ID
	if len(parts) == 1 {
		switch r.Method {
		case http.MethodGet:
			h.GetListing(w, r)
		case http.MethodDelete:
			h.DeleteListing(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
		return
	}

	writeError(w, http.StatusNotFound, "Endpoint not found")
}

func (h *Handler) routeUsers(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	userPath := strings.TrimPrefix(path, "/v3/application/users/")
	parts := strings.Split(userPath, "/")

	if len(parts) >= 2 {
		switch parts[1] {
		case "addresses":
			if len(parts) == 3 && r.Method == http.MethodDelete {
				h.DeleteUserAddress(w, r)
				return
			}
			h.GetUserAddresses(w, r)
			return
		case "shops":
			h.GetUserShops(w, r)
			return
		}
	}

	h.GetUser(w, r)
}

func (h *Handler) routeShops(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// GET /v3/application/shops?shop_name=xxx (shop search)
	if path == "/v3/application/shops" || path == "/v3/application/shops/" {
		if r.Method == http.MethodGet {
			h.FindShops(w, r)
			return
		}
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Routes under /v3/application/shops/{shop_id}/...
	shopPath := strings.TrimPrefix(path, "/v3/application/shops/")
	parts := strings.Split(shopPath, "/")

	if len(parts) < 1 {
		writeError(w, http.StatusBadRequest, "Invalid path")
		return
	}

	// /v3/application/shops/{shop_id}
	if len(parts) == 1 {
		switch r.Method {
		case http.MethodGet:
			h.GetShop(w, r)
		case http.MethodPut:
			h.UpdateShop(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
		return
	}

	resource := parts[1]

	switch resource {
	case "listings":
		h.routeShopListings(w, r, parts)
	case "receipts":
		h.routeShopReceipts(w, r, parts)
	case "transactions":
		h.routeShopTransactions(w, r, parts)
	case "payments":
		h.GetShopPayments(w, r)
	case "reviews":
		h.GetShopReviews(w, r)
	case "sections":
		h.routeShopSections(w, r, parts)
	case "shop-sections":
		// /shops/{id}/shop-sections/listings
		if len(parts) >= 3 && parts[2] == "listings" {
			h.GetShopSectionListings(w, r)
			return
		}
		writeError(w, http.StatusNotFound, "Endpoint not found")
	case "return-policies":
		h.routeReturnPolicies(w, r, parts)
	case "shipping-profiles":
		h.routeShippingProfiles(w, r, parts)
	case "payment-account":
		if len(parts) >= 3 && parts[2] == "ledger-entries" {
			h.GetLedgerEntries(w, r)
			return
		}
		writeError(w, http.StatusNotFound, "Endpoint not found")
	case "production-partners":
		h.GetProductionPartners(w, r)
	case "holiday-preferences":
		switch r.Method {
		case http.MethodGet:
			h.GetHolidayPreferences(w, r)
		case http.MethodPut:
			h.UpdateHolidayPreferences(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	case "readiness-state-definitions":
		h.GetReadinessStateDefinitions(w, r)
	default:
		writeError(w, http.StatusNotFound, "Endpoint not found")
	}
}

func (h *Handler) routeShopListings(w http.ResponseWriter, r *http.Request, parts []string) {
	// /shops/{id}/listings
	if len(parts) == 2 {
		switch r.Method {
		case http.MethodGet:
			h.GetShopListings(w, r)
		case http.MethodPost:
			h.CreateListing(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
		return
	}

	// /shops/{id}/listings/active
	if len(parts) == 3 && parts[2] == "active" {
		h.GetShopActiveListings(w, r)
		return
	}

	// /shops/{id}/listings/featured
	if len(parts) == 3 && parts[2] == "featured" {
		h.GetFeaturedListings(w, r)
		return
	}

	// /shops/{id}/listings/{listing_id}
	if len(parts) == 3 {
		switch r.Method {
		case http.MethodGet:
			h.GetListing(w, r)
		case http.MethodPut, http.MethodPatch:
			h.UpdateListing(w, r)
		case http.MethodDelete:
			h.DeleteListing(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
		return
	}

	// /shops/{id}/listings/{listing_id}/images
	if len(parts) >= 4 && parts[3] == "images" {
		if len(parts) == 4 {
			switch r.Method {
			case http.MethodGet:
				h.GetListingImages(w, r)
			case http.MethodPost:
				h.UploadListingImage(w, r)
			default:
				writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
			return
		}
		// /shops/{id}/listings/{listing_id}/images/{image_id}
		if len(parts) == 5 && r.Method == http.MethodDelete {
			h.DeleteListingImage(w, r)
			return
		}
	}

	// /shops/{id}/listings/{listing_id}/files
	if len(parts) >= 4 && parts[3] == "files" {
		if len(parts) == 4 {
			switch r.Method {
			case http.MethodGet:
				h.GetListingFiles(w, r)
			case http.MethodPost:
				h.UploadListingFile(w, r)
			default:
				writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
			return
		}
		if len(parts) == 5 {
			switch r.Method {
			case http.MethodGet:
				h.GetListingFile(w, r)
			case http.MethodDelete:
				h.DeleteListingFile(w, r)
			default:
				writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
			return
		}
	}

	// /shops/{id}/listings/{listing_id}/inventory
	if len(parts) == 4 && parts[3] == "inventory" {
		h.GetListingInventory(w, r)
		return
	}

	// /shops/{id}/listings/{listing_id}/personalization
	if len(parts) == 4 && parts[3] == "personalization" {
		switch r.Method {
		case http.MethodGet:
			h.GetListingPersonalization(w, r)
		case http.MethodPut:
			h.UpdateListingPersonalization(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
		return
	}

	// /shops/{id}/listings/{listing_id}/videos[/{video_id}]
	if len(parts) >= 4 && parts[3] == "videos" {
		if len(parts) == 4 {
			switch r.Method {
			case http.MethodGet:
				h.GetListingVideos(w, r)
			case http.MethodPost:
				h.UploadListingVideo(w, r)
			default:
				writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
			return
		}
		if len(parts) == 5 {
			switch r.Method {
			case http.MethodGet:
				h.GetListingVideo(w, r)
			case http.MethodDelete:
				h.DeleteListingVideo(w, r)
			default:
				writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
			return
		}
	}

	// /shops/{id}/listings/{listing_id}/translations/{language}
	if len(parts) == 5 && parts[3] == "translations" {
		switch r.Method {
		case http.MethodGet:
			h.GetListingTranslation(w, r)
		case http.MethodPut:
			h.UpdateListingTranslation(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
		return
	}

	// /shops/{id}/listings/{listing_id}/variation-images
	if len(parts) == 4 && parts[3] == "variation-images" {
		switch r.Method {
		case http.MethodGet:
			h.GetListingVariationImages(w, r)
		case http.MethodPost:
			h.UpdateListingVariationImages(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
		return
	}

	// /shops/{id}/listings/{listing_id}/properties
	if len(parts) == 4 && parts[3] == "properties" {
		h.GetListingProperties(w, r)
		return
	}

	writeError(w, http.StatusNotFound, "Endpoint not found")
}

func (h *Handler) routeShopReceipts(w http.ResponseWriter, r *http.Request, parts []string) {
	// /shops/{id}/receipts
	if len(parts) == 2 {
		h.GetShopReceipts(w, r)
		return
	}

	// /shops/{id}/receipts/{receipt_id}/payments
	if len(parts) == 4 && parts[3] == "payments" {
		h.GetReceiptPayments(w, r)
		return
	}

	// /shops/{id}/receipts/{receipt_id}/transactions
	if len(parts) == 4 && parts[3] == "transactions" {
		h.GetReceiptTransactions(w, r)
		return
	}

	// /shops/{id}/receipts/{receipt_id}/tracking
	if len(parts) == 4 && parts[3] == "tracking" && r.Method == http.MethodPost {
		h.CreateReceiptTracking(w, r)
		return
	}

	// /shops/{id}/receipts/{receipt_id}
	if len(parts) == 3 {
		switch r.Method {
		case http.MethodGet:
			h.GetReceipt(w, r)
		case http.MethodPut:
			h.UpdateReceipt(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
		return
	}

	writeError(w, http.StatusNotFound, "Endpoint not found")
}

func (h *Handler) routeShopTransactions(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 2 {
		h.GetShopTransactions(w, r)
		return
	}
	if len(parts) == 3 {
		h.GetTransaction(w, r)
		return
	}
	writeError(w, http.StatusNotFound, "Endpoint not found")
}

func (h *Handler) routeShopSections(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 2 {
		switch r.Method {
		case http.MethodGet:
			h.GetShopSections(w, r)
		case http.MethodPost:
			h.CreateShopSection(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
		return
	}
	if len(parts) == 3 {
		h.GetShopSection(w, r)
		return
	}
	writeError(w, http.StatusNotFound, "Endpoint not found")
}

func (h *Handler) routeReturnPolicies(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 2 {
		switch r.Method {
		case http.MethodGet:
			h.GetReturnPolicies(w, r)
		case http.MethodPost:
			h.CreateReturnPolicy(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
		return
	}
	if len(parts) == 3 {
		h.GetReturnPolicy(w, r)
		return
	}
	writeError(w, http.StatusNotFound, "Endpoint not found")
}

func (h *Handler) routeShippingProfiles(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 2 {
		switch r.Method {
		case http.MethodGet:
			h.GetShopShippingProfiles(w, r)
		case http.MethodPost:
			h.CreateShippingProfile(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
		return
	}
	if len(parts) == 3 {
		switch r.Method {
		case http.MethodGet:
			h.GetShippingProfile(w, r)
		case http.MethodDelete:
			h.DeleteShippingProfile(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
		return
	}
	writeError(w, http.StatusNotFound, "Endpoint not found")
}
