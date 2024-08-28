package usecases

import (
	"context"
	"strconv"
	"strings"
	"utility/internal/consts"
	"utility/internal/entities"
	"utility/internal/repo"
	"utility/utilities"

	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils"
)

// CountryUseCases represents the use cases for country-related operations.
type CountryUseCases struct {
	repo repo.CountryRepoImply
}

// CountryUseCaseImply is an interface for the CountryUseCases.
type CountryUseCaseImply interface {
	GetCountries(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errMap map[string]models.ErrorResponse) (*entities.Response, map[string]models.ErrorResponse, error)
	GetStatesOfCountry(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errMap map[string]models.ErrorResponse) (*entities.Response, map[string]models.ErrorResponse, error)
	CheckCountryExists(ctx context.Context, params entities.IsoParam, validation entities.Validation, errMap map[string]models.ErrorResponse) (entities.CountryExists, map[string]models.ErrorResponse, error)
	GetAllCountryCodes(ctx context.Context) (entities.IsoList, error)
	CheckStateExists(ctx context.Context, iso string, countryCode string, validation entities.Validation, errMap map[string]models.ErrorResponse) (bool, map[string]models.ErrorResponse, error)
}

// NewCountryUseCases creates a new instance of CountryUseCases.
func NewCountryUseCases(CountryRepo repo.CountryRepoImply) CountryUseCaseImply {
	return &CountryUseCases{
		repo: CountryRepo,
	}
}

func (country *CountryUseCases) CheckStateExists(ctx context.Context, iso string, countryCode string, validation entities.Validation, errMap map[string]models.ErrorResponse) (bool, map[string]models.ErrorResponse, error) {

	var (
		isExists bool
		log      = logger.Log().WithContext(ctx)
	)

	//Validatecode function is used to validate the ISO param
	validation.ID = iso
	errMap = utilities.Validatecode(ctx, validation, consts.StateIsoLength, consts.IsoField, errMap)

	validation.ID = countryCode
	errMap = utilities.Validatecode(ctx, validation, consts.CodeLength, consts.CodeField, errMap)

	if len(errMap) > 0 {
		log.Errorf("[CountryUseCases][CheckStateExists] ValidationError")
		return false, errMap, nil
	}

	isExists, err := country.repo.CheckStateExists(ctx, iso, countryCode)

	if err != nil {
		log.Errorf("[CountryUseCases][CheckStateExists], Error : %s", err.Error())
		return isExists, nil, err
	}

	return isExists, nil, nil
}

// GetCountries retrieves a list of countries based on provided parameters.
func (country *CountryUseCases) GetCountries(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errMap map[string]models.ErrorResponse) (*entities.Response, map[string]models.ErrorResponse, error) {

	log := logger.Log().WithContext(ctx)

	if params.Iso != "" {
		//Validatecode function is used to validate the ISO param
		validation.ID = params.Iso
		errMap = utilities.Validatecode(ctx, validation, consts.StateIsoLength, consts.IsoField, errMap)

		if len(errMap) > 0 {
			log.Errorf("[CountryUseCases][CheckStateExists] ValidationError")
			return nil, errMap, nil
		}
	}

	//get data
	countries, totalRecords, err := country.repo.GetCountries(ctx, params, pagination, validation, errMap)

	if err != nil {
		log.Errorf("[CountryUseCases][GetCountries], Error : %s", err.Error())
		return nil, nil, err
	}

	if len(errMap) > 0 {
		log.Errorf("[CountryUseCases][GetCountries] ValidationError")
		return nil, errMap, nil
	}

	metaData := utils.MetaDataInfo(&models.MetaData{
		Total:       totalRecords,
		PerPage:     pagination.Limit,
		CurrentPage: pagination.Page,
	})

	resp := &entities.Response{
		MetaData: metaData,
		Data:     countries,
	}
	return resp, nil, nil
}

