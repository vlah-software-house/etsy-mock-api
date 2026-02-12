package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vlah-software-house/etsy-mock-api/internal/models"
)

// GET /v3/application/listings/{listing_id}/personalization
func (h *Handler) GetListingPersonalization(w http.ResponseWriter, r *http.Request) {
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
	writeJSON(w, http.StatusOK, models.ListingPersonalization{
		IsPersonalizable:            listing.IsPersonalizable,
		PersonalizationIsRequired:   listing.PersonalizationIsRequired,
		PersonalizationCharCountMax: listing.PersonalizationCharCountMax,
		PersonalizationInstructions: listing.PersonalizationInstructions,
	})
}

// PUT /v3/application/shops/{shop_id}/listings/{listing_id}/personalization
func (h *Handler) UpdateListingPersonalization(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "listings_w") {
		return
	}
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

	var body struct {
		IsPersonalizable          *bool   `json:"is_personalizable"`
		PersonalizationIsRequired *bool   `json:"personalization_is_required"`
		PersonalizationCharCountMax *int  `json:"personalization_char_count_max"`
		PersonalizationInstructions *string `json:"personalization_instructions"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if body.IsPersonalizable != nil {
		listing.IsPersonalizable = *body.IsPersonalizable
	}
	if body.PersonalizationIsRequired != nil {
		listing.PersonalizationIsRequired = *body.PersonalizationIsRequired
	}
	if body.PersonalizationCharCountMax != nil {
		listing.PersonalizationCharCountMax = body.PersonalizationCharCountMax
	}
	if body.PersonalizationInstructions != nil {
		listing.PersonalizationInstructions = body.PersonalizationInstructions
	}
	listing.UpdatedTimestamp = time.Now().Unix()
	listing.LastModifiedTimestamp = time.Now().Unix()

	writeJSON(w, http.StatusOK, models.ListingPersonalization{
		IsPersonalizable:            listing.IsPersonalizable,
		PersonalizationIsRequired:   listing.PersonalizationIsRequired,
		PersonalizationCharCountMax: listing.PersonalizationCharCountMax,
		PersonalizationInstructions: listing.PersonalizationInstructions,
	})
}

// GET /v3/application/listings/{listing_id}/videos
func (h *Handler) GetListingVideos(w http.ResponseWriter, r *http.Request) {
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
	videos := listing.Videos
	if videos == nil {
		videos = []models.ListingVideo{}
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   len(videos),
		Results: videos,
	})
}

// GET /v3/application/listings/{listing_id}/videos/{video_id}
func (h *Handler) GetListingVideo(w http.ResponseWriter, r *http.Request) {
	listingID, ok := extractPathID(r.URL.Path, "listings")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid listing_id")
		return
	}
	videoID, ok := extractPathID(r.URL.Path, "videos")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid video_id")
		return
	}
	listing, found := h.Store.GetListing(listingID)
	if !found {
		writeError(w, http.StatusNotFound, "Listing not found")
		return
	}
	for _, v := range listing.Videos {
		if v.VideoID == videoID {
			writeJSON(w, http.StatusOK, v)
			return
		}
	}
	writeError(w, http.StatusNotFound, "Video not found")
}

// POST /v3/application/shops/{shop_id}/listings/{listing_id}/videos
func (h *Handler) UploadListingVideo(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "listings_w") {
		return
	}
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
	id := h.Store.NextID()
	video := models.ListingVideo{
		VideoID:      id,
		Height:       1080,
		Width:        1920,
		ThumbnailURL: fmt.Sprintf("https://mock.etsy.com/videos/%d_thumb.jpg", id),
		VideoURL:     fmt.Sprintf("https://mock.etsy.com/videos/%d.mp4", id),
		VideoState:   "active",
	}
	listing.Videos = append(listing.Videos, video)
	writeJSON(w, http.StatusCreated, video)
}

// DELETE /v3/application/shops/{shop_id}/listings/{listing_id}/videos/{video_id}
func (h *Handler) DeleteListingVideo(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "listings_w") {
		return
	}
	listingID, ok := extractPathID(r.URL.Path, "listings")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid listing_id")
		return
	}
	videoID, ok := extractPathID(r.URL.Path, "videos")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid video_id")
		return
	}
	listing, found := h.Store.GetListing(listingID)
	if !found {
		writeError(w, http.StatusNotFound, "Listing not found")
		return
	}
	for i, v := range listing.Videos {
		if v.VideoID == videoID {
			listing.Videos = append(listing.Videos[:i], listing.Videos[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	writeError(w, http.StatusNotFound, "Video not found")
}

// GET /v3/application/shops/{shop_id}/listings/{listing_id}/translations/{language}
func (h *Handler) GetListingTranslation(w http.ResponseWriter, r *http.Request) {
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
	lang := extractPathSegment(r.URL.Path, "translations")
	writeJSON(w, http.StatusOK, models.ListingTranslation{
		ListingID:   listingID,
		Language:    lang,
		Title:       &listing.Title,
		Description: &listing.Description,
		Tags:        listing.Tags,
	})
}

// PUT /v3/application/shops/{shop_id}/listings/{listing_id}/translations/{language}
func (h *Handler) UpdateListingTranslation(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "listings_w") {
		return
	}
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
	lang := extractPathSegment(r.URL.Path, "translations")

	var body struct {
		Title       *string  `json:"title"`
		Description *string  `json:"description"`
		Tags        []string `json:"tags"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	title := listing.Title
	desc := listing.Description
	tags := listing.Tags
	if body.Title != nil {
		title = *body.Title
	}
	if body.Description != nil {
		desc = *body.Description
	}
	if body.Tags != nil {
		tags = body.Tags
	}

	writeJSON(w, http.StatusOK, models.ListingTranslation{
		ListingID:   listingID,
		Language:    lang,
		Title:       &title,
		Description: &desc,
		Tags:        tags,
	})
}

