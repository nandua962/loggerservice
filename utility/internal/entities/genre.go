package entities

import "github.com/google/uuid"

// Genre represents a genre entity
type Genre struct {
	ID   uuid.UUID `json:"id,omitempty"`
	Name string    `json:"name"`
}

type GenreDetails struct {
	Name                    string `json:"name"`
	LanguageLabelIdentifier string `json:"language_label_identifier"`
	IsDeleted               bool   `json:"is_deleted"`
}
