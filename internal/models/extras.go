package models

// ListingPersonalization represents personalization settings for a listing.
type ListingPersonalization struct {
	IsPersonalizable          bool    `json:"is_personalizable"`
	PersonalizationIsRequired bool    `json:"personalization_is_required"`
	PersonalizationCharCountMax *int  `json:"personalization_char_count_max"`
	PersonalizationInstructions *string `json:"personalization_instructions"`
}

// ListingTranslation represents a translated version of a listing.
type ListingTranslation struct {
	ListingID   int64   `json:"listing_id"`
	Language    string  `json:"language"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Tags        []string `json:"tags"`
}

// ListingVariationImage maps a variation property value to an image.
type ListingVariationImage struct {
	PropertyID   int64 `json:"property_id"`
	ValueID      int64 `json:"value_id"`
	ImageID      int64 `json:"image_id"`
}

// ShopProductionPartner represents a production partner for a shop.
type ShopProductionPartner struct {
	ProductionPartnerID int64  `json:"production_partner_id"`
	PartnerName         string `json:"partner_name"`
	Location            string `json:"location"`
}

// ShopHolidayPreferences represents a shop's holiday/vacation schedule.
type ShopHolidayPreferences struct {
	ShopID         int64   `json:"shop_id"`
	IsVacation     bool    `json:"is_vacation"`
	VacationStart  *int64  `json:"vacation_start"`
	VacationEnd    *int64  `json:"vacation_end"`
	VacationMessage *string `json:"vacation_message"`
}

// ReadinessStateDefinition represents a processing/readiness state.
type ReadinessStateDefinition struct {
	ReadinessStateID    int64  `json:"readiness_state_id"`
	Name                string `json:"name"`
	Label               string `json:"label"`
	ProcessingTimeUnit  string `json:"processing_time_unit"`
	ProcessingMin       int    `json:"processing_min"`
	ProcessingMax       int    `json:"processing_max"`
}

// ShippingCarrier represents a supported shipping carrier.
type ShippingCarrier struct {
	ShippingCarrierID     int    `json:"shipping_carrier_id"`
	Name                  string `json:"name"`
	DomesticClasses       []ShippingCarrierMailClass `json:"domestic_classes"`
	InternationalClasses  []ShippingCarrierMailClass `json:"international_classes"`
}

type ShippingCarrierMailClass struct {
	MailClassKey string `json:"mail_class_key"`
	Name         string `json:"name"`
}

// ReceiptShipmentTracking represents tracking info for a receipt.
type ReceiptShipmentTracking struct {
	CarrierName  string `json:"carrier_name"`
	TrackingCode string `json:"tracking_code"`
}
