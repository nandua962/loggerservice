package entities

// Theme represents a theme entity in the database.
type Theme struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Value    string `json:"code"`
	LayoutID int64  `json:"layout_id"`
}
