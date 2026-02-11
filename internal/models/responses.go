package models

// Standard paginated response wrapper used across Etsy API
type PaginatedResponse struct {
	Count   int         `json:"count"`
	Results interface{} `json:"results"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
