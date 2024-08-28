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

// RoleController represents a controller responsible for handling role-related API requests.
type RoleController struct {
	router   *gin.RouterGroup
	useCases usecases.RoleUseCaseImply
	cfg      *entities.EnvConfig
}

// NewRoleController creates a new RoleController instance.
func NewRoleController(router *gin.RouterGroup, roleUseCase usecases.RoleUseCaseImply, cfg *entities.EnvConfig) *RoleController {
	return &RoleController{
		router:   router,
		useCases: roleUseCase,
		cfg:      cfg,
	}
}

// InitRoutes initializes and configures the role-related routes for the RoleController.
func (role *RoleController) InitRoutes() {

	role.router.GET("/:version/roles", func(ctx *gin.Context) {
		version.RenderHandler(ctx, role, "GetRoles")
	})
	role.router.GET("/:version/roles/:id", func(ctx *gin.Context) {
		version.RenderHandler(ctx, role, "GetRoleByID")
	})
	role.router.DELETE("/:version/roles/:id", func(ctx *gin.Context) {
		version.RenderHandler(ctx, role, "DeleteRoles")
	})
	role.router.POST("/:version/roles", func(ctx *gin.Context) {
		version.RenderHandler(ctx, role, "CreateRole")
	})
	role.router.PATCH("/:version/roles/:id", func(ctx *gin.Context) {
		version.RenderHandler(ctx, role, "UpdateRole")
	})
}

// GetRoleByID handles the retrieval of role of specified role ID
func (role *RoleController) GetRoleByID(ctx *gin.Context) {

	var (
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		errMap        = utilities.NewErrorMap()
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)
	validation.HelpLink = consts.ErrorHelpLink
	log.Info("[RoleController][GetRoleByID] Processing GetRoleByID request")

	validation.ID = ctx.Param(consts.IDKey)

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.GetRoleByIDIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	resp, errMap, err := role.useCases.GetRoleByID(ctxt, validation, errMap)

	if err != nil {
		if err == sql.ErrNoRows {
			utilities.HandleNotFoundError(ctx, consts.RoleNotExist, http.StatusNotFound, consts.NotFound)
			return
		}

		log.Errorf("[RoleController][GetRoleByID] Error : %s", err.Error())
		utilities.HandleError(ctx, err, validation)
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[RoleController][GetRoleByID] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}

	log.Info("[RoleController][GetRoleByID] role code fetched successfully")

	result := utilities.SuccessResponseGenerator("role retrieved successfully", http.StatusOK, resp)
	ctx.JSON(http.StatusOK, result)

}

