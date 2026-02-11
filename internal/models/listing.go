package models

type ShopListing struct {
	ListingID                  int64    `json:"listing_id"`
	UserID                     int64    `json:"user_id"`
	ShopID                     int64    `json:"shop_id"`
	Title                      string   `json:"title"`
	Description                string   `json:"description"`
	State                      string   `json:"state"`
	CreationTimestamp           int64    `json:"creation_timestamp"`
	CreatedTimestamp            int64    `json:"created_timestamp"`
	EndingTimestamp             int64    `json:"ending_timestamp"`
	OriginalCreationTimestamp   int64    `json:"original_creation_timestamp"`
	LastModifiedTimestamp       int64    `json:"last_modified_timestamp"`
	UpdatedTimestamp            int64    `json:"updated_timestamp"`
	StateTimestamp              *int64   `json:"state_timestamp"`
	Quantity                   int      `json:"quantity"`
	ShopSectionID              *int64   `json:"shop_section_id"`
	FeaturedRank               int      `json:"featured_rank"`
	URL                        string   `json:"url"`
	NumFavorers                int      `json:"num_favorers"`
	NonTaxable                 bool     `json:"non_taxable"`
	IsTaxable                  bool     `json:"is_taxable"`
	IsCustomizable             bool     `json:"is_customizable"`
	IsPersonalizable           bool     `json:"is_personalizable"`
	PersonalizationIsRequired  bool     `json:"personalization_is_required"`
	PersonalizationCharCountMax *int    `json:"personalization_char_count_max"`
	PersonalizationInstructions *string `json:"personalization_instructions"`
	ListingType                string   `json:"listing_type"`
	Tags                       []string `json:"tags"`
	Materials                  []string `json:"materials"`
	ShippingProfileID          *int64   `json:"shipping_profile_id"`
	ReturnPolicyID             *int64   `json:"return_policy_id"`
	ProcessingMin              *int     `json:"processing_min"`
	ProcessingMax              *int     `json:"processing_max"`
	WhoMade                    *string  `json:"who_made"`
	WhenMade                   *string  `json:"when_made"`
	IsSupply                   *bool    `json:"is_supply"`
	ItemWeight                 *float32 `json:"item_weight"`
	ItemWeightUnit             *string  `json:"item_weight_unit"`
	ItemLength                 *float32 `json:"item_length"`
	ItemWidth                  *float32 `json:"item_width"`
	ItemHeight                 *float32 `json:"item_height"`
	ItemDimensionsUnit         *string  `json:"item_dimensions_unit"`
	IsPrivate                  bool     `json:"is_private"`
	Style                      []string `json:"style"`
	FileData                   *string  `json:"file_data"`
	HasVariations              bool     `json:"has_variations"`
	ShouldAutoRenew            bool     `json:"should_auto_renew"`
	Language                   *string  `json:"language"`
	Price                      Money    `json:"price"`
	TaxonomyID                 *int     `json:"taxonomy_id"`
	Views                      int      `json:"views"`
	// Associations (optional, included when requested)
	Images    []ListingImage     `json:"images,omitempty"`
	Videos    []ListingVideo     `json:"videos,omitempty"`
	Inventory *ListingInventory  `json:"inventory,omitempty"`
	SKUs      []string           `json:"skus,omitempty"`
}

type ListingImage struct {
	ListingID        int64   `json:"listing_id"`
	ListingImageID   int64   `json:"listing_image_id"`
	HexCode          *string `json:"hex_code"`
	Red              *int    `json:"red"`
	Green            *int    `json:"green"`
	Blue             *int    `json:"blue"`
	Hue              *int    `json:"hue"`
	Saturation       *int    `json:"saturation"`
	Brightness       *int    `json:"brightness"`
	IsBlackAndWhite  *bool   `json:"is_black_and_white"`
	CreationTsz      int64   `json:"creation_tsz"`
	CreatedTimestamp  int64   `json:"created_timestamp"`
	Rank             int     `json:"rank"`
	URL75x75         string  `json:"url_75x75"`
	URL170x135       string  `json:"url_170x135"`
	URL570xN         string  `json:"url_570xN"`
	URLFullxfull     string  `json:"url_fullxfull"`
	FullHeight       *int    `json:"full_height"`
	FullWidth        *int    `json:"full_width"`
	AltText          *string `json:"alt_text"`
}

