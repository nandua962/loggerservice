package entities

// Role represents a role entity
type Role struct {
	ID                      string `json:"id,omitempty"`
	Name                    string `json:"name"`
	CustomName              string `json:"custom_name"`
	LanguageLabelIdentifier string `json:"language_label_identifier"`
	IsDefault               bool   `json:"default"`
}
