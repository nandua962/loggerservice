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

// PaymentGatewayRepo is responsible for handling PaymentGateway-related data.
type PaymentGatewayRepo struct {
	db *sql.DB
}

// PaymentGatewayRepoImply is the interface defining the methods for working with PaymentGateway data.
type PaymentGatewayRepoImply interface {
	GetPaymentGatewayByID(ctx context.Context, id string) (entities.PaymentGatewayName, error)
	GetAllPaymentGateway(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errs map[string]models.ErrorResponse) ([]entities.PaymentGateway, int64, error)
}

// NewPaymentGatewayRepo creates a new PaymentGatewayRepo instance.
func NewPaymentGatewayRepo(db *sql.DB) PaymentGatewayRepoImply {
	return &PaymentGatewayRepo{db: db}
}

// GetPaymentGateways retrieves PaymentGateway data based on specified payment gateway id.
func (paymentGatewayRepo *PaymentGatewayRepo) GetPaymentGatewayByID(ctx context.Context, id string) (entities.PaymentGatewayName, error) {

	var (
		res entities.PaymentGatewayName
		log = logger.Log().WithContext(ctx)
	)

	query := `
		SELECT 
			name
		FROM payment_gateway 
		WHERE id= $1
		`

	row := paymentGatewayRepo.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&res.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Errorf("[paymentGatewayRepo][GetPaymentGatewayByID], No Content, Error : %s", err.Error())
			return entities.PaymentGatewayName{}, err
		}
		log.Errorf("[paymentGatewayRepo][GetPaymentGatewayByID], Error : %s", err.Error())
		return entities.PaymentGatewayName{}, err
	}
	return res, nil
}

// GetPaymentGatewayByID retrieves PaymentGateway data based on specified query paramaters
func (paymentGatewayRepo *PaymentGatewayRepo) GetAllPaymentGateway(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errs map[string]models.ErrorResponse) ([]entities.PaymentGateway, int64, error) {

	var (
		payments     []entities.PaymentGateway
		totalRecords int64
		log          = logger.Log().WithContext(ctx)
	)

	query := `
		SELECT 
			id,
			name,
			COUNT(id) OVER() as total_records
		FROM payment_gateway 
	`

	if !utils.IsEmpty(params.Name) {
		query = utilities.Search(query, params.Name, "name")
	}

	query = utilities.GroupBy(query, "id, name")
	sortQ, err := utilities.OrderBy(
		ctx,
		params.Sort, params.Order,
		consts.SortOptns,
		validation.Endpoint,
		validation.Method,
		errs,
	)
	if err != nil {
		log.Errorf("[paymentGatewayRepo][GetAllPaymentGateway] Error : %s", err.Error())
		return nil, 0, err
	}
	query = fmt.Sprintf("%s %s", query, sortQ)
	query = fmt.Sprintf("%s %s", query, utils.CalculateOffset(pagination.Page, pagination.Limit))
	if len(errs) > 0 {
		return nil, 0, nil
	}

	res, err := paymentGatewayRepo.db.QueryContext(ctx, query)

	if err != nil {
		log.Errorf("[paymentGatewayRepo][GetPaymentGatewayByID], Error : %s", err.Error())
		return nil, 0, err
	}

	defer res.Close()
	for res.Next() {
		var payment entities.PaymentGateway
		err := res.Scan(
			&payment.ID,
			&payment.Name,
			&totalRecords,
		)
		if err != nil {
			log.Errorf("[paymentGatewayRepo][GetPaymentGatewayByID],Error while Scan, Error : %s", err.Error())
			return nil, 0, err
		}

		payments = append(payments, payment)
	}

	return payments, totalRecords, nil
}