// GET /v3/application/shops/{shop_id}/listings/{listing_id}/variation-images
func (h *Handler) GetListingVariationImages(w http.ResponseWriter, r *http.Request) {
	_, ok := extractPathID(r.URL.Path, "listings")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid listing_id")
		return
	}
	// Return empty set by default
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   0,
		Results: []models.ListingVariationImage{},
	})
}

// POST /v3/application/shops/{shop_id}/listings/{listing_id}/variation-images
func (h *Handler) UpdateListingVariationImages(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "listings_w") {
		return
	}
	_, ok := extractPathID(r.URL.Path, "listings")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid listing_id")
		return
	}
	var body struct {
		VariationImages []models.ListingVariationImage `json:"variation_images"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if body.VariationImages == nil {
		body.VariationImages = []models.ListingVariationImage{}
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   len(body.VariationImages),
		Results: body.VariationImages,
	})
}

// GET /v3/application/shops/{shop_id}/listings/{listing_id}/properties
func (h *Handler) GetListingProperties(w http.ResponseWriter, r *http.Request) {
	_, ok := extractPathID(r.URL.Path, "listings")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid listing_id")
		return
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   0,
		Results: []models.ListingPropertyValue{},
	})
}

// GET /v3/application/listings/{listing_id}/products/{product_id}/offerings/{offering_id}
func (h *Handler) GetListingOffering(w http.ResponseWriter, r *http.Request) {
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
	offeringID, ok := extractPathID(r.URL.Path, "offerings")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid offering_id")
		return
	}
	for _, p := range inv.Products {
		for _, o := range p.Offerings {
			if o.OfferingID == offeringID {
				writeJSON(w, http.StatusOK, o)
				return
			}
		}
	}
	writeError(w, http.StatusNotFound, "Offering not found")
}

// GET /v3/application/shops/{shop_id}/production-partners
func (h *Handler) GetProductionPartners(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "shops_r") {
		return
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   0,
		Results: []models.ShopProductionPartner{},
	})
}

// GET /v3/application/shops/{shop_id}/holiday-preferences
func (h *Handler) GetHolidayPreferences(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "shops_r") {
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
	writeJSON(w, http.StatusOK, models.ShopHolidayPreferences{
		ShopID:          shopID,
		IsVacation:      shop.IsVacation,
		VacationMessage: shop.VacationMessage,
	})
}

// PUT /v3/application/shops/{shop_id}/holiday-preferences
func (h *Handler) UpdateHolidayPreferences(w http.ResponseWriter, r *http.Request) {
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
		IsVacation      *bool   `json:"is_vacation"`
		VacationMessage *string `json:"vacation_message"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if body.IsVacation != nil {
		shop.IsVacation = *body.IsVacation
	}
	if body.VacationMessage != nil {
		shop.VacationMessage = body.VacationMessage
	}
	h.Store.UpdateShop(shop)

	writeJSON(w, http.StatusOK, models.ShopHolidayPreferences{
		ShopID:          shopID,
		IsVacation:      shop.IsVacation,
		VacationMessage: shop.VacationMessage,
	})
}

