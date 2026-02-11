package handlers

import (
	"net/http"
)

// GET /v3/application/shops/{shop_id}
func (h *Handler) GetShop(w http.ResponseWriter, r *http.Request) {
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
	writeJSON(w, http.StatusOK, shop)
}

// GET /v3/application/shops?shop_name=xxx
func (h *Handler) FindShops(w http.ResponseWriter, r *http.Request) {
	shopName := r.URL.Query().Get("shop_name")
	if shopName == "" {
		writeError(w, http.StatusBadRequest, "shop_name parameter is required")
		return
	}
	shop, found := h.Store.GetShopByName(shopName)
	if !found {
		writeError(w, http.StatusNotFound, "Shop not found")
		return
	}
	writeJSON(w, http.StatusOK, shop)
}

// PUT /v3/application/shops/{shop_id}
func (h *Handler) UpdateShop(w http.ResponseWriter, r *http.Request) {
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

	var updates map[string]interface{}
	if err := decodeJSON(r, &updates); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if v, ok := updates["title"].(string); ok {
		shop.Title = &v
	}
	if v, ok := updates["announcement"].(string); ok {
		shop.Announcement = &v
	}
	if v, ok := updates["sale_message"].(string); ok {
		shop.SaleMessage = &v
	}
	if v, ok := updates["is_vacation"].(bool); ok {
		shop.IsVacation = v
	}
	if v, ok := updates["vacation_message"].(string); ok {
		shop.VacationMessage = &v
	}

	h.Store.UpdateShop(shop)
	writeJSON(w, http.StatusOK, shop)
}

// GET /v3/application/shops/{shop_id}/sections
func (h *Handler) GetShopSections(w http.ResponseWriter, r *http.Request) {
	shopID, ok := extractPathID(r.URL.Path, "shops")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid shop_id")
		return
	}
	sections := h.Store.GetShopSections(shopID)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"count":   len(sections),
		"results": sections,
	})
}

// POST /v3/application/shops/{shop_id}/sections
func (h *Handler) CreateShopSection(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "shops_w") {
		return
	}
	shopID, ok := extractPathID(r.URL.Path, "shops")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid shop_id")
		return
	}
	var body struct {
		Title string `json:"title"`
		Rank  int    `json:"rank"`
	}
	if err := decodeJSON(r, &body); err != nil || body.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	sec := h.Store.CreateShopSection(shopID, body.Title, body.Rank)
	writeJSON(w, http.StatusCreated, sec)
}

// GET /v3/application/shops/{shop_id}/return-policies
func (h *Handler) GetReturnPolicies(w http.ResponseWriter, r *http.Request) {
	shopID, ok := extractPathID(r.URL.Path, "shops")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid shop_id")
		return
	}
	policies := h.Store.GetReturnPolicies(shopID)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"count":   len(policies),
		"results": policies,
	})
}

// POST /v3/application/shops/{shop_id}/return-policies
func (h *Handler) CreateReturnPolicy(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "shops_w") {
		return
	}
	shopID, ok := extractPathID(r.URL.Path, "shops")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid shop_id")
		return
	}
	var body struct {
		AcceptsReturns   bool `json:"accepts_returns"`
		AcceptsExchanges bool `json:"accepts_exchanges"`
		ReturnDeadline   *int `json:"return_deadline"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	p := h.Store.CreateReturnPolicy(shopID, body.AcceptsReturns, body.AcceptsExchanges, body.ReturnDeadline)
	writeJSON(w, http.StatusCreated, p)
}

// GET /v3/application/shops/{shop_id}/return-policies/{return_policy_id}
func (h *Handler) GetReturnPolicy(w http.ResponseWriter, r *http.Request) {
	policyID, ok := extractPathID(r.URL.Path, "return-policies")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid return_policy_id")
		return
	}
	p, found := h.Store.GetReturnPolicy(policyID)
	if !found {
		writeError(w, http.StatusNotFound, "Return policy not found")
		return
	}
	writeJSON(w, http.StatusOK, p)
}
