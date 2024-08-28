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

// CurrencyController represents the controller for currency-related operations.
type CurrencyController struct {
	router   *gin.RouterGroup              // The Gin router group for currency routes.
	useCases usecases.CurrencyUseCaseImply // The use case interface for currency operations.
	cfg      *entities.EnvConfig
}

// NewCurrencyController creates a new instance of CurrencyController.
func NewCurrencyController(router *gin.RouterGroup, currencyUseCase usecases.CurrencyUseCaseImply, cfg *entities.EnvConfig) *CurrencyController {
	return &CurrencyController{
		router:   router,
		useCases: currencyUseCase,
		cfg:      cfg,
	}
}

// InitRoutes initializes the currency-related routes for the CurrencyController.
func (currencyCtrl *CurrencyController) InitRoutes() {
	currencyCtrl.router.GET("/:version/currencies", func(ctx *gin.Context) {
		version.RenderHandler(ctx, currencyCtrl, "GetAllCurrency")
	})
	currencyCtrl.router.GET("/:version/currencies/:id", func(ctx *gin.Context) {
		version.RenderHandler(ctx, currencyCtrl, "GetCurrencyByID")
	})
	currencyCtrl.router.HEAD("/:version/currencies/exists/:iso", func(ctx *gin.Context) {
		version.RenderHandler(ctx, currencyCtrl, "GetCurrencyByISO")
	})
	currencyCtrl.router.GET("/:version/currencies/exists/:iso", func(ctx *gin.Context) {
		version.RenderHandler(ctx, currencyCtrl, "GetCurrencyByISO")
	})

}

// GetAllCurrency retrieves currencies based on the provided parameters and returns them as JSON.
func (currencyCtrl *CurrencyController) GetAllCurrency(ctx *gin.Context) {
	var (
		req           entities.Params
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		errMap        = utilities.NewErrorMap()
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)

	validation.HelpLink = consts.ErrorHelpLink
	log.Info("[CurrencyController][GetAllCurrency] Processing GetAllCurrency request")

	if err := ctx.BindQuery(&req); err != nil {
		log.Errorf("[CurrencyController][GetAllCurrency], Invalid query params, Error : %s", err.Error())
		result := utilities.ErrorResponseGenerator(consts.BindingError, http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.GetAllCurrencyIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	paginationInfo, _ := utils.GetContext[entities.Pagination](ctx, consts.PaginationKey)
	resp, errMap, err := currencyCtrl.useCases.GetCurrencies(ctx, req, paginationInfo, validation, errMap)

	if err != nil {
		log.Errorf("[CurrencyController][GetAllCurrency] Error : %s", err.Error())
		utilities.HandleError(ctx, err, validation)
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[CurrencyController][GetAllCurrency] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}
	var output entities.Result

	if resp.MetaData == nil {
		result := utilities.SuccessResponseGenerator("Currency listed successfully", http.StatusNoContent, "")
		ctx.JSON(http.StatusNoContent, result)
		return
	}

	log.Info("[CurrencyController][GetAllCurrency] Currency data fetched successfully")

	output.Data = resp.Data
	output.Metadata = resp.MetaData
	result := utilities.SuccessResponseGenerator("Currency retrieved successfully", http.StatusOK, output)
	ctx.JSON(http.StatusOK, result)
}

// GetCurrencyByID retrieves currencies based on the provided ID .
func (currencyCtrl *CurrencyController) GetCurrencyByID(ctx *gin.Context) {
	var (
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		errMap        = utilities.NewErrorMap()
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)

	validation.HelpLink = consts.ErrorHelpLink
	log.Info("[CurrencyController][GetCurrencyByID] Processing GetCurrencyByID request")

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.GetCurrencyByIDIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	validation.ID = ctx.Param(consts.IDKey)

	resp, errMap, err := currencyCtrl.useCases.GetCurrencyByID(ctx, validation, errMap)

	if err != nil {
		if err == sql.ErrNoRows {
			result := utilities.ErrorResponseGenerator("Currency not found", http.StatusNotFound, "")
			ctx.JSON(http.StatusNotFound, result)
			return
		}
		log.Errorf("[CurrencyController][GetCurrencyByID], Error : %s", err.Error())
		utilities.HandleError(ctx, err, validation)
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[CurrencyController][GetCurrencyByID] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}

	log.Info("[CurrencyController][GetCurrencyByID] Currency data fetched successfully")

	result := utilities.SuccessResponseGenerator("Currency retrieved successfully", http.StatusOK, resp)
	ctx.JSON(http.StatusOK, result)
}

// GetCurrencyByISO retrieves currencies based on the provided ISO .
func (currencyCtrl *CurrencyController) GetCurrencyByISO(ctx *gin.Context) {
	var (
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		errMap        = utilities.NewErrorMap()
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)

	validation.HelpLink = consts.ErrorHelpLink
	log.Info("[CurrencyController][GetCurrencyByISO] Processing GetCurrencyByISO request")

	validation.ID = ctx.Param((consts.IsoField))

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.GetCurrencyByISOIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	resp, errMap, err := currencyCtrl.useCases.GetCurrencyByISO(ctx, validation, errMap)

	if err != nil {
		if err == sql.ErrNoRows {
			result := utilities.ErrorResponseGenerator("Currency not found", http.StatusNotFound, "")
			ctx.JSON(http.StatusNotFound, result)
			return
		}
		log.Errorf("[CurrencyController][GetCurrencyByISO], Error : %s", err.Error())
		utilities.HandleError(ctx, err, validation)
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[CurrencyController][GetCurrencyByISO] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}

	log.Info("[CurrencyController][GetCurrencyByISO] Currency data fetched successfully")

	result := utilities.SuccessResponseGenerator("Currency retrieved successfully", http.StatusOK, resp)
	ctx.JSON(http.StatusOK, result)

}
