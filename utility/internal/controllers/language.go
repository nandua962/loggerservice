package controllers

import (
	"net/http"
	"strings"
	"utility/internal/entities"
	"utility/internal/usecases"
	"utility/utilities"

	"utility/internal/consts"

	constants "gitlab.com/tuneverse/toolkit/consts"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/core/version"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils"

	"github.com/gin-gonic/gin"
)

// LanguageController handles HTTP requests related to languages.
type LanguageController struct {
	router   *gin.RouterGroup
	useCases usecases.LanguageUseCaseImply
	cfg      *entities.EnvConfig
}

// NewLanguageController creates a new instance of LanguageController.
func NewLanguageController(router *gin.RouterGroup, languageUseCase usecases.LanguageUseCaseImply, cfg *entities.EnvConfig) *LanguageController {
	return &LanguageController{
		router:   router,
		useCases: languageUseCase,
		cfg:      cfg,
	}
}

// InitRoutes initializes and configures the language-related routes for the LanguageController.
func (language *LanguageController) InitRoutes() {

	// Define and configure the route for getting languages.
	language.router.GET("/:version/languages", func(ctx *gin.Context) {
		version.RenderHandler(ctx, language, "GetLanguages")
	})
	language.router.HEAD("/:version/languages/exists/:code", func(ctx *gin.Context) {
		version.RenderHandler(ctx, language, "GetLanguageCodeExists")
	})
}

// GetLanguages handles the HTTP GET request to retrieve languages.
func (language *LanguageController) GetLanguages(ctx *gin.Context) {
	var (
		req           entities.LangParams
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		errMap        = utilities.NewErrorMap()
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)

	validation.HelpLink = consts.ErrorHelpLink
	log.Info("[LanguageController][GetLanguages] Processing GetLanguages request")

	// Bind query parameters to the 'req' variable.
	if err := ctx.BindQuery(&req); err != nil {
		log.Errorf("[LanguageController][GetLanguages], Invalid query params, Error : %s", err.Error())
		result := utilities.ErrorResponseGenerator(consts.BindingError, http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.GetLanguagesIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)
	paginationInfo, _ := utils.GetContext[entities.Pagination](ctx, consts.PaginationKey)

	resp, errMap, err := language.useCases.GetLanguages(ctxt, req, paginationInfo, validation, errMap)

	if err != nil {
		log.Errorf("[LanguageController][GetLanguages] Error : %s", err.Error())
		utilities.HandleError(ctx, err, validation)
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[LanguageController][GetLanguages] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}

	var output entities.Result

	if resp.MetaData == nil {
		result := utilities.SuccessResponseGenerator("Languages listed successfully", http.StatusNoContent, "")
		ctx.JSON(http.StatusNoContent, result)
		return
	}

	log.Info("[LanguageController][GetLanguages] Languages data fetched successfully")

	output.Data = resp.Data
	output.Metadata = resp.MetaData
	result := utilities.SuccessResponseGenerator("Languages retrieved successfully", http.StatusOK, output)
	ctx.JSON(http.StatusOK, result)
}

// GetLanguageCodeExists retrieves currencies based on the provided Code .
func (language *LanguageController) GetLanguageCodeExists(ctx *gin.Context) {
	var (
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		errMap        = utilities.NewErrorMap()
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)

	validation.HelpLink = consts.ErrorHelpLink
	log.Info("[LanguageController][GetLanguageCodeExists] Processing GetLanguageCodeExists request")

	validation.ID = ctx.Param(consts.CodeField)

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.GetLanguageCodeExistsIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	isExists, errMap, err := language.useCases.GetLanguageCodeExists(ctx, validation, errMap)

	if err != nil {
		log.Errorf("[LanguageController][GetLanguageCodeExists] Error : %s", err.Error())
		utilities.HandleError(ctx, err, validation)
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[LanguageController][GetLanguageCodeExists] Validation Error")
		result := utilities.SuccessResponseGenerator("Data retrieved successfully", http.StatusBadRequest, "")
		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	if !isExists {
		result := utilities.SuccessResponseGenerator("Data retrieved successfully", http.StatusNotFound, "")
		ctx.JSON(http.StatusNotFound, result)
		return
	}
	log.Info("Language data fetched successfully")

	result := utilities.SuccessResponseGenerator("Data retrieved successfully", http.StatusOK, "")
	ctx.JSON(http.StatusOK, result)
}
