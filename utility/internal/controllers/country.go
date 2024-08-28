package controllers

import (
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

// CountryController handles HTTP requests related to countries and states.
type CountryController struct {
	router   *gin.RouterGroup
	useCases usecases.CountryUseCaseImply
	cfg      *entities.EnvConfig
}

// NewCountryController creates a new instance of CountryController.
func NewCountryController(router *gin.RouterGroup, CountryUseCase usecases.CountryUseCaseImply, cfg *entities.EnvConfig) *CountryController {
	return &CountryController{
		router:   router,
		useCases: CountryUseCase,
		cfg:      cfg,
	}
}

// InitRoutes initializes and configures the country-related routes for the CountryController.
func (countryCtrl *CountryController) InitRoutes() {

	countryCtrl.router.GET("/:version/countries", func(ctx *gin.Context) {
		version.RenderHandler(ctx, countryCtrl, "GetCountries")
	})
	countryCtrl.router.GET("/:version/countries/:country_id/states", func(ctx *gin.Context) {
		version.RenderHandler(ctx, countryCtrl, "GetStatesOfCountry")
	})
	countryCtrl.router.GET("/:version/countries/exists", func(ctx *gin.Context) {
		version.RenderHandler(ctx, countryCtrl, "CheckCountryExists")
	})
	countryCtrl.router.GET("/:version/countries/iso", func(ctx *gin.Context) {
		version.RenderHandler(ctx, countryCtrl, "GetAllCountryCodes")
	})
	countryCtrl.router.HEAD("/:version/countries/:country_code/states/:iso", func(ctx *gin.Context) {
		version.RenderHandler(ctx, countryCtrl, "CheckStateExists")
	})

}

// CheckStateExists checks the existence of states based on ISO and country_code parameters.
func (countryCtrl *CountryController) CheckStateExists(ctx *gin.Context) {

	var (
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		errMap        = utilities.NewErrorMap()
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)
	iso := ctx.Param(consts.IsoField)
	countryCode := ctx.Param(consts.CountryCode)
	validation.HelpLink = consts.ErrorHelpLink
	log.Info("[CountryController][CheckStateExists] Processing CheckStateExists request")

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.CheckStateExistsIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	isExists, errMap, err := countryCtrl.useCases.CheckStateExists(ctxt, iso, countryCode, validation, errMap)

	if err != nil {
		log.Errorf("[CountryController][CheckStateExists], Error : %s", err.Error())
		utilities.HandleError(ctx, err, validation)
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[CountryController][CheckStateExists] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}

	if !isExists {
		result := utilities.ErrorResponseGenerator("State retrieved successfully", http.StatusNotFound, "")
		ctx.JSON(http.StatusNotFound, result)
		return
	}
	log.Info("[CountryController][CheckStateExists] State data fetched successfully")

	result := utilities.SuccessResponseGenerator("State retrieved successfully", http.StatusOK, "")
	ctx.JSON(http.StatusOK, result)

}

// GetCountries retrieves a list of countries based on provided query parameters.
func (countryCtrl *CountryController) GetCountries(ctx *gin.Context) {

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
	log.Info("[CountryController][GetCountries] Processing GetCountries request")

	if err := ctx.BindQuery(&req); err != nil {
		log.Errorf("[CountryController][GetCountries], Invalid query params, Error : %s", err.Error())
		result := utilities.ErrorResponseGenerator(consts.BindingError, http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.GetCountriesIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	paginationInfo, _ := utils.GetContext[entities.Pagination](ctx, consts.PaginationKey)
	resp, errMap, err := countryCtrl.useCases.GetCountries(ctxt, req, paginationInfo, validation, errMap)

	if err != nil {
		log.Errorf("[CountryController][GetCountries], Error : %s", err.Error())
		utilities.HandleError(ctx, err, validation)
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[CountryController][GetCountries] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}

	var output entities.Result

	if resp.MetaData == nil {
		result := utilities.SuccessResponseGenerator("Countries listed successfully", http.StatusNoContent, "")
		ctx.JSON(http.StatusNoContent, result)
		return
	}

	log.Info("Countries fetched successfully")

	output.Data = resp.Data
	output.Metadata = resp.MetaData
	result := utilities.SuccessResponseGenerator("Data retrieved successfully", http.StatusOK, output)
	ctx.JSON(http.StatusOK, result)
}

// GetStatesOfCountry retrieves a list of states for a specific country based on provided query parameters.
func (countryCtrl *CountryController) GetStatesOfCountry(ctx *gin.Context) {

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
	log.Info("[CountryController][GetStatesOfCountry] Processing GetStatesOfCountry request")

	if err := ctx.BindQuery(&req); err != nil {
		log.Errorf("[CountryController][GetStatesOfCountry], Invalid query params, Error : %s", err.Error())
		result := utilities.ErrorResponseGenerator(consts.BindingError, http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.GetStatesOfCountryIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	paginationInfo, _ := utils.GetContext[entities.Pagination](ctx, consts.PaginationKey)

	validation.ID = ctx.Param(consts.Country)
	if utils.IsEmpty(validation.ID) {
		result := utilities.ErrorResponseGenerator(consts.CountryNotFound, http.StatusBadRequest, "")
		ctx.JSON(result.Code, result)
		return
	}

	resp, errMap, err := countryCtrl.useCases.GetStatesOfCountry(ctxt, req, paginationInfo, validation, errMap)

	if err != nil {
		if err == consts.ErrCountryNotFound {
			result := utilities.ErrorResponseGenerator(consts.CountryNotFound, http.StatusBadRequest, "")
			ctx.JSON(result.Code, result)
			return
		} else {
			utilities.HandleError(ctx, err, validation)
			return
		}
	}

	if len(errMap) != 0 {
		log.Errorf("[CountryController][GetStatesOfCountry] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}

	var output entities.Result

	if resp.MetaData == nil {
		result := utilities.SuccessResponseGenerator("StatesOfCountry listed successfully", http.StatusNoContent, "")
		ctx.JSON(http.StatusNoContent, result)
		return
	}

	log.Info("States fetched successfully")

	output.Data = resp.Data
	output.Metadata = resp.MetaData
	result := utilities.SuccessResponseGenerator("Data retrieved successfully", http.StatusOK, output)
	ctx.JSON(http.StatusOK, result)
}

// CheckCountryExists checks the existence of countries based on ISO parameters.
func (countryCtrl *CountryController) CheckCountryExists(ctx *gin.Context) {

	var (
		req           entities.IsoParam
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		errMap        = utilities.NewErrorMap()
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)
	validation.HelpLink = consts.ErrorHelpLink
	log.Info("[CountryController][CheckCountryExists] Processing CheckCountryExists request")

	if err := ctx.BindQuery(&req); err != nil {
		log.Errorf("[CountryController][CheckCountryExists], Invalid query params, Error : %s", err.Error())
		result := utilities.ErrorResponseGenerator(consts.BindingError, http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}
	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.CheckCountryExistsIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	resp, errMap, err := countryCtrl.useCases.CheckCountryExists(ctxt, req, validation, errMap)

	if err != nil {
		log.Errorf("[CountryController][CheckCountryExists], Error : %s", err.Error())
		utilities.HandleError(ctx, err, validation)
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[CountryController][CheckCountryExists] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}

	log.Info("Result fetched successfully")

	result := utilities.SuccessResponseGenerator("Data retrieved successfully", http.StatusOK, resp)
	ctx.JSON(http.StatusOK, result)
}

// GetCountries retrieves a list of countries based on provided query parameters.
func (countryCtrl *CountryController) GetAllCountryCodes(ctx *gin.Context) {

	var (
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)
	validation.HelpLink = consts.ErrorHelpLink
	log.Info("[CountryController][GetAllCountryCodes] Processing GetAllCountryCodes request")

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.GetAllCountryCodesIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	resp, err := countryCtrl.useCases.GetAllCountryCodes(ctxt)

	if err != nil {
		log.Errorf("[CountryController][GetAllCountryCodes], Error : %s", err.Error())
		utilities.HandleError(ctx, err, validation)
		return
	}

	log.Info("Country List fetched successfully")

	result := utilities.SuccessResponseGenerator("Data retrieved successfully", http.StatusOK, resp)
	ctx.JSON(http.StatusOK, result)
}
