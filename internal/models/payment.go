package models

type Payment struct {
	PaymentID            int64               `json:"payment_id"`
	BuyerUserID          int64               `json:"buyer_user_id"`
	ShopID               int64               `json:"shop_id"`
	ReceiptID            int64               `json:"receipt_id"`
	AmountGross          Money               `json:"amount_gross"`
	AmountFees           Money               `json:"amount_fees"`
	AmountNet            Money               `json:"amount_net"`
	PostedGross          *Money              `json:"posted_gross"`
	PostedFees           *Money              `json:"posted_fees"`
	PostedNet            *Money              `json:"posted_net"`
	AdjustedGross        *Money              `json:"adjusted_gross"`
	AdjustedFees         *Money              `json:"adjusted_fees"`
	AdjustedNet          *Money              `json:"adjusted_net"`
	Currency             string              `json:"currency"`
	ShopCurrency         *string             `json:"shop_currency"`
	BuyerCurrency        *string             `json:"buyer_currency"`
	ShippingUserID       *int64              `json:"shipping_user_id"`
	ShippingAddressID    int64               `json:"shipping_address_id"`
	BillingAddressID     int                 `json:"billing_address_id"`
	Status               string              `json:"status"`
	ShippedTimestamp      *int64              `json:"shipped_timestamp"`
	CreateTimestamp      int64               `json:"create_timestamp"`
	CreatedTimestamp      int64               `json:"created_timestamp"`
	UpdateTimestamp       int64               `json:"update_timestamp"`
	UpdatedTimestamp      int64               `json:"updated_timestamp"`
	PaymentAdjustments   []PaymentAdjustment `json:"payment_adjustments"`
}

type PaymentAdjustment struct {
	PaymentAdjustmentID        int64                   `json:"payment_adjustment_id"`
	PaymentID                  int64                   `json:"payment_id"`
	Status                     string                  `json:"status"`
	IsSuccess                  bool                    `json:"is_success"`
	UserID                     int64                   `json:"user_id"`
	ReasonCode                 string                  `json:"reason_code"`
	TotalAdjustmentAmount      *int                    `json:"total_adjustment_amount"`
	ShopTotalAdjustmentAmount  *int                    `json:"shop_total_adjustment_amount"`
	BuyerTotalAdjustmentAmount *int                    `json:"buyer_total_adjustment_amount"`
	TotalFeeAdjustmentAmount   *int                    `json:"total_fee_adjustment_amount"`
	CreateTimestamp            int64                   `json:"create_timestamp"`
	CreatedTimestamp            int64                   `json:"created_timestamp"`
	UpdateTimestamp             int64                   `json:"update_timestamp"`
	UpdatedTimestamp            int64                   `json:"updated_timestamp"`
	PaymentAdjustmentItems     []PaymentAdjustmentItem `json:"payment_adjustment_items"`
}

type PaymentAdjustmentItem struct {
	PaymentAdjustmentID     int64  `json:"payment_adjustment_id"`
	PaymentAdjustmentItemID int64  `json:"payment_adjustment_item_id"`
	AdjustmentType          *string `json:"adjustment_type"`
	Amount                  int    `json:"amount"`
	ShopAmount              int    `json:"shop_amount"`
	TransactionID           *int64 `json:"transaction_id"`
	BillPaymentID           *int64 `json:"bill_payment_id"`
	CreatedTimestamp         int64  `json:"created_timestamp"`
	UpdatedTimestamp         int64  `json:"updated_timestamp"`
}

type PaymentAccountLedgerEntry struct {
	EntryID            int64               `json:"entry_id"`
	LedgerID           int64               `json:"ledger_id"`
	SequenceNumber     int                 `json:"sequence_number"`
	Amount             int                 `json:"amount"`
	Currency           string              `json:"currency"`
	Description        string              `json:"description"`
	Balance            int                 `json:"balance"`
	CreateDate         int64               `json:"create_date"`
	CreatedTimestamp    int64               `json:"created_timestamp"`
	LedgerType         string              `json:"ledger_type"`
	ReferenceType      string              `json:"reference_type"`
	ReferenceID        *string             `json:"reference_id"`
	ParentEntryID      int                 `json:"parent_entry_id"`
	PaymentAdjustments []PaymentAdjustment `json:"payment_adjustments"`
}
