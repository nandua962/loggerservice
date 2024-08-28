package entities

import "gitlab.com/tuneverse/toolkit/models"

// Response represents a standard response structure for API responses.
type Response struct {
	MetaData *models.MetaData `json:"meta_data,omitempty"`
	Data     interface{}      `json:"data,omitempty"`
}

type Result struct {
	Metadata any `json:"metadata"`
	Data     any `json:"records"`
}

type ResponseData struct {
	Data any `json:"records"`
}
