package usecases

import (
	"context"
	"utility/internal/consts"
	"utility/internal/entities"
	"utility/internal/repo"

	"github.com/google/uuid"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils"
)

// PaymentGatewayUseCases represents use cases for handling PaymentGateway-related operations.
type PaymentGatewayUseCases struct {
	repo repo.PaymentGatewayRepoImply
}

// PaymentGatewayUseCaseImply is an interface defining the methods for working with PaymentGateway use cases.
type PaymentGatewayUseCaseImply interface {
	GetPaymentGatewayByID(ctx context.Context, validation entities.Validation, errMap map[string]models.ErrorResponse) (entities.PaymentGatewayName, map[string]models.ErrorResponse, error)
	GetAllPaymentGateway(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errMap map[string]models.ErrorResponse) (*entities.Response, map[string]models.ErrorResponse, error)
}

// NewPaymentGatewayUseCases creates a new PaymentGatewayUseCases instance.
func NewPaymentGatewayUseCases(paymentGatewayRepo repo.PaymentGatewayRepoImply) PaymentGatewayUseCaseImply {
	return &PaymentGatewayUseCases{
		repo: paymentGatewayRepo,
	}
}

// GetPaymentGatewayByID retrieves PaymentGateway data based on specified payment gateway id.
func (paymentGateway *PaymentGatewayUseCases) GetPaymentGatewayByID(ctx context.Context, validation entities.Validation, errMap map[string]models.ErrorResponse) (entities.PaymentGatewayName, map[string]models.ErrorResponse, error) {

	log := logger.Log().WithContext(ctx)
	var result entities.PaymentGatewayName

	_, err := uuid.Parse(validation.ID)

	if err != nil {

		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.PaymentIdField, consts.InvalidKey)
		if err != nil {
			log.Errorf("[PaymentGatewayUseCases][GetPaymentGatewayByID] Error while loading service code Error : %s", err.Error())
		}
		errMap[consts.PaymentIdField] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}
	}

	if len(errMap) > 0 {
		return entities.PaymentGatewayName{}, errMap, nil
	}

	result, err = paymentGateway.repo.GetPaymentGatewayByID(ctx, validation.ID)

	if err != nil {
		log.Errorf("[PaymentGatewayUseCases][GetPaymentGatewayByID] Error : %s", err.Error())
		return entities.PaymentGatewayName{}, nil, err
	}

	return result, nil, nil
}

// GetPaymentGatewayByID retrieves PaymentGateway data based on specified query paramaters
func (paymentGateway *PaymentGatewayUseCases) GetAllPaymentGateway(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errMap map[string]models.ErrorResponse) (*entities.Response, map[string]models.ErrorResponse, error) {

	log := logger.Log().WithContext(ctx)

	payments, totalRecords, err := paymentGateway.repo.GetAllPaymentGateway(ctx, params, pagination, validation, errMap)

	if len(errMap) > 0 {
		return nil, errMap, nil
	}

	if err != nil {
		log.Errorf("[PaymentGatewayUseCases][GetAllPaymentGateway] Error : %s", err.Error())
		return nil, nil, err
	}

	metaData := utils.MetaDataInfo(&models.MetaData{
		Total:       totalRecords,
		PerPage:     pagination.Limit,
		CurrentPage: pagination.Page,
	})

	result := &entities.Response{
		MetaData: metaData,
		Data:     payments,
	}

	return result, nil, nil
}
