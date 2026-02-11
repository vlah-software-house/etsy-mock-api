package models

type User struct {
	UserID        int64   `json:"user_id"`
	PrimaryEmail  *string `json:"primary_email"`
	FirstName     *string `json:"first_name"`
	LastName      *string `json:"last_name"`
	ImageURL75x75 *string `json:"image_url_75x75"`
}

type UserAddress struct {
	UserAddressID           int64   `json:"user_address_id"`
	UserID                  int64   `json:"user_id"`
	Name                    string  `json:"name"`
	FirstLine               string  `json:"first_line"`
	SecondLine              *string `json:"second_line"`
	City                    string  `json:"city"`
	State                   *string `json:"state"`
	Zip                     *string `json:"zip"`
	ISOCountryCode          *string `json:"iso_country_code"`
	CountryName             *string `json:"country_name"`
	IsDefaultShippingAddress bool   `json:"is_default_shipping_address"`
}
