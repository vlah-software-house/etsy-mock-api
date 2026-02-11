package models

type ShopShippingProfile struct {
	ShippingProfileID         int64                            `json:"shipping_profile_id"`
	Title                     *string                          `json:"title"`
	UserID                    int64                            `json:"user_id"`
	OriginCountryISO          string                           `json:"origin_country_iso"`
	IsDeleted                 bool                             `json:"is_deleted"`
	ShippingProfileDestinations []ShopShippingProfileDestination `json:"shipping_profile_destinations"`
	ShippingProfileUpgrades   []ShopShippingProfileUpgrade     `json:"shipping_profile_upgrades"`
	OriginPostalCode          *string                          `json:"origin_postal_code"`
	ProfileType               string                           `json:"profile_type"`
	DomesticHandlingFee       float64                          `json:"domestic_handling_fee"`
	InternationalHandlingFee  float64                          `json:"international_handling_fee"`
}

type ShopShippingProfileDestination struct {
	ShippingProfileDestinationID int64  `json:"shipping_profile_destination_id"`
	ShippingProfileID            int64  `json:"shipping_profile_id"`
	OriginCountryISO             string `json:"origin_country_iso"`
	DestinationCountryISO        string `json:"destination_country_iso"`
	DestinationRegion            string `json:"destination_region"`
	PrimaryCost                  Money  `json:"primary_cost"`
	SecondaryCost                Money  `json:"secondary_cost"`
	ShippingCarrierID            *int   `json:"shipping_carrier_id"`
	MailClass                    *string `json:"mail_class"`
	MinDeliveryDays              *int   `json:"min_delivery_days"`
	MaxDeliveryDays              *int   `json:"max_delivery_days"`
}

type ShopShippingProfileUpgrade struct {
	ShippingProfileID int64   `json:"shipping_profile_id"`
	UpgradeID         int64   `json:"upgrade_id"`
	UpgradeName       string  `json:"upgrade_name"`
	Type              int     `json:"type"`
	Rank              int     `json:"rank"`
	Language          string  `json:"language"`
	Price             Money   `json:"price"`
	SecondaryPrice    Money   `json:"secondary_price"`
	ShippingCarrierID *int    `json:"shipping_carrier_id"`
	MailClass         *string `json:"mail_class"`
	MinDeliveryDays   *int    `json:"min_delivery_days"`
	MaxDeliveryDays   *int    `json:"max_delivery_days"`
}
