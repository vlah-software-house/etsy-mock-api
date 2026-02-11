package models

type Shop struct {
	ShopID                        int64    `json:"shop_id"`
	UserID                        int64    `json:"user_id"`
	ShopName                      string   `json:"shop_name"`
	CreateDate                    int64    `json:"create_date"`
	CreatedTimestamp               int64    `json:"created_timestamp"`
	Title                         *string  `json:"title"`
	Announcement                  *string  `json:"announcement"`
	CurrencyCode                  string   `json:"currency_code"`
	IsVacation                    bool     `json:"is_vacation"`
	VacationMessage               *string  `json:"vacation_message"`
	SaleMessage                   *string  `json:"sale_message"`
	DigitalSaleMessage            *string  `json:"digital_sale_message"`
	UpdateDate                    int64    `json:"update_date"`
	UpdatedTimestamp               int64    `json:"updated_timestamp"`
	ListingActiveCount            int      `json:"listing_active_count"`
	DigitalListingCount           int      `json:"digital_listing_count"`
	LoginName                     string   `json:"login_name"`
	AcceptsCustomRequests         bool     `json:"accepts_custom_requests"`
	PolicyWelcome                 *string  `json:"policy_welcome"`
	PolicyPayment                 *string  `json:"policy_payment"`
	PolicyShipping                *string  `json:"policy_shipping"`
	PolicyRefunds                 *string  `json:"policy_refunds"`
	PolicyAdditional              *string  `json:"policy_additional"`
	PolicySellerInfo              *string  `json:"policy_seller_info"`
	PolicyUpdateDate              int64    `json:"policy_update_date"`
	PolicyHasPrivateReceiptInfo   bool     `json:"policy_has_private_receipt_info"`
	HasUnstructuredPolicies       bool     `json:"has_unstructured_policies"`
	PolicyPrivacy                 *string  `json:"policy_privacy"`
	VacationAutoreply             *string  `json:"vacation_autoreply"`
	URL                           string   `json:"url"`
	ImageURL760x100               *string  `json:"image_url_760x100"`
	NumFavorers                   int      `json:"num_favorers"`
	Languages                     []string `json:"languages"`
	IconURLFullxfull              *string  `json:"icon_url_fullxfull"`
	IsUsingStructuredPolicies     bool     `json:"is_using_structured_policies"`
	HasOnboardedStructuredPolicies bool    `json:"has_onboarded_structured_policies"`
	IncludeDisputeFormLink        bool     `json:"include_dispute_form_link"`
	IsDirectCheckoutOnboarded     bool     `json:"is_direct_checkout_onboarded"`
	IsEtsyPaymentsOnboarded       bool     `json:"is_etsy_payments_onboarded"`
	IsCalculatedEligible          bool     `json:"is_calculated_eligible"`
	IsOptedInToBuyerPromise       bool     `json:"is_opted_in_to_buyer_promise"`
	IsShopUSBased                 bool     `json:"is_shop_us_based"`
	TransactionSoldCount          int      `json:"transaction_sold_count"`
	ShippingFromCountryISO        *string  `json:"shipping_from_country_iso"`
	ShopLocationCountryISO        *string  `json:"shop_location_country_iso"`
	ReviewCount                   *int     `json:"review_count"`
	ReviewAverage                 *float32 `json:"review_average"`
}

type ShopSection struct {
	ShopSectionID      int64  `json:"shop_section_id"`
	Title              string `json:"title"`
	Rank               int    `json:"rank"`
	UserID             int64  `json:"user_id"`
	ActiveListingCount int    `json:"active_listing_count"`
}

type ShopReturnPolicy struct {
	ReturnPolicyID   int64 `json:"return_policy_id"`
	ShopID           int64 `json:"shop_id"`
	AcceptsReturns   bool  `json:"accepts_returns"`
	AcceptsExchanges bool  `json:"accepts_exchanges"`
	ReturnDeadline   *int  `json:"return_deadline"`
}
