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

// LookupController represents a controller responsible for handling lookup-related API requests.
type LookupController struct {
	router   *gin.RouterGroup
	useCases usecases.LookupUseCaseImply
	cfg      *entities.EnvConfig
}

// NewLookupController creates a new LookupController instance.
func NewLookupController(router *gin.RouterGroup, lookupUseCase usecases.LookupUseCaseImply, cfg *entities.EnvConfig) *LookupController {
	return &LookupController{
		router:   router,
		useCases: lookupUseCase,
		cfg:      cfg,
	}
}

// InitRoutes initializes and configures the lookup-related routes for the LookupController.
func (lookup *LookupController) InitRoutes() {
	lookup.router.POST("/:version/lookup", func(ctx *gin.Context) {
		version.RenderHandler(ctx, lookup, "GetLookupByIdList")
	})
	lookup.router.GET("/:version/lookup/type/:name", func(ctx *gin.Context) {
		version.RenderHandler(ctx, lookup, "GetLookupByTypeName")
	})

}

// GetLookupByTypeName handles the retrieval of lookup.
func (lookup *LookupController) GetLookupByTypeName(ctx *gin.Context) {

	var (
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		errMap        = utilities.NewErrorMap()
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)

	validation.HelpLink = consts.ErrorHelpLink
	log.Info("[LookupController][GetLookupByTypeName] Processing GetLookupByTypeName request")

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.GetLookupByTypeNameIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	// Set the type name for validation
	validation.ID = ctx.Param(consts.NameKey) // name from lookup_type table

	filter := make(map[string]string)
	filterValue := ctx.DefaultQuery("value", "") // value from lookup table
	if filterValue != "" {
		filter["value"] = filterValue
	}
	// Call the use case to get a lookup
	resp, errMap, err := lookup.useCases.GetLookupByTypeName(ctxt, filter, validation, errMap)

	if err != nil {
		log.Errorf("[LookupController][GetLookupByTypeName], Error : %s", err.Error())
		utilities.HandleError(ctx, err, validation)
		return
	}

	// Handle validation errors if any
	if len(errMap) != 0 {
		log.Errorf("[LookupController][GetLookupByTypeName] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}

	// Handle case where no lookup is returned
	if resp == nil {
		result := utilities.SuccessResponseGenerator("Lookup not found", http.StatusNotFound, "")
		ctx.JSON(http.StatusNotFound, result)
		return
	}

	log.Info("[LookupController][GetLookupByTypeName] Lookup data fetched successfully")

	// Generate a success response
	result := utilities.SuccessResponseGenerator("Lookup retrieved successfully", http.StatusOK, resp)
	ctx.JSON(http.StatusOK, result)
}

// GetLookupByIdList handles the retrieval of lookup by array of lookup Ids
func (lookup *LookupController) GetLookupByIdList(ctx *gin.Context) {
	var (
		req           entities.LookupIDs
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)

	validation.HelpLink = consts.ErrorHelpLink
	if err := ctx.BindJSON(&req); err != nil {
		log.Errorf("[LookupController][GetLookupByIdList], Invalid query params, Error : %s", err.Error())
		result := utilities.ErrorResponseGenerator(consts.BindingError, http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}

	log.Info("[LookupController][GetLookupByIdList] Processing GetLookupByIdList request")

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.GetLookupByIdListIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	// Call the use case to get a lookup
	resp, err := lookup.useCases.GetLookupByIdList(ctxt, req)

	if err != nil {
		if err == sql.ErrNoRows {
			result := utilities.ErrorResponseGenerator("Lookup not found", http.StatusNotFound, "")
			ctx.JSON(http.StatusNotFound, result)
			return
		}
		log.Errorf("[LookupController][GetLookupByIdList], Error : %s", err.Error())
		utilities.HandleError(ctx, err, validation)
		return
	}

	log.Info("[LookupController][GetLookupByIdList] Lookup data fetched successfully")

	// Generate a success response
	result := utilities.SuccessResponseGenerator("Lookup retrieved successfully", http.StatusOK, resp)
	ctx.JSON(http.StatusOK, result)
}