type ListingVideo struct {
	VideoID      int64  `json:"video_id"`
	Height       int    `json:"height"`
	Width        int    `json:"width"`
	ThumbnailURL string `json:"thumbnail_url"`
	VideoURL     string `json:"video_url"`
	VideoState   string `json:"video_state"`
}

type ListingFile struct {
	ListingFileID   int64  `json:"listing_file_id"`
	ListingID       int64  `json:"listing_id"`
	Rank            int    `json:"rank"`
	Filename        string `json:"filename"`
	Filesize        string `json:"filesize"`
	SizeBytes       int    `json:"size_bytes"`
	Filetype        string `json:"filetype"`
	CreateTimestamp int64  `json:"create_timestamp"`
	CreatedTimestamp int64 `json:"created_timestamp"`
}

type ListingInventory struct {
	Products             []ListingInventoryProduct `json:"products"`
	PriceOnProperty      []int                     `json:"price_on_property"`
	QuantityOnProperty   []int                     `json:"quantity_on_property"`
	SKUOnProperty        []int                     `json:"sku_on_property"`
}

type ListingInventoryProduct struct {
	ProductID      int64                          `json:"product_id"`
	SKU            string                         `json:"sku"`
	IsDeleted      bool                           `json:"is_deleted"`
	Offerings      []ListingInventoryProductOffering `json:"offerings"`
	PropertyValues []ListingPropertyValue          `json:"property_values"`
}

type ListingInventoryProductOffering struct {
	OfferingID int64 `json:"offering_id"`
	Quantity   int   `json:"quantity"`
	IsEnabled  bool  `json:"is_enabled"`
	IsDeleted  bool  `json:"is_deleted"`
	Price      Money `json:"price"`
}

type ListingPropertyValue struct {
	PropertyID   int64    `json:"property_id"`
	PropertyName *string  `json:"property_name"`
	ScaleID      *int64   `json:"scale_id"`
	ScaleName    *string  `json:"scale_name"`
	ValueIDs     []int    `json:"value_ids"`
	Values       []string `json:"values"`
}

type CreateListingRequest struct {
	Quantity           int      `json:"quantity"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Price              float64  `json:"price"`
	WhoMade            string   `json:"who_made"`
	WhenMade           string   `json:"when_made"`
	TaxonomyID         int      `json:"taxonomy_id"`
	ShippingProfileID  *int64   `json:"shipping_profile_id"`
	ReturnPolicyID     *int64   `json:"return_policy_id"`
	Materials          []string `json:"materials"`
	Tags               []string `json:"tags"`
	Styles             []string `json:"styles"`
	ListingType        string   `json:"listing_type"`
	IsSupply           *bool    `json:"is_supply"`
	IsCustomizable     bool     `json:"is_customizable"`
	IsPersonalizable   bool     `json:"is_personalizable"`
	ItemWeight         *float32 `json:"item_weight"`
	ItemWeightUnit     *string  `json:"item_weight_unit"`
	ItemLength         *float32 `json:"item_length"`
	ItemWidth          *float32 `json:"item_width"`
	ItemHeight         *float32 `json:"item_height"`
	ItemDimensionsUnit *string  `json:"item_dimensions_unit"`
}

type UpdateListingRequest struct {
	Quantity           *int      `json:"quantity"`
	Title              *string   `json:"title"`
	Description        *string   `json:"description"`
	Price              *float64  `json:"price"`
	WhoMade            *string   `json:"who_made"`
	WhenMade           *string   `json:"when_made"`
	TaxonomyID         *int      `json:"taxonomy_id"`
	ShippingProfileID  *int64    `json:"shipping_profile_id"`
	ReturnPolicyID     *int64    `json:"return_policy_id"`
	Materials          []string  `json:"materials"`
	Tags               []string  `json:"tags"`
	State              *string   `json:"state"`
	IsSupply           *bool     `json:"is_supply"`
	ItemWeight         *float32  `json:"item_weight"`
	ItemWeightUnit     *string   `json:"item_weight_unit"`
	ItemLength         *float32  `json:"item_length"`
	ItemWidth          *float32  `json:"item_width"`
	ItemHeight         *float32  `json:"item_height"`
	ItemDimensionsUnit *string   `json:"item_dimensions_unit"`
}