// GetRoles handles the retrieval of roles.
func (role *RoleController) GetRoles(ctx *gin.Context) {
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
	log.Info("[RoleController][GetRoles] Processing GetRoles request")

	if err := ctx.BindQuery(&req); err != nil {
		log.Errorf("[RoleController][GetRoles], Invalid query params, Error : %s", err.Error())
		result := utilities.ErrorResponseGenerator(consts.BindingError, http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.GetRolesIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)
	paginationInfo, _ := utils.GetContext[entities.Pagination](ctx, consts.PaginationKey)
	resp, errMap, err := role.useCases.GetRoles(ctx, req, paginationInfo, validation, errMap)

	if err != nil {
		log.Errorf("[RoleController][GetRoles] Error : %s", err.Error())
		utilities.HandleError(ctx, err, validation)
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[RoleController][GetRoles] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}

	var output entities.Result

	if resp.MetaData == nil {
		result := utilities.SuccessResponseGenerator("Roles listed successfully", http.StatusNoContent, "")
		ctx.JSON(http.StatusNoContent, result)
		return
	}

	log.Info("[RoleController][GetRoles] Roles fetched successfully")

	output.Data = resp.Data
	output.Metadata = resp.MetaData
	result := utilities.SuccessResponseGenerator("Roles retrieved successfully", http.StatusOK, output)
	ctx.JSON(http.StatusOK, result)
}

// DeleteRoles handles the deletion of a role.
func (role *RoleController) DeleteRoles(ctx *gin.Context) {

	var (
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		errMap        = utilities.NewErrorMap()
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)

	validation.HelpLink = consts.ErrorHelpLink
	log.Info("[RoleController][DeleteRoles] Processing DeleteRoles request")

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.DeleteRolesIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	validation.ID = ctx.Param(consts.IDKey)

	errMap, err := role.useCases.DeleteRoles(ctxt, validation, errMap)

	if err != nil {
		if err == consts.ErrNotExist {
			log.Errorf("[RoleController][DeleteRoles] Error : %s", err.Error())
			utilities.HandleNotFoundError(ctx, consts.RoleNotExist, http.StatusNotFound, consts.NotFound)
			return
		} else {
			log.Errorf("[RoleController][DeleteRoles] Error : %s", err.Error())
			utilities.HandleError(ctx, err, validation)
		}
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[RoleController][GetRoles] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}

	result := utilities.SuccessResponseGenerator("Roles deleted successfully", http.StatusOK, "")
	// Data deleted successfully
	ctx.JSON(http.StatusOK, result)
}

// CreateRole handles the creation of a role.
func (role *RoleController) CreateRole(ctx *gin.Context) {

	var (
		req           entities.Role
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		errMap        = utilities.NewErrorMap()
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)

	validation.HelpLink = consts.ErrorHelpLink
	log.Info("[RoleController][CreateRole] Processing DeleteRoles request")
	if err := ctx.BindJSON(&req); err != nil {
		log.Errorf("[RoleController][CreateRole], Invalid query params, Error : %s", err.Error())
		result := utilities.ErrorResponseGenerator(consts.BindingError, http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.CreateRoleIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	errMap, err := role.useCases.CreateRole(ctxt, req, validation, errMap)

	if err != nil {
		log.Errorf("[RoleController][CreateRole] Error : %s", err.Error())
		utilities.HandleError(ctx, err, validation)
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[RoleController][CreateRole] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}

	result := utilities.SuccessResponseGenerator("Role created successfully", http.StatusCreated, "")
	// Data  added successfully
	ctx.JSON(http.StatusCreated, result)
}

// UpdateRole handles the updating of a role.
func (role *RoleController) UpdateRole(ctx *gin.Context) {

	var (
		req           entities.Role
		ctxt          = ctx.Request.Context()
		log           = logger.Log().WithContext(ctxt)
		errMap        = utilities.NewErrorMap()
		validation    entities.Validation
		endpointURL   string
		contextStatus bool
	)

	validation.HelpLink = consts.ErrorHelpLink
	log.Info("[RoleController][UpdateRole] Processing UpdateRole request")

	if err := ctx.BindJSON(&req); err != nil {
		log.Errorf("[RoleController][UpdateRole], Invalid json data, Error : %s", err.Error())
		result := utilities.ErrorResponseGenerator(consts.BindingError, http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	validation.Method, endpointURL = strings.ToLower(ctx.Request.Method), ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	validation.ContextError, contextStatus = utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
	//check isEndpointExists and contextStatus
	utilities.IsEndpointExists(ctx, isEndpointExists, validation.ContextError, validation.HelpLink, contextStatus, consts.UpdateRoleIdentifier)
	validation.Endpoint = utils.GetEndPoints(contextEndpoints, endpointURL, validation.Method)

	validation.ID = ctx.Param(consts.IDKey)

	errMap, err := role.useCases.UpdateRole(ctxt, req, validation, errMap)

	if err != nil {
		if err == consts.ErrNotExist {
			log.Errorf("[RoleController][UpdateRole] Error : %s", err.Error())
			utilities.HandleNotFoundError(ctx, consts.RoleNotExist, http.StatusNotFound, consts.NotFound)
			return
		} else {
			log.Errorf("[RoleController][UpdateRole] Error : %s", err.Error())
			utilities.HandleError(ctx, err, validation)
		}
		return
	}

	if len(errMap) != 0 {
		log.Errorf("[RoleController][UpdateRole] ValidationError")
		utilities.HandleValidationError(ctx, errMap, validation)
		return
	}
	log.Info("[RoleController][UpdateRole] Roles fetched successfully")
	result := utilities.SuccessResponseGenerator("Role updated successfully", http.StatusOK, "")
	// Data  added successfully
	ctx.JSON(http.StatusOK, result)
}
