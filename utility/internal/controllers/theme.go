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

// ThemeController represents the controller for theme-related operations.
type ThemeController struct {
	router   *gin.RouterGroup           // The Gin router group for theme routes.
	useCases usecases.ThemeUseCaseImply // The use case interface for theme operations.
	cfg      *entities.EnvConfig
}

// NewThemeController creates a new instance of ThemeController.
func NewThemeController(router *gin.RouterGroup, themeUseCase usecases.ThemeUseCaseImply, cfg *entities.EnvConfig) *ThemeController {
	return &ThemeController{
		router:   router,
		useCases: themeUseCase,
		cfg:      cfg,
	}
}

// InitRoutes initializes and configures the theme-related routes for the ThemeController.
func (theme *ThemeController) InitRoutes() {

	// Define and configure the route for getting themes.
	theme.router.GET("/:version/theme/:id", func(ctx *gin.Context) {
		version.RenderHandler(ctx, theme, "GetThemeByID")
	})
}

// GetThemeByID handles the HTTP GET request to retrieve theme by specified ID.
func (theme *ThemeController) GetThemeByID(ctx *gin.Context) {

	var (
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		errMap        = utilities.NewErrorMap()
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)

	validation.HelpLink = consts.ErrorHelpLink
	log.Info("[ThemeController][GetThemeByID] Processing GetThemeByID request")

	validation.ID = ctx.Param(consts.IDKey)

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.GetThemeByIDIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	resp, errMap, err := theme.useCases.GetThemeByID(ctxt, validation, errMap)

	if err != nil {
		if err == sql.ErrNoRows {
			result := utilities.ErrorResponseGenerator("Theme not found", http.StatusNotFound, "")
			ctx.JSON(http.StatusNotFound, result)
			return
		}
		log.Errorf("[ThemeController][GetThemeByID] Error : %s", err.Error())
		utilities.HandleError(ctx, err, validation)
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[ThemeController][GetThemeByID] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}

	log.Info("[ThemeController][GetThemeByID] Theme fetched successfully")

	result := utilities.SuccessResponseGenerator("Theme retrieved successfully", http.StatusOK, resp)
	ctx.JSON(http.StatusOK, result)
}
