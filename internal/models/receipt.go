package models

type ShopReceipt struct {
	ReceiptID          int64                    `json:"receipt_id"`
	ReceiptType        int                      `json:"receipt_type"`
	SellerUserID       int64                    `json:"seller_user_id"`
	SellerEmail        *string                  `json:"seller_email"`
	BuyerUserID        int64                    `json:"buyer_user_id"`
	BuyerEmail         *string                  `json:"buyer_email"`
	Name               string                   `json:"name"`
	FirstLine          *string                  `json:"first_line"`
	SecondLine         *string                  `json:"second_line"`
	City               *string                  `json:"city"`
	State              *string                  `json:"state"`
	Zip                *string                  `json:"zip"`
	Status             string                   `json:"status"`
	FormattedAddress   *string                  `json:"formatted_address"`
	CountryISO         *string                  `json:"country_iso"`
	PaymentMethod      string                   `json:"payment_method"`
	PaymentEmail       *string                  `json:"payment_email"`
	MessageFromSeller  *string                  `json:"message_from_seller"`
	MessageFromBuyer   *string                  `json:"message_from_buyer"`
	MessageFromPayment *string                  `json:"message_from_payment"`
	IsPaid             bool                     `json:"is_paid"`
	IsShipped          bool                     `json:"is_shipped"`
	CreateTimestamp    int64                    `json:"create_timestamp"`
	CreatedTimestamp    int64                    `json:"created_timestamp"`
	UpdateTimestamp     int64                    `json:"update_timestamp"`
	UpdatedTimestamp    int64                    `json:"updated_timestamp"`
	IsGift             bool                     `json:"is_gift"`
	GiftMessage        string                   `json:"gift_message"`
	GiftSender         string                   `json:"gift_sender"`
	Grandtotal         Money                    `json:"grandtotal"`
	Subtotal           Money                    `json:"subtotal"`
	TotalPrice         Money                    `json:"total_price"`
	TotalShippingCost  Money                    `json:"total_shipping_cost"`
	TotalTaxCost       Money                    `json:"total_tax_cost"`
	TotalVatCost       Money                    `json:"total_vat_cost"`
	DiscountAmt        Money                    `json:"discount_amt"`
	GiftWrapPrice      Money                    `json:"gift_wrap_price"`
	Shipments          []ShopReceiptShipment    `json:"shipments"`
	Transactions       []ShopReceiptTransaction `json:"transactions"`
	Refunds            []ShopRefund             `json:"refunds"`
}

type ShopReceiptShipment struct {
	ReceiptShippingID              *int64 `json:"receipt_shipping_id"`
	ShipmentNotificationTimestamp  int64  `json:"shipment_notification_timestamp"`
	CarrierName                    string `json:"carrier_name"`
	TrackingCode                   string `json:"tracking_code"`
}

type ShopReceiptTransaction struct {
	TransactionID     int64                  `json:"transaction_id"`
	Title             *string                `json:"title"`
	Description       *string                `json:"description"`
	SellerUserID      int64                  `json:"seller_user_id"`
	BuyerUserID       int64                  `json:"buyer_user_id"`
	CreateTimestamp   int64                  `json:"create_timestamp"`
	CreatedTimestamp   int64                  `json:"created_timestamp"`
	PaidTimestamp     *int64                 `json:"paid_timestamp"`
	ShippedTimestamp  *int64                 `json:"shipped_timestamp"`
	Quantity          int                    `json:"quantity"`
	ListingImageID    *int64                 `json:"listing_image_id"`
	ReceiptID         int64                  `json:"receipt_id"`
	IsDigital         bool                   `json:"is_digital"`
	FileData          string                 `json:"file_data"`
	ListingID         *int                   `json:"listing_id"`
	TransactionType   string                 `json:"transaction_type"`
	ProductID         *int64                 `json:"product_id"`
	SKU               *string                `json:"sku"`
	Price             Money                  `json:"price"`
	ShippingCost      Money                  `json:"shipping_cost"`
	Variations        []TransactionVariation `json:"variations"`
	ShippingProfileID *int64                 `json:"shipping_profile_id"`
	MinProcessingDays *int                   `json:"min_processing_days"`
	MaxProcessingDays *int                   `json:"max_processing_days"`
	ShippingMethod    *string                `json:"shipping_method"`
	ShippingUpgrade   *string                `json:"shipping_upgrade"`
	ExpectedShipDate  *int64                 `json:"expected_ship_date"`
	BuyerCoupon       float64                `json:"buyer_coupon"`
	ShopCoupon        float64                `json:"shop_coupon"`
}

type TransactionVariation struct {
	PropertyID   int64  `json:"property_id"`
	ValueID      int64  `json:"value_id"`
	FormattedName  string `json:"formatted_name"`
	FormattedValue string `json:"formatted_value"`
}

type ShopRefund struct {
	Amount           Money   `json:"amount"`
	CreatedTimestamp  int64   `json:"created_timestamp"`
	Reason           *string `json:"reason"`
	NoteFromIssuer   *string `json:"note_from_issuer"`
	Status           *string `json:"status"`
}
