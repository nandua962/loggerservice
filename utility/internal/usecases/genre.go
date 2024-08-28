package usecases

import (
	"context"
	"utility/internal/consts"
	"utility/internal/entities"
	"utility/internal/repo"
	"utility/utilities"

	"github.com/google/uuid"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils"
)

// GenreUseCases represents use cases for handling genre-related operations.
type GenreUseCases struct {
	repo repo.GenreRepoImply
}

// GenreUseCaseImply is an interface defining the methods for working with genre use cases.
type GenreUseCaseImply interface {
	CreateGenre(ctx context.Context, req entities.Genre, validation entities.Validation, errMap map[string]models.ErrorResponse) (map[string]models.ErrorResponse, error)
	GetGenres(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errMap map[string]models.ErrorResponse) (*entities.Response, map[string]models.ErrorResponse, error)
	DeleteGenre(ctx context.Context, validation entities.Validation, errMap map[string]models.ErrorResponse) (map[string]models.ErrorResponse, error)
	UpdateGenre(ctx context.Context, genreInfo entities.Genre, validation entities.Validation, errMap map[string]models.ErrorResponse) (map[string]models.ErrorResponse, error)
	GetGenresByID(ctx context.Context, validation entities.Validation, errMap map[string]models.ErrorResponse) (entities.GenreDetails, map[string]models.ErrorResponse, error)
}

// NewGenreUseCases creates a new GenreUseCases instance.
func NewGenreUseCases(GenreRepo repo.GenreRepoImply) GenreUseCaseImply {
	return &GenreUseCases{
		repo: GenreRepo,
	}
}

// CreateGenre handles the creation of a genre and does necessary validations.
func (genre *GenreUseCases) CreateGenre(ctx context.Context, req entities.Genre, validation entities.Validation, errMap map[string]models.ErrorResponse) (map[string]models.ErrorResponse, error) {

	log := logger.Log().WithContext(ctx)

	if utils.IsEmpty(req.Name) {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.NameKey, consts.Required)
		if err != nil {
			log.Errorf("[GenreUseCases][CreateGenre],Error while loading service code, Error : %s", err.Error())

		}

		errMap[consts.NameKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.Required},
		}

	} else if err := utils.ValidateStringLength(req.Name, consts.MinimumLength, consts.MaximumLength); err != nil {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.NameKey, consts.LengthKey)
		if err != nil {
			log.Errorf("[GenreUseCases][CreateGenre][ValidateStringLength],Error while loading service code, Error : %s", err.Error())
		}

		errMap[consts.NameKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.LengthKey},
		}

		return errMap, nil
	} else {
		isExist, err := genre.repo.GenreNameExists(ctx, req.Name)

		if err != nil {
			log.Errorf("[GenreUseCases][CreateGenre][GenreNameExists], Error : %s", err.Error())
			return nil, err
		}
		if isExist {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.NameKey, consts.GenreExists)
			if err != nil {
				log.Errorf("[GenreUseCases][CreateGenre][GenreNameExists],Error while loading service code, Error : %s", err.Error())
			}

			errMap[consts.NameKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.GenreExists},
			}

		}
	}

	if len(errMap) > 0 {
		log.Errorf("[GenreUseCases][CreateGenre], ValidationError")
		return errMap, nil
	}

	langIdentifier, err := utilities.GenerateLangIdentifier(req.Name, consts.Genre, consts.LabelIdentifier)

	if err != nil {
		log.Errorf("[GenreUseCases][CreateGenre][GenerateLangIdentifier], Error : %s", err.Error())
		return nil, err
	}

	err = genre.repo.CreateGenre(ctx, req.Name, langIdentifier)
	if err != nil {
		log.Errorf("[GenreUseCases][CreateGenre], Error : %s", err.Error())
		return nil, err
	}
	return nil, nil
}

// GetGenres retrieves genre data based on specified parameters.
func (genre *GenreUseCases) GetGenres(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errMap map[string]models.ErrorResponse) (*entities.Response, map[string]models.ErrorResponse, error) {

	log := logger.Log().WithContext(ctx)

	genres, totalRecords, err := genre.repo.GetGenres(ctx, params, pagination, validation, errMap)
	if len(errMap) > 0 {
		log.Errorf("[GenreUseCases][CreateGenre], ValidationError")
		return nil, errMap, nil
	}

	if err != nil {
		log.Errorf("[GenreUseCases][GetGenres], Error : %s", err.Error())
		return nil, nil, err
	}

	metaData := utils.MetaDataInfo(&models.MetaData{
		Total:       totalRecords,
		PerPage:     pagination.Limit,
		CurrentPage: pagination.Page,
	})

	resp := &entities.Response{
		MetaData: metaData,
		Data:     genres,
	}
	return resp, nil, nil

}

