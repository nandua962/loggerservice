package repo

import (
	"context"
	"database/sql"
	"fmt"
	"utility/internal/consts"
	"utility/internal/entities"
	"utility/utilities"

	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils"
)

// CurrencyRepo represents the repository for currency-related operations.
type CurrencyRepo struct {
	db *sql.DB
}

// CurrencyRepoImply is an interface for the CurrencyRepo.
type CurrencyRepoImply interface {
	GetCurrencies(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errs map[string]models.ErrorResponse) ([]entities.Currency, int64, error)
	GetCurrencyByISO(ctx context.Context, iso string) (entities.Currency, error)
	GetCurrencyByID(ctx context.Context, id string) (entities.GeographicInfo, error)
}

// NewCurrencyRepo creates a new instance of CurrencyRepo.
func NewCurrencyRepo(db *sql.DB) CurrencyRepoImply {
	return &CurrencyRepo{db: db}
}

// GetCurrencies retrieves a list of currencies based on provided parameters.
func (currencyRepo *CurrencyRepo) GetCurrencies(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errs map[string]models.ErrorResponse) ([]entities.Currency, int64, error) {
	var (
		currencies   []entities.Currency
		totalRecords int64
		log          = logger.Log().WithContext(ctx)
	)

	getCurrencies := `
		SELECT
			id,
			name,
			iso,
			symbol,
			SUM(COUNT(id)) OVER() as total_records
		FROM currency
	`
	if !utils.IsEmpty(params.Name) {
		getCurrencies = utilities.Search(getCurrencies, params.Name, "name")
	}
	getCurrencies = utilities.GroupBy(getCurrencies, "id, name")
	sortQ, err := utilities.OrderBy(
		ctx,
		params.Sort, params.Order,
		consts.SortOptns,
		validation.Endpoint,
		validation.Method,
		errs,
	)
	if err != nil {
		log.Errorf("[CurrencyRepo][GetCurrencies] Error : %s", err.Error())
		return nil, 0, err
	}
	getCurrencies = fmt.Sprintf("%s %s", getCurrencies, sortQ)
	getCurrencies = fmt.Sprintf("%s %s", getCurrencies, utils.CalculateOffset(pagination.Page, pagination.Limit))
	if len(errs) > 0 {
		return nil, 0, nil
	}

	res, err := currencyRepo.db.QueryContext(ctx, getCurrencies)
	if err != nil {
		log.Errorf("[CurrencyRepo][GetCurrencies], Error : %s", err.Error())
		return nil, 0, err
	}
	defer res.Close()
	for res.Next() {
		var currency entities.Currency
		err := res.Scan(
			&currency.ID,
			&currency.Name,
			&currency.ISO,
			&currency.Symbol,
			&totalRecords,
		)

		if err != nil {
			log.Errorf("[CurrencyRepo][GetCurrencies],error while Scan, Error : %s", err.Error())
			return nil, 0, err
		}
		currencies = append(currencies, currency)
	}

	return currencies, totalRecords, nil
}

func (currencyRepo *CurrencyRepo) GetCurrencyByID(ctx context.Context, id string) (entities.GeographicInfo, error) {

	var (
		res entities.GeographicInfo
		log = logger.Log().WithContext(ctx)
	)

	query := `
		SELECT 
			id,
			name,
			iso
		FROM currency 
		WHERE id= $1
		`

	row := currencyRepo.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&res.ID, &res.Name, &res.ISO)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Errorf("[CurrencyRepo][GetCurrencyByID], No Content, Error : %s", err.Error())
			return entities.GeographicInfo{}, err
		}
		log.Errorf("[CurrencyRepo][GetCurrencyByID], Error : %s", err.Error())
		return entities.GeographicInfo{}, err
	}
	return res, nil

}

func (currencyRepo *CurrencyRepo) GetCurrencyByISO(ctx context.Context, iso string) (entities.Currency, error) {

	var (
		currency entities.Currency
		log      = logger.Log().WithContext(ctx)
	)

	query := `
		SELECT 
			id,
			name 
		FROM currency 
		WHERE LOWER(iso)= LOWER($1)
		`

	row := currencyRepo.db.QueryRowContext(ctx, query, iso)
	err := row.Scan(&currency.ID, &currency.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Errorf("[CurrencyRepo][GetCurrencyByISO], No Content, Error : %s", err.Error())
			return entities.Currency{}, err
		}
		log.Errorf("[CurrencyRepo][GetCurrencyByISO], Error : %s", err.Error())
		return entities.Currency{}, err
	}

	return currency, nil

}
