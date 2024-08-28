package entities

// Lookup represents a lookup entity
type Lookup struct {
	ID                 int64  `json:"id"`
	Name               string `json:"name,omitempty"`
	Value              string `json:"value,omitempty"`
	Description        string `json:"description,omitempty"`
	Position           int    `json:"position,omitempty"`
	LookupTypeId       int    `json:"lookup_type_id,omitempty"`
	LanguageIdentifier string `json:"language_identifier,omitempty"`
}

// LookupData represents the lookup and Ids
type LookupData struct {
	Lookup           []Lookup `json:"records"`
	InvalidLookupIds []int64  `json:"invalid_lookup_ids,omitempty"`
}

type LookupIDs struct {
	ID []int64 `json:"ids"`
}
