package entities

// Language represents a Language entity in the database.
type Language struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	IsActive bool   `json:"is_active"`
}

type LanguageCode struct {
	Code string `json:"code"`
}
