package seed

import (
	"encoding/json"
	"os"
)

type Range struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type SeedConfig struct {
	Shops                     int      `json:"shops"`
	ListingsPerShop           Range    `json:"listings_per_shop"`
	ReviewsPerShop            Range    `json:"reviews_per_shop"`
	ReceiptsPerShop           Range    `json:"receipts_per_shop"`
	IncludeDigitalListings    bool     `json:"include_digital_listings"`
	IncludePersonalizedListings bool   `json:"include_personalized_listings"`
	ListingStates             []string `json:"listing_states"`
}

func DefaultConfig() SeedConfig {
	return SeedConfig{
		Shops:                       5,
		ListingsPerShop:             Range{Min: 8, Max: 25},
		ReviewsPerShop:              Range{Min: 3, Max: 15},
		ReceiptsPerShop:             Range{Min: 2, Max: 10},
		IncludeDigitalListings:      true,
		IncludePersonalizedListings: true,
		ListingStates:               []string{"active", "active", "active", "draft", "sold_out"},
	}
}

func LoadConfig(path string) (SeedConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return SeedConfig{}, err
	}
	cfg := DefaultConfig()
	if err := json.Unmarshal(data, &cfg); err != nil {
		return SeedConfig{}, err
	}
	if cfg.Shops < 1 {
		cfg.Shops = 1
	}
	if cfg.ListingsPerShop.Min < 1 {
		cfg.ListingsPerShop.Min = 1
	}
	if cfg.ListingsPerShop.Max < cfg.ListingsPerShop.Min {
		cfg.ListingsPerShop.Max = cfg.ListingsPerShop.Min
	}
	if len(cfg.ListingStates) == 0 {
		cfg.ListingStates = []string{"active"}
	}
	return cfg, nil
}
