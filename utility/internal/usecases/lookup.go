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

// LookupUseCases represents use cases for handling lookup-related operations.
type LookupUseCases struct {
	repo repo.LookupRepoImply
}

// LookupUseCaseImply is an interface defining the methods for working with lookup use cases.
type LookupUseCaseImply interface {
	GetLookupByIdList(ctx context.Context, idList entities.LookupIDs) (entities.LookupData, error)
	GetLookupByTypeName(ctx context.Context, filter map[string]string, validation entities.Validation, errMap map[string]models.ErrorResponse) ([]entities.Lookup, map[string]models.ErrorResponse, error)
}

// NewLookupUseCases creates a new LookupUseCases instance.
func NewLookupUseCases(LookupRepo repo.LookupRepoImply) LookupUseCaseImply {
	return &LookupUseCases{
		repo: LookupRepo,
	}
}

// GetLookupByTypeName retrieves lookup data based on specified parameters.
func (lookup *LookupUseCases) GetLookupByTypeName(ctx context.Context, filter map[string]string, validation entities.Validation, errMap map[string]models.ErrorResponse) ([]entities.Lookup, map[string]models.ErrorResponse, error) {

	log := logger.Log().WithContext(ctx)

	//IsValidValueWithUnderscore checks whether the id is valid or not

	isTypeNameValid := utilities.IsValidValueWithUnderscore(validation.ID)

	if !isTypeNameValid {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.NameKey, consts.InvalidKey)

		if err != nil {
			log.Errorf("[LookupUseCases][GetLookupByTypeName] Error while loading service code %s", err)
		}
		errMap[consts.NameKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}
	}

	if len(errMap) > 0 {
		log.Errorf("[LookupUseCases][GetLookupByTypeName] ValidationError")
		return nil, errMap, nil
	}

	result, err := lookup.repo.GetLookupByTypeName(ctx, validation.ID, filter)

	if err != nil {
		log.Errorf("[LookupUseCases][GetLookupByTypeName] Error : %s", err.Error())
		return nil, nil, err
	}

	return result, nil, nil
}

// GetLookup retrieves lookup data based by array of lookup Ids
func (lookup *LookupUseCases) GetLookupByIdList(ctx context.Context, idList entities.LookupIDs) (entities.LookupData, error) {

	log := logger.Log().WithContext(ctx)

	result, err := lookup.repo.GetLookupByIdList(ctx, idList)

	if err != nil {
		log.Errorf("[LookupUseCases][GetLookupByIdList] Error : %s", err.Error())
		return entities.LookupData{}, err
	}

	return result, nil
}
