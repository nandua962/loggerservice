package controllers

import (
	"database/sql"
	"net/http"
	"strings"
	"utility/internal/consts"
	"utility/internal/entities"
	"utility/internal/usecases"
	"utility/utilities"

	"github.com/gin-gonic/gin"
	constants "gitlab.com/tuneverse/toolkit/consts"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/core/version"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils"
)

// PaymentGatewayController represents a controller responsible for handling PaymentGateway-related API requests.
type PaymentGatewayController struct {
	router   *gin.RouterGroup
	useCases usecases.PaymentGatewayUseCaseImply
	cfg      *entities.EnvConfig
}

// NewPaymentGatewayController creates a new PaymentGatewayController instance.
func NewPaymentGatewayController(router *gin.RouterGroup, paymentGatewayUseCase usecases.PaymentGatewayUseCaseImply, cfg *entities.EnvConfig) *PaymentGatewayController {
	return &PaymentGatewayController{
		router:   router,
		useCases: paymentGatewayUseCase,
		cfg:      cfg,
	}
}

// InitRoutes initializes and configures the PaymentGateway-related routes for the PaymentGatewayController.
func (paymentGateway *PaymentGatewayController) InitRoutes() {
	paymentGateway.router.GET("/:version/payment_gateway/:id", func(ctx *gin.Context) {
		version.RenderHandler(ctx, paymentGateway, "GetPaymentGatewayByID")
	})
	paymentGateway.router.GET("/:version/payment_gateway/all", func(ctx *gin.Context) {
		version.RenderHandler(ctx, paymentGateway, "GetAllPaymentGateway")
	})
}

// GetPaymentGatewayByID handles the retrieval of PaymentGateways by unique ID.
func (paymentGateway *PaymentGatewayController) GetPaymentGatewayByID(ctx *gin.Context) {

	var (
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		errMap        = utilities.NewErrorMap()
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)
	validation.HelpLink = consts.ErrorHelpLink
	log.Info("[PaymentGatewayController][GetPaymentGatewayByID] Processing GetPaymentGatewayByID request")

	validation.ID = ctx.Param(consts.IDKey)

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.GetPaymentGatewayByIDIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	resp, errMap, err := paymentGateway.useCases.GetPaymentGatewayByID(ctx, validation, errMap)

	if err != nil {
		if err == sql.ErrNoRows {
			result := utilities.ErrorResponseGenerator("Payment Gateway not found", http.StatusNotFound, "")
			ctx.JSON(http.StatusNotFound, result)
			return
		}
		log.Errorf("[PaymentGatewayController][GetPaymentGatewayByID], Error : %s", err.Error())
		utilities.HandleError(ctx, err, validation)
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[PaymentGatewayController][GetPaymentGatewayByID] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}

	log.Info("[PaymentGatewayController][GetPaymentGatewayByID] Payment Gateway data fetched successfully")

	result := utilities.SuccessResponseGenerator("Payment Gateway retrieved successfully", http.StatusOK, resp)
	ctx.JSON(http.StatusOK, result)

}

// GetAllPaymentGateway handles the retrieval of PaymentGateways based on the provided parameters and returns them as JSON.
func (paymentGateway *PaymentGatewayController) GetAllPaymentGateway(ctx *gin.Context) {

	var (
		req           entities.Params
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		errMap        = utilities.NewErrorMap()
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)
	validation.HelpLink = paymentGateway.cfg.ErrorHelpLink
	log.Info("[PaymentGatewayController][GetAllPaymentGateway] Processing GetAllPaymentGateway request")

	if err := ctx.BindQuery(&req); err != nil {
		log.Errorf("[PaymentGatewayController][GetAllPaymentGateway], Invalid query params, Error : %s", err.Error())
		result := utilities.ErrorResponseGenerator(consts.BindingError, http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.GetAllPaymentGatewayIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	paginationInfo, _ := utils.GetContext[entities.Pagination](ctx, consts.PaginationKey)

	resp, errMap, err := paymentGateway.useCases.GetAllPaymentGateway(ctx, req, paginationInfo, validation, errMap)

	if err != nil {
		log.Errorf("[PaymentGatewayController][GetAllPaymentGateway], Error : %s", err.Error())
		utilities.HandleError(ctx, err, validation)
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[PaymentGatewayController][GetAllPaymentGateway] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}

	var output entities.Result
	if resp.MetaData == nil {
		result := utilities.SuccessResponseGenerator("Payment gateway listed successfully", http.StatusNoContent, "")
		ctx.JSON(http.StatusNoContent, result)
		return
	}

	log.Info("[PaymentGatewayController][GetAllPaymentGateway] Payment Gateway data fetched successfully")

	output.Data = resp.Data
	output.Metadata = resp.MetaData
	result := utilities.SuccessResponseGenerator("Payment Gateway retrieved successfully", http.StatusOK, output)
	ctx.JSON(http.StatusOK, result)

}
