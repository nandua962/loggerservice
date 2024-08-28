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

// LanguageUseCases represents the use cases for language-related operations.
type LanguageUseCases struct {
	repo repo.LanguageRepoImply
}

// LanguageUseCaseImply is an interface for the LanguageUseCases.
type LanguageUseCaseImply interface {
	GetLanguages(ctx context.Context, params entities.LangParams, pagination entities.Pagination, validation entities.Validation, errMap map[string]models.ErrorResponse) (*entities.Response, map[string]models.ErrorResponse, error)
	GetLanguageCodeExists(ctx context.Context, validation entities.Validation, errMap map[string]models.ErrorResponse) (bool, map[string]models.ErrorResponse, error)
}

// NewLanguageUseCases creates a new instance of LanguageUseCases.
func NewLanguageUseCases(LanguageRepo repo.LanguageRepoImply) LanguageUseCaseImply {
	return &LanguageUseCases{
		repo: LanguageRepo,
	}
}

// GetLanguages retrieves a list of languages based on provided parameters.
func (language *LanguageUseCases) GetLanguages(ctx context.Context, params entities.LangParams, pagination entities.Pagination, validation entities.Validation, errMap map[string]models.ErrorResponse) (*entities.Response, map[string]models.ErrorResponse, error) {

	log := logger.Log().WithContext(ctx)

	//Validatecode function is used to validate the code param
	if !utils.IsEmpty(params.Params.Code) {
		validation.ID = params.Params.Code
		errMap = utilities.Validatecode(ctx, validation, consts.CodeLength, consts.CodeField, errMap)
		if len(errMap) > 0 {
			log.Errorf("[LanguageUseCases][GetLanguages] ValidationError")
			return nil, errMap, nil
		}
	}

	languages, totalRecords, err := language.repo.GetLanguages(ctx, params, pagination, validation, errMap)
	if len(errMap) > 0 {
		log.Errorf("[LanguageUseCases][GetLanguages] ValidationError")
		return nil, errMap, nil
	}
	if err != nil {
		log.Errorf("[LanguageUseCases][GetLanguages], Error : %s", err.Error())
		return nil, nil, err
	}

	metaData := utils.MetaDataInfo(&models.MetaData{
		Total:       totalRecords,
		PerPage:     pagination.Limit,
		CurrentPage: pagination.Page,
	})

	resp := &entities.Response{
		MetaData: metaData,
		Data:     languages,
	}
	return resp, nil, nil
}

// GetLanguageCodeExists to check language code is exists or not
func (language *LanguageUseCases) GetLanguageCodeExists(ctx context.Context, validation entities.Validation, errMap map[string]models.ErrorResponse) (bool, map[string]models.ErrorResponse, error) {

	var (
		isExists bool
		log      = logger.Log().WithContext(ctx)
	)

	//Validatecode function is used to validate the language code param
	errMap = utilities.Validatecode(ctx, validation, consts.CodeLength, consts.CodeField, errMap)

	if len(errMap) > 0 {
		log.Errorf("[LanguageUseCases][GetLanguageCodeExists] ValidationError")
		return false, errMap, nil
	}

	isExists, err := language.repo.GetLanguageCodeExists(ctx, validation.ID)

	if err != nil {
		log.Errorf("[LanguageUseCases][GetLanguageCodeExists], Error :%s", err.Error())
		return isExists, nil, err
	}

	return isExists, nil, nil

}