// GetStatesOfCountry retrieves a list of states for a specific country based on provided parameters.
func (country *CountryUseCases) GetStatesOfCountry(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errMap map[string]models.ErrorResponse) (*entities.Response, map[string]models.ErrorResponse, error) {

	log := logger.Log().WithContext(ctx)
	ID := validation.ID
	id, err := strconv.ParseInt(ID, consts.DecimalBase, consts.BitSize64)
	if err != nil || id <= 0 {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.Country, consts.InvalidKey)

		if err != nil {
			log.Errorf("[CountryUseCases][GetStatesOfCountry], Error :  %s", err)
		}
		errMap[consts.Country] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}
		return nil, errMap, nil
	}

	states, totalRecords, err := country.repo.GetStatesOfCountry(ctx, params, pagination, errMap, id, validation)

	if err != nil {
		log.Errorf("[CountryUseCases][GetStatesOfCountry], Error :%s", err.Error())
		return nil, nil, err
	}

	if len(errMap) > 0 {
		log.Errorf("[CountryUseCases][GetStatesOfCountry] ValidationError")
		return nil, errMap, nil
	}

	metaData := utils.MetaDataInfo(&models.MetaData{
		Total:       totalRecords,
		PerPage:     pagination.Limit,
		CurrentPage: pagination.Page,
	})

	resp := &entities.Response{
		MetaData: metaData,
		Data:     states,
	}
	return resp, nil, nil
}

// CheckCountryExists checks the existence of countries based on ISO parameters.
func (country *CountryUseCases) CheckCountryExists(ctx context.Context, params entities.IsoParam, validation entities.Validation, errMap map[string]models.ErrorResponse) (entities.CountryExists, map[string]models.ErrorResponse, error) {

	log := logger.Log().WithContext(ctx)

	//IsoValidation function is used to validate the ISO param
	errMap = IsoValidation(ctx, params, validation, errMap)

	if len(errMap) > 0 {
		log.Errorf("[CountryUseCases][CheckCountryExists] ValidationError")
		return entities.CountryExists{}, errMap, nil
	}

	result, err := country.repo.CheckCountryExists(ctx, params)

	if err != nil {
		log.Errorf("[CountryUseCases][CheckCountryExists], Error :%s", err.Error())
		return entities.CountryExists{}, nil, err
	}

	return result, nil, nil
}

// IsoValidation validates ISO parameters.
func IsoValidation(ctx context.Context, params entities.IsoParam, validation entities.Validation, errs map[string]models.ErrorResponse) map[string]models.ErrorResponse {

	isoCodes := params.Iso
	var log = logger.Log().WithContext(ctx)
	if utils.IsEmpty(isoCodes) {

		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.IsoField, consts.Required)

		if err != nil {
			log.Errorf("[CountryUseCases][IsoValidation], Error : %s", err)
		}
		errs[consts.IsoField] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.Required},
		}

	} else {
		iso := strings.Split(isoCodes, ",")

		for _, value := range iso {

			if len(value) != 2 {

				code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.IsoField, consts.LengthKey)

				if err != nil {
					log.Errorf("[CountryUseCases][IsoValidation], Error : %s", err)

				}
				errs[consts.IsoField] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.LengthKey},
				}

			}

			isISOValid := utilities.IsValidValue(value)

			if !isISOValid {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.IsoField, consts.InvalidKey)

				if err != nil {
					log.Errorf("[CountryUseCases][IsValidValue], Error : %s", err)
				}
				errs[consts.IsoField] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.InvalidKey},
				}
			}

		}
	}

	return errs
}

// GetAllCountryCodes List all country codes.
func (country *CountryUseCases) GetAllCountryCodes(ctx context.Context) (entities.IsoList, error) {

	log := logger.Log().WithContext(ctx)

	result, err := country.repo.GetAllCountryCodes(ctx)

	if err != nil {
		log.Errorf("[CountryUseCases][GetAllCountryCodes], Error :%s", err.Error())
		return entities.IsoList{}, err
	}

	return result, nil
}
