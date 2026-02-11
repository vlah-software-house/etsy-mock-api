package models

type ListingReview struct {
	ShopID            int64   `json:"shop_id"`
	ListingID         int64   `json:"listing_id"`
	TransactionID     int64   `json:"transaction_id"`
	BuyerUserID       *int64  `json:"buyer_user_id"`
	Rating            int     `json:"rating"`
	Review            string  `json:"review"`
	Language          string  `json:"language"`
	ImageURLFullxfull *string `json:"image_url_fullxfull"`
	CreateTimestamp   int64   `json:"create_timestamp"`
	CreatedTimestamp   int64   `json:"created_timestamp"`
	UpdateTimestamp    int64   `json:"update_timestamp"`
	UpdatedTimestamp   int64   `json:"updated_timestamp"`
}
