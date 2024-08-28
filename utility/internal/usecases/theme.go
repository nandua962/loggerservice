package usecases

import (
	"context"
	"utility/internal/consts"
	"utility/internal/entities"
	"utility/internal/repo"
	"utility/utilities"

	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils"
)

// ThemeUseCases represents the use cases for Theme-related operations.
type ThemeUseCases struct {
	repo repo.ThemeRepoImply
}

// ThemeUseCaseImply is an interface for the ThemeUseCases.
type ThemeUseCaseImply interface {
	GetThemeByID(ctx context.Context, validation entities.Validation, errMap map[string]models.ErrorResponse) (entities.Theme, map[string]models.ErrorResponse, error)
}

// NewThemeUseCases creates a new instance of ThemeUseCases.
func NewThemeUseCases(ThemeRepo repo.ThemeRepoImply) ThemeUseCaseImply {
	return &ThemeUseCases{
		repo: ThemeRepo,
	}
}

func (theme *ThemeUseCases) GetThemeByID(ctx context.Context, validation entities.Validation, errMap map[string]models.ErrorResponse) (entities.Theme, map[string]models.ErrorResponse, error) {

	log := logger.Log().WithContext(ctx)

	isValidID := utilities.IsValidID(validation.ID)

	if !isValidID {

		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.ThemeIdField, "invalid")

		if err != nil {
			log.Errorf("[CountryUseCases][IsValidValue] Error while loading service code,  Error : %s", err.Error())
		}
		errMap[consts.ThemeIdField] = models.ErrorResponse{
			Code:    code,
			Message: []string{"invalid"},
		}

	}

	if len(errMap) > 0 {
		log.Errorf("[CountryUseCases][IsValidValue] ValidationError")
		return entities.Theme{}, errMap, nil
	}

	result, err := theme.repo.GetThemeByID(ctx, validation.ID)

	if err != nil {
		log.Errorf("[ThemeUseCases][GetThemeByID] Error : %s", err.Error())
		return entities.Theme{}, nil, err
	}

	return result, nil, nil

}
