package handlers

import (
	"net/http"

	"github.com/vlah-software-house/etsy-mock-api/internal/models"
)

// POST /v3/application/shops/{shop_id}/listings
func (h *Handler) CreateListing(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "listings_w") {
		return
	}
	shopID, ok := extractPathID(r.URL.Path, "shops")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid shop_id")
		return
	}
	if _, found := h.Store.GetShop(shopID); !found {
		writeError(w, http.StatusNotFound, "Shop not found")
		return
	}

	var req models.CreateListingRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.Title == "" || req.Quantity <= 0 || req.Price <= 0 || req.WhoMade == "" || req.WhenMade == "" || req.TaxonomyID == 0 {
		writeError(w, http.StatusBadRequest, "Missing required fields: title, quantity, price, who_made, when_made, taxonomy_id")
		return
	}

	listing := h.Store.CreateListing(shopID, req)
	writeJSON(w, http.StatusCreated, listing)
}

// GET /v3/application/shops/{shop_id}/listings
func (h *Handler) GetShopListings(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "listings_r") {
		return
	}
	shopID, ok := extractPathID(r.URL.Path, "shops")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid shop_id")
		return
	}
	state := queryString(r, "state", "")
	limit := queryInt(r, "limit", 25)
	offset := queryInt(r, "offset", 0)

	listings, total := h.Store.GetShopListings(shopID, state, limit, offset)
	if listings == nil {
		listings = []models.ShopListing{}
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   total,
		Results: listings,
	})
}

// GET /v3/application/shops/{shop_id}/listings/active
func (h *Handler) GetShopActiveListings(w http.ResponseWriter, r *http.Request) {
	shopID, ok := extractPathID(r.URL.Path, "shops")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid shop_id")
		return
	}
	limit := queryInt(r, "limit", 25)
	offset := queryInt(r, "offset", 0)

	listings, total := h.Store.GetShopListings(shopID, "active", limit, offset)
	if listings == nil {
		listings = []models.ShopListing{}
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   total,
		Results: listings,
	})
}

// GET /v3/application/listings/{listing_id}
func (h *Handler) GetListing(w http.ResponseWriter, r *http.Request) {
	listingID, ok := extractPathID(r.URL.Path, "listings")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid listing_id")
		return
	}
	listing, found := h.Store.GetListing(listingID)
	if !found {
		writeError(w, http.StatusNotFound, "Listing not found")
		return
	}

	// Attach images if available
	listing.Images = h.Store.GetListingImages(listingID)

	writeJSON(w, http.StatusOK, listing)
}

// PUT/PATCH /v3/application/shops/{shop_id}/listings/{listing_id}
func (h *Handler) UpdateListing(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "listings_w") {
		return
	}
	listingID, ok := extractPathID(r.URL.Path, "listings")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid listing_id")
		return
	}

	var req models.UpdateListingRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	listing, found := h.Store.UpdateListing(listingID, req)
	if !found {
		writeError(w, http.StatusNotFound, "Listing not found")
		return
	}
	writeJSON(w, http.StatusOK, listing)
}

// DELETE /v3/application/listings/{listing_id}
func (h *Handler) DeleteListing(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "listings_d") {
		return
	}
	listingID, ok := extractPathID(r.URL.Path, "listings")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid listing_id")
		return
	}
	if !h.Store.DeleteListing(listingID) {
		writeError(w, http.StatusNotFound, "Listing not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GET /v3/application/listings/active
func (h *Handler) FindAllActiveListings(w http.ResponseWriter, r *http.Request) {
	keyword := queryString(r, "keywords", "")
	limit := queryInt(r, "limit", 25)
	offset := queryInt(r, "offset", 0)
	sortOn := queryString(r, "sort_on", "created")
	sortOrder := queryString(r, "sort_order", "desc")

	var taxonomyID *int
	if tid := r.URL.Query().Get("taxonomy_id"); tid != "" {
		if v, ok := parseID(tid); ok {
			intV := int(v)
			taxonomyID = &intV
		}
	}

	listings, total := h.Store.GetActiveListings(keyword, taxonomyID, limit, offset, sortOn, sortOrder)
	if listings == nil {
		listings = []models.ShopListing{}
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   total,
		Results: listings,
	})
}

// GET /v3/application/shops/{shop_id}/listings/{listing_id}/images
func (h *Handler) GetListingImages(w http.ResponseWriter, r *http.Request) {
	listingID, ok := extractPathID(r.URL.Path, "listings")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid listing_id")
		return
	}
	images := h.Store.GetListingImages(listingID)
	if images == nil {
		images = []models.ListingImage{}
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   len(images),
		Results: images,
	})
}