// GET /v3/application/shops/{shop_id}/readiness-state-definitions
func (h *Handler) GetReadinessStateDefinitions(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "shops_r") {
		return
	}
	defs := []models.ReadinessStateDefinition{
		{ReadinessStateID: 1, Name: "ready_to_ship", Label: "Ready to ship", ProcessingTimeUnit: "business_days", ProcessingMin: 1, ProcessingMax: 3},
		{ReadinessStateID: 2, Name: "made_to_order", Label: "Made to order", ProcessingTimeUnit: "business_days", ProcessingMin: 3, ProcessingMax: 7},
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   len(defs),
		Results: defs,
	})
}

// GET /v3/application/shipping-carriers
func (h *Handler) GetShippingCarriers(w http.ResponseWriter, r *http.Request) {
	carriers := []models.ShippingCarrier{
		{ShippingCarrierID: 1, Name: "USPS", DomesticClasses: []models.ShippingCarrierMailClass{{MailClassKey: "usps_first_class", Name: "First Class"}, {MailClassKey: "usps_priority", Name: "Priority Mail"}, {MailClassKey: "usps_priority_express", Name: "Priority Mail Express"}}, InternationalClasses: []models.ShippingCarrierMailClass{{MailClassKey: "usps_first_class_international", Name: "First Class International"}, {MailClassKey: "usps_priority_international", Name: "Priority Mail International"}}},
		{ShippingCarrierID: 2, Name: "UPS", DomesticClasses: []models.ShippingCarrierMailClass{{MailClassKey: "ups_ground", Name: "Ground"}, {MailClassKey: "ups_2day", Name: "2nd Day Air"}, {MailClassKey: "ups_next_day", Name: "Next Day Air"}}, InternationalClasses: []models.ShippingCarrierMailClass{{MailClassKey: "ups_worldwide_express", Name: "Worldwide Express"}}},
		{ShippingCarrierID: 3, Name: "FedEx", DomesticClasses: []models.ShippingCarrierMailClass{{MailClassKey: "fedex_ground", Name: "Ground"}, {MailClassKey: "fedex_2day", Name: "2Day"}, {MailClassKey: "fedex_overnight", Name: "Standard Overnight"}}, InternationalClasses: []models.ShippingCarrierMailClass{{MailClassKey: "fedex_international_economy", Name: "International Economy"}}},
		{ShippingCarrierID: 4, Name: "Canada Post", DomesticClasses: []models.ShippingCarrierMailClass{{MailClassKey: "canadapost_regular", Name: "Regular Parcel"}, {MailClassKey: "canadapost_expedited", Name: "Expedited Parcel"}}, InternationalClasses: []models.ShippingCarrierMailClass{{MailClassKey: "canadapost_international", Name: "International Parcel"}}},
		{ShippingCarrierID: 5, Name: "Royal Mail", DomesticClasses: []models.ShippingCarrierMailClass{{MailClassKey: "royalmail_first", Name: "1st Class"}, {MailClassKey: "royalmail_second", Name: "2nd Class"}}, InternationalClasses: []models.ShippingCarrierMailClass{{MailClassKey: "royalmail_international_standard", Name: "International Standard"}}},
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   len(carriers),
		Results: carriers,
	})
}

// GET /v3/application/listings/batch
func (h *Handler) GetListingsBatch(w http.ResponseWriter, r *http.Request) {
	idsParam := r.URL.Query().Get("listing_ids")
	if idsParam == "" {
		writeError(w, http.StatusBadRequest, "listing_ids parameter is required")
		return
	}

	var listings []models.ShopListing
	for _, idStr := range splitCSV(idsParam) {
		id, ok := parseID(idStr)
		if !ok {
			continue
		}
		if listing, found := h.Store.GetListing(id); found {
			listing.Images = h.Store.GetListingImages(id)
			listings = append(listings, *listing)
		}
	}
	if listings == nil {
		listings = []models.ShopListing{}
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   len(listings),
		Results: listings,
	})
}

// GET /v3/application/shops/{shop_id}/listings/featured
func (h *Handler) GetFeaturedListings(w http.ResponseWriter, r *http.Request) {
	shopID, ok := extractPathID(r.URL.Path, "shops")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid shop_id")
		return
	}
	limit := queryInt(r, "limit", 25)
	offset := queryInt(r, "offset", 0)

	listings, total := h.Store.GetShopListings(shopID, "active", limit, offset)
	// Filter to featured (featured_rank > 0) - in mock, return first few actives
	var featured []models.ShopListing
	for _, l := range listings {
		if l.FeaturedRank > 0 || len(featured) < 3 {
			featured = append(featured, l)
		}
	}
	if featured == nil {
		featured = []models.ShopListing{}
	}
	_ = total
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   len(featured),
		Results: featured,
	})
}

