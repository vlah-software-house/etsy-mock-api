package handlers

import (
	"net/http"

	"github.com/gabriel/etsy-mock/internal/models"
)

// GET /v3/application/shops/{shop_id}/receipts
func (h *Handler) GetShopReceipts(w http.ResponseWriter, r *http.Request) {
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

	receipts, total := h.Store.GetShopReceipts(shopID, limit, offset)
	if receipts == nil {
		receipts = []models.ShopReceipt{}
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   total,
		Results: receipts,
	})
}

// GET /v3/application/shops/{shop_id}/receipts/{receipt_id}
func (h *Handler) GetReceipt(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "transactions_r") {
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
	writeJSON(w, http.StatusOK, receipt)
}

// PUT /v3/application/shops/{shop_id}/receipts/{receipt_id}
func (h *Handler) UpdateReceipt(w http.ResponseWriter, r *http.Request) {
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

	var updates map[string]interface{}
	if err := decodeJSON(r, &updates); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if v, ok := updates["message_from_seller"].(string); ok {
		receipt.MessageFromSeller = &v
	}
	if v, ok := updates["is_shipped"].(bool); ok {
		receipt.IsShipped = v
	}
	if v, ok := updates["is_paid"].(bool); ok {
		receipt.IsPaid = v
	}

	h.Store.UpdateReceipt(receipt)
	writeJSON(w, http.StatusOK, receipt)
}

// GET /v3/application/shops/{shop_id}/transactions
func (h *Handler) GetShopTransactions(w http.ResponseWriter, r *http.Request) {
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

	txns, total := h.Store.GetShopTransactions(shopID, limit, offset)
	if txns == nil {
		txns = []models.ShopReceiptTransaction{}
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   total,
		Results: txns,
	})
}

// GET /v3/application/shops/{shop_id}/transactions/{transaction_id}
func (h *Handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "transactions_r") {
		return
	}
	txnID, ok := extractPathID(r.URL.Path, "transactions")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid transaction_id")
		return
	}
	txn, found := h.Store.GetTransaction(txnID)
	if !found {
		writeError(w, http.StatusNotFound, "Transaction not found")
		return
	}
	writeJSON(w, http.StatusOK, txn)
}

// GET /v3/application/shops/{shop_id}/receipts/{receipt_id}/transactions
func (h *Handler) GetReceiptTransactions(w http.ResponseWriter, r *http.Request) {
	if !requireScope(w, r, "transactions_r") {
		return
	}
	receiptID, ok := extractPathID(r.URL.Path, "receipts")
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid receipt_id")
		return
	}
	txns := h.Store.GetReceiptTransactions(receiptID)
	if txns == nil {
		txns = []models.ShopReceiptTransaction{}
	}
	writeJSON(w, http.StatusOK, models.PaginatedResponse{
		Count:   len(txns),
		Results: txns,
	})
}
