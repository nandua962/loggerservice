package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"utility/internal/consts"
	"utility/internal/entities"
	"utility/utilities"

	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils"
)

// LanguageRepo represents the repository for language-related operations.
type LanguageRepo struct {
	db *sql.DB
}

// LanguageRepoImply is an interface for the LanguageRepo.
type LanguageRepoImply interface {
	GetLanguages(ctx context.Context, params entities.LangParams, pagination entities.Pagination, Validation entities.Validation, errs map[string]models.ErrorResponse) ([]*entities.Language, int64, error)
	GetLanguageCodeExists(ctx context.Context, code string) (bool, error)
}

// NewLanguageRepo creates a new instance of LanguageRepo.
func NewLanguageRepo(db *sql.DB) LanguageRepoImply {
	return &LanguageRepo{db: db}
}

// GetLanguages retrieves a list of languages based on provided parameters.
func (language *LanguageRepo) GetLanguages(ctx context.Context, params entities.LangParams, pagination entities.Pagination, validation entities.Validation, errs map[string]models.ErrorResponse) ([]*entities.Language, int64, error) {

	var (
		languages    []*entities.Language
		totalRecords int64
		log          = logger.Log().WithContext(ctx)
	)
	getLanguages := `
		SELECT 
			id, 
			name, 
			code,
			is_active,
			COUNT(id) OVER() as total_records
		FROM language
	`

	switch strings.ToLower(params.Status) {
	case "", "active":
		getLanguages += "WHERE is_active=true"
	case "inactive":
		getLanguages += "WHERE is_active=false"
	case "all":
	default:

		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, "status", "invalid")
		if err != nil {
			log.Errorf("[LanguageRepo][GetLanguages] Error while loading service code, Error : %s", err.Error())
		}

		errs["status"] = models.ErrorResponse{
			Code:    code,
			Message: []string{"invalid"},
		}
	}

	if !utils.IsEmpty(params.Params.Name) {
		getLanguages = utilities.Search(getLanguages, params.Params.Name, "name")
	}
	if !utils.IsEmpty(params.Params.Id) {
		getLanguages = utilities.SearchById(getLanguages, params.Params.Id, "id")
	}
	if !utils.IsEmpty(params.Params.Code) {
		getLanguages = utilities.Search(getLanguages, params.Params.Code, "code")
	}
	getLanguages = utilities.GroupBy(getLanguages, "id, name")
	sortQ, err := utilities.OrderBy(
		ctx,
		params.Params.Sort, params.Params.Order,
		consts.SortOptns,
		validation.Endpoint,
		validation.Method,
		errs,
	)
	if err != nil {
		log.Errorf("[LanguageRepo][GetLanguages] Error : %s", err.Error())
		return nil, 0, err
	}
	getLanguages = fmt.Sprintf("%s %s", getLanguages, sortQ)
	getLanguages = fmt.Sprintf("%s %s", getLanguages, utils.CalculateOffset(pagination.Page, pagination.Limit))
	if len(errs) > 0 {
		return nil, 0, nil
	}

	res, err := language.db.QueryContext(ctx, getLanguages)
	if err != nil {
		log.Errorf("[LanguageRepo][GetLanguages], Error : %s", err.Error())
		return nil, 0, err
	}
	defer res.Close()
	for res.Next() {
		var language entities.Language
		err := res.Scan(
			&language.ID,
			&language.Name,
			&language.Code,
			&language.IsActive,
			&totalRecords,
		)

		if err != nil {
			log.Errorf("[LanguageRepo][GetLanguages] Error while Scan, Error : %s", err.Error())
			return nil, 0, err
		}

		languages = append(languages, &language)
	}

	return languages, totalRecords, nil
}

func (language *LanguageRepo) GetLanguageCodeExists(ctx context.Context, code string) (bool, error) {

	var (
		isExists bool
		log      = logger.Log().WithContext(ctx)
	)
	code = strings.ToLower(code)

	query := `
		SELECT EXISTS (
			SELECT 1 FROM language WHERE code= $1 AND is_active = $2 
		)`

	if err := language.db.QueryRowContext(ctx, query, code, true).Scan(&isExists); err != nil {
		log.Errorf("[LanguageRepo][GetLanguageCodeExists] Error : %s", err.Error())
		return false, err
	}

	return isExists, nil

}