// POST /v3/application/shops/{shop_id}/receipts/{receipt_id}/tracking
func (h *Handler) CreateReceiptTracking(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "transactions_w") {
		return
	}
	receiptID, ok := extractPathID(r.URL.Path, "receipts")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid receipt_id")
		return
	}
	receipt, found := h.Store.GetReceipt(receiptID)
	if !found {
		writeError(w, http.StatusNotFound, "Receipt not found")
		return
	}

	var body models.ReceiptShipmentTracking
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	shipID := h.Store.NextID()
	receipt.Shipments = append(receipt.Shipments, models.ShopReceiptShipment{
		ReceiptShippingID:             &shipID,
		ShipmentNotificationTimestamp: time.Now().Unix(),
		CarrierName:                   body.CarrierName,
		TrackingCode:                  body.TrackingCode,
	})
	receipt.IsShipped = true
	h.Store.UpdateReceipt(receipt)

	writeJSON(w, http.StatusCreated, receipt)
}

// GET /v3/application/shops/{shop_id}/sections/{shop_section_id}
func (h *Handler) GetShopSection(w http.ResponseWriter, r *http.Request) {
	shopID, ok := extractPathID(r.URL.Path, "shops")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid shop_id")
		return
	}
	sectionID, ok := extractPathID(r.URL.Path, "sections")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid section_id")
		return
	}
	sec, found := h.Store.GetShopSection(shopID, sectionID)
	if !found {
		writeError(w, http.StatusNotFound, "Section not found")
		return
	}
	writeJSON(w, http.StatusOK, sec)
}

// GET /v3/application/shops/{shop_id}/shop-sections/listings
func (h *Handler) GetShopSectionListings(w http.ResponseWriter, r *http.Request) {
	shopID, ok := extractPathID(r.URL.Path, "shops")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid shop_id")
		return
	}
	limit := queryInt(r, "limit", 25)
	offset := queryInt(r, "offset", 0)
	sectionIDs := r.URL.Query().Get("shop_section_ids")

	listings, total := h.Store.GetShopListings(shopID, "active", limit, offset)
	if sectionIDs != "" {
		var filtered []models.ShopListing
		for _, idStr := range splitCSV(sectionIDs) {
			secID, ok := parseID(idStr)
			if !ok {
				continue
			}
			for _, l := range listings {
				if l.ShopSectionID != nil && *l.ShopSectionID == secID {
					filtered = append(filtered, l)
				}
			}
		}
		listings = filtered
		total = len(filtered)
	}
	if listings == nil {
		listings = []models.ShopListing{}
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   total,
		Results: listings,
	})
}

// Seller Taxonomy (mirrors Buyer Taxonomy endpoints)

// GET /v3/application/seller-taxonomy/nodes
func (h *Handler) GetSellerTaxonomyNodes(w http.ResponseWriter, r *http.Request) {
	nodes := h.Store.GetTaxonomyNodes()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"count":   len(nodes),
		"results": nodes,
	})
}

// GET /v3/application/seller-taxonomy/nodes/{taxonomy_id}/properties
func (h *Handler) GetSellerTaxonomyProperties(w http.ResponseWriter, r *http.Request) {
	taxonomyID, ok := extractPathID(r.URL.Path, "nodes")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid taxonomy_id")
		return
	}
	props, found := h.Store.GetTaxonomyProperties(taxonomyID)
	if !found {
		writeError(w, http.StatusNotFound, "Taxonomy not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"count":   len(props),
		"results": props,
	})
}

// DELETE /v3/application/users/{user_id}/addresses/{user_address_id}
func (h *Handler) DeleteUserAddress(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "address_r") {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Helper: extract a string path segment after a given key
func extractPathSegment(path string, after string) string {
	parts := splitPath(path)
	for i, p := range parts {
		if p == after && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

func splitPath(path string) []string {
	var parts []string
	for _, p := range splitCSV(path) {
		parts = append(parts, p)
	}
	return strings.Split(path, "/")
}

func splitCSV(s string) []string {
	var result []string
	for _, part := range strings.Split(s, ",") {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
