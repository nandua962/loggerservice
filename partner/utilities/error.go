package utilities

import (
	"gitlab.com/tuneverse/toolkit/models"
)

type errorMap map[string]models.ErrorResponse

func NewErrorMap() errorMap {
	return make(errorMap)
}
