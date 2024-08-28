package entities

// Common struct used for Country, State and Currency entities
type GeographicInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	ISO  string `json:"iso"`
}

// Validation represents a structure for storing validation-related information.
type Validation struct {
	ID           string
	Endpoint     string
	Method       string
	Field        string
	Key          string
	HelpLink     string
	ContextError map[string]any
}
