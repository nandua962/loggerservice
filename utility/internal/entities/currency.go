package entities

// Currency represents a currency entity.
type Currency struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	ISO    string `json:"iso,omitempty"`
	Symbol string `json:"symbol,omitempty"`
}
