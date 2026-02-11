package handlers

import (
	"net/http"

	"github.com/gabriel/etsy-mock/internal/models"
)

// GET /v3/application/shops/{shop_id}/receipts/{receipt_id}/payments
func (h *Handler) GetReceiptPayments(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "transactions_r") {
		return
	}
	receiptID, ok := extractPathID(r.URL.Path, "receipts")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid receipt_id")
		return
	}
	payments := h.Store.GetPaymentsByReceipt(receiptID)
	if payments == nil {
		payments = []models.Payment{}
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   len(payments),
		Results: payments,
	})
}

// GET /v3/application/shops/{shop_id}/payments
func (h *Handler) GetShopPayments(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "transactions_r") {
		return
	}
	shopID, ok := extractPathID(r.URL.Path, "shops")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid shop_id")
		return
	}
	payments := h.Store.GetPaymentsByShop(shopID)
	if payments == nil {
		payments = []models.Payment{}
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   len(payments),
		Results: payments,
	})
}

// GET /v3/application/shops/{shop_id}/payment-account/ledger-entries
func (h *Handler) GetLedgerEntries(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "transactions_r") {
		return
	}
	shopID, ok := extractPathID(r.URL.Path, "shops")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid shop_id")
		return
	}
	limit := queryInt(r, "limit", 25)
	offset := queryInt(r, "offset", 0)

	entries, total := h.Store.GetLedgerEntries(shopID, limit, offset)
	if entries == nil {
		entries = []models.PaymentAccountLedgerEntry{}
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   total,
		Results: entries,
	})
}
