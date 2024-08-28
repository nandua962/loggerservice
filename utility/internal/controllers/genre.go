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

// GenreController represents a controller responsible for handling genre-related API requests.
type GenreController struct {
	router   *gin.RouterGroup
	useCases usecases.GenreUseCaseImply
	cfg      *entities.EnvConfig
}

// NewGenreController creates a new GenreController instance.
func NewGenreController(router *gin.RouterGroup, genreUseCase usecases.GenreUseCaseImply, cfg *entities.EnvConfig) *GenreController {
	return &GenreController{
		router:   router,
		useCases: genreUseCase,
		cfg:      cfg,
	}
}

// InitRoutes initializes and configures the genre-related routes for the GenreController.
func (genre *GenreController) InitRoutes() {
	genre.router.GET("/:version/health", func(ctx *gin.Context) {
		version.RenderHandler(ctx, genre, "HealthHandler")
	})
	genre.router.POST("/:version/genres", func(ctx *gin.Context) {
		version.RenderHandler(ctx, genre, "CreateGenre")
	})
	genre.router.GET("/:version/genres", func(ctx *gin.Context) {
		version.RenderHandler(ctx, genre, "GetGenres")
	})
	genre.router.GET("/:version/genres/:id", func(ctx *gin.Context) {
		version.RenderHandler(ctx, genre, "GetGenresByID")
	})
	genre.router.DELETE("/:version/genres/:id", func(ctx *gin.Context) {
		version.RenderHandler(ctx, genre, "DeleteGenre")
	})
	genre.router.PATCH("/:version/genres/:id", func(ctx *gin.Context) {
		version.RenderHandler(ctx, genre, "UpdateGenre")
	})
}

// CreateGenre handles the creation of a genre.
func (genre *GenreController) CreateGenre(ctx *gin.Context) {

	var (
		req           entities.Genre
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		errMap        = utilities.NewErrorMap()
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)

	validation.HelpLink = consts.ErrorHelpLink
	log.Info("[GenreController][CreateGenre] Processing CreateGenre request")

	if err := ctx.BindJSON(&req); err != nil {
		log.Errorf("[GenreController][CreateGenre], Invalid query params, Error : %s", err.Error())
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
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.CreateGenreIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	errMap, err := genre.useCases.CreateGenre(ctxt, req, validation, errMap)

	if err != nil {
		log.Errorf("[GenreController][CreateGenre] Error : %s", err.Error())
		utilities.HandleError(ctx, err, validation)
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[GenreController][CreateGenre] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}

	result := utilities.SuccessResponseGenerator("Genre created successfully", http.StatusCreated, "")
	// Data  added successfully
	ctx.JSON(http.StatusCreated, result)

}

// GetGenres handles the retrieval of genres.
func (genre *GenreController) GetGenres(ctx *gin.Context) {
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
	if err := ctx.BindQuery(&req); err != nil {
		log.Errorf("[GenreController][GetGenres], Invalid query params, Error : %s", err.Error())
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
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.GetGenresIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)
	paginationInfo, _ := utils.GetContext[entities.Pagination](ctx, consts.PaginationKey)
	resp, errMap, err := genre.useCases.GetGenres(ctxt, req, paginationInfo, validation, errMap)

	if err != nil {
		log.Errorf("[GenreController][GetGenres] Error : %s", err.Error())
		utilities.HandleError(ctx, err, validation)
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[GenreController][GetGenres] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}

	var output entities.Result

	if resp.MetaData == nil {
		result := utilities.SuccessResponseGenerator("Genre listed successfully", http.StatusNoContent, "")
		ctx.JSON(http.StatusNoContent, result)
		return
	}

	log.Info("genre data fetched successfully")

	output.Data = resp.Data
	output.Metadata = resp.MetaData
	result := utilities.SuccessResponseGenerator("Data retrieved successfully", http.StatusOK, output)
	ctx.JSON(http.StatusOK, result)
}

// DeleteGenre handles the deletion of a genre.
func (genre *GenreController) DeleteGenre(ctx *gin.Context) {

	var (
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		errMap        = utilities.NewErrorMap()
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)
	validation.HelpLink = consts.ErrorHelpLink
	log.Info("[GenreController][DeleteGenre] Processing DeleteGenre request")

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.DeleteGenreIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	validation.ID = ctx.Param(consts.IDKey)
	errMap, err := genre.useCases.DeleteGenre(ctxt, validation, errMap)

	if err != nil {
		log.Errorf("[GenreController][DeleteGenre] Error : %s", err.Error())
		if err == consts.ErrNotExist {
			utilities.HandleNotFoundError(ctx, consts.GenreNotExist, http.StatusNotFound, consts.NotFound)
			return
		} else {
			utilities.HandleError(ctx, err, validation)
		}
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[GenreController][DeleteGenre] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}

	result := utilities.SuccessResponseGenerator("Genre deleted successfully", http.StatusOK, "")
	// Data deleted successfully
	ctx.JSON(http.StatusOK, result)

}

// UpdateGenre handles the updating of a genre.
func (genre *GenreController) UpdateGenre(ctx *gin.Context) {

	var (
		req           entities.Genre
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		errMap        = utilities.NewErrorMap()
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)

	validation.HelpLink = consts.ErrorHelpLink
	log.Info("[GenreController][UpdateGenre] Processing UpdateGenre request")

	if err := ctx.BindJSON(&req); err != nil {
		log.Errorf("[GenreController][UpdateGenre], Invalid json data, err=%s", err.Error())
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
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.UpdateGenreIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	validation.ID = ctx.Param(consts.IDKey)

	errMap, err := genre.useCases.UpdateGenre(ctxt, req, validation, errMap)

	if err != nil {
		log.Errorf("[GenreController][UpdateGenre], Error : %s", err.Error())
		if err == consts.ErrNotExist {
			utilities.HandleNotFoundError(ctx, consts.GenreNotExist, http.StatusNotFound, consts.NotFound)
			return
		} else {
			utilities.HandleError(ctx, err, validation)
		}
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[GenreController][UpdateGenre] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}

	result := utilities.SuccessResponseGenerator("Genre updated successfully", http.StatusOK, "")
	// Data  added successfully
	ctx.JSON(http.StatusOK, result)
}

// GetGenresByID handles the retrieval of genre of specified genre ID
func (genre *GenreController) GetGenresByID(ctx *gin.Context) {

	var (
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		errMap        = utilities.NewErrorMap()
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)
	validation.HelpLink = consts.ErrorHelpLink
	log.Info("[GenreController][GetGenresByID] Processing GetGenresByID request")

	validation.ID = ctx.Param(consts.IDKey)

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.GetGenresByIDIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	resp, errMap, err := genre.useCases.GetGenresByID(ctxt, validation, errMap)

	if err != nil {
		if err == sql.ErrNoRows {
			utilities.HandleNotFoundError(ctx, consts.GenreNotExist, http.StatusNotFound, consts.NotFound)
			return
		}

		log.Errorf("[GenreController][GetGenresByID] Error : %s", err.Error())
		utilities.HandleError(ctx, err, validation)
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[GenreController][GetGenresByID] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}

	log.Info("[GenreController][GetGenresByID] Genre code fetched successfully")

	result := utilities.SuccessResponseGenerator("Genre retrieved successfully", http.StatusOK, resp)
	ctx.JSON(http.StatusOK, result)

}
