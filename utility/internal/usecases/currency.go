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

// CurrencyUseCases represents the use case for currency-related operations.
type CurrencyUseCases struct {
	repo repo.CurrencyRepoImply
}

// CurrencyUseCaseImply is an interface for the CurrencyUseCases.
type CurrencyUseCaseImply interface {
	GetCurrencies(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errMap map[string]models.ErrorResponse) (*entities.Response, map[string]models.ErrorResponse, error)
	GetCurrencyByID(ctx context.Context, validation entities.Validation, errMap map[string]models.ErrorResponse) (entities.GeographicInfo, map[string]models.ErrorResponse, error)
	GetCurrencyByISO(ctx context.Context, validation entities.Validation, errMap map[string]models.ErrorResponse) (entities.Currency, map[string]models.ErrorResponse, error)
}

// NewCurrencyUseCases creates a new instance of CurrencyUseCases.
func NewCurrencyUseCases(currencyRepo repo.CurrencyRepoImply) CurrencyUseCaseImply {
	return &CurrencyUseCases{
		repo: currencyRepo,
	}
}

// GetCurrencies retrieves a list of currencies based on provided parameters.
func (currency *CurrencyUseCases) GetCurrencies(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errMap map[string]models.ErrorResponse) (*entities.Response, map[string]models.ErrorResponse, error) {

	log := logger.Log().WithContext(ctx)

	currencies, totalRecords, err := currency.repo.GetCurrencies(ctx, params, pagination, validation, errMap)
	if len(errMap) > 0 {
		log.Errorf("[CurrencyUseCases][GetCurrencies] ValidationError")
		return nil, errMap, nil
	}

	if err != nil {
		log.Errorf("[CurrencyUseCases][GetCurrencies] Error : %s", err.Error())
		return nil, nil, err
	}

	metaData := utils.MetaDataInfo(&models.MetaData{
		Total:       totalRecords,
		PerPage:     pagination.Limit,
		CurrentPage: pagination.Page,
	})

	resp := &entities.Response{
		MetaData: metaData,
		Data:     currencies,
	}
	return resp, nil, nil
}

func (currency *CurrencyUseCases) GetCurrencyByID(ctx context.Context, validation entities.Validation, errMap map[string]models.ErrorResponse) (entities.GeographicInfo, map[string]models.ErrorResponse, error) {

	log := logger.Log().WithContext(ctx)
	var result entities.GeographicInfo

	errMap = utilities.IDValidation(ctx, validation, errMap, consts.CurrencyIdField)

	if len(errMap) > 0 {
		log.Errorf("[CurrencyUseCases][IDValidation], ValidationError")
		return entities.GeographicInfo{}, errMap, nil
	}

	result, err := currency.repo.GetCurrencyByID(ctx, validation.ID)

	if err != nil {
		log.Errorf("[CurrencyUseCases][GetCurrencyByID] Error : %s", err.Error())
		return entities.GeographicInfo{}, nil, err
	}

	return result, nil, nil
}

func (currency *CurrencyUseCases) GetCurrencyByISO(ctx context.Context, validation entities.Validation, errMap map[string]models.ErrorResponse) (entities.Currency, map[string]models.ErrorResponse, error) {

	var (
		result entities.Currency
		log    = logger.Log().WithContext(ctx)
	)

	//Validatecode function is used to validate the ISO param
	errMap = utilities.Validatecode(ctx, validation, consts.IsoLength, consts.IsoField, errMap)

	if len(errMap) > 0 {
		log.Errorf("[CurrencyUseCases][Validatecode], ValidationError")
		return entities.Currency{}, errMap, nil
	}

	result, err := currency.repo.GetCurrencyByISO(ctx, validation.ID)

	if err != nil {
		log.Errorf("[CurrencyUseCases][GetCurrencyByISO], Error : %s", err.Error())
		return entities.Currency{}, nil, err
	}

	return result, nil, nil
}
