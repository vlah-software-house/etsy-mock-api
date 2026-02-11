package models

type BuyerTaxonomyNode struct {
	ID                   int64               `json:"id"`
	Level                int                 `json:"level"`
	Name                 string              `json:"name"`
	ParentID             *int64              `json:"parent_id"`
	Children             []BuyerTaxonomyNode `json:"children"`
	FullPathTaxonomyIDs  []int               `json:"full_path_taxonomy_ids"`
}

type BuyerTaxonomyNodeProperty struct {
	PropertyID          int64                         `json:"property_id"`
	Name                string                        `json:"name"`
	DisplayName         string                        `json:"display_name"`
	Scales              []BuyerTaxonomyPropertyScale  `json:"scales"`
	IsRequired          bool                          `json:"is_required"`
	SupportsAttributes  bool                          `json:"supports_attributes"`
	SupportsVariations  bool                          `json:"supports_variations"`
	IsMultivalued       bool                          `json:"is_multivalued"`
	MaxValuesAllowed    *int                          `json:"max_values_allowed"`
	PossibleValues      []BuyerTaxonomyPropertyValue  `json:"possible_values"`
	SelectedValues      []BuyerTaxonomyPropertyValue  `json:"selected_values"`
}

type BuyerTaxonomyPropertyScale struct {
	ScaleID     int64  `json:"scale_id"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

type BuyerTaxonomyPropertyValue struct {
	ValueID   int64   `json:"value_id"`
	Name      string  `json:"name"`
	ScaleID   *int64  `json:"scale_id"`
	EqualTo   []int64 `json:"equal_to"`
}