// POST /v3/application/shops/{shop_id}/listings/{listing_id}/images
func (h *Handler) UploadListingImage(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "listings_w") {
		return
	}
	listingID, ok := extractPathID(r.URL.Path, "listings")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid listing_id")
		return
	}
	if _, found := h.Store.GetListing(listingID); !found {
		writeError(w, http.StatusNotFound, "Listing not found")
		return
	}

	altText := r.FormValue("alt_text")
	rank := queryInt(r, "rank", 1)

	img := h.Store.AddListingImage(listingID, altText, rank)
	writeJSON(w, http.StatusCreated, img)
}

// DELETE /v3/application/shops/{shop_id}/listings/{listing_id}/images/{listing_image_id}
func (h *Handler) DeleteListingImage(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "listings_w") {
		return
	}
	listingID, ok := extractPathID(r.URL.Path, "listings")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid listing_id")
		return
	}
	imageID, ok := extractPathID(r.URL.Path, "images")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid listing_image_id")
		return
	}
	if !h.Store.DeleteListingImage(listingID, imageID) {
		writeError(w, http.StatusNotFound, "Image not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GET /v3/application/shops/{shop_id}/listings/{listing_id}/files
func (h *Handler) GetListingFiles(w http.ResponseWriter, r *http.Request) {
	listingID, ok := extractPathID(r.URL.Path, "listings")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid listing_id")
		return
	}
	files := h.Store.GetListingFiles(listingID)
	if files == nil {
		files = []models.ListingFile{}
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   len(files),
		Results: files,
	})
}

// GET /v3/application/shops/{shop_id}/listings/{listing_id}/files/{listing_file_id}
func (h *Handler) GetListingFile(w http.ResponseWriter, r *http.Request) {
	listingID, ok := extractPathID(r.URL.Path, "listings")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid listing_id")
		return
	}
	fileID, ok := extractPathID(r.URL.Path, "files")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid listing_file_id")
		return
	}
	f, found := h.Store.GetListingFile(listingID, fileID)
	if !found {
		writeError(w, http.StatusNotFound, "File not found")
		return
	}
	writeJSON(w, http.StatusOK, f)
}

// POST /v3/application/shops/{shop_id}/listings/{listing_id}/files
func (h *Handler) UploadListingFile(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "listings_w") {
		return
	}
	listingID, ok := extractPathID(r.URL.Path, "listings")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid listing_id")
		return
	}
	if _, found := h.Store.GetListing(listingID); !found {
		writeError(w, http.StatusNotFound, "Listing not found")
		return
	}

	filename := r.FormValue("filename")
	if filename == "" {
		filename = "file.pdf"
	}
	filetype := r.FormValue("filetype")
	if filetype == "" {
		filetype = "application/pdf"
	}

	f := h.Store.AddListingFile(listingID, filename, filetype, 1024)
	writeJSON(w, http.StatusCreated, f)
}

// DELETE /v3/application/shops/{shop_id}/listings/{listing_id}/files/{listing_file_id}
func (h *Handler) DeleteListingFile(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "listings_w") {
		return
	}
	listingID, ok := extractPathID(r.URL.Path, "listings")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid listing_id")
		return
	}
	fileID, ok := extractPathID(r.URL.Path, "files")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid listing_file_id")
		return
	}
	if !h.Store.DeleteListingFile(listingID, fileID) {
		writeError(w, http.StatusNotFound, "File not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GET /v3/application/listings/{listing_id}/inventory
func (h *Handler) GetListingInventory(w http.ResponseWriter, r *http.Request) {
	listingID, ok := extractPathID(r.URL.Path, "listings")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid listing_id")
		return
	}
	inv, found := h.Store.GetListingInventory(listingID)
	if !found {
		writeError(w, http.StatusNotFound, "Listing not found")
		return
	}
	writeJSON(w, http.StatusOK, inv)
}
