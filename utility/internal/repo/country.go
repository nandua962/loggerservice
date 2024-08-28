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

// CountryRepo represents the repository for country-related operations.
type CountryRepo struct {
	db *sql.DB
}

// CountryRepoImply is an interface for the CountryRepo.
type CountryRepoImply interface {
	GetCountries(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errs map[string]models.ErrorResponse) ([]*entities.GeographicInfo, int64, error)
	GetStatesOfCountry(ctx context.Context, params entities.Params, pagination entities.Pagination, errs map[string]models.ErrorResponse, id int64, validation entities.Validation) ([]*entities.GeographicInfo, int64, error)
	CheckCountryExists(ctx context.Context, params entities.IsoParam) (entities.CountryExists, error)
	GetAllCountryCodes(ctx context.Context) (entities.IsoList, error)
	CheckStateExists(ctx context.Context, iso string, countryCode string) (bool, error)
}

// NewCountryRepo creates a new instance of CountryRepo.
func NewCountryRepo(db *sql.DB) CountryRepoImply {
	return &CountryRepo{db: db}
}

func (country *CountryRepo) CheckStateExists(ctx context.Context, iso string, countryCode string) (bool, error) {

	var (
		isExists bool
		log      = logger.Log().WithContext(ctx)
	)

	query := `
		SELECT EXISTS (
			 SELECT 1 
            FROM public.country_state 
            WHERE LOWER(iso) = LOWER($1) AND LOWER(country_code) = LOWER($2)
		)`

	if err := country.db.QueryRowContext(ctx, query, iso, countryCode).Scan(&isExists); err != nil {
		log.Errorf("[CountryRepo][CheckStateExists], Error : %s", err.Error())
		return false, err
	}

	return isExists, nil

}

// GetCountries retrieves a list of countries based on provided parameters.
func (country *CountryRepo) GetCountries(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errs map[string]models.ErrorResponse) ([]*entities.GeographicInfo, int64, error) {

	var (
		countries    []*entities.GeographicInfo
		totalRecords int64
		log          = logger.Log().WithContext(ctx)
	)

	getCountries := `
		SELECT 
			id, 
			name,
			iso,
			SUM(COUNT(id)) OVER() as total_records
		FROM country
	`

	if !utils.IsEmpty(params.Name) {
		getCountries = utilities.Search(getCountries, params.Name, "name")
	}
	if !utils.IsEmpty(params.Iso) {
		getCountries = utilities.SearchIso(getCountries, params.Iso, "iso")
	}

	getCountries = utilities.GroupBy(getCountries, "id, name")
	sortQ, err := utilities.OrderBy(
		ctx,
		params.Sort,
		params.Order,
		consts.SortOptns,
		validation.Endpoint,
		validation.Method,
		errs,
	)
	if err != nil {
		log.Errorf("[CountryRepo][GetCountries] Error : %s", err.Error())
		return nil, 0, err
	}
	if len(errs) > 0 {
		return nil, 0, nil
	}

	getCountries = fmt.Sprintf("%s %s", getCountries, sortQ)
	getCountries = fmt.Sprintf("%s %s", getCountries, utils.CalculateOffset(pagination.Page, pagination.Limit))

	res, err := country.db.QueryContext(ctx, getCountries)
	if err != nil {
		log.Errorf("[CountryRepo][GetCountries] Error : %s", err.Error())
		return nil, 0, err
	}
	defer res.Close()
	for res.Next() {
		var country entities.GeographicInfo
		err := res.Scan(
			&country.ID,
			&country.Name,
			&country.ISO,
			&totalRecords,
		)

		if err != nil {
			log.Errorf("[CountryRepo][GetCountries] Error :%s", err.Error())
			return nil, 0, err
		}

		countries = append(countries, &country)
	}

	return countries, totalRecords, nil
}

