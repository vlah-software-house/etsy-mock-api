package handlers

import (
	"net/http"
)

// GET /v3/application/buyer-taxonomy/nodes
func (h *Handler) GetBuyerTaxonomyNodes(w http.ResponseWriter, r *http.Request) {
	nodes := h.Store.GetTaxonomyNodes()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"count":   len(nodes),
		"results": nodes,
	})
}

// GET /v3/application/buyer-taxonomy/nodes/{taxonomy_id}/properties
func (h *Handler) GetBuyerTaxonomyProperties(w http.ResponseWriter, r *http.Request) {
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
