// Package controllers contains the controller implementations for partner-related HTTP endpoints.
package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"partner/internal/consts"
	"partner/internal/entities"
	"partner/internal/usecases"
	"partner/utilities"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/tuneverse/toolkit/core/activitylog"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/core/version"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils"
)

// PartnerController is responsible for handling partner-related HTTP requests.
type PartnerController struct {
	router      *gin.RouterGroup
	Cfg         *entities.EnvConfig
	useCases    usecases.PartnerUseCaseImply
	activitylog *activitylog.ActivityLogOptions
}

// NewPartnerController creates a new PartnerController instance with the given router and partner use case.
func NewPartnerController(router *gin.RouterGroup, partnerUseCase usecases.PartnerUseCaseImply, cfg *entities.EnvConfig, activitylog *activitylog.ActivityLogOptions) *PartnerController {
	return &PartnerController{
		router:      router,
		Cfg:         cfg,
		useCases:    partnerUseCase,
		activitylog: activitylog,
	}
}

// InitRoutes initializes routes for partner-related endpoints.
func (partner *PartnerController) InitRoutes() {

	// health handler
	partner.router.GET("/:version/health", func(ctx *gin.Context) {
		version.RenderHandler(ctx, partner, "HealthHandler")
	})
	// Create partner Oauth-credentials
	partner.router.GET("/:version/partners/:partner_id/oauth-credentials", func(ctx *gin.Context) {
		version.RenderHandler(ctx, partner, "GetPartnerOauthCredential")
	})
	// Update partner
	partner.router.PATCH("/:version/partners/:partner_id", func(ctx *gin.Context) {
		version.RenderHandler(ctx, partner, "UpdatePartner")
	})

	// Create partner
	partner.router.POST("/:version/partners", func(ctx *gin.Context) {
		version.RenderHandler(ctx, partner, "CreatePartner")
	})

	// Get partner by id
	partner.router.GET("/:version/partners/:partner_id", func(ctx *gin.Context) {
		version.RenderHandler(ctx, partner, "GetPartnerById")
	})

	// Delete partner
	partner.router.DELETE("/:version/partners/:partner_id", func(ctx *gin.Context) {
		version.RenderHandler(ctx, partner, "DeletePartner")
	})

	//Get all partners
	partner.router.GET("/:version/partners", func(ctx *gin.Context) {
		version.RenderHandler(ctx, partner, "GetAllPartners")
	})

	// get terms and conditions of a partner
	partner.router.GET("/:version/terms-and-conditions", func(ctx *gin.Context) {
		version.RenderHandler(ctx, partner, "GetAllTermsAndConditions")
	})

	// update terms and conditions
	partner.router.PATCH("/:version/partners/:partner_id/terms-and-conditions", func(ctx *gin.Context) {
		version.RenderHandler(ctx, partner, "UpdateTermsAndConditions")
	})

	// get patner payment gateway details
	partner.router.GET("/:version/partners/:partner_id/payment-gateways", func(ctx *gin.Context) {
		version.RenderHandler(ctx, partner, "GetPartnerPaymentGateways")
	})

	// get patner stores
	partner.router.GET("/:version/partners/:partner_id/stores", func(ctx *gin.Context) {
		version.RenderHandler(ctx, partner, "GetPartnerStores")
	})
	// to check the partner exists in partner table
	partner.router.HEAD("/:version/partners/:partner_id", func(ctx *gin.Context) {
		version.RenderHandler(ctx, partner, "IsPartnerExists")
	})
	// Delete  partner genre language
	partner.router.DELETE("/:version/partners/:partner_id/genres/:genre_id", func(ctx *gin.Context) {
		version.RenderHandler(ctx, partner, "DeletePartnerGenreLanguage")
	})

	// Delete  partner artist role language
	partner.router.DELETE("/:version/partners/:partner_id/artist-role/:role_id", func(ctx *gin.Context) {
		version.RenderHandler(ctx, partner, "DeletePartnerArtistRoleLanguage")
	})
	// create partner store
	partner.router.PATCH("/:version/partners/:partner_id/stores", func(ctx *gin.Context) {
		version.RenderHandler(ctx, partner, "CreatePartnerStores")
	})

	partner.router.PATCH("/:version/partners/:partner_id/status", func(ctx *gin.Context) {
		version.RenderHandler(ctx, partner, "UpdatePartnerStatus")
	})

}