// DeleteGenre handles the deletion of a genre.
func (genre *GenreUseCases) DeleteGenre(ctx context.Context, validation entities.Validation, errMap map[string]models.ErrorResponse) (map[string]models.ErrorResponse, error) {

	log := logger.Log().WithContext(ctx)

	uuid, err := uuid.Parse(validation.ID)

	if err != nil {

		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.CommonMsg, consts.InvalidGenre)
		if err != nil {
			log.Errorf("[GenreUseCases][DeleteGenre],Error while loading service code, Error : %s", err)
		}

		errMap[consts.CommonMsg] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidGenre},
		}

		return errMap, nil
	}

	err = genre.repo.IsGenreExists(ctx, uuid)
	if err != nil {

		if err.Error() == "genre already deleted" {
			log.Errorf("[GenreUseCases][DeleteGenre], Error : %s", err.Error())
			return nil, consts.ErrNotExist
		}
		log.Errorf("[GenreUseCases][DeleteGenre], Error : %s", err.Error())
		return nil, err
	}

	rowsAffected, err := genre.repo.DeleteGenre(ctx, uuid)

	if err != nil {
		log.Errorf("[GenreUseCases][DeleteGenre], Error : %s", err.Error())
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, consts.ErrNotFound
	}
	return nil, nil
}

// UpdateGenre handles the updating of a genre and does necessary validattions for the input data.
func (genre *GenreUseCases) UpdateGenre(ctx context.Context, genreInfo entities.Genre, validation entities.Validation, errMap map[string]models.ErrorResponse) (map[string]models.ErrorResponse, error) {

	log := logger.Log().WithContext(ctx)
	uuid, err := uuid.Parse(validation.ID)

	if err != nil {

		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.CommonMsg, consts.InvalidGenre)
		if err != nil {
			log.Errorf("[GenreUseCases][UpdateGenre],Error while loading service code, Error : %s", err.Error())
		}

		errMap[consts.CommonMsg] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidGenre},
		}

		return errMap, nil
	}

	err = genre.repo.IsGenreExists(ctx, uuid)
	if err != nil {

		if err.Error() == "genre already deleted" {
			log.Errorf("[GenreUseCases][UpdateGenre] UpdateGenre failed, validation error, err = genre already inactive")
			return nil, consts.ErrNotExist
		}
		log.Errorf("[GenreUseCases][UpdateGenre][IsGenreExists], Error : %s", err.Error())
		return nil, err
	}

	if utils.IsEmpty(genreInfo.Name) {

		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.NameKey, consts.Required)
		if err != nil {
			log.Errorf("[GenreUseCases][UpdateGenre],Error while loading service code, Error : %s", err)
		}

		errMap[consts.NameKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.Required},
		}

	} else if err := utils.ValidateStringLength(genreInfo.Name, consts.MinimumLength, consts.MaximumLength); err != nil {

		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.NameKey, consts.LengthKey)
		if err != nil {
			log.Errorf("[GenreUseCases][UpdateGenre],Error while loading service code, Error : %s", err.Error())
		}

		errMap[consts.NameKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.LengthKey},
		}

		return errMap, nil
	} else {
		isNameExists, err := genre.repo.IsDuplicateExists(ctx, consts.NameKey, genreInfo.Name, validation.ID)
		if err != nil {
			log.Errorf("[GenreUseCases][UpdateGenre][IsDuplicateExists], Error :  %s", err.Error())
			return nil, err
		}

		if isNameExists {
			log.Errorf("UpdateGenre failed, validation error, Error : duplicate name.")
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.NameKey, consts.GenreExists)
			if err != nil {
				log.Errorf("[GenreUseCases][UpdateGenre], Error while loading service code, Error : %s", err.Error())
			}

			errMap[consts.NameKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.GenreExists},
			}

		}
	}

	if len(errMap) > 0 {
		log.Errorf("[GenreUseCases][UpdateGenre], ValidationError")
		return errMap, nil
	}

	langIdentifier, err := utilities.GenerateLangIdentifier(genreInfo.Name, consts.Genre)
	if err != nil {
		log.Errorf("[GenreUseCases][UpdateGenre][GenerateLangIdentifier], Error : %s", err.Error())
		return nil, err
	}
	genreInfo.ID = uuid
	err = genre.repo.UpdateGenre(ctx, genreInfo, langIdentifier)
	if err != nil {
		log.Errorf("[GenreUseCases][UpdateGenre], Error : %s", err.Error())
		return nil, err
	}

	return nil, nil
}

// GetGenresByID to retreive genre details
func (genre *GenreUseCases) GetGenresByID(ctx context.Context, validation entities.Validation, errMap map[string]models.ErrorResponse) (entities.GenreDetails, map[string]models.ErrorResponse, error) {

	log := logger.Log().WithContext(ctx)

	_, err := uuid.Parse(validation.ID)

	if err != nil {

		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.GenreIdField, consts.InvalidKey)
		if err != nil {
			log.Errorf("[GenreUseCases][GetGenresByID],Error while loading service code, Error : %s", err.Error())
		}

		errMap[consts.GenreIdField] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}

		return entities.GenreDetails{}, errMap, nil
	}

	if len(errMap) > 0 {
		log.Errorf("[GenreUseCases][GetGenresByID], ValidationError")
		return entities.GenreDetails{}, errMap, nil
	}

	result, err := genre.repo.GetGenresByID(ctx, validation.ID)

	if err != nil {
		log.Errorf("[GenreUseCases][GetGenresByID] Error : %s", err.Error())
		return entities.GenreDetails{}, nil, err
	}

	return result, nil, nil

}