// GetStatesOfCountry retrieves a list of states for a specific country based on provided parameters.
func (country *CountryRepo) GetStatesOfCountry(ctx context.Context, params entities.Params, pagination entities.Pagination, errs map[string]models.ErrorResponse, id int64, validation entities.Validation) ([]*entities.GeographicInfo, int64, error) {

	var (
		states       []*entities.GeographicInfo
		totalRecords int64
		log          = logger.Log().WithContext(ctx)
	)

	getStatesForCountry := `
		SELECT
			cs.id,
			cs.name,
			cs.iso,
			SUM(COUNT(cs.id)) OVER() AS total_records
		FROM country_state cs
		JOIN country coun
		ON cs.country_code = coun.iso
		WHERE coun.id = $1
	`

	if !utils.IsEmpty(params.Name) {
		getStatesForCountry = utilities.Search(getStatesForCountry, params.Name, "cs.name")
	}

	getStatesForCountry = utilities.GroupBy(getStatesForCountry, "cs.id")
	sortQ, err := utilities.OrderBy(
		ctx,
		params.Sort,
		params.Order,
		consts.StateSortOpns,
		validation.Endpoint,
		validation.Method,
		errs,
	)
	if err != nil {
		log.Errorf("[CountryRepo][GetStatesOfCountry] Error : %s", err.Error())
		return nil, 0, err
	}
	getStatesForCountry = fmt.Sprintf("%s %s", getStatesForCountry, sortQ)
	getStatesForCountry = fmt.Sprintf("%s %s", getStatesForCountry, utils.CalculateOffset(pagination.Page, pagination.Limit))

	if len(errs) > 0 {
		return nil, 0, nil
	}

	res, err := country.db.QueryContext(ctx, getStatesForCountry, id)
	if err != nil {
		log.Errorf("[CountryRepo][GetStatesOfCountry], Error :%s", err.Error())
		return nil, 0, err
	}
	defer res.Close()
	for res.Next() {
		var state entities.GeographicInfo
		err := res.Scan(
			&state.ID,
			&state.Name,
			&state.ISO,
			&totalRecords,
		)

		if err != nil {
			log.Errorf("[CountryRepo][GetStatesOfCountry], Error :%s", err.Error())
			return nil, 0, err
		}

		states = append(states, &state)
	}

	return states, totalRecords, nil
}

// CheckCountryExists retrieves the existence of countries based on provided ISO parameters.
func (country *CountryRepo) CheckCountryExists(ctx context.Context, params entities.IsoParam) (entities.CountryExists, error) {

	var (
		output   entities.CountryExists
		isExists bool
		log      = logger.Log().WithContext(ctx)

		missingCountryCodes []string
	)

	isoCodes := params.Iso
	iso := strings.Split(isoCodes, ",")

	for _, value := range iso {

		query := `SELECT EXISTS(SELECT 1 FROM country where LOWER(iso)=LOWER($1))`

		if err := country.db.QueryRowContext(ctx, query, value).Scan(&isExists); err != nil {
			log.Errorf("[CountryRepo][CheckCountryExists], Error : %s", err.Error())
			return entities.CountryExists{}, err
		}

		// If country does not exist, add ISO code to missingCountryCodes list
		if !isExists {
			missingCountryCodes = append(missingCountryCodes, value)
		}
	}

	if missingCountryCodes != nil {
		output.Exists = false
		output.MissingCountryCodes = missingCountryCodes
	} else {
		output.Exists = true
	}

	return output, nil
}

// GetAllCountryCodes List all country codes.
func (country *CountryRepo) GetAllCountryCodes(ctx context.Context) (entities.IsoList, error) {

	var (
		countryCodes []string
		log          = logger.Log().WithContext(ctx)
		output       entities.IsoList
	)

	query := `SELECT iso FROM country`

	// Execute query to retrieve country codes
	rows, err := country.db.QueryContext(ctx, query)
	if err != nil {
		log.Errorf("[CountryRepo][GetAllCountryCodes], Error : %s", err.Error())
		return entities.IsoList{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var isoCode string
		if err := rows.Scan(&isoCode); err != nil {
			log.Errorf("[CountryRepo][GetAllCountryCodes] failed to scan row, Error : %s", err.Error())
			return entities.IsoList{}, err
		}
		countryCodes = append(countryCodes, isoCode)
	}

	output.Iso = countryCodes

	return output, nil
}