// HealthHandler handles the health check endpoint.
func (partner *PartnerController) HealthHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "server run with base version",
	})
}
func (partner *PartnerController) IsPartnerExists(ctx *gin.Context) {
	var log = logger.Log().WithContext(ctx)
	errMap := utilities.NewErrorMap()
	partnerId := ctx.Param(consts.PartnerIDKey)
	method := consts.IsPartnerExistMethod
	endpoint := consts.IsPartnerExistEndpoint

	// to check whether partner is valid and already exist in partner table
	partnerExists, err := partner.useCases.IsPartnerExists(ctx, partnerId, endpoint, method, errMap)
	if !partnerExists && len(errMap) != 0 {
		log.Error(consts.IsPartnerExistErrMsg, errMap)
		result := utilities.ErrorResponseGenerator(consts.PathParameterErrMsg, http.StatusNotFound, errMap)
		ctx.JSON(http.StatusNotFound,
			result,
		)
		return
	}
	if err != nil {
		log.Error(consts.IsPartnerExistErrMsg, err)
		result := utilities.ErrorResponseGenerator(consts.InternalServerErr, http.StatusInternalServerError, "")
		ctx.JSON(http.StatusInternalServerError, result)
		return
	}

	result := utilities.SuccessResponseGenerator(consts.IsPartnerExistSuccessMsg, http.StatusOK, "")
	ctx.JSON(http.StatusOK, result)

}
func (partner *PartnerController) UpdatePartnerStatus(ctx *gin.Context) {
	var (
		log           = logger.Log().WithContext(ctx)
		partnerStatus entities.UpdatePartnerStatus
	)
	errMap := utilities.NewErrorMap()
	serviceCode := make(map[string]string)
	helpLink := consts.ErrorHelpLink
	loggerVar := logger.Log().WithContext(ctx.Request.Context())
	method := strings.ToLower(ctx.Request.Method)
	endpointURL := ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	contextError, errVal := utils.GetContext[map[string]interface{}](ctx, consts.ContextErrorResponses)
	if !errVal {
		loggerVar.Error(consts.UpdatePartnerStatusErrMsg, consts.ContextErrMsg)
		return
	}
	endpoint := utils.GetEndPoints(contextEndpoints, endpointURL, method)
	if !isEndpointExists {
		val, _, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		log.Errorf(consts.UpdatePartnerStatusErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")
		ctx.JSON(http.StatusBadRequest, result)
		return
	}
	partnerID := ctx.Param(consts.PartnerIDKey)
	// to check whether partner is valid and already exist in partner table
	partnerExists, err := partner.useCases.IsExists(ctx, consts.PartnerTable, consts.IDKey, partnerID)
	if !partnerExists {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, _, _ := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.UpdatePartnerStatusErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(consts.PathParameterErrMsg, http.StatusNotFound, val)
		ctx.JSON(http.StatusNotFound,
			result,
		)
		return
	}
	if err != nil {
		loggerVar.Errorf(consts.UpdatePartnerStatusErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		if !hasError {
			loggerVar.Error(consts.UpdatePartnerStatusErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		loggerVar.Error(consts.UpdatePartnerStatusErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

		ctx.JSON(http.StatusInternalServerError, result)
		return
	}
	if err := ctx.BindJSON(&partnerStatus); err != nil {
		log.Errorf(consts.UpdatePartnerStatusErrMsg, err.Error())
		result := utilities.ErrorResponseGenerator(consts.BindingErrorErrMsg, http.StatusBadRequest, consts.ParsingErrMsg)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}
	err = partner.useCases.UpdatePartnerStatus(ctx, partnerID, partnerStatus, endpoint, method, errMap)
	if err != nil {
		loggerVar.Errorf(consts.UpdatePartnerStatusErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		if !hasError {
			loggerVar.Error(consts.UpdatePartnerStatusErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		loggerVar.Error(consts.UpdatePartnerStatusErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

		ctx.JSON(http.StatusInternalServerError, result)
		return

	}
	result := utilities.SuccessResponseGenerator(consts.UpdatePartnerStatusSuccessMsg, http.StatusOK, "")
	ctx.JSON(http.StatusOK, result)
}

// function to get partner stores
func (partner *PartnerController) GetPartnerStores(ctx *gin.Context) {

	var log = logger.Log().WithContext(ctx)
	errMap := utilities.NewErrorMap()
	serviceCode := make(map[string]string)
	helpLink := consts.ErrorHelpLink
	loggerVar := logger.Log().WithContext(ctx.Request.Context())
	method := strings.ToLower(ctx.Request.Method)
	endpointURL := ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	contextError, errVal := utils.GetContext[map[string]interface{}](ctx, consts.ContextErrorResponses)
	if !errVal {
		loggerVar.Error(consts.GetPartnerStoresErrMsg, consts.ContextErrMsg)
		return
	}
	endpoint := utils.GetEndPoints(contextEndpoints, endpointURL, method)
	if !isEndpointExists {
		val, _, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		log.Errorf(consts.GetPartnerStoresErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")
		ctx.JSON(http.StatusBadRequest, result)
		return
	}
	partnerID := ctx.Param(consts.PartnerIDKey)
	// to check whether partner is valid and already exist in partner table
	partnerExists, err := partner.useCases.IsPartnerExists(ctx, partnerID, endpoint, method, errMap)
	if !partnerExists && len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, _, _ := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.GetPartnerStoresErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(consts.PathParameterErrMsg, http.StatusNotFound, val)
		ctx.JSON(http.StatusNotFound,
			result,
		)
		return
	}
	if err != nil {
		loggerVar.Errorf(consts.GetPartnerStoresErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		if !hasError {
			loggerVar.Error(consts.GetPartnerStoresErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		loggerVar.Error(consts.GetPartnerStoresErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

		ctx.JSON(http.StatusInternalServerError, result)
		return
	}
	data, err := partner.useCases.GetPartnerStores(ctx, partnerID, endpoint, method, errMap)
	if len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.GetPartnerStoresErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), val)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}
	if err != nil {
		switch {

		case errors.Is(err, consts.ErrMemberServiceConnectionLost):
			loggerVar.Error(consts.GetPartnerStoresErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrMemberServiceConnectionLost)
			ctx.JSON(http.StatusBadRequest, result)
			return

		case errors.Is(err, consts.ErrUtilityServiceConnectionLost):
			loggerVar.Error(consts.GetPartnerStoresErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrUtilityServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrSubscriptionServiceConnectionLost):
			loggerVar.Error(consts.GetPartnerStoresErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrSubscriptionServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrStoreServiceConnectionLost):
			loggerVar.Error(consts.GetPartnerStoresErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrStoreServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrOauthServiceConnectionLost):
			loggerVar.Error(consts.GetPartnerStoresErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrOauthServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return

		default:
			loggerVar.Errorf(consts.GetPartnerStoresErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
			val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
				"", contextError, "", "", nil, helpLink)
			if !hasError {
				loggerVar.Error(consts.GetPartnerStoresErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
			}
			loggerVar.Error(consts.GetPartnerStoresErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
			result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

			ctx.JSON(http.StatusInternalServerError, result)
			return

		}
	}
	result := utilities.SuccessResponseGenerator(consts.GetPartnerStoresSuccessMsg, http.StatusOK, data)
	ctx.JSON(http.StatusOK, result)

}

// Get Partner Payment Gateways
func (partner *PartnerController) GetPartnerPaymentGateways(ctx *gin.Context) {
	var log = logger.Log().WithContext(ctx)
	errMap := utilities.NewErrorMap()
	serviceCode := make(map[string]string)
	helpLink := consts.ErrorHelpLink
	loggerVar := logger.Log().WithContext(ctx.Request.Context())
	method := strings.ToLower(ctx.Request.Method)
	endpointURL := ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	contextError, errVal := utils.GetContext[map[string]interface{}](ctx, consts.ContextErrorResponses)
	if !errVal {
		loggerVar.Error(consts.GetPartnerPaymentGatewaysErrMsg, consts.ContextErrMsg)
		return
	}
	endpoint := utils.GetEndPoints(contextEndpoints, endpointURL, method)
	if !isEndpointExists {
		val, _, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		log.Errorf(consts.GetPartnerPaymentGatewaysErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")
		ctx.JSON(http.StatusBadRequest, result)
		return
	}
	partnerId := ctx.Param(consts.PartnerIDKey)

	// to check whether partner is valid and already exist in partner table
	partnerExists, err := partner.useCases.IsPartnerExists(ctx, partnerId, endpoint, method, errMap)
	if !partnerExists && len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, _, _ := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.GetPartnerPaymentGatewaysErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(consts.PathParameterErrMsg, http.StatusNotFound, val)
		ctx.JSON(http.StatusNotFound,
			result,
		)
		return
	}
	if err != nil {
		loggerVar.Errorf(consts.GetPartnerPaymentGatewaysErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		if !hasError {
			loggerVar.Error(consts.GetPartnerPaymentGatewaysErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		loggerVar.Error(consts.GetPartnerPaymentGatewaysErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

		ctx.JSON(http.StatusInternalServerError, result)
		return
	}

	// function to get partner payment gateway details
	data, err := partner.useCases.GetPartnerPaymentGateways(ctx, partnerId, endpoint, method, errMap)
	if len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.GetPartnerPaymentGatewaysErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), val)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}
	if err != nil {
		switch {

		case errors.Is(err, consts.ErrMemberServiceConnectionLost):
			loggerVar.Error(consts.GetPartnerPaymentGatewaysErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrMemberServiceConnectionLost)
			ctx.JSON(http.StatusBadRequest, result)
			return

		case errors.Is(err, consts.ErrUtilityServiceConnectionLost):
			loggerVar.Error(consts.GetPartnerPaymentGatewaysErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrUtilityServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrSubscriptionServiceConnectionLost):
			loggerVar.Error(consts.GetPartnerPaymentGatewaysErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrSubscriptionServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrStoreServiceConnectionLost):
			loggerVar.Error(consts.GetPartnerPaymentGatewaysErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrStoreServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrOauthServiceConnectionLost):
			loggerVar.Error(consts.GetPartnerPaymentGatewaysErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrOauthServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return

		default:
			loggerVar.Errorf(consts.GetPartnerPaymentGatewaysErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
			val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
				"", contextError, "", "", nil, helpLink)
			if !hasError {
				loggerVar.Error(consts.GetPartnerPaymentGatewaysErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
			}
			loggerVar.Error(consts.GetPartnerPaymentGatewaysErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
			result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

			ctx.JSON(http.StatusInternalServerError, result)
			return

		}
	}

	result := utilities.SuccessResponseGenerator(consts.GetPartnerPaymentGatewaysSuccessMsg, http.StatusOK, data)
	ctx.JSON(http.StatusOK, result)

}

// func tion to Get All Partners
func (partner *PartnerController) GetAllPartners(ctx *gin.Context) {

	var (
		params entities.Params
		output entities.ResponseData
		ctxt   = ctx.Request.Context()
		log    = logger.Log().WithContext(ctxt)
	)
	errMap := utilities.NewErrorMap()
	serviceCode := make(map[string]string)
	helpLink := consts.ErrorHelpLink
	loggerVar := logger.Log().WithContext(ctx.Request.Context())

	if err := ctx.ShouldBindQuery(&params); err != nil {
		logger.Log().WithContext(ctx).Errorf(consts.GetAllPartnersErrMsg, err.Error())
		result := utilities.ErrorResponseGenerator(consts.BindingErrorErrMsg, http.StatusBadRequest, consts.ParsingErrMsg)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}
	method := strings.ToLower(ctx.Request.Method)
	endpointURL := ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)

	contextError, errVal := utils.GetContext[map[string]interface{}](ctx, consts.ContextErrorResponses)
	if !errVal {
		loggerVar.Error(consts.GetAllPartnersErrMsg, consts.ContextErrMsg)
		return
	}
	endpoint := utils.GetEndPoints(contextEndpoints, endpointURL, method)
	if !isEndpointExists {
		val, _, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		log.Errorf(consts.GetAllPartnersErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")
		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	data, metadata, errMap, err := partner.useCases.GetAllPartners(ctx, params, endpoint, method, errMap)
	for key, value := range errMap {
		serviceCode[key] = value.Code
	}
	// validation error check
	if len(errMap) != 0 {
		fields := utils.FieldMapping(errMap)
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.GetAllPartnersErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), val)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}
	if err != nil {
		switch {

		case errors.Is(err, consts.ErrMemberServiceConnectionLost):
			loggerVar.Error(consts.GetAllPartnersErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrMemberServiceConnectionLost)
			ctx.JSON(http.StatusBadRequest, result)
			return

		case errors.Is(err, consts.ErrUtilityServiceConnectionLost):
			loggerVar.Error(consts.GetAllPartnersErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrUtilityServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrSubscriptionServiceConnectionLost):
			loggerVar.Error(consts.GetAllPartnersErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrSubscriptionServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrStoreServiceConnectionLost):
			loggerVar.Error(consts.GetAllPartnersErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrStoreServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrOauthServiceConnectionLost):
			loggerVar.Error(consts.GetAllPartnersErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrOauthServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrMaximumRequest):
			loggerVar.Error(consts.GetAllPartnersErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.MaximumRequestErrMsg, http.StatusTooManyRequests, consts.ErrMaximumRequest)
			ctx.JSON(http.StatusTooManyRequests, result)
			return

		default:
			loggerVar.Errorf(consts.GetAllPartnersErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
			val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
				"", contextError, "", "", nil, helpLink)
			if !hasError {
				loggerVar.Error(consts.GetAllPartnersErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
			}
			loggerVar.Error(consts.GetAllPartnersErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
			result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

			ctx.JSON(http.StatusInternalServerError, result)
			return

		}
	}

	if metadata.Total == 0 {
		result := utilities.SuccessResponseGenerator("", http.StatusNoContent, "")
		ctx.JSON(http.StatusNoContent, result)
		return
	}

	if reflect.DeepEqual(metadata, models.MetaData{}) {
		result := utilities.SuccessResponseGenerator(consts.GetAllPartnersSuccessMsg, http.StatusOK, data)
		ctx.JSON(http.StatusOK, result)
		return
	}
	output.Data = data
	output.Metadata = metadata
	result := utilities.SuccessResponseGenerator(consts.GetAllPartnersSuccessMsg, http.StatusOK, output)
	ctx.JSON(http.StatusOK, result)
}

// function to get partner oauth credential
func (partner *PartnerController) GetPartnerOauthCredential(ctx *gin.Context) {
	ctxt := ctx.Request.Context()
	var (
		log         = logger.Log().WithContext(ctxt)
		oauthHeader entities.PartnerOAuthHeader
	)
	errMap := utilities.NewErrorMap()
	serviceCode := make(map[string]string)
	helpLink := consts.ErrorHelpLink
	loggerVar := logger.Log().WithContext(ctx.Request.Context())
	method := strings.ToLower(ctx.Request.Method)
	endpointURL := ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)

	contextError, errVal := utils.GetContext[map[string]interface{}](ctx, consts.ContextErrorResponses)
	if !errVal {
		loggerVar.Errorf(consts.GetPartnerOauthCredentialErrMsg, consts.ContextErrMsg)
		return
	}
	endpoint := utils.GetEndPoints(contextEndpoints, endpointURL, method)
	if !isEndpointExists {
		val, _, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		log.Errorf(consts.GetPartnerOauthCredentialErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")
		ctx.JSON(http.StatusBadRequest, result)
		return
	}
	partnerId := ctx.Param(consts.PartnerIDKey)
	// to check whether partner is valid and already exist in partner table
	partnerExists, err := partner.useCases.IsPartnerExists(ctx, partnerId, endpoint, method, errMap)
	if !partnerExists && len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, _, _ := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.GetPartnerOauthCredentialErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(consts.PathParameterErrMsg, http.StatusNotFound, "")
		ctx.JSON(http.StatusNotFound,
			result,
		)
		return
	}

	if err != nil {
		loggerVar.Errorf(consts.GetPartnerOauthCredentialErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		if !hasError {
			loggerVar.Error(consts.GetPartnerOauthCredentialErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		loggerVar.Error(consts.GetPartnerOauthCredentialErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

		ctx.JSON(http.StatusInternalServerError, result)
		return
	}
	err = ctx.BindHeader(&oauthHeader)
	if err != nil {
		log.Errorf(consts.GetPartnerOauthCredentialErrMsg, err)
		result := utilities.ErrorResponseGenerator(consts.BindingErrorErrMsg, http.StatusBadRequest, consts.InvalidHeaderData)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}

	data, err := partner.useCases.GetPartnerOauthCredential(ctxt, partnerId, oauthHeader, endpoint, method, errMap)
	// check validation errors
	if len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.GetPartnerOauthCredentialErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), val)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}
	if err != nil {
		switch {

		case errors.Is(err, consts.ErrMemberServiceConnectionLost):
			loggerVar.Error(consts.GetPartnerOauthCredentialErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrMemberServiceConnectionLost)
			ctx.JSON(http.StatusBadRequest, result)
			return

		case errors.Is(err, consts.ErrUtilityServiceConnectionLost):
			loggerVar.Error(consts.GetPartnerOauthCredentialErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrUtilityServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrSubscriptionServiceConnectionLost):
			loggerVar.Error(consts.GetPartnerOauthCredentialErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrSubscriptionServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrStoreServiceConnectionLost):
			loggerVar.Error(consts.GetPartnerOauthCredentialErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrStoreServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrOauthServiceConnectionLost):
			loggerVar.Error(consts.GetPartnerOauthCredentialErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrOauthServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrMaximumRequest):
			loggerVar.Error(consts.GetPartnerOauthCredentialErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.MaximumRequestErrMsg, http.StatusTooManyRequests, consts.ErrMaximumRequest)
			ctx.JSON(http.StatusTooManyRequests, result)
			return
		case errors.Is(err, sql.ErrNoRows):
			loggerVar.Errorf(consts.GetPartnerOauthCredentialErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
			val, hasError, _, _ := utils.ParseFields(ctx, consts.InternalServerErr,
				"", contextError, "", "", nil, helpLink)
			if !hasError {
				loggerVar.Error(consts.GetPartnerOauthCredentialErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
			}
			loggerVar.Error(consts.GetPartnerOauthCredentialErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
			result := utilities.ErrorResponseGenerator(consts.NotFoundErrMsg, http.StatusNotFound, val)
			ctx.JSON(http.StatusNotFound,
				result,
			)
			return

		default:
			loggerVar.Errorf(consts.GetPartnerOauthCredentialErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
			val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
				"", contextError, "", "", nil, helpLink)
			if !hasError {
				loggerVar.Error(consts.GetPartnerOauthCredentialErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
			}
			loggerVar.Error(consts.GetPartnerOauthCredentialErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
			result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

			ctx.JSON(http.StatusInternalServerError, result)
			return

		}
	}

	result := utilities.SuccessResponseGenerator(consts.GetPartnerOauthCredentialSuccessMsg, http.StatusOK, data)
	ctx.JSON(http.StatusOK, result)
}

// function to create a partner
func (partner *PartnerController) CreatePartner(ctx *gin.Context) {

	var (
		ctxt       = ctx.Request.Context()
		newPartner entities.Partner
		log        = logger.Log().WithContext(ctxt)
	)

	errMap := utilities.NewErrorMap()
	serviceCode := make(map[string]string)
	helpLink := consts.ErrorHelpLink
	loggerVar := logger.Log().WithContext(ctx.Request.Context())

	if err := ctx.BindJSON(&newPartner); err != nil {
		log.Errorf(consts.CreatePartnerErrMsg, err.Error())
		result := utilities.ErrorResponseGenerator(consts.BindingErrorErrMsg, http.StatusBadRequest, consts.ParsingErrMsg)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}

	method := strings.ToLower(ctx.Request.Method)
	endpointURL := ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)

	contextError, errVal := utils.GetContext[map[string]interface{}](ctx, consts.ContextErrorResponses)
	if !errVal {
		loggerVar.Error(consts.CreatePartnerErrMsg, consts.ContextErrMsg)
		return
	}
	endpoint := utils.GetEndPoints(contextEndpoints, endpointURL, method)
	if !isEndpointExists {
		val, _, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		log.Errorf(consts.CreatePartnerErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")
		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	// function to create a partner
	errMap, _, err := partner.useCases.CreatePartner(ctxt, newPartner, endpoint, method, errMap)

	// check validation errors
	if len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.CreatePartnerErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), val)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}
	if err != nil {
		switch {

		case errors.Is(err, consts.ErrMemberServiceConnectionLost):
			loggerVar.Error(consts.CreatePartnerErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrMemberServiceConnectionLost)
			ctx.JSON(http.StatusBadRequest, result)
			return

		case errors.Is(err, consts.ErrUtilityServiceConnectionLost):
			loggerVar.Error(consts.CreatePartnerErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrUtilityServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrSubscriptionServiceConnectionLost):
			loggerVar.Error(consts.CreatePartnerErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrSubscriptionServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrStoreServiceConnectionLost):
			loggerVar.Error(consts.CreatePartnerErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrStoreServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrOauthServiceConnectionLost):
			loggerVar.Error(consts.CreatePartnerErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrOauthServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrMaximumRequest):
			loggerVar.Error(consts.CreatePartnerErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.MaximumRequestErrMsg, http.StatusTooManyRequests, consts.ErrMaximumRequest)
			ctx.JSON(http.StatusTooManyRequests, result)
			return

		default:
			loggerVar.Errorf(consts.CreatePartnerErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
			val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
				"", contextError, "", "", nil, helpLink)
			if !hasError {
				loggerVar.Error(consts.CreatePartnerErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
			}
			loggerVar.Error(consts.CreatePartnerErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
			result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

			ctx.JSON(http.StatusInternalServerError, result)
			return
		}
	}

	result := utilities.SuccessResponseGenerator(consts.CreatePartnerSuccessMsg, http.StatusCreated, "")
	ctx.JSON(http.StatusCreated, result)
}

// Get PartnerById retrieves partner details by ID and responds with the partner data.
func (partner *PartnerController) GetPartnerById(ctx *gin.Context) {
	var (
		ctxt   = ctx.Request.Context()
		log    = logger.Log().WithContext(ctxt)
		params entities.QueryParams
		output entities.ResponseData
	)
	errMap := utilities.NewErrorMap()
	serviceCode := make(map[string]string)
	helpLink := consts.ErrorHelpLink
	loggerVar := logger.Log().WithContext(ctx.Request.Context())

	method := strings.ToLower(ctx.Request.Method)
	endpointURL := ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)

	contextError, errVal := utils.GetContext[map[string]interface{}](ctx, consts.ContextErrorResponses)
	if !errVal {
		loggerVar.Error(consts.GetPartnerByIdErrMsg, consts.ContextErrMsg)
		return
	}
	endpoint := utils.GetEndPoints(contextEndpoints, endpointURL, method)
	if !isEndpointExists {
		val, _, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		log.Errorf(consts.GetPartnerByIdErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")
		ctx.JSON(http.StatusBadRequest, result)
		return
	}
	partnerId := ctx.Param(consts.PartnerIDKey)
	// to check whether partner is valid and already exist in partner table
	partnerExists, err := partner.useCases.IsPartnerExists(ctx, partnerId, endpoint, method, errMap)
	if !partnerExists && len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, _, _ := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.GetPartnerByIdErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(consts.PathParameterErrMsg, http.StatusNotFound, val)
		ctx.JSON(http.StatusNotFound,
			result,
		)
		return
	}
	if err != nil {
		loggerVar.Errorf(consts.GetPartnerByIdErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		if !hasError {
			loggerVar.Error(consts.GetPartnerByIdErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		loggerVar.Error(consts.GetPartnerByIdErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

		ctx.JSON(http.StatusInternalServerError, result)
		return
	}
	err = ctx.BindQuery(&params)
	if err != nil {
		log.Errorf(consts.GetPartnerByIdErrMsg, err.Error())
		result := utilities.ErrorResponseGenerator(consts.QueryBindingErrorErrMsg, http.StatusBadRequest, consts.QueryParsingError)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}

	data, metadata, errMap, err := partner.useCases.GetPartnerById(ctx, params, partnerId, endpoint, method, errMap)
	// check validation error
	if len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.GetPartnerByIdErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), val)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}
	if err != nil {
		switch err {

		case consts.ErrMemberServiceConnectionLost:
			loggerVar.Error(consts.GetPartnerByIdErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrMemberServiceConnectionLost)
			ctx.JSON(http.StatusBadRequest, result)
			return

		case consts.ErrUtilityServiceConnectionLost:
			loggerVar.Error(consts.GetPartnerByIdErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrUtilityServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case consts.ErrSubscriptionServiceConnectionLost:
			loggerVar.Error(consts.GetPartnerByIdErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrSubscriptionServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case consts.ErrStoreServiceConnectionLost:
			loggerVar.Error(consts.GetPartnerByIdErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrStoreServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case consts.ErrOauthServiceConnectionLost:
			loggerVar.Error(consts.GetPartnerByIdErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrOauthServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return

		default:
			loggerVar.Errorf(consts.GetPartnerByIdErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
			val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
				"", contextError, "", "", nil, helpLink)
			if !hasError {
				loggerVar.Error(consts.GetPartnerByIdErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
			}
			loggerVar.Error(consts.GetPartnerByIdErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
			result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

			ctx.JSON(http.StatusInternalServerError, result)
			return

		}
	}
	if reflect.DeepEqual(metadata, models.MetaData{}) {
		result := utilities.SuccessResponseGenerator(consts.GetPartnerByIdSuccessMsg, http.StatusOK, data)
		ctx.JSON(http.StatusOK, result)
	}
	if metadata.Total == 0 {
		result := utilities.SuccessResponseGenerator("", http.StatusNoContent, "")
		ctx.JSON(http.StatusNoContent, result)
		return
	}

	output.Data = data
	output.Metadata = metadata
	ctx.JSON(http.StatusOK, utilities.SuccessResponseGenerator(consts.GetPartnerByIdSuccessMsg, http.StatusOK, output))

}

// function to  get all terms and conditions of a partner
func (partner *PartnerController) GetAllTermsAndConditions(ctx *gin.Context) {

	var (
		ctxt = ctx.Request.Context()
		log  = logger.Log().WithContext(ctxt)
	)
	errMap := utilities.NewErrorMap()
	serviceCode := make(map[string]string)
	helpLink := consts.ErrorHelpLink
	loggerVar := logger.Log().WithContext(ctx.Request.Context())

	method := strings.ToLower(ctx.Request.Method)
	endpointURL := ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)

	contextError, errVal := utils.GetContext[map[string]interface{}](ctx, consts.ContextErrorResponses)
	if !errVal {
		loggerVar.Error(consts.GetAllTermsAndConditionErrMsg, consts.ContextErrMsg)
		return
	}
	endpoint := utils.GetEndPoints(contextEndpoints, endpointURL, method)
	if !isEndpointExists {
		val, _, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		log.Errorf(consts.GetAllTermsAndConditionErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), val)
		ctx.JSON(http.StatusBadRequest, result)
		return
	}
	partnerId := ctx.Query(consts.PartnerIDKey)
	// to check whether partner is valid and already exist in partner table
	partnerExists, err := partner.useCases.IsPartnerExists(ctx, partnerId, endpoint, method, errMap)
	if !partnerExists && len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, _, _ := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.GetAllTermsAndConditionErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(consts.PathParameterErrMsg, http.StatusNotFound, val)
		ctx.JSON(http.StatusNotFound,
			result,
		)
		return
	}

	if err != nil {
		loggerVar.Errorf(consts.GetAllTermsAndConditionErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		if !hasError {
			loggerVar.Error(consts.GetAllTermsAndConditionErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		loggerVar.Error(consts.GetAllTermsAndConditionErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

		ctx.JSON(http.StatusInternalServerError, result)
		return
	}

	data, err := partner.useCases.GetAllTermsAndConditions(ctx, partnerId, endpoint, method, errMap)
	// check validation error
	if len(errMap) != 0 {
		fields := utils.FieldMapping(errMap)
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.GetAllTermsAndConditionErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), val)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}
	if err != nil {
		switch {

		case errors.Is(err, consts.ErrMemberServiceConnectionLost):
			loggerVar.Error(consts.GetAllTermsAndConditionErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrMemberServiceConnectionLost)
			ctx.JSON(http.StatusBadRequest, result)
			return

		case errors.Is(err, consts.ErrUtilityServiceConnectionLost):
			loggerVar.Error(consts.GetAllTermsAndConditionErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrUtilityServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrSubscriptionServiceConnectionLost):
			loggerVar.Error(consts.GetAllTermsAndConditionErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrSubscriptionServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrStoreServiceConnectionLost):
			loggerVar.Error(consts.GetAllTermsAndConditionErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrStoreServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrOauthServiceConnectionLost):
			loggerVar.Error(consts.GetAllTermsAndConditionErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrOauthServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return

		default:
			loggerVar.Errorf(consts.GetAllTermsAndConditionErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
			val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
				"", contextError, "", "", nil, helpLink)
			if !hasError {
				loggerVar.Error(consts.GetAllTermsAndConditionErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
			}
			loggerVar.Error(consts.GetAllTermsAndConditionErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
			result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

			ctx.JSON(http.StatusInternalServerError, result)
			return

		}
	}

	result := utilities.SuccessResponseGenerator(consts.GetAllTermsAndConditionSuccessMsg, http.StatusOK, data)
	ctx.JSON(http.StatusOK, result)
}

// function to update terms and conditions
func (partner *PartnerController) UpdateTermsAndConditions(ctx *gin.Context) {
	var (
		log                    = logger.Log().WithContext(ctx)
		termsAndConditionsData entities.UpdateTermsAndConditions
	)

	errMap := utilities.NewErrorMap()
	serviceCode := make(map[string]string)
	helpLink := consts.ErrorHelpLink
	loggerVar := logger.Log().WithContext(ctx.Request.Context())

	method := strings.ToLower(ctx.Request.Method)
	endpointURL := ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)

	contextError, errVal := utils.GetContext[map[string]interface{}](ctx, consts.ContextErrorResponses)
	if !errVal {
		loggerVar.Error(consts.UpdateTermsAndConditionsErrMsg, consts.ContextErrMsg)
		return
	}
	endpoint := utils.GetEndPoints(contextEndpoints, endpointURL, method)
	if !isEndpointExists {
		val, _, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		log.Errorf(consts.UpdateTermsAndConditionsErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")
		ctx.JSON(http.StatusBadRequest, result)
		return
	}

	partnerId := ctx.Param(consts.PartnerIDKey)
	// to check whether partner is valid and already exist in partner table
	partnerExists, err := partner.useCases.IsPartnerExists(ctx, partnerId, endpoint, method, errMap)
	if !partnerExists && len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, _, _ := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.UpdateTermsAndConditionsErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(consts.PathParameterErrMsg, http.StatusNotFound, val)
		ctx.JSON(http.StatusNotFound,
			result,
		)
		return
	}

	if err != nil {
		loggerVar.Errorf(consts.UpdateTermsAndConditionsErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		if !hasError {
			loggerVar.Error(consts.UpdateTermsAndConditionsErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		loggerVar.Error(consts.UpdateTermsAndConditionsErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

		ctx.JSON(http.StatusInternalServerError, result)
		return
	}
	memberIDStr := ctx.Request.Header.Get(consts.MemberIdKey)
	memberID, err := uuid.Parse(memberIDStr)
	if err != nil {
		log.Errorf(consts.UpdateTermsAndConditionsErrMsg, err.Error())
	}
	err = partner.useCases.IsMemberExists(ctx, partnerId, memberIDStr, endpoint, method, errMap)
	if len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.UpdateTermsAndConditionsErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), val)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}
	if err != nil {
		loggerVar.Errorf(consts.UpdateTermsAndConditionsErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		if !hasError {
			loggerVar.Error(consts.UpdateTermsAndConditionsErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		loggerVar.Error(consts.UpdateTermsAndConditionsErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

		ctx.JSON(http.StatusInternalServerError, result)
		return

	}

	err = ctx.BindJSON(&termsAndConditionsData)
	if err != nil {
		log.Errorf(consts.UpdateTermsAndConditionsErrMsg, err.Error())
		result := utilities.ErrorResponseGenerator(consts.BindingErrorErrMsg, http.StatusBadRequest, consts.ParsingErrMsg)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}

	// function to update terms and conditions
	errMap, err = partner.useCases.UpdateTermsAndConditions(ctx, partnerId, memberID, termsAndConditionsData, endpoint, method, errMap)
	if len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.UpdateTermsAndConditionsErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), val)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}
	if err != nil {
		switch {

		case errors.Is(err, consts.ErrMemberServiceConnectionLost):
			loggerVar.Error(consts.UpdateTermsAndConditionsErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrMemberServiceConnectionLost)
			ctx.JSON(http.StatusBadRequest, result)
			return

		case errors.Is(err, consts.ErrUtilityServiceConnectionLost):
			loggerVar.Error(consts.UpdateTermsAndConditionsErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrUtilityServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrSubscriptionServiceConnectionLost):
			loggerVar.Error(consts.UpdateTermsAndConditionsErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrSubscriptionServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrStoreServiceConnectionLost):
			loggerVar.Error(consts.UpdateTermsAndConditionsErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrStoreServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrOauthServiceConnectionLost):
			loggerVar.Error(consts.UpdateTermsAndConditionsErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrOauthServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return

		default:
			loggerVar.Errorf(consts.UpdateTermsAndConditionsErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
			val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
				"", contextError, "", "", nil, helpLink)
			if !hasError {
				loggerVar.Error(consts.UpdateTermsAndConditionsErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
			}
			loggerVar.Error(consts.UpdateTermsAndConditionsErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
			result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

			ctx.JSON(http.StatusInternalServerError, result)
			return

		}
	}

	// ***************Activity code starts here********************
	name, err := partner.useCases.GetPartnerName(ctx, partnerId)
	if err != nil {
		log.Errorf(consts.UpdateTermsAndConditionsErrMsg, err.Error())
	}

	data := map[string]interface{}{
		consts.PartnerBaseUrlKey: consts.PartnerServiceURL,
		consts.PartnerIDKey:      partnerId,
		consts.PartnerNameKey:    name,
		consts.DateTimeKey:       time.Now(),
	}

	_, ok := data[consts.PartnerIDKey].(string)
	if ok {
		activityDetails := utilities.NewActivityLog(memberIDStr, consts.PartnerTermsAndConditionsUpdatedActvityLogKey, data)

		resp, err := partner.activitylog.Log(activityDetails)
		if err != nil {
			logger.Log().WithContext(ctx.Request.Context()).Error(consts.ActivityLogErrMsg, err.Error(), resp)

		}
	}

	// ***************Activity code ends here********************
	result := utilities.SuccessResponseGenerator(consts.UpdateTermsAndConditionsSuccessMsg, http.StatusOK, "")
	ctx.JSON(http.StatusOK, result)

}

// function to Update a Partner
func (partner *PartnerController) UpdatePartner(ctx *gin.Context) {

	var (
		partnerData entities.PartnerProperties
		ctxt        = ctx.Request.Context()
		log         = logger.Log().WithContext(ctxt)
	)
	errMap := utilities.NewErrorMap()
	serviceCode := make(map[string]string)
	helpLink := consts.ErrorHelpLink
	loggerVar := logger.Log().WithContext(ctx.Request.Context())

	method := strings.ToLower(ctx.Request.Method)
	endpointURL := ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)

	contextError, errVal := utils.GetContext[map[string]interface{}](ctx, consts.ContextErrorResponses)
	if !errVal {
		loggerVar.Error(consts.UpdatePartnerErrMsg, consts.ContextErrMsg)
		return
	}
	endpoint := utils.GetEndPoints(contextEndpoints, endpointURL, method)
	if !isEndpointExists {
		val, _, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		log.Errorf(consts.UpdatePartnerErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")
		ctx.JSON(http.StatusBadRequest, result)
		return
	}
	partnerId := ctx.Param(consts.PartnerIDKey)
	// to check whether partner is valid and already exist in partner table
	partnerExists, err := partner.useCases.IsPartnerExists(ctx, partnerId, endpoint, method, errMap)
	if !partnerExists && len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, _, _ := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.UpdatePartnerErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(consts.PathParameterErrMsg, http.StatusNotFound, val)
		ctx.JSON(http.StatusNotFound,
			result,
		)
		return
	}
	if err != nil {
		log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		if !hasError {
			loggerVar.Error(consts.UpdatePartnerErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		loggerVar.Error(consts.UpdatePartnerErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

		ctx.JSON(http.StatusInternalServerError, result)
		return

	}
	memberIDStr := ctx.Request.Header.Get(consts.MemberIdKey)
	memberID, err := uuid.Parse(memberIDStr)
	if err != nil {
		log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
	}
	err = partner.useCases.IsMemberExists(ctx, partnerId, memberIDStr, endpoint, method, errMap)
	if len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.UpdatePartnerErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), val)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}
	if err != nil {
		loggerVar.Errorf(consts.UpdatePartnerErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		if !hasError {
			loggerVar.Error(consts.UpdatePartnerErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		loggerVar.Error(consts.UpdatePartnerErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

		ctx.JSON(http.StatusInternalServerError, result)
		return

	}
	err = ctx.BindJSON(&partnerData)

	if err != nil {
		log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
		result := utilities.ErrorResponseGenerator(consts.BindingErrorErrMsg, http.StatusBadRequest, consts.ParsingErrMsg)
		ctx.JSON(http.StatusBadRequest,
			result,
		)

	}
	// function to update a partner data
	errMap, err = partner.useCases.UpdatePartner(ctx, partnerId, memberID, partnerData, endpoint, method, errMap)

	// validation error check
	if len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.UpdatePartnerErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), val)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}
	if err != nil {
		switch {

		case errors.Is(err, consts.ErrMemberServiceConnectionLost):
			loggerVar.Error(consts.UpdatePartnerErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrMemberServiceConnectionLost)
			ctx.JSON(http.StatusBadRequest, result)
			return

		case errors.Is(err, consts.ErrUtilityServiceConnectionLost):
			loggerVar.Error(consts.UpdatePartnerErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrUtilityServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrSubscriptionServiceConnectionLost):
			loggerVar.Error(consts.UpdatePartnerErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrSubscriptionServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrStoreServiceConnectionLost):
			loggerVar.Error(consts.UpdatePartnerErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrStoreServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrOauthServiceConnectionLost):
			loggerVar.Error(consts.UpdatePartnerErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrOauthServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return

		default:
			loggerVar.Errorf(consts.UpdatePartnerErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
			val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
				"", contextError, "", "", nil, helpLink)
			if !hasError {
				loggerVar.Error(consts.UpdatePartnerErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
			}
			loggerVar.Error(consts.UpdatePartnerErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
			result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

			ctx.JSON(http.StatusInternalServerError, result)
			return

		}
	}

	result := utilities.SuccessResponseGenerator(consts.UpdatePartnerSuccessMsg, http.StatusOK, "")

	// ***************Activity code starts here********************
	name, err := partner.useCases.GetPartnerName(ctx, partnerId)
	if err != nil {
		log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
	}

	data := map[string]interface{}{
		consts.PartnerBaseUrlKey: consts.PartnerServiceURL,
		consts.PartnerIDKey:      partnerId,
		consts.PartnerNameKey:    name,
		consts.DateTimeKey:       time.Now(),
	}

	_, ok := data[consts.PartnerIDKey].(string)
	if ok {
		activityDetails := utilities.NewActivityLog(memberIDStr, consts.PartnerUpdatedActivityLogKey, data)

		resp, err := partner.activitylog.Log(activityDetails)
		if err != nil {
			logger.Log().WithContext(ctx.Request.Context()).Error(consts.ActivityLogErrMsg, err.Error(), resp)

		}
	}

	// ***************Activity code ends here********************

	ctx.JSON(http.StatusOK, result)
}

// function to delete a Partner
func (partner *PartnerController) DeletePartner(ctx *gin.Context) {

	var (
		ctxt = ctx.Request.Context()
		log  = logger.Log().WithContext(ctxt)
	)
	errMap := utilities.NewErrorMap()
	serviceCode := make(map[string]string)
	helpLink := consts.ErrorHelpLink
	loggerVar := logger.Log().WithContext(ctx.Request.Context())
	method := strings.ToLower(ctx.Request.Method)
	endpointURL := ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)

	// contextEndpoints, isEndpointExists := utils.GetContext[models.ResponseData](ctx, consts.ContextEndPoints)
	contextError, errVal := utils.GetContext[map[string]interface{}](ctx, consts.ContextErrorResponses)
	if !errVal {
		loggerVar.Error(consts.DeletePartnerErrMsg, consts.ContextErrMsg)
		return
	}
	endpoint := utils.GetEndPoints(contextEndpoints, endpointURL, method)
	if !isEndpointExists {
		val, _, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		log.Errorf(consts.DeletePartnerErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")
		ctx.JSON(http.StatusBadRequest, result)
		return
	}
	partnerId := ctx.Param(consts.PartnerIDKey)
	// to check whether partner is valid and already exist in partner table
	partnerExists, err := partner.useCases.IsPartnerExists(ctx, partnerId, endpoint, method, errMap)
	if !partnerExists && len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, _, _ := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.DeletePartnerErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(consts.PathParameterErrMsg, http.StatusNotFound, val)
		ctx.JSON(http.StatusNotFound,
			result,
		)
		return
	}

	if err != nil {
		log.Errorf(consts.DeletePartnerErrMsg, err.Error())
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		if !hasError {
			loggerVar.Errorf(consts.DeletePartnerErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		loggerVar.Error(consts.DeletePartnerErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

		ctx.JSON(http.StatusInternalServerError, result)
		return

	}

	// function to delete a partner
	err = partner.useCases.DeletePartner(ctx, partnerId, endpoint, method, errMap)

	if err != nil {
		loggerVar.Errorf(consts.DeletePartnerErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		if !hasError {
			loggerVar.Error(consts.DeletePartnerErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		loggerVar.Error(consts.DeletePartnerErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

		ctx.JSON(http.StatusInternalServerError, result)
		return

	}
	result := utilities.SuccessResponseGenerator(consts.DeletePartnerSuccessMsg, http.StatusOK, "")
	ctx.JSON(http.StatusOK, result)

}

// function to delete a Partner genre language
func (partner *PartnerController) DeletePartnerGenreLanguage(ctx *gin.Context) {

	var (
		ctxt = ctx.Request.Context()
		log  = logger.Log().WithContext(ctxt)
	)
	errMap := utilities.NewErrorMap()
	serviceCode := make(map[string]string)
	helpLink := consts.ErrorHelpLink
	loggerVar := logger.Log().WithContext(ctx.Request.Context())
	genreId := ctx.Param(consts.GenreIdKey)
	_, err := uuid.Parse(genreId)
	if err != nil {
		log.Errorf(consts.DeletePartnerGenreLanguageErrMsg, err.Error())
		result := utilities.ErrorResponseGenerator(consts.InvalidPartner, http.StatusBadRequest, consts.ParsingErrMsg)
		ctx.JSON(result.Code, result)
		return
	}

	method := strings.ToLower(ctx.Request.Method)
	endpointUrl := ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	contextError, errVal := utils.GetContext[map[string]interface{}](ctx, consts.ContextErrorResponses)
	if !errVal {
		loggerVar.Error(consts.DeletePartnerGenreLanguageErrMsg, consts.ContextErrMsg)
		return
	}
	endpoint := utils.GetEndPoints(contextEndpoints, endpointUrl, method)
	if !isEndpointExists {
		val, _, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		log.Errorf(consts.DeletePartnerGenreLanguageErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")
		ctx.JSON(result.Code, result)
		return
	}
	partnerId := ctx.Param(consts.PartnerIDKey)
	// to check whether partner is valid and already exist in partner table
	partnerExists, err := partner.useCases.IsPartnerExists(ctx, partnerId, endpoint, method, errMap)
	if !partnerExists && len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, _, _ := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.DeletePartnerGenreLanguageErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(consts.PathParameterErrMsg, http.StatusNotFound, val)
		ctx.JSON(http.StatusNotFound,
			result,
		)
		return
	}

	if err != nil {
		loggerVar.Errorf(consts.DeletePartnerGenreLanguageErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		if !hasError {
			loggerVar.Error(consts.DeletePartnerGenreLanguageErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		loggerVar.Error(consts.DeletePartnerGenreLanguageErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

		ctx.JSON(http.StatusInternalServerError, result)
		return
	}
	memberIDStr := ctx.Request.Header.Get(consts.MemberIdKey)
	err = partner.useCases.IsMemberExists(ctx, partnerId, memberIDStr, endpoint, method, errMap)
	if len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.DeletePartnerGenreLanguageErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), val)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}
	if err != nil {
		loggerVar.Errorf(consts.DeletePartnerGenreLanguageErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		if !hasError {
			loggerVar.Error(consts.DeletePartnerGenreLanguageErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		loggerVar.Error(consts.DeletePartnerGenreLanguageErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

		ctx.JSON(http.StatusInternalServerError, result)
		return

	}
	err = partner.useCases.DeletePartnerGenreLanguage(ctx, genreId, endpoint, method, errMap)

	if len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, _, _ := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.DeletePartnerGenreLanguageErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(consts.PathParameterErrMsg, http.StatusNotFound, val)
		ctx.JSON(http.StatusNotFound,
			result,
		)
		return
	}
	if err != nil {
		switch {

		case errors.Is(err, consts.ErrMemberServiceConnectionLost):
			loggerVar.Error(consts.DeletePartnerGenreLanguageErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrMemberServiceConnectionLost)
			ctx.JSON(http.StatusBadRequest, result)
			return

		case errors.Is(err, consts.ErrUtilityServiceConnectionLost):
			loggerVar.Error(consts.DeletePartnerGenreLanguageErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrUtilityServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrSubscriptionServiceConnectionLost):
			loggerVar.Error(consts.DeletePartnerGenreLanguageErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrSubscriptionServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrStoreServiceConnectionLost):
			loggerVar.Error(consts.DeletePartnerGenreLanguageErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrStoreServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrOauthServiceConnectionLost):
			loggerVar.Error(consts.DeletePartnerGenreLanguageErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrOauthServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return

		default:
			loggerVar.Errorf(consts.DeletePartnerGenreLanguageErrMsg, err)
			val, hasError, _, _ := utils.ParseFields(ctx, consts.InternalServerErr,
				"", contextError, "", "", nil, helpLink)
			if !hasError {
				loggerVar.Error(consts.DeletePartnerGenreLanguageErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
				return
			}
			result := utilities.ErrorResponseGenerator(consts.GenreNotExist, http.StatusNotFound, consts.NotFoundKey)
			ctx.JSON(http.StatusNotFound,
				result)
			return

		}
	}

	// ***************Activity code starts here********************
	name, err := partner.useCases.GetPartnerName(ctx, partnerId)
	if err != nil {
		log.Errorf(consts.DeletePartnerGenreLanguageErrMsg, err.Error())
	}
	genreName, err := utilities.GetGenre(ctx, genreId, consts.UtilityServiceURL)
	if err != nil {
		log.Errorf(consts.DeletePartnerGenreLanguageErrMsg, err.Error())
	}

	data := map[string]interface{}{
		consts.PartnerBaseUrlKey: consts.PartnerServiceURL,
		consts.PartnerIDKey:      partnerId,
		consts.PartnerNameKey:    name,
		consts.GenreIdKey:        genreId,
		consts.GenreNameKey:      genreName,
		consts.DateTimeKey:       time.Now(),
	}

	_, ok := data[consts.PartnerIDKey].(string)
	if ok {
		activityDetails := utilities.NewActivityLog(memberIDStr, consts.PartnerGenreDeletedActivityLog, data)

		resp, err := partner.activitylog.Log(activityDetails)
		if err != nil {
			logger.Log().WithContext(ctx.Request.Context()).Error(consts.ActivityLogErrMsg, err.Error(), resp)

		}
	}

	// ***************Activity code ends here********************

	result := utilities.SuccessResponseGenerator(consts.DeletePartnerGenreLanguageSuccessMsg, http.StatusOK, "")
	ctx.JSON(result.Code, result)
}

// function to delete Partner artist role language
func (partner *PartnerController) DeletePartnerArtistRoleLanguage(ctx *gin.Context) {

	var (
		ctxt = ctx.Request.Context()
		log  = logger.Log().WithContext(ctxt)
	)
	errMap := utilities.NewErrorMap()
	serviceCode := make(map[string]string)
	helpLink := consts.ErrorHelpLink
	loggerVar := logger.Log().WithContext(ctx.Request.Context())
	roleId := ctx.Param(consts.RoleIdKey)
	_, err := uuid.Parse(roleId)
	if err != nil {
		log.Errorf(consts.DeletePartnerArtistRoleLanguageErrMsg, err.Error())
		result := utilities.ErrorResponseGenerator(consts.InvalidPartner, http.StatusBadRequest, consts.ParsingErrMsg)
		ctx.JSON(result.Code,
			result,
		)
		return
	}

	method := strings.ToLower(ctx.Request.Method)
	endpointURL := ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	contextError, errVal := utils.GetContext[map[string]interface{}](ctx, consts.ContextErrorResponses)
	if !errVal {
		loggerVar.Error(consts.DeletePartnerArtistRoleLanguageErrMsg, consts.ContextErrMsg)
		return
	}
	endpoint := utils.GetEndPoints(contextEndpoints, endpointURL, method)
	if !isEndpointExists {
		val, _, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		log.Errorf(consts.DeletePartnerArtistRoleLanguageErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, ""))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), val)
		ctx.JSON(http.StatusBadRequest, result)
		return
	}
	partnerId := ctx.Param(consts.PartnerIDKey)
	// to check whether partner is valid and already exist in partner table
	partnerExists, err := partner.useCases.IsPartnerExists(ctx, partnerId, endpoint, method, errMap)
	if !partnerExists && len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, _, _ := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.DeletePartnerArtistRoleLanguageErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(consts.PathParameterErrMsg, http.StatusNotFound, val)
		// result := api.Response{Status: consts.FailureKey, Message: consts.PathParameterErrMsg, Code: http.StatusNotFound, Data: map[string]string{}, Errors: val}
		ctx.JSON(http.StatusNotFound,
			result,
		)
		return
	}

	if err != nil {
		loggerVar.Errorf(consts.DeletePartnerArtistRoleLanguageErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		if !hasError {
			loggerVar.Error(consts.DeletePartnerArtistRoleLanguageErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		loggerVar.Error(consts.DeletePartnerArtistRoleLanguageErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

		ctx.JSON(http.StatusInternalServerError, result)
		return
	}
	memberIDStr := ctx.Request.Header.Get(consts.MemberIdKey)
	err = partner.useCases.IsMemberExists(ctx, partnerId, memberIDStr, endpoint, method, errMap)
	if len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.DeletePartnerArtistRoleLanguageErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), val)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}
	if err != nil {
		loggerVar.Errorf(consts.DeletePartnerArtistRoleLanguageErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		if !hasError {
			loggerVar.Error(consts.DeletePartnerArtistRoleLanguageErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		loggerVar.Error(consts.DeletePartnerArtistRoleLanguageErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

		ctx.JSON(http.StatusInternalServerError, result)
		return

	}

	err = partner.useCases.DeletePartnerArtistRoleLanguage(ctx, roleId, endpoint, method, errMap)
	if len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, _, _ := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.DeletePartnerArtistRoleLanguageErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(consts.PathParameterErrMsg, http.StatusNotFound, val)
		ctx.JSON(http.StatusNotFound,
			result,
		)
		return
	}
	if err != nil {
		switch {

		case errors.Is(err, consts.ErrMemberServiceConnectionLost):
			loggerVar.Error(consts.DeletePartnerArtistRoleLanguageErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrMemberServiceConnectionLost)
			ctx.JSON(http.StatusBadRequest, result)
			return

		case errors.Is(err, consts.ErrUtilityServiceConnectionLost):
			loggerVar.Error(consts.DeletePartnerArtistRoleLanguageErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrUtilityServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrSubscriptionServiceConnectionLost):
			loggerVar.Error(consts.DeletePartnerArtistRoleLanguageErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrSubscriptionServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrStoreServiceConnectionLost):
			loggerVar.Error(consts.DeletePartnerArtistRoleLanguageErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrStoreServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrOauthServiceConnectionLost):
			loggerVar.Error(consts.DeletePartnerArtistRoleLanguageErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrOauthServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return

		default:
			loggerVar.Errorf(consts.DeletePartnerArtistRoleLanguageErrMsg, err.Error())
			val, hasError, _, _ := utils.ParseFields(ctx, consts.InternalServerErr,
				"", contextError, "", "", nil, helpLink)
			if !hasError {
				loggerVar.Error(consts.DeletePartnerArtistRoleLanguageErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
				return
			}
			result := utilities.ErrorResponseGenerator(consts.RoleNotExist, http.StatusNotFound, consts.NotFoundKey)
			ctx.JSON(http.StatusNotFound,
				result)
			return

		}
	}

	// ***************Activity code starts here********************
	name, err := partner.useCases.GetPartnerName(ctx, partnerId)
	if err != nil {
		log.Errorf(consts.DeletePartnerArtistRoleLanguageErrMsg, err.Error())
	}
	roleName, err := utilities.GetAristRole(ctx, roleId, consts.UtilityServiceURL)
	if err != nil {
		log.Errorf(consts.DeletePartnerArtistRoleLanguageErrMsg, err.Error())
	}
	data := map[string]interface{}{
		consts.PartnerBaseUrlKey: consts.PartnerServiceURL,
		consts.PartnerIDKey:      partnerId,
		consts.PartnerNameKey:    name,
		consts.RoleIdKey:         roleId,
		consts.RoleNameKey:       roleName,
		consts.DateTimeKey:       time.Now(),
	}

	_, ok := data[consts.PartnerIDKey].(string)
	if ok {
		activityDetails := utilities.NewActivityLog(memberIDStr, consts.PartnerRoleDeletedActivityLog, data)

		resp, err := partner.activitylog.Log(activityDetails)
		if err != nil {
			logger.Log().WithContext(ctx.Request.Context()).Error(consts.ActivityLogErrMsg, err.Error(), resp)

		}
	}

	// ***************Activity code ends here********************

	result := utilities.SuccessResponseGenerator(consts.DeletePartnerArtistRoleLanguageSuccessMsg, http.StatusOK, "")
	ctx.JSON(http.StatusOK, result)

}

// function to create partner stores
func (partner *PartnerController) CreatePartnerStores(ctx *gin.Context) {

	var (
		ctxt             = ctx.Request.Context()
		newPartnerStores entities.PartnerStores
		log              = logger.Log().WithContext(ctxt)
	)
	partnerId := ctx.Param(consts.PartnerIDKey)
	errMap := utilities.NewErrorMap()
	serviceCode := make(map[string]string)
	helpLink := consts.ErrorHelpLink

	method := strings.ToLower(ctx.Request.Method)
	endpointURL := ctx.FullPath()
	contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
	contextError, errVal := utils.GetContext[map[string]interface{}](ctx, consts.ContextErrorResponses)
	loggerVar := logger.Log().WithContext(ctx.Request.Context())
	if !errVal {
		loggerVar.Error(consts.CreatePartnerStoresErrMsg, consts.ContextErrMsg)
		return
	}
	endpoint := utils.GetEndPoints(contextEndpoints, endpointURL, method)
	if !isEndpointExists {
		val, _, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		log.Errorf(consts.CreatePartnerStoresErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), val)
		ctx.JSON(http.StatusBadRequest, result)
		return
	}
	if err := ctx.BindJSON(&newPartnerStores); err != nil {
		log.Errorf(consts.CreatePartnerStoresErrMsg, err.Error())
		result := utilities.ErrorResponseGenerator(consts.BindingErrorErrMsg, http.StatusBadRequest, consts.ParsingErrMsg)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}

	partnerExists, err := partner.useCases.IsPartnerExists(ctx, partnerId, endpoint, method, errMap)
	if !partnerExists && len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, _, _ := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.CreatePartnerStoresErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(consts.PathParameterErrMsg, http.StatusNotFound, val)
		ctx.JSON(http.StatusNotFound,
			result,
		)
		return
	}
	if err != nil {
		loggerVar.Errorf(consts.CreatePartnerStoresErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		if !hasError {
			loggerVar.Error(consts.CreatePartnerStoresErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		loggerVar.Error(consts.CreatePartnerStoresErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

		ctx.JSON(http.StatusInternalServerError, result)
		return
	}
	memberIDStr := ctx.Request.Header.Get(consts.MemberIdKey)
	err = partner.useCases.IsMemberExists(ctx, partnerId, memberIDStr, endpoint, method, errMap)
	if len(errMap) != 0 {
		for key, value := range errMap {
			serviceCode[key] = value.Code
		}
		fields := utils.FieldMapping(errMap)
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.CreatePartnerStoresErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), val)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}
	if err != nil {
		loggerVar.Errorf(consts.CreatePartnerStoresErrMsg, fmt.Sprintf(consts.LogErrMsg, err.Error()))
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		if !hasError {
			loggerVar.Error(consts.CreatePartnerStoresErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		loggerVar.Error(consts.CreatePartnerStoresErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

		ctx.JSON(http.StatusInternalServerError, result)
		return

	}
	validationErrors, storeIds, err := partner.useCases.CreatePartnerStores(ctxt, newPartnerStores, partnerId, endpoint, method, errMap)
	for key, value := range validationErrors {
		serviceCode[key] = value.Code
	}

	// validation error check
	if len(errMap) != 0 {
		fields := utils.FieldMapping(errMap)
		val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.ValidationErr,
			fields, contextError, endpoint, method, serviceCode, helpLink)
		if hasError {
			loggerVar.Error(consts.CreatePartnerStoresErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
		}
		result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), val)
		ctx.JSON(http.StatusBadRequest,
			result,
		)
		return
	}
	if err != nil {
		switch {

		case errors.Is(err, consts.ErrMemberServiceConnectionLost):
			loggerVar.Error(consts.CreatePartnerStoresErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrMemberServiceConnectionLost)
			ctx.JSON(http.StatusBadRequest, result)
			return

		case errors.Is(err, consts.ErrUtilityServiceConnectionLost):
			loggerVar.Error(consts.CreatePartnerStoresErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrUtilityServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrSubscriptionServiceConnectionLost):
			loggerVar.Error(consts.CreatePartnerStoresErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrSubscriptionServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrStoreServiceConnectionLost):
			loggerVar.Error(consts.CreatePartnerStoresErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrStoreServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return
		case errors.Is(err, consts.ErrOauthServiceConnectionLost):
			loggerVar.Error(consts.CreatePartnerStoresErrMsg, err)
			result := utilities.ErrorResponseGenerator(consts.ServiceUnavailableErrMsg, http.StatusServiceUnavailable, consts.ErrOauthServiceConnectionLost)
			ctx.JSON(http.StatusServiceUnavailable, result)
			return

		default:
			log.Errorf(consts.CreatePartnerStoresErrMsg, err.Error())
			val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
				"", contextError, "", "", nil, helpLink)
			if !hasError {
				loggerVar.Error(consts.CreatePartnerStoresErrMsg, fmt.Sprintf(consts.LocalizationModuleErrMsg, val))
			}
			loggerVar.Error(consts.CreatePartnerStoresErrMsg, fmt.Sprintf(consts.InvalidEndpointErrMsg, val))
			result := utilities.ErrorResponseGenerator(errDet.Message, int(errorCode), "")

			ctx.JSON(http.StatusBadRequest, result)
			return

		}
	}

	// ***************Activity code starts here********************
	name, err := partner.useCases.GetPartnerName(ctx, partnerId)
	if err != nil {
		log.Errorf(consts.CreatePartnerStoresErrMsg, err.Error())
	}
	storeName, err := partner.useCases.GetStoreName(ctx, storeIds)
	if err != nil {
		log.Errorf(consts.CreatePartnerStoresErrMsg, err.Error())
	}

	data := map[string]interface{}{
		consts.PartnerBaseUrlKey: consts.PartnerServiceURL,
		consts.PartnerIDKey:      partnerId,
		consts.PartnerNameKey:    name,
		consts.StoresKey:         strings.Join(storeName, ","),
		consts.DateTimeKey:       time.Now(),
	}

	_, ok := data[consts.PartnerIDKey].(string)
	if ok {
		activityDetails := utilities.NewActivityLog(memberIDStr, consts.PartnerStoreCreatedActivityLog, data)

		resp, err := partner.activitylog.Log(activityDetails)
		if err != nil {
			logger.Log().WithContext(ctx.Request.Context()).Error(consts.ActivityLogErrMsg, err.Error(), resp)

		}
	}

	// ***************Activity code ends here********************

	result := utilities.SuccessResponseGenerator(consts.CreatePartnerStoreSuccessMsg, http.StatusCreated, "")
	ctx.JSON(http.StatusCreated, result)

}
