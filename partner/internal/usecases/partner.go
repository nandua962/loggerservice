// Package usecases contains the use case implementations related to partner operations.
package usecases

import (
	"context"
	"errors"
	"html/template"
	"partner/internal/consts"
	"partner/internal/entities"
	"partner/internal/repo"
	"slices"
	"strings"

	"partner/utilities"

	"github.com/gin-gonic/gin"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils/crypto"

	"github.com/google/uuid"

	cacheConf "gitlab.com/tuneverse/toolkit/core/cache"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/utils"
)

// PartnerUseCases defines the use cases for partner operations.
type PartnerUseCases struct {
	repo  repo.PartnerRepoImply
	cache cacheConf.Cache
}

// PartnerUseCaseImply defines the interface for partner use cases.
type PartnerUseCaseImply interface {
	DeletePartner(*gin.Context, string, string, string, map[string]models.ErrorResponse) error
	GetAllPartners(*gin.Context, entities.Params, string, string, map[string]models.ErrorResponse) ([]entities.ListAllPartners, models.MetaData, map[string]models.ErrorResponse, error)
	CreatePartner(context.Context, entities.Partner, string, string, map[string]models.ErrorResponse) (map[string]models.ErrorResponse, string, error)
	GetPartnerById(*gin.Context, entities.QueryParams, string, string, string, map[string]models.ErrorResponse) (interface{}, models.MetaData, map[string]models.ErrorResponse, error)
	GetPartnerOauthCredential(context.Context, string, entities.PartnerOAuthHeader, string, string, map[string]models.ErrorResponse) (entities.GetPartnerOauthCredential, error)
	UpdateTermsAndConditions(context.Context, string, uuid.UUID, entities.UpdateTermsAndConditions, string, string, map[string]models.ErrorResponse) (map[string]models.ErrorResponse, error)
	UpdatePartner(context.Context, string, uuid.UUID, entities.PartnerProperties, string, string, map[string]models.ErrorResponse) (map[string]models.ErrorResponse, error)
	GetAllTermsAndConditions(context.Context, string, string, string, map[string]models.ErrorResponse) (entities.TermsAndConditions, error)
	IsPartnerExists(context.Context, string, string, string, map[string]models.ErrorResponse) (bool, error)
	GetPartnerPaymentGateways(*gin.Context, string, string, string, map[string]models.ErrorResponse) (entities.GetPaymentGateways, error)
	GetPartnerStores(context.Context, string, string, string, map[string]models.ErrorResponse) (entities.GetPartnerStores, error)
	DeletePartnerGenreLanguage(*gin.Context, string, string, string, map[string]models.ErrorResponse) error
	DeletePartnerArtistRoleLanguage(*gin.Context, string, string, string, map[string]models.ErrorResponse) error
	IsMemberExists(*gin.Context, string, string, string, string, map[string]models.ErrorResponse) error
	GetPartnerName(*gin.Context, string) (string, error)
	CreatePartnerStores(context.Context, entities.PartnerStores, string, string, string, map[string]models.ErrorResponse) (map[string]models.ErrorResponse, []string, error)
	GetStoreName(context.Context, []string) ([]string, error)
	IsExists(context.Context, string, string, string) (bool, error)
	UpdatePartnerStatus(*gin.Context, string, entities.UpdatePartnerStatus, string, string, map[string]models.ErrorResponse) error
}

// NewPartnerUseCases
func NewPartnerUseCases(partnerRepo repo.PartnerRepoImply, cache cacheConf.Cache) PartnerUseCaseImply {
	return &PartnerUseCases{
		repo:  partnerRepo,
		cache: cache,
	}
}

// function to get store name
func (partner *PartnerUseCases) GetStoreName(ctx context.Context, storeIds []string) ([]string, error) {
	var storeName []string
	data, err := utilities.GetStores(ctx, partner.cache, consts.StoreServiceURL)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf(consts.CreatePartnerStoresErrMsg, err.Error())
	}
	for _, v := range data {
		if slices.Contains(storeIds, v.Id) {
			storeName = append(storeName, v.Name)
		}

	}
	return storeName, nil
}

// function to check a value exists in a table
func (partner *PartnerUseCases) IsExists(ctx context.Context, tableName string, fieldName string, fieldValue string) (bool, error) {
	return partner.repo.IsExists(ctx, tableName, fieldName, fieldValue)
}

// function to get partner name
func (partner *PartnerUseCases) GetPartnerName(ctx *gin.Context, partnerID string) (string, error) {
	return partner.repo.GetPartnerName(ctx, partnerID)
}

// function to update the status of  a partner
func (partner *PartnerUseCases) UpdatePartnerStatus(ctx *gin.Context, partnerID string, partnerStatus entities.UpdatePartnerStatus, endpoint string, method string, errMap map[string]models.ErrorResponse) error {
	return partner.repo.UpdatePartnerStatus(ctx, partnerID, partnerStatus)
}

// function to check whether a member exists under a partner based on partner id
func (partner *PartnerUseCases) IsMemberExists(ctx *gin.Context, partnerID string, memberID string, endpoint string, method string, errMap map[string]models.ErrorResponse) error {
	if utils.IsEmpty(memberID) {
		logger.Log().WithContext(ctx).Error("MemberID is epmty")
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.MemberIdKey, consts.RequiredKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error("IsMemberExists failed", err)
		}

		errMap[consts.MemberIdKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.RequiredKey},
		}
		return err
	}
	_, err := uuid.Parse(memberID)
	if err != nil {
		logger.Log().WithContext(ctx).Error("IsMemberExists failed", err)
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.MemberIdKey, consts.InvalidKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error("IsMemberExists failed", err)
		}

		errMap[consts.MemberIdKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}
		return err
	}

	exists, err := utilities.IsMemberExists(ctx, memberID, partnerID, consts.MemberServiceURL)
	if err != nil {
		logger.Log().WithContext(ctx).Error("IsMemberExists failed", err)
		return err
	}
	if !exists {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.MemberIdKey, consts.NotFoundKey)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
		}
		errMap[consts.MemberIdKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.NotFoundKey},
		}
	}

	return nil
}

// function to get partner stores based on partner id
func (partner *PartnerUseCases) GetPartnerStores(ctx context.Context, PartnerID string, endpoint string, method string, errMap map[string]models.ErrorResponse) (entities.GetPartnerStores, error) {
	return partner.repo.GetPartnerStores(ctx, PartnerID)
}

// funtion to get partner payment gateways based on partner id
func (partner *PartnerUseCases) GetPartnerPaymentGateways(ctx *gin.Context, PartnerID string, endpoint string, method string, errMap map[string]models.ErrorResponse) (entities.GetPaymentGateways, error) {
	return partner.repo.GetPartnerPaymentGateways(ctx, PartnerID)
}

// function to Get All Partners
func (partner *PartnerUseCases) GetAllPartners(ctx *gin.Context, params entities.Params, endpoint string, method string, errMap map[string]models.ErrorResponse) ([]entities.ListAllPartners, models.MetaData, map[string]models.ErrorResponse, error) {

	validKeys := map[string]bool{
		consts.NameKey:    true,
		consts.CountryKey: true,
		consts.SortKey:    true,
		consts.OrderKey:   true,
		consts.StatusKey:  true,
		consts.PageKey:    true,
		consts.LimitKey:   true,
	}

	// Check if the keys in the query parameters are valid
	invalidKeys := []string{}
	for key := range ctx.Request.URL.Query() {
		if _, ok := validKeys[key]; !ok {
			invalidKeys = append(invalidKeys, key)
		}
	}

	// If invalid keys are found, return an error response
	if len(invalidKeys) > 0 {
		err := errors.New("Invalid query parameter key(s): " + strings.Join(invalidKeys, ", "))
		logger.Log().WithContext(ctx).Errorf(consts.GetAllPartnersErrMsg, err.Error())
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.Key, consts.InvalidKey)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf(consts.GetAllPartnersErrMsg, err)
		}
		errMap[consts.Key] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}

	}
	//set default value for page parameter
	if params.Page == 0 {
		params.Page = consts.PageDefaultVal
	}
	//set default value for limit parameter
	if params.Limit == 0 {
		params.Limit = consts.LimitDefaultVal
	}
	//Set default value for sort parameter
	if utils.IsEmpty(params.Sort) {
		params.Sort = consts.NameKey

	}

	// check Page paramater is an integer
	if params.Page <= 0 {
		logger.Log().WithContext(ctx).Errorf(consts.GetAllPartnersErrMsg, errors.New("page number is not an integer"))
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PageKey, consts.InvalidKey)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf(consts.GetAllPartnersErrMsg, err)
		}

		errMap[consts.PageKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}

	}
	//check Limit parameter is an  integer
	if params.Limit <= 0 {
		logger.Log().WithContext(ctx).Errorf(consts.GetAllPartnersErrMsg, errors.New("limit is not an integer"))
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.LimitKey, consts.InvalidKey)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf(consts.GetAllPartnersErrMsg, err)
		}
		errMap[consts.LimitKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}

	}
	// check the Limit parameter not exceeds the Maximum limit
	if params.Limit > consts.MaxLimit {
		err := errors.New("cannot exceed maximum page limit")
		logger.Log().WithContext(ctx).Errorf(consts.GetAllPartnersErrMsg, err.Error())
		return []entities.ListAllPartners{}, models.MetaData{}, nil, consts.ErrMaximumRequest
	}

	//check status parameter holds an valid value
	if params.Status != "" {
		if strings.ToLower(params.Status) != consts.StatusActive && strings.ToLower(params.Status) != consts.StatusInActive && strings.ToLower(params.Status) != consts.StatusAll {
			logger.Log().WithContext(ctx).Errorf(consts.GetAllPartnersErrMsg, errors.New("invalid status"))
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.StatusKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Errorf(consts.GetAllPartnersErrMsg, err)
			}
			errMap[consts.StatusKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}

		}
	}
	//check sort parameter holds an valid value
	sortFields := strings.Split(params.Sort, ",")
	for i := range sortFields {
		if sortFields[i] != "" {
			if strings.ToLower(sortFields[i]) != consts.NameKey && strings.ToLower(sortFields[i]) != consts.EmailKey {
				logger.Log().WithContext(ctx).Errorf(consts.GetAllPartnersErrMsg, errors.New("invalid sort field"))
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.SortKey, consts.InvalidKey)
				if err != nil {
					logger.Log().WithContext(ctx).Errorf(consts.GetAllPartnersErrMsg, err)
				}
				errMap[consts.SortKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.InvalidKey},
				}

			}
		}
	}
	//check order parameter holds an valid value
	sortOrders := strings.Split(params.Order, ",")
	for i := range sortOrders {
		if sortOrders[i] != "" {
			if strings.ToLower(sortOrders[i]) != consts.Ascending && strings.ToLower(sortOrders[i]) != consts.Descending {
				logger.Log().WithContext(ctx).Errorf(consts.GetAllPartnersErrMsg, errors.New("invalid sort order"))
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.OrderKey, consts.InvalidKey)
				if err != nil {
					logger.Log().WithContext(ctx).Errorf(consts.GetAllPartnersErrMsg, err)
				}
				errMap[consts.OrderKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.InvalidKey},
				}
				break

			}
		}
	}
	if len(errMap) != 0 {
		logger.Log().WithContext(ctx).Errorf(consts.GetAllPartnersErrMsg, errMap)
		return []entities.ListAllPartners{}, models.MetaData{}, errMap, nil
	}

	params.Page, params.Limit = utils.Paginate(params.Page, params.Limit, consts.LimitDefaultVal)

	data, recordCount, err := partner.repo.GetAllPartners(ctx, params, endpoint, method, errMap)

	if err != nil {
		logger.Log().WithContext(ctx).Errorf(consts.GetAllPartnersErrMsg, err)
		return []entities.ListAllPartners{}, models.MetaData{}, nil, err
	}
	if (params.Page*params.Limit)-int32(recordCount) >= params.Limit {
		return []entities.ListAllPartners{}, models.MetaData{}, nil, nil
	}

	metadata := &models.MetaData{
		CurrentPage: params.Page,
		PerPage:     params.Limit,
		Total:       recordCount,
	}

	metadata = utils.MetaDataInfo(metadata)

	return data, *metadata, nil, nil
}

// function to get partner Oauth Credentials based on partner id
func (partner *PartnerUseCases) GetPartnerOauthCredential(ctx context.Context, partnerId string, oauthHeader entities.PartnerOAuthHeader, endpoint string, method string, errMap map[string]models.ErrorResponse) (entities.GetPartnerOauthCredential, error) {
	var (
		log             = logger.Log().WithContext(ctx)
		err             error
		oauthProviderID string
	)
	if utils.IsEmpty(oauthHeader.ProviderName) {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.OauthProviderKey, consts.RequiredKey)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf(consts.GetPartnerOauthCredentialErrMsg, err)
		}
		errMap[consts.OauthProviderKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.RequiredKey},
		}
		return entities.GetPartnerOauthCredential{}, nil
	} else {
		oauthProviderID, err = utilities.GetOauthProviderId(ctx, partner.cache, oauthHeader.ProviderName, consts.OAuthServiceURL)
		if err != nil {
			log.Errorf(consts.GetPartnerOauthCredentialErrMsg, err.Error())
			return entities.GetPartnerOauthCredential{}, err
		}
		if utils.IsEmpty(oauthProviderID) {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.OauthProviderKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Errorf(consts.GetPartnerOauthCredentialErrMsg, err)
			}
			errMap[consts.OauthProviderKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}
			return entities.GetPartnerOauthCredential{}, nil
		}

	}

	data, err := partner.repo.GetPartnerOauthCredential(ctx, partnerId, oauthHeader, oauthProviderID)
	if err != nil {
		log.Errorf(consts.GetPartnerOauthCredentialErrMsg, err.Error())
		return entities.GetPartnerOauthCredential{}, err
	}
	encryptionKey := consts.ClientCredentialEncryptionKey
	key := []byte(encryptionKey)
	clientSecret, err := crypto.Decrypt(data.ClientSecret, key)
	if err != nil {
		log.Errorf(consts.GetPartnerOauthCredentialErrMsg, err.Error())
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.EncryptionKey, consts.InvalidKey)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf(consts.GetPartnerOauthCredentialErrMsg, err)
		}
		errMap[consts.EncryptionKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}
		return entities.GetPartnerOauthCredential{}, err
	}
	data.ClientSecret = clientSecret
	clientId, err := crypto.Decrypt(data.ClientId, key)
	if err != nil {
		log.Errorf(consts.GetPartnerOauthCredentialErrMsg, err.Error())
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.EncryptionKey, consts.InvalidKey)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf(consts.GetPartnerOauthCredentialErrMsg, err)
		}
		errMap[consts.EncryptionKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}
		return entities.GetPartnerOauthCredential{}, err
	}
	data.ClientId = clientId
	return data, nil
}

// fnction to get a partner details based on partner id
func (partner *PartnerUseCases) GetPartnerById(ctx *gin.Context, params entities.QueryParams, partnerId string, endpoint string, method string, errMap map[string]models.ErrorResponse) (interface{}, models.MetaData, map[string]models.ErrorResponse, error) {
	var log = logger.Log().WithContext(ctx)
	if !utils.IsEmpty(params.Fields) {
		switch params.Fields {
		case consts.ProductTypeKey:
			errMap, err := QueryParamValidations(ctx, &params, endpoint, method, errMap)
			if len(errMap) != 0 {
				logger.Log().WithContext(ctx).Errorf(consts.GetPartnerProductTypesErrMsg, errMap)
				return nil, models.MetaData{}, errMap, nil
			}
			if err != nil {
				logger.Log().WithContext(ctx).Errorf(consts.GetPartnerProductTypesErrMsg, err)
				return nil, models.MetaData{}, map[string]models.ErrorResponse{}, err
			}

			params.Page, params.Limit = utils.Paginate(params.Page, params.Limit, consts.LimitDefaultVal)

			recordCount, data, err := partner.repo.GetPartnerProductTypes(ctx, partnerId, params)

			if err != nil {
				logger.Log().WithContext(ctx).Errorf(consts.GetPartnerProductTypesErrMsg, err)
				return nil, models.MetaData{}, nil, err
			}
			if (params.Page*params.Limit)-int32(recordCount) >= params.Limit {
				return nil, models.MetaData{}, nil, nil
			}

			metadata := &models.MetaData{
				CurrentPage: params.Page,
				PerPage:     params.Limit,
				Total:       recordCount,
			}

			metadata = utils.MetaDataInfo(metadata)

			return data, *metadata, nil, nil
		case consts.TrackFileQualityKey:
			errMap, err := QueryParamValidations(ctx, &params, endpoint, method, errMap)
			if len(errMap) != 0 {
				logger.Log().WithContext(ctx).Errorf(consts.GetPartnerTrackQualityErrMsg, errMap)
				return nil, models.MetaData{}, errMap, nil
			}
			if err != nil {
				logger.Log().WithContext(ctx).Errorf(consts.GetPartnerTrackQualityErrMsg, err)
				return nil, models.MetaData{}, errMap, err
			}

			params.Page, params.Limit = utils.Paginate(params.Page, params.Limit, consts.LimitDefaultVal)

			recordCount, data, err := partner.repo.GetPartnerTrackFileQuality(ctx, partnerId, params)
			if err != nil {
				log.Errorf(consts.GetPartnerTrackQualityErrMsg, err.Error())
				return []entities.GetPartnerProdTypesAndTrackQuality{}, models.MetaData{}, errMap, err
			}
			if (params.Page*params.Limit)-int32(recordCount) >= params.Limit {
				return nil, models.MetaData{}, nil, nil
			}

			metadata := &models.MetaData{
				CurrentPage: params.Page,
				PerPage:     params.Limit,
				Total:       recordCount,
			}

			metadata = utils.MetaDataInfo(metadata)

			return data, *metadata, nil, nil
		default:
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.FieldKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Errorf(consts.GetPartnerByIdErrMsg, err)
			}
			errMap[consts.FieldKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}
			return entities.GetPartner{}, models.MetaData{}, errMap, nil
		}
	}

	data, err := partner.repo.GetPartnerById(ctx, partnerId)
	if err != nil {
		log.Errorf(consts.GetPartnerByIdErrMsg, err.Error())
		return entities.GetPartner{}, models.MetaData{}, errMap, err
	}
	return data, models.MetaData{}, nil, nil

}

func QueryParamValidations(ctx *gin.Context, params *entities.QueryParams, endpoint string, method string, errMap map[string]models.ErrorResponse) (map[string]models.ErrorResponse, error) {
	validKeys := map[string]bool{
		consts.SortKey:  true,
		consts.OrderKey: true,
		consts.PageKey:  true,
		consts.LimitKey: true,
		consts.FieldKey: true,
	}

	// Check if the keys in the query parameters are valid
	invalidKeys := []string{}
	for key := range ctx.Request.URL.Query() {
		if _, ok := validKeys[key]; !ok {
			invalidKeys = append(invalidKeys, key)
		}
	}

	// If invalid keys are found, return an error response
	if len(invalidKeys) > 0 {
		err := errors.New("Invalid query parameter key(s): " + strings.Join(invalidKeys, ", "))
		logger.Log().WithContext(ctx).Errorf("err=%s", err)
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.Key, consts.InvalidKey)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf("err=%s", err)
		}
		errMap[consts.Key] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}

	}
	//set default value for page parameter
	if params.Page == 0 {
		params.Page = consts.PageDefaultVal
	}
	//set default value for limit parameter
	if params.Limit == 0 {
		params.Limit = consts.LimitDefaultVal
	}
	//Set default value for sort parameter
	if utils.IsEmpty(params.Sort) {
		params.Sort = consts.NameKey

	}

	// check Page paramater is an integer
	if params.Page <= 0 {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PageKey, consts.InvalidKey)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf("err=%s", err)
		}

		errMap[consts.PageKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}

	}
	//check Limit parameter is an  integer
	if params.Limit <= 0 {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.LimitKey, consts.InvalidKey)
		if err != nil {
			logger.Log().WithContext(ctx).Errorf("err=%s", err)
		}
		errMap[consts.LimitKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}

	}
	// check the Limit parameter not exceeds the Maximum limit
	if params.Limit > consts.MaxLimit {
		return map[string]models.ErrorResponse{}, consts.ErrMaximumRequest
	}

	//check sort parameter holds an valid value
	if !utils.IsEmpty(params.Sort) {
		if strings.ToLower(params.Sort) != consts.NameKey {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.SortKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Errorf("err=%s", err)
			}
			errMap[consts.SortKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}

		}
	}

	//check order parameter holds an valid value
	if !utils.IsEmpty(params.Order) {
		if strings.ToLower(params.Order) != consts.Ascending && strings.ToLower(params.Order) != consts.Descending {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.OrderKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Errorf("err=%s", err)
			}
			errMap[consts.OrderKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}

		}
	}

	return errMap, nil
}

// function to get all terms and conditions of partner based on partner id
func (partner *PartnerUseCases) GetAllTermsAndConditions(ctx context.Context, partnerId string, endpoint string, method string, errMap map[string]models.ErrorResponse) (entities.TermsAndConditions, error) {
	var log = logger.Log().WithContext(ctx)
	data, err := partner.repo.GetAllTermsAndConditions(ctx, partnerId)

	if err != nil {
		log.Errorf(consts.GetAllTermsAndConditionErrMsg, err.Error())
		return entities.TermsAndConditions{}, err
	}

	return data, nil

}

// function to validate partner update fields
func (partner *PartnerUseCases) ValidatePartnerUpdate(ctx context.Context, partnerID string, arg entities.PartnerProperties, data entities.Partner, endpoint string, method string, errMap map[string]models.ErrorResponse) (map[string]models.ErrorResponse, entities.Partner, error) {

	var log = logger.Log().WithContext(ctx)

	contactDetails, err := arg.GetContactDetails()
	if err != nil {
		log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
		return nil, entities.Partner{}, err
	}

	addressDetails, err := arg.GetAddressDetails()
	if err != nil {
		log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
		return nil, entities.Partner{}, err
	}

	subscriptionDetails, err := arg.GetSubscriptionPlanDetails()
	if err != nil {
		log.Errorf(consts.UpdatePartnerErrMsg, err)
		return nil, entities.Partner{}, err
	}

	gatewayDetails, err := arg.GetPaymentDetails()
	if err != nil {
		log.Errorf(consts.UpdatePartnerErrMsg, err)
		return nil, entities.Partner{}, err
	}

	for argField := range arg {
		switch argField {

		case consts.NameKey:

			partnerName := arg.GetPartnerName()
			if utils.IsEmpty(partnerName) {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.NameKey, consts.RequiredKey)
				if err != nil {
					logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
				}
				errMap[consts.NameKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.RequiredKey},
				}

			} else {
				isValid := utilities.IsValidPartnerName(partnerName)
				if !isValid {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.NameKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.NameKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				} else {
					name := strings.TrimSpace(partnerName)
					if len(name) > consts.PartnerNameMaxLength {
						code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.NameKey, consts.LimitExceedsKey)
						if err != nil {
							logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
						}
						errMap[consts.NameKey] = models.ErrorResponse{
							Code:    code,
							Message: []string{consts.LimitExceedsKey},
						}
					} else {
						unique, err := partner.repo.IsFieldValueUnique(ctx, consts.NameKey, name, partnerID)
						if err != nil {
							log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
							return nil, entities.Partner{}, err
						}

						if !unique {
							code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.NameKey, consts.AlreadyExistsKey)
							if err != nil {
								logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
							}
							errMap[consts.NameKey] = models.ErrorResponse{
								Code:    code,
								Message: []string{consts.AlreadyExistsKey},
							}
						}
					}
				}

			}
			data.Name = partnerName

		case consts.URLKey:

			url := arg.GetURL()
			if utils.IsEmpty(url) {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.URLKey, consts.RequiredKey)
				if err != nil {
					logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
				}
				errMap[consts.URLKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.RequiredKey},
				}

			} else {
				if !utilities.IsValidURL(url) {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.URLKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.URLKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				} else {
					if len(url) > consts.URLMaxLength {
						code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.URLKey, consts.LimitExceedsKey)
						if err != nil {
							logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
						}
						errMap[consts.URLKey] = models.ErrorResponse{
							Code:    code,
							Message: []string{consts.LimitExceedsKey},
						}
					} else {
						unique, err := partner.repo.IsFieldValueUnique(ctx, consts.URLKey, url, partnerID)
						if err != nil {
							log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
							return nil, entities.Partner{}, err
						}

						if !unique {
							code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.URLKey, consts.AlreadyExistsKey)
							if err != nil {
								logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
							}
							errMap[consts.URLKey] = models.ErrorResponse{
								Code:    code,
								Message: []string{consts.AlreadyExistsKey},
							}
						}
					}
				}

			}
			data.URL = url
		case consts.LanguageKey:
			language := arg.GetLang()
			if !utils.IsEmpty(language) {
				exists, err := utilities.IsLanguageIsoExists(ctx, partner.cache, language, consts.UtilityServiceURL)
				if err != nil {
					log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
					return nil, entities.Partner{}, err
				}
				if !exists {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.LanguageKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.LanguageKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				}
				data.Language = language
			}

		case consts.DefaultPriceCodeCurrencyKey:
			defaultPriceCodeCurrency, err := arg.GetDefaultPriceCodeCurrency()
			if err != nil {
				log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
			}

			if !utils.IsEmpty(defaultPriceCodeCurrency.Name) {
				defaultPriceCodeCurrencyID, err := utilities.GetCurrencyId(ctx, partner.cache, defaultPriceCodeCurrency.Name, consts.UtilityServiceURL)
				if err != nil {
					log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
					return nil, entities.Partner{}, err
				}
				if defaultPriceCodeCurrencyID == 0 {

					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.DefaultPriceCodeCurrencyKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.DefaultPriceCodeCurrencyKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				}
				data.DefaultPriceCodeCurrencyID = defaultPriceCodeCurrencyID
			}

		case consts.PayoutTargetCurrencyKey:
			payoutTargetCurrency := arg.GetPayoutTargetCurren()
			if !utils.IsEmpty(payoutTargetCurrency) {
				payoutTargetCurrencyID, err := utilities.GetCurrencyId(ctx, partner.cache, payoutTargetCurrency, consts.UtilityServiceURL)
				if err != nil {
					log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
					return nil, entities.Partner{}, err
				}
				if payoutTargetCurrencyID == 0 {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PayoutTargetCurrencyKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.PayoutTargetCurrencyKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				}
				data.PayoutTargetCurrencyID = payoutTargetCurrencyID
			}
		case consts.MemberGracePeriodKey:
			memberGracePeriod := arg.GetMemberGracePeriod()
			if !utils.IsEmpty(memberGracePeriod) {
				id, err := utilities.GetSubscriptionDurationId(ctx, partner.cache, memberGracePeriod, consts.SubcriptionServiceApiURL)
				if err != nil {
					log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
					return nil, entities.Partner{}, err
				}

				if id == 0 {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.MemberGracePeriodKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.MemberGracePeriodKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				}
				data.MemberGracePeriodID = id
			}
		case consts.BusinessModelKey:
			businessModel := arg.GetBusinessModel()
			if utils.IsEmpty(businessModel) {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.BusinessModelKey, consts.RequiredKey)
				if err != nil {
					logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
				}
				errMap[consts.BusinessModelKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.RequiredKey},
				}
			} else {
				data.BusinessModelID, err = utilities.GetLookupId(ctx, partner.cache, consts.BusinessLookupTypeName, businessModel, consts.UtilityServiceURL)

				if err != nil {
					log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
					return nil, entities.Partner{}, err

				}
				if data.BusinessModelID == 0 {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.BusinessModelKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.BusinessModelKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				}
			}
		case consts.ProductReviewKey:
			productReview := arg.GetProductReview()
			if !utils.IsEmpty(productReview) {
				data.ProductReviewID, err = utilities.GetLookupId(ctx, partner.cache, consts.ProductReviewLookupTypeName, productReview, consts.UtilityServiceURL)
				if err != nil {
					log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
					return nil, entities.Partner{}, err
				}
				if data.ProductReviewID == 0 {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.ProductReviewKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.ProductReviewKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				}
			}

		case consts.LoginTypeKey:
			loginType := arg.GetLoginType()
			if !utils.IsEmpty(loginType) {
				data.LoginTypeID, err = utilities.GetLookupId(ctx, partner.cache, consts.LoginLookupTypeName, loginType, consts.UtilityServiceURL)
				if err != nil {
					log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
					return nil, entities.Partner{}, err
				}
				if data.LoginTypeID == 0 {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.LoginTypeKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.LoginTypeKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				}
			}
			data.LoginType = loginType

		case consts.LogoKey:
			logo := arg.GetLogo()
			if !utils.IsEmpty(logo) {
				if !utilities.IsValidURL(logo) {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.LogoKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.LogoKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				}
			}
			data.Logo = logo

		case consts.ContactDetailsKey:

			for val := range contactDetails {
				switch val {
				case consts.EmailKey:
					email, err := arg.GetEmail()
					if err != nil {
						log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
						return nil, entities.Partner{}, err
					}

					if utils.IsEmpty(email) {
						code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.EmailKey, consts.RequiredKey)
						if err != nil {
							logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
						}
						errMap[consts.EmailKey] = models.ErrorResponse{
							Code:    code,
							Message: []string{consts.RequiredKey},
						}

					} else {
						if !utils.IsValidEmail(email) {
							code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.EmailKey, consts.InvalidKey)
							if err != nil {
								logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
							}
							errMap[consts.EmailKey] = models.ErrorResponse{
								Code:    code,
								Message: []string{consts.InvalidKey},
							}
						} else {
							if len(email) > consts.EmailMaxLength {
								code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.EmailKey, consts.LimitExceedsKey)
								if err != nil {
									logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
								}
								errMap[consts.EmailKey] = models.ErrorResponse{
									Code:    code,
									Message: []string{consts.LimitExceedsKey},
								}
							} else {
								unique, err := partner.repo.IsFieldValueUnique(ctx, consts.EmailKey, email, partnerID)

								if err != nil {
									log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
									return nil, entities.Partner{}, err
								}

								if !unique {
									code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.EmailKey, consts.AlreadyExistsKey)
									if err != nil {
										logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
									}
									errMap[consts.EmailKey] = models.ErrorResponse{
										Code:    code,
										Message: []string{consts.AlreadyExistsKey},
									}
								}
							}
						}

					}

					data.ContactDetails.Email = email

				case consts.NoReplyEmailKey:
					noreplyEmail, err := arg.GetNoReplyEmail()
					if err != nil {
						log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
						return nil, entities.Partner{}, err
					}

					if utils.IsEmpty(noreplyEmail) {
						code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.NoReplyEmailKey, consts.RequiredKey)
						if err != nil {
							logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
						}
						errMap[consts.NoReplyEmailKey] = models.ErrorResponse{
							Code:    code,
							Message: []string{consts.RequiredKey},
						}

					} else {
						if !utils.IsValidEmail(noreplyEmail) {
							code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.NoReplyEmailKey, consts.InvalidKey)
							if err != nil {
								logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
							}
							errMap[consts.NoReplyEmailKey] = models.ErrorResponse{
								Code:    code,
								Message: []string{consts.InvalidKey},
							}
						} else {
							if len(noreplyEmail) > consts.EmailMaxLength {
								code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.NoReplyEmailKey, consts.LimitExceedsKey)
								if err != nil {
									logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
								}
								errMap[consts.NoReplyEmailKey] = models.ErrorResponse{
									Code:    code,
									Message: []string{consts.LimitExceedsKey},
								}
							} else {
								unique, err := partner.repo.IsFieldValueUnique(ctx, consts.NoReplyEmailKey, noreplyEmail, partnerID)

								if err != nil {
									log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
									return nil, entities.Partner{}, err
								}

								if !unique {
									code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.NoReplyEmailKey, consts.AlreadyExistsKey)
									if err != nil {
										logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
									}
									errMap[consts.NoReplyEmailKey] = models.ErrorResponse{
										Code:    code,
										Message: []string{consts.AlreadyExistsKey},
									}
								}
							}
						}

					}
					data.ContactDetails.NoReplyEmail = noreplyEmail

				case consts.FeedbackEmailKey:
					feedbackEmail, err := arg.GetFeedBackEmail()
					if err != nil {
						log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
						return nil, entities.Partner{}, err
					}

					if utils.IsEmpty(feedbackEmail) {
						code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.FeedbackEmailKey, consts.RequiredKey)
						if err != nil {
							logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
						}
						errMap[consts.FeedbackEmailKey] = models.ErrorResponse{
							Code:    code,
							Message: []string{consts.RequiredKey},
						}

					} else {
						if !utils.IsValidEmail(feedbackEmail) {
							code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.FeedbackEmailKey, consts.InvalidKey)
							if err != nil {
								logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
							}
							errMap[consts.FeedbackEmailKey] = models.ErrorResponse{
								Code:    code,
								Message: []string{consts.InvalidKey},
							}
						} else {
							if len(feedbackEmail) > consts.EmailMaxLength {
								code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.FeedbackEmailKey, consts.LimitExceedsKey)
								if err != nil {
									logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
								}
								errMap[consts.FeedbackEmailKey] = models.ErrorResponse{
									Code:    code,
									Message: []string{consts.LimitExceedsKey},
								}
							} else {
								unique, err := partner.repo.IsFieldValueUnique(ctx, consts.FeedbackEmailKey, feedbackEmail, partnerID)

								if err != nil {
									log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
									return nil, entities.Partner{}, err
								}

								if !unique {
									code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.FeedbackEmailKey, consts.AlreadyExistsKey)
									if err != nil {
										logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
									}
									errMap[consts.FeedbackEmailKey] = models.ErrorResponse{
										Code:    code,
										Message: []string{consts.AlreadyExistsKey},
									}
								}
							}
						}

					}
					data.ContactDetails.FeedbackEmail = feedbackEmail
				case consts.SupportEmailKey:
					supportEmail, err := arg.GetSupportEmail()
					if err != nil {
						log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
						return nil, entities.Partner{}, err
					}

					if utils.IsEmpty(supportEmail) {
						code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.SupportEmailKey, consts.RequiredKey)
						if err != nil {
							logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
						}
						errMap[consts.SupportEmailKey] = models.ErrorResponse{
							Code:    code,
							Message: []string{consts.RequiredKey},
						}

					} else {
						if !utils.IsValidEmail(supportEmail) {
							code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.SupportEmailKey, consts.InvalidKey)
							if err != nil {
								logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
							}
							errMap[consts.SupportEmailKey] = models.ErrorResponse{
								Code:    code,
								Message: []string{consts.InvalidKey},
							}
						} else {
							if len(supportEmail) > consts.EmailMaxLength {
								code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.SupportEmailKey, consts.LimitExceedsKey)
								if err != nil {
									logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
								}
								errMap[consts.SupportEmailKey] = models.ErrorResponse{
									Code:    code,
									Message: []string{consts.LimitExceedsKey},
								}
							} else {
								unique, err := partner.repo.IsFieldValueUnique(ctx, consts.SupportEmailKey, supportEmail, partnerID)
								if err != nil {
									log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
									return nil, entities.Partner{}, err
								}

								if !unique {
									code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.SupportEmailKey, consts.AlreadyExistsKey)
									if err != nil {
										logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
									}
									errMap[consts.SupportEmailKey] = models.ErrorResponse{
										Code:    code,
										Message: []string{consts.AlreadyExistsKey},
									}
								}
							}
						}

					}
					data.ContactDetails.SupportEmail = supportEmail
				case consts.ContactPersonKey:
					contactPerson, err := arg.GetContactPerson()
					if err != nil {
						log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
						return nil, entities.Partner{}, err
					}
					contactPersonName := strings.TrimSpace(contactPerson)
					if utils.IsEmpty(contactPersonName) {
						code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.ContactPersonKey, consts.RequiredKey)
						if err != nil {
							logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
						}
						errMap[consts.ContactPersonKey] = models.ErrorResponse{
							Code:    code,
							Message: []string{consts.RequiredKey},
						}

					} else {
						isValid := utilities.IsValidPartnerName(contactPersonName)

						if !isValid {
							code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.ContactPersonKey, consts.InvalidKey)
							if err != nil {
								logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
							}
							errMap[consts.ContactPersonKey] = models.ErrorResponse{
								Code:    code,
								Message: []string{consts.InvalidKey},
							}
						} else {
							if len(contactPersonName) > consts.ContactPersonMaxLength {
								code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.ContactPersonKey, consts.LimitExceedsKey)
								if err != nil {
									logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
								}
								errMap[consts.ContactPersonKey] = models.ErrorResponse{
									Code:    code,
									Message: []string{consts.LimitExceedsKey},
								}
							}
						}

					}
					data.ContactDetails.ContactPerson = contactPersonName
				}

			}
		case consts.SubscriptionDetailsKey:
			for val := range subscriptionDetails {
				switch val {
				case consts.PlanIDKey:
					planID, err := arg.GetSubscriptionId()
					if err != nil {
						log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
						return nil, entities.Partner{}, err
					}

					if !utils.IsEmpty(planID) {
						data.SubscriptionPlanDetails.ID, err = partner.repo.GetID(ctx, consts.PartnerPlanTable, consts.IDKey, planID, consts.IDKey)
						if utils.IsEmpty(data.SubscriptionPlanDetails.ID) {
							code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PlanIDKey, consts.InvalidKey)
							if err != nil {
								logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
							}
							errMap[consts.PlanIDKey] = models.ErrorResponse{
								Code:    code,
								Message: []string{consts.InvalidKey},
							}
						}
						if err != nil {
							log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
							return nil, entities.Partner{}, err
						}
						data.SubscriptionPlanDetails.ID = planID
					}

				case consts.PlanStartDateKey:
					startDate, err := arg.GetStartDate()
					if err != nil {
						log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
						return nil, entities.Partner{}, err
					}
					data.SubscriptionPlanDetails.StartDate = startDate
				case consts.PlanLaunchDateKey:
					launchDate, err := arg.GetLaunchDate()
					if err != nil {
						log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
						return nil, entities.Partner{}, err
					}
					data.SubscriptionPlanDetails.LaunchDate = launchDate

				}

			}

		case consts.BackgroundColorKey:

			bgColor := arg.GetBGColor()
			if !utils.IsEmpty(bgColor) {
				if !utilities.IsValidHexColorCode(bgColor) {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.BackgroundColorKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.BackgroundColorKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				}
			}
			data.BackgroundColor = bgColor

		case consts.BackgroundImageKey:

			bgImage := arg.GetBGImage()
			if !utils.IsEmpty(bgImage) {
				if !utilities.IsValidURL(bgImage) {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.BackgroundImageKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.BackgroundImageKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				}
			}
			data.BackgroundImage = bgImage

		case consts.WebsiteURLKey:

			websiteURL := arg.GetWebsiteURL()
			if !utils.IsEmpty(websiteURL) {
				if !utilities.IsValidURL(websiteURL) {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.WebsiteURLKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.WebsiteURLKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				} else {
					if len(websiteURL) > consts.URLMaxLength {
						code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.WebsiteURLKey, consts.LimitExceedsKey)
						if err != nil {
							logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
						}
						errMap[consts.WebsiteURLKey] = models.ErrorResponse{
							Code:    code,
							Message: []string{consts.LimitExceedsKey},
						}
					} else {
						unique, err := partner.repo.IsFieldValueUnique(ctx, consts.WebsiteURLKey, websiteURL, partnerID)

						if err != nil {
							log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
							return nil, entities.Partner{}, err
						}

						if !unique {
							code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.WebsiteURLKey, consts.AlreadyExistsKey)
							if err != nil {
								logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
							}
							errMap[consts.WebsiteURLKey] = models.ErrorResponse{
								Code:    code,
								Message: []string{consts.AlreadyExistsKey},
							}
						}
					}
				}
			}
			data.WebsiteURL = websiteURL

		case consts.BrowserTitleKey:

			browserTitle := arg.GetBrowserTitle()
			browserTitleURL := strings.TrimSpace(browserTitle)
			if !utils.IsEmpty(browserTitleURL) {
				if len(browserTitleURL) > consts.BrowserTitleMaxLength {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.BrowserTitleKey, consts.LimitExceedsKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.BrowserTitleKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.LimitExceedsKey},
					}
				}
			}
			data.BrowserTitle = browserTitle

		case consts.ProfileURLKey:

			profileURL := arg.GetProfileURL()
			if !utils.IsEmpty(profileURL) {
				if !utilities.IsValidURL(profileURL) {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.ProfileURLKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.ProfileURLKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				} else {
					if len(profileURL) > consts.URLMaxLength {
						code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.ProfileURLKey, consts.LimitExceedsKey)
						if err != nil {
							logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
						}
						errMap[consts.ProfileURLKey] = models.ErrorResponse{
							Code:    code,
							Message: []string{consts.LimitExceedsKey},
						}
					}
				}

			}
			data.ProfileURL = profileURL

		case consts.PaymentURLKey:

			paymentURL := arg.GetPaymentURL()
			if !utils.IsEmpty(paymentURL) {
				if !utilities.IsValidURL(paymentURL) {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PaymentURLKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.PaymentURLKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				} else {
					if len(paymentURL) > consts.URLMaxLength {
						code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PaymentURLKey, consts.LimitExceedsKey)
						if err != nil {
							logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
						}
						errMap[consts.PaymentURLKey] = models.ErrorResponse{
							Code:    code,
							Message: []string{consts.LimitExceedsKey},
						}
					}
				}

			}
			data.PaymentURL = paymentURL

		case consts.SiteInfoKey:

			siteInfo := arg.GetSiteInfo()
			siteInfo = template.HTMLEscaper(siteInfo)
			siteInf := strings.TrimSpace(siteInfo)
			if !utils.IsEmpty(siteInf) {
				data.SiteInfo = siteInf
			}
			if !utils.IsEmpty(siteInf) && len(siteInf) > consts.SiteInfoMaxLength {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.SiteInfoKey, consts.LimitExceedsKey)
				if err != nil {
					logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
				}
				errMap[consts.SiteInfoKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.LimitExceedsKey},
				}
			}

		case consts.MusicLanguageKey:

			musicLanguage := arg.GetMusicLang()

			if !utils.IsEmpty(musicLanguage) {
				exists, err := utilities.IsLanguageIsoExists(ctx, partner.cache, musicLanguage, consts.UtilityServiceURL)
				if err != nil {
					log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
					return nil, entities.Partner{}, err
				}

				if !exists {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.MusicLanguageKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.MusicLanguageKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				}
				data.MusicLanguage = musicLanguage
			}

		case consts.MemberDefaultCountryKey:
			memberCountry := arg.GetMemberDefaultCountry()
			if !utils.IsEmpty(memberCountry) {
				exists, err := utilities.IsCountryExists(ctx, partner.cache, memberCountry, consts.UtilityServiceURL)
				if !exists {
					code, _ := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.MemberDefaultCountryKey, consts.InvalidKey)
					errMap[consts.MemberDefaultCountryKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				}

				if err != nil {
					log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
					return nil, entities.Partner{}, err
				}
				data.MemberDefaultCountry = memberCountry
			}
		case consts.ThemeKey:

			themeID := arg.GetThemeID()
			if themeID != 0 && themeID < 0 {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.ThemeKey, consts.InvalidKey)
				if err != nil {
					logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
				}
				errMap[consts.ThemeKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.InvalidKey},
				}
			}
			data.ThemeID = int(themeID)

		case consts.MobileVerifyIntervalKey:
			mobileVerifyInterval := arg.GetMobileVerifyInterval()
			if mobileVerifyInterval != 0 && mobileVerifyInterval < 0 {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.MobileVerifyIntervalKey, consts.InvalidKey)
				if err != nil {
					logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
				}
				errMap[consts.MobileVerifyIntervalKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.InvalidKey},
				}
			}
			data.MobileVerifyInterval = int(mobileVerifyInterval)

		case consts.LandingPageKey:

			landingPage := arg.GetLandingPage()
			if !utils.IsEmpty(landingPage) {
				if !utilities.IsValidURL(landingPage) {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.LandingPageKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.LandingPageKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				} else {
					if len(landingPage) > consts.URLMaxLength {
						code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.LandingPageKey, consts.LimitExceedsKey)
						if err != nil {
							logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
						}
						errMap[consts.LandingPageKey] = models.ErrorResponse{
							Code:    code,
							Message: []string{consts.LimitExceedsKey},
						}
					} else {
						unique, err := partner.repo.IsFieldValueUnique(ctx, consts.LandingPageKey, landingPage, partnerID)

						if err != nil {
							log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
							return nil, entities.Partner{}, err
						}

						if !unique {
							code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.LandingPageKey, consts.AlreadyExistsKey)
							if err != nil {
								logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
							}
							errMap[consts.LandingPageKey] = models.ErrorResponse{
								Code:    code,
								Message: []string{consts.AlreadyExistsKey},
							}
						}
					}
				}

			}
			data.LandingPage = landingPage

		case consts.AddressDetailsKey:
			for val := range addressDetails {
				switch val {
				case consts.AddressKey:
					address, err := arg.GetAddress()
					if err != nil {
						log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
						return nil, entities.Partner{}, err
					}
					if utils.IsEmpty(address) {
						code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.AddressKey, consts.RequiredKey)
						if err != nil {
							logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
						}
						errMap[consts.AddressKey] = models.ErrorResponse{
							Code:    code,
							Message: []string{consts.RequiredKey},
						}

					} else {
						address := strings.TrimSpace(address)
						if len(address) > consts.AddressMaxLength {
							code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.AddressKey, consts.LimitExceedsKey)
							if err != nil {
								logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
							}
							errMap[consts.AddressKey] = models.ErrorResponse{
								Code:    code,
								Message: []string{consts.LimitExceedsKey},
							}
						}
					}
					data.AddressDetails.Address = address

				case consts.StreetKey:
					street, err := arg.GetStreet()
					if err != nil {
						log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
						return nil, entities.Partner{}, err
					}
					partnerStreet := strings.TrimSpace(street)

					if !utils.IsEmpty(partnerStreet) {
						if len(partnerStreet) > consts.StreetMaxLength {
							code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.StreetKey, consts.LimitExceedsKey)
							if err != nil {
								logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
							}
							errMap[consts.StreetKey] = models.ErrorResponse{
								Code:    code,
								Message: []string{consts.LimitExceedsKey},
							}
						}
					}
					data.AddressDetails.Street = partnerStreet
				case consts.CountryKey:
					country, err := arg.GetCountry()
					if err != nil {
						log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
						return nil, entities.Partner{}, err
					}
					if !utils.IsEmpty(country) {
						exists, err := utilities.IsCountryExists(ctx, partner.cache, country, consts.UtilityServiceURL)
						if err != nil {
							log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
							return nil, entities.Partner{}, err
						}
						if !exists {
							code, _ := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.CountryKey, consts.InvalidKey)
							errMap[consts.CountryKey] = models.ErrorResponse{
								Code:    code,
								Message: []string{consts.InvalidKey},
							}
						}
						data.AddressDetails.Country = country
					}
					data.AddressDetails.Country = country
				case consts.StateKey:
					state, err := arg.GetState()
					if err != nil {
						log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
						return nil, entities.Partner{}, err
					}
					country, err := arg.GetCountry()
					if err != nil {
						log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
						return nil, entities.Partner{}, err
					}
					if !utils.IsEmpty(state) {
						exists, err := utilities.IsStateIsoExists(ctx, partner.cache, country, state, consts.UtilityServiceURL)
						if err != nil {
							log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
							return nil, entities.Partner{}, err

						}
						if !exists {
							code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.StateKey, consts.InvalidKey)
							if err != nil {
								logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
							}
							errMap[consts.StateKey] = models.ErrorResponse{
								Code:    code,
								Message: []string{consts.InvalidKey},
							}
						}

						data.AddressDetails.State = state
					}
				case consts.CityKey:
					city, err := arg.GetCity()
					if err != nil {
						log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
						return nil, entities.Partner{}, err
					}
					partnerCity := strings.TrimSpace(city)
					if len(partnerCity) > consts.CityMaxLength {
						code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.CityKey, consts.LimitExceedsKey)
						if err != nil {
							logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
						}

						errMap[consts.CityKey] = models.ErrorResponse{
							Code:    code,
							Message: []string{consts.LimitExceedsKey},
						}
					}
					data.AddressDetails.City = partnerCity
				case consts.PostalCodeKey:
					postalCode, err := arg.GetPostalCode()
					if err != nil {
						log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
						return nil, entities.Partner{}, err
					}
					if !utils.IsEmpty(postalCode) {
						isValid := utilities.IsValidPostalCode(postalCode)
						if !isValid {
							code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PostalCodeKey, consts.InvalidKey)
							if err != nil {
								logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
							}
							errMap[consts.PostalCodeKey] = models.ErrorResponse{
								Code:    code,
								Message: []string{consts.InvalidKey},
							}
						}
					}
					data.AddressDetails.PostalCode = postalCode

				}
			}

		case consts.AlbumReviewEmailKey:

			albumEmail := arg.GetAlbumEmail()
			if utils.IsEmpty(albumEmail) {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.AlbumReviewEmailKey, consts.RequiredKey)
				if err != nil {
					logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
				}
				errMap[consts.AlbumReviewEmailKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.RequiredKey},
				}

			} else {
				if !utils.IsValidEmail(albumEmail) {

					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.AlbumReviewEmailKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.AlbumReviewEmailKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				} else {
					if len(albumEmail) > consts.EmailMaxLength {
						code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.AlbumReviewEmailKey, consts.LimitExceedsKey)
						if err != nil {
							logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
						}
						errMap[consts.AlbumReviewEmailKey] = models.ErrorResponse{
							Code:    code,
							Message: []string{consts.LimitExceedsKey},
						}
					} else {
						unique, err := partner.repo.IsFieldValueUnique(ctx, consts.AlbumReviewEmailKey, albumEmail, partnerID)
						if err != nil {
							log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
							return nil, entities.Partner{}, err
						}

						if !unique {
							code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.AlbumReviewEmailKey, consts.AlreadyExistsKey)
							if err != nil {
								logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
							}
							errMap[consts.AlbumReviewEmailKey] = models.ErrorResponse{
								Code:    code,
								Message: []string{consts.AlreadyExistsKey},
							}
						}
					}
				}

			}
			data.AlbumReviewEmail = albumEmail

		case consts.FreePlanLimitKey:

			freePlanLimit := arg.GetFreePlanLimit()
			if freePlanLimit != 0 {
				if freePlanLimit > consts.MaxFreePlanLimit || freePlanLimit < 0 {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.FreePlanLimitKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.FreePlanLimitKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				} else {
					data.FreePlanLimit = int(freePlanLimit)

				}
			}

		case consts.ExpiryWarningCountKey:

			expiryWarningCount := arg.GetExpiryWarningCount()
			if expiryWarningCount != 0 {
				if expiryWarningCount > consts.MaxExpiryWarningCount || expiryWarningCount < 0 {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.ExpiryWarningCountKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.ExpiryWarningCountKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				} else {
					data.ExpiryWarningCount = int(expiryWarningCount)

				}
			}

		case consts.OutletsProcessingDurationKey:

			outletsProcessingDuration := arg.GetOutletsProcessingDuration()
			if outletsProcessingDuration != 0 {
				if outletsProcessingDuration < consts.MinOutletsProcessingDuration || outletsProcessingDuration > consts.MaxOutletsProcessingDuration {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.OutletsProcessingDurationKey, consts.LimitExceedsKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.OutletsProcessingDurationKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.LimitExceedsKey},
					}
				}
			}
			data.OutletsProcessingDuration = int(outletsProcessingDuration)

		case consts.PaymentKey:
			var gateways []string

			paymentGateways, err := arg.GetPaymentGateways()
			if err != nil {
				log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
				return nil, entities.Partner{}, err

			}
			for i := range paymentGateways {
				gateways = append(gateways, paymentGateways[i].Gateway)
				if paymentGateways[i].Gateway == "" {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PaymentGatewayKey, consts.RequiredKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.PaymentGatewayKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.RequiredKey},
					}
					break

				}
				paymentGateways[i].GatewayId, err = utilities.GetPaymentGatewayId(ctx, partner.cache, paymentGateways[i].Gateway, consts.UtilityServiceURL)
				if err != nil {
					logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					return nil, entities.Partner{}, err
				}
				if paymentGateways[i].GatewayId == "" {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PaymentGatewayKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.PaymentGatewayKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
				}

			}

			for i := range paymentGateways {
				if paymentGateways[i].ClientSecret == "" {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.ClientSecretKey, consts.RequiredKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.ClientSecretKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.RequiredKey},
					}
					break

				}
			}
			for i := range paymentGateways {
				if paymentGateways[i].Email == "" {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PaymentGatewayEmailKey, consts.RequiredKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.PaymentGatewayEmailKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.RequiredKey},
					}
					break

				}
				if !utils.IsValidEmail(paymentGateways[i].Email) {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PaymentGatewayEmailKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.PaymentGatewayEmailKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
					break

				}
			}
			for i := range paymentGateways {
				if paymentGateways[i].ClientId == "" {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.ClientIdKey, consts.RequiredKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.ClientIdKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.RequiredKey},
					}
					break

				}
			}
			for i := range paymentGateways {
				if paymentGateways[i].DefaultPayoutCurrency == "" {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.DefaultPayoutCurrencyKey, consts.RequiredKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.DefaultPayoutCurrencyKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.RequiredKey},
					}
					break
				}

				if !slices.Contains(consts.SupportedCurrencies, paymentGateways[i].DefaultPayoutCurrency) {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.DefaultPayoutCurrencyKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.DefaultPayoutCurrencyKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
					break

				}
			}
			for i := range paymentGateways {
				if paymentGateways[i].DefaultPayinCurrency == "" {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.DefaultPayinCurrencyKey, consts.RequiredKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.DefaultPayinCurrencyKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.RequiredKey},
					}
					break

				}

				if !slices.Contains(consts.SupportedCurrencies, paymentGateways[i].DefaultPayinCurrency) {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.DefaultPayinCurrencyKey, consts.InvalidKey)
					if err != nil {
						logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
					}
					errMap[consts.DefaultPayinCurrencyKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.InvalidKey},
					}
					break

				}
			}

			data.PaymentGatewayDetails.PaymentGateways = paymentGateways
			for val := range gatewayDetails {
				switch val {

				case consts.PayoutMinLimitKey:
					payoutMinLimit, err := arg.GetPayoutMinLmit()
					if err != nil {
						log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
						return nil, entities.Partner{}, err

					}

					if payoutMinLimit == 0 {
						data.PaymentGatewayDetails.PayoutMinLimit = consts.PayoutMinLimitDefaultVal
					} else {
						if payoutMinLimit < 0 {
							code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PayoutMinLimitKey, consts.InvalidKey)
							if err != nil {
								logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
							}
							errMap[consts.PayoutMinLimitKey] = models.ErrorResponse{
								Code:    code,
								Message: []string{consts.InvalidKey},
							}
						}
					}

					data.PaymentGatewayDetails.PayoutMinLimit = int(payoutMinLimit)

				case consts.MaxRemittancePerMonthKey:
					maxRemittance, err := arg.GetMaxRemittance()
					if err != nil {
						log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
						return nil, entities.Partner{}, err

					}

					if maxRemittance != 0 {
						if maxRemittance > consts.MaxRemittancePerMonth || maxRemittance < 0 {
							code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.MaxRemittancePerMonthKey, consts.InvalidKey)
							if err != nil {
								logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
							}
							errMap[consts.MaxRemittancePerMonthKey] = models.ErrorResponse{
								Code:    code,
								Message: []string{consts.InvalidKey},
							}
						}
					} else {
						data.PaymentGatewayDetails.MaxRemittancePerMonth = int(maxRemittance)
					}

				case consts.DefaultCurrencyKey:
					defaultCurrency, err := arg.GetDefaultCurrency()
					if err != nil {
						log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
						return nil, entities.Partner{}, err

					}
					if !utils.IsEmpty(defaultCurrency) {
						data.DefaultCurrencyID, err = utilities.GetCurrencyId(ctx, partner.cache, defaultCurrency, consts.UtilityServiceURL)
						if err != nil {
							log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
							return nil, entities.Partner{}, err
						}
						if data.DefaultCurrencyID == 0 {
							code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PayoutCurrencyKey, consts.InvalidKey)
							if err != nil {
								logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
							}
							errMap[consts.PayoutCurrencyKey] = models.ErrorResponse{
								Code:    code,
								Message: []string{consts.InvalidKey},
							}
						}
					}
				case consts.DefaultPaymentGatewayKey:
					defaultGateway, err := arg.GetDefaultGateway()
					if err != nil {
						log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
						return nil, entities.Partner{}, err

					}
					if utils.IsEmpty(defaultGateway) {
						code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.DefaultPaymentGatewayKey, consts.RequiredKey)
						if err != nil {
							logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
						}
						errMap[consts.DefaultPaymentGatewayKey] = models.ErrorResponse{
							Code:    code,
							Message: []string{consts.RequiredKey},
						}

					} else {
						if !slices.Contains(gateways, defaultGateway) {
							code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.DefaultPaymentGatewayKey, consts.InvalidKey)
							if err != nil {
								logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerErrMsg, err)
							}
							errMap[consts.DefaultPaymentGatewayKey] = models.ErrorResponse{
								Code:    code,
								Message: []string{consts.InvalidKey},
							}

						} else {
							data.DefaultPaymentGatewayId, err = utilities.GetPaymentGatewayId(ctx, partner.cache, defaultGateway, consts.UtilityServiceURL)
							if utils.IsEmpty(data.DefaultPaymentGatewayId) {
								code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.DefaultPaymentGatewayKey, consts.NotFoundKey)
								if err != nil {
									logger.Log().WithContext(ctx).Error(consts.UpdatePartnerErrMsg, err)
								}
								errMap[consts.DefaultPaymentGatewayKey] = models.ErrorResponse{
									Code:    code,
									Message: []string{consts.NotFoundKey},
								}
							}
							if err != nil {
								log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
								return nil, entities.Partner{}, err
							}
						}
					}
					data.PaymentGatewayDetails.DefaultPaymentGateway = defaultGateway
				}

			}

		}
	}
	return errMap, data, nil
}

// function to update partner details based on partner id
func (partner *PartnerUseCases) UpdatePartner(ctx context.Context, partnerID string, memberID uuid.UUID, partnerData entities.PartnerProperties, endpoint string, method string, errMap map[string]models.ErrorResponse) (map[string]models.ErrorResponse, error) {

	var (
		log               = logger.Log().WithContext(ctx)
		updatePartnerData entities.Partner
	)
	errMap, data, err := partner.ValidatePartnerUpdate(ctx, partnerID, partnerData, updatePartnerData, endpoint, method, errMap)
	if err != nil {
		log.Errorf(consts.UpdatePartnerErrMsg, err)
		return nil, err
	}
	if len(errMap) != 0 {
		log.Errorf(consts.UpdatePartnerErrMsg, "validation error")
		return errMap, nil
	}
	_, err = partner.repo.UpdatePartner(ctx, partnerID, memberID, &data, &partnerData, endpoint, method, errMap)

	if err != nil {
		log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
		return nil, err
	}

	return nil, nil
}

// function to check partner existence based on partner id
func (partner *PartnerUseCases) IsPartnerExists(ctx context.Context, partnerID string, endpoint string, method string, errMap map[string]models.ErrorResponse) (bool, error) {

	_, err := uuid.Parse(partnerID)

	if err != nil {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PartnerIDKey, consts.InvalidKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error("IsPartnerExists failed", err)
		}

		errMap[consts.PartnerIDKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}
		return false, err
	}
	exists, err := partner.repo.IsPartnerExists(ctx, partnerID)
	if exists {
		return true, nil
	}

	if err != nil {
		return false, err
	}
	code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PartnerIDKey, consts.NotFoundKey)
	if err != nil {
		logger.Log().WithContext(ctx).Error("IsPartnerExists failed", err)
	}
	errMap[consts.PartnerIDKey] = models.ErrorResponse{
		Code:    code,
		Message: []string{consts.NotFoundKey},
	}
	return false, nil
}

// function to update terms and conditions of a partner by partner id
func (partner *PartnerUseCases) UpdateTermsAndConditions(ctx context.Context, partnerID string, memberID uuid.UUID, termsAndConditionsData entities.UpdateTermsAndConditions, endpoint string, method string, errMap map[string]models.ErrorResponse) (map[string]models.ErrorResponse, error) {
	var (
		log = logger.Log().WithContext(ctx)
		err error
	)
	errMap, err = partner.TermsAndConditionValidations(ctx, endpoint, method, errMap, termsAndConditionsData)
	if len(errMap) != 0 {
		log.Errorf(consts.UpdateTermsAndConditionsErrMsg, "validation error")
		return errMap, nil
	}

	if err != nil {
		log.Errorf(consts.UpdateTermsAndConditionsErrMsg, err.Error())
		return nil, err
	}
	err = partner.repo.UpdateTermsAndConditions(ctx, partnerID, memberID, termsAndConditionsData, endpoint, method, errMap)
	if len(errMap) != 0 {
		log.Errorf(consts.UpdateTermsAndConditionsErrMsg, "validation error")
		return errMap, nil
	}

	if err != nil {
		log.Errorf(consts.UpdateTermsAndConditionsErrMsg, err.Error())
		return nil, err
	}
	return errMap, nil
}

func (partner *PartnerUseCases) TermsAndConditionValidations(ctx context.Context, endpoint string, method string, errMap map[string]models.ErrorResponse, data entities.UpdateTermsAndConditions) (map[string]models.ErrorResponse, error) {

	for argField := range data {
		switch argField {
		case consts.NameKey:
			name := data.GetData(consts.NameKey)
			if utils.IsEmpty(name) {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.TermsAndConditionsNameKey, consts.RequiredKey)
				if err != nil {
					logger.Log().WithContext(ctx).Error(consts.UpdateTermsAndConditionsErrMsg, err)
				}
				errMap[consts.TermsAndConditionsNameKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.RequiredKey},
				}

			} else {
				if len(name) > consts.TermsAndConditionsNameMaxLength {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.TermsAndConditionsNameKey, consts.LimitExceedsKey)
					if err != nil {
						logger.Log().WithContext(ctx).Error(consts.UpdateTermsAndConditionsErrMsg, err)
					}
					errMap[consts.TermsAndConditionsNameKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.LimitExceedsKey},
					}
				}
			}
		case consts.DescriptionKey:
			description := data.GetData(consts.DescriptionKey)
			if utils.IsEmpty(description) {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.TermsAndConditionsDescriptionKey, consts.RequiredKey)
				if err != nil {
					logger.Log().WithContext(ctx).Error(consts.UpdateTermsAndConditionsErrMsg, err)
				}
				errMap[consts.TermsAndConditionsDescriptionKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.RequiredKey},
				}

			}
		case consts.LanguageKey:
			language := data.GetData(consts.LanguageKey)
			if utils.IsEmpty(language) {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.LanguageKey, consts.RequiredKey)
				if err != nil {
					logger.Log().WithContext(ctx).Error(consts.UpdateTermsAndConditionsErrMsg, err)
				}
				errMap[consts.LanguageKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.RequiredKey},
				}

			}
		}

	}

	return errMap, nil
}

// function to create a partner
func (partner *PartnerUseCases) CreatePartner(ctx context.Context, newPartner entities.Partner, endpoint string, method string, errMap map[string]models.ErrorResponse) (map[string]models.ErrorResponse, string, error) {
	var (
		log = logger.Log().WithContext(ctx)
		err error
	)
	SetPartnerDefaultVal(&newPartner)
	errMap, err = partner.PartnerValidations(ctx, endpoint, method, errMap, &newPartner)
	if err != nil {
		log.Errorf(consts.CreatePartnerErrMsg, err)
		return nil, "", err
	}
	if len(errMap) != 0 {
		log.Errorf(consts.CreatePartnerErrMsg, "validation error")
		return errMap, "", nil
	}
	OauthCredentials, err := partner.GeneratePartnerOauthcredential(ctx)
	if err != nil {
		log.Errorf(consts.CreatePartnerErrMsg, err)
		return nil, "", err
	}
	id, err := partner.repo.CreatePartner(ctx, endpoint, method, errMap, newPartner, OauthCredentials)

	if err != nil {
		log.Errorf(consts.CreatePartnerErrMsg, err.Error())
		return nil, "", err
	}
	if len(errMap) != 0 {
		log.Errorf(consts.CreatePartnerErrMsg, "validation error")
		return errMap, "", nil
	}

	return nil, id, nil
}

// function generate client id and client secret for a partner
func (partner *PartnerUseCases) GeneratePartnerOauthcredential(ctx context.Context) (entities.PartnerOauthCredential, error) {
	var (
		oauthCredentials entities.PartnerOauthCredential
		err              error
		log              = logger.Log().WithContext(ctx)
	)
	encryptionKey := consts.ClientCredentialEncryptionKey
	key := []byte(encryptionKey)

	// oauthCredentials.PartnerId = partnerId
	oauthCredentials.ClientId, err = utilities.GenerateClientCredentials()
	if err != nil {
		log.Errorf(consts.PartnerOauthCredentialErrMsg, err.Error())
		return entities.PartnerOauthCredential{}, err
	}
	oauthCredentials.ClientSecret, err = utilities.GenerateClientCredentials()
	if err != nil {
		log.Errorf(consts.PartnerOauthCredentialErrMsg, err.Error())
		return entities.PartnerOauthCredential{}, err
	}
	clientId, err := crypto.Encrypt(oauthCredentials.ClientId, key)
	if err != nil {
		log.Errorf(consts.PartnerOauthCredentialErrMsg, err.Error())
		return entities.PartnerOauthCredential{}, err
	}
	oauthCredentials.ClientId = clientId
	secret, err := crypto.Encrypt(oauthCredentials.ClientSecret, key)
	if err != nil {
		log.Errorf(consts.PartnerOauthCredentialErrMsg, err.Error())
		return entities.PartnerOauthCredential{}, err
	}
	oauthCredentials.ClientSecret = secret
	oauthCredentials.ProviderId, err = utilities.GetOauthProviderId(ctx, partner.cache, consts.InternalVal, consts.OAuthServiceURL)
	if err != nil {
		log.Errorf(consts.PartnerOauthCredentialErrMsg, err.Error())
		return entities.PartnerOauthCredential{}, err
	}
	return oauthCredentials, nil
}

// function to validate partner fields while ceating a partner
func (partner *PartnerUseCases) PartnerValidations(ctx context.Context, endpoint string, method string, errMap map[string]models.ErrorResponse, data *entities.Partner) (map[string]models.ErrorResponse, error) {

	var (
		log = logger.Log().WithContext(ctx)
		err error
	)
	//Name field validation
	if utils.IsEmpty(data.Name) {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.NameKey, consts.RequiredKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.NameKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.RequiredKey},
		}
	} else {
		isExists, err := partner.repo.IsExists(ctx, consts.PartnerTable, consts.NameKey, data.Name)
		if isExists {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.NameKey, consts.AlreadyExistsKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.NameKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.AlreadyExistsKey},
			}
		} else {
			if !utilities.IsValidPartnerName(data.Name) {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.NameKey, consts.InvalidKey)
				if err != nil {
					logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
				}
				errMap[consts.NameKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.InvalidKey},
				}
			} else {
				name := strings.TrimSpace(data.Name)
				if len(name) > consts.PartnerNameMaxLength {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.NameKey, consts.LimitExceedsKey)
					if err != nil {
						logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
					}
					errMap[consts.NameKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.LimitExceedsKey},
					}
				}
			}

		}
		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err)
			return nil, err
		}

	}

	// Url field validation
	if utils.IsEmpty(data.URL) {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.URLKey, consts.RequiredKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.URLKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.RequiredKey},
		}

	} else {
		isExists, err := partner.repo.IsExists(ctx, consts.PartnerTable, consts.URLKey, data.URL)
		if isExists {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.URLKey, consts.AlreadyExistsKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.URLKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.AlreadyExistsKey},
			}

		} else {
			if !utilities.IsValidURL(data.URL) {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.URLKey, consts.InvalidKey)
				if err != nil {
					logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
				}
				errMap[consts.URLKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.InvalidKey},
				}

			} else {
				if len(data.URL) > consts.URLMaxLength {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.URLKey, consts.LimitExceedsKey)
					if err != nil {
						logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
					}
					errMap[consts.URLKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.LimitExceedsKey},
					}

				}
			}

		}

		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err)
			return nil, err
		}

	}

	// Logo field validation
	if utils.IsEmpty(data.Logo) {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.LogoKey, consts.RequiredKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.LogoKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.RequiredKey},
		}

	} else {
		if !utilities.IsValidURL(data.Logo) {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.LogoKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.LogoKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}

		}
	}

	//language validation

	if !utils.IsEmpty(data.Language) {
		exists, err := utilities.IsLanguageIsoExists(ctx, partner.cache, data.Language, consts.UtilityServiceURL)
		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err.Error())
			return nil, err

		}
		if !exists {
			code, _ := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.LanguageKey, consts.InvalidKey)
			errMap[consts.LanguageKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}
		}
	}

	// payout target currency validation
	if !utils.IsEmpty(data.PayoutTargetCurrency) {
		data.PayoutTargetCurrencyID, err = utilities.GetCurrencyId(ctx, partner.cache, data.PayoutTargetCurrency, consts.UtilityServiceURL)
		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err.Error())
			return nil, err
		}
		if data.PayoutTargetCurrencyID == 0 {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PayoutTargetCurrencyKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.PayoutTargetCurrencyKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}

		}
	}

	// BusinessModel field validation
	if utils.IsEmpty(data.BusinessModel) {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.BusinessModelKey, consts.RequiredKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.BusinessModelKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.RequiredKey},
		}

	} else {

		data.BusinessModelID, err = utilities.GetLookupId(ctx, partner.cache, consts.BusinessLookupTypeName, data.BusinessModel, consts.UtilityServiceURL)
		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err.Error())
			return nil, err

		}
		if data.BusinessModelID == 0 {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.BusinessModelKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.BusinessModelKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}

		}
	}

	// ContactPerson field validation
	if utils.IsEmpty(data.ContactDetails.ContactPerson) {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.ContactPersonKey, consts.RequiredKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.ContactPersonKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.RequiredKey},
		}

	} else {
		contactPerson := strings.TrimSpace(data.ContactDetails.ContactPerson)
		if len(contactPerson) > consts.ContactPersonMaxLength {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.ContactPersonKey, consts.LimitExceedsKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.ContactPersonKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.LimitExceedsKey},
			}
		}
	}

	// Email field validation
	if utils.IsEmpty(data.ContactDetails.Email) {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.EmailKey, consts.RequiredKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.EmailKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.RequiredKey},
		}

	} else {
		isExists, err := partner.repo.IsExists(ctx, consts.PartnerTable, consts.EmailKey, data.ContactDetails.Email)
		if isExists {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.EmailKey, consts.AlreadyExistsKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.EmailKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.AlreadyExistsKey},
			}

		} else {
			if !utils.IsValidEmail(data.ContactDetails.Email) {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.EmailKey, consts.InvalidKey)
				if err != nil {
					logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
				}
				errMap[consts.EmailKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.InvalidKey},
				}

			} else {
				if len(data.ContactDetails.Email) > consts.EmailMaxLength {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.EmailKey, consts.LimitExceedsKey)
					if err != nil {
						logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
					}
					errMap[consts.EmailKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.LimitExceedsKey},
					}
				}
			}

		}
		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err)
			return nil, err
		}

	}

	// NoReplyEmail field validation
	if utils.IsEmpty(data.ContactDetails.NoReplyEmail) {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.NoReplyEmailKey, consts.RequiredKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.NoReplyEmailKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.RequiredKey},
		}

	} else {
		isExists, err := partner.repo.IsExists(ctx, consts.PartnerTable, consts.NoReplyEmailKey, data.ContactDetails.NoReplyEmail)
		if isExists {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.NoReplyEmailKey, consts.AlreadyExistsKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.NoReplyEmailKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.AlreadyExistsKey},
			}

		} else {
			if !utils.IsValidEmail(data.ContactDetails.NoReplyEmail) {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.NoReplyEmailKey, consts.InvalidKey)
				if err != nil {
					logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
				}
				errMap[consts.NoReplyEmailKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.InvalidKey},
				}

			} else {
				if len(data.ContactDetails.NoReplyEmail) > consts.EmailMaxLength {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.NoReplyEmailKey, consts.LimitExceedsKey)
					if err != nil {
						logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
					}
					errMap[consts.NoReplyEmailKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.LimitExceedsKey},
					}
				}
			}

		}
		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err)
			return nil, err
		}

	}

	// SupportEmail field validation
	if utils.IsEmpty(data.ContactDetails.SupportEmail) {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.SupportEmailKey, consts.RequiredKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.SupportEmailKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.RequiredKey},
		}

	} else {
		isExists, err := partner.repo.IsExists(ctx, consts.PartnerTable, consts.SupportEmailKey, data.ContactDetails.SupportEmail)
		if isExists {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.SupportEmailKey, consts.AlreadyExistsKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.SupportEmailKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.AlreadyExistsKey},
			}

		} else {
			if !utils.IsValidEmail(data.ContactDetails.SupportEmail) {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.SupportEmailKey, consts.InvalidKey)
				if err != nil {
					logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
				}
				errMap[consts.SupportEmailKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.InvalidKey},
				}

			} else {
				if len(data.ContactDetails.SupportEmail) > consts.EmailMaxLength {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.SupportEmailKey, consts.LimitExceedsKey)
					if err != nil {
						logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
					}
					errMap[consts.SupportEmailKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.LimitExceedsKey},
					}

				}
			}

		}
		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err)
			return nil, err
		}

	}
	// FeedBack Email field validation
	if utils.IsEmpty(data.ContactDetails.FeedbackEmail) {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.FeedbackEmailKey, consts.RequiredKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.FeedbackEmailKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.RequiredKey},
		}

	} else {
		isExists, err := partner.repo.IsExists(ctx, consts.PartnerTable, consts.FeedbackEmailKey, data.ContactDetails.FeedbackEmail)
		if isExists {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.FeedbackEmailKey, consts.AlreadyExistsKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.FeedbackEmailKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.AlreadyExistsKey},
			}

		} else {
			if !utils.IsValidEmail(data.ContactDetails.FeedbackEmail) {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.FeedbackEmailKey, consts.InvalidKey)
				if err != nil {
					logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
				}
				errMap[consts.FeedbackEmailKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.InvalidKey},
				}

			} else {
				if len(data.ContactDetails.FeedbackEmail) > consts.EmailMaxLength {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.FeedbackEmailKey, consts.LimitExceedsKey)
					if err != nil {
						logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
					}
					errMap[consts.FeedbackEmailKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.LimitExceedsKey},
					}
				}
			}

		}
		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err)
			return nil, err
		}

	}

	// AlbumReviewEmail field validation
	if utils.IsEmpty(data.AlbumReviewEmail) {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.AlbumReviewEmailKey, consts.RequiredKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.AlbumReviewEmailKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.RequiredKey},
		}

	} else {
		isExists, err := partner.repo.IsExists(ctx, consts.PartnerTable, consts.AlbumReviewEmailKey, data.AlbumReviewEmail)
		if isExists {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.AlbumReviewEmailKey, consts.AlreadyExistsKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.AlbumReviewEmailKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.AlreadyExistsKey},
			}

		} else {
			if !utils.IsValidEmail(data.AlbumReviewEmail) {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.AlbumReviewEmailKey, consts.InvalidKey)
				if err != nil {
					logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
				}
				errMap[consts.AlbumReviewEmailKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.InvalidKey},
				}

			} else {
				if len(data.AlbumReviewEmail) > consts.EmailMaxLength {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.AlbumReviewEmailKey, consts.LimitExceedsKey)
					if err != nil {
						logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
					}
					errMap[consts.AlbumReviewEmailKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.LimitExceedsKey},
					}
				}
			}

		}
		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err)
			return nil, err
		}

	}

	// website url validation
	if !utils.IsEmpty(data.WebsiteURL) {
		isExists, err := partner.repo.IsExists(ctx, consts.PartnerTable, consts.WebsiteURLKey, data.WebsiteURL)
		if isExists {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.WebsiteURLKey, consts.AlreadyExistsKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.WebsiteURLKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.AlreadyExistsKey},
			}

		} else {
			if !utilities.IsValidURL(data.WebsiteURL) {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.WebsiteURLKey, consts.InvalidKey)
				if err != nil {
					logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
				}
				errMap[consts.WebsiteURLKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.InvalidKey},
				}

			} else {
				if len(data.WebsiteURL) > consts.URLMaxLength {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.WebsiteURLKey, consts.LimitExceedsKey)
					if err != nil {
						logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
					}
					errMap[consts.WebsiteURLKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.LimitExceedsKey},
					}
				}
			}

		}
		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err)
			return nil, err
		}

	}

	// landing page validation
	if !utils.IsEmpty(data.LandingPage) {
		isExists, err := partner.repo.IsExists(ctx, consts.PartnerTable, consts.LandingPageKey, data.LandingPage)
		if isExists {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.LandingPageKey, consts.AlreadyExistsKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.LandingPageKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.AlreadyExistsKey},
			}

		} else {
			if !utilities.IsValidURL(data.LandingPage) {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.LandingPageKey, consts.InvalidKey)
				if err != nil {
					logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
				}
				errMap[consts.LandingPageKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.InvalidKey},
				}

			} else {
				if len(data.LandingPage) > consts.URLMaxLength {
					code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.LandingPageKey, consts.LimitExceedsKey)
					if err != nil {
						logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
					}
					errMap[consts.LandingPageKey] = models.ErrorResponse{
						Code:    code,
						Message: []string{consts.LimitExceedsKey},
					}
				}
			}

		}
		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err)
			return nil, err
		}

	}

	// Profile url validation
	if !utils.IsEmpty(data.ProfileURL) {
		if !utilities.IsValidURL(data.ProfileURL) {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.ProfileURLKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.ProfileURLKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}

		} else {
			if len(data.ProfileURL) > consts.URLMaxLength {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.ProfileURLKey, consts.LimitExceedsKey)
				if err != nil {
					logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
				}
				errMap[consts.ProfileURLKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.LimitExceedsKey},
				}
			}
		}

	}

	// Payment url validation
	if !utils.IsEmpty(data.PaymentURL) {
		if !utilities.IsValidURL(data.PaymentURL) {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PaymentURLKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.PaymentURLKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}

		} else {
			if len(data.PaymentURL) > consts.URLMaxLength {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PaymentURLKey, consts.LimitExceedsKey)
				if err != nil {
					logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
				}
				errMap[consts.PaymentURLKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.LimitExceedsKey},
				}
			}
		}

	}

	// Address field validation
	if utils.IsEmpty(data.AddressDetails.Address) {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.AddressKey, consts.RequiredKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.AddressKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.RequiredKey},
		}

	} else {
		address := strings.TrimSpace(data.AddressDetails.Address)
		if len(address) > consts.AddressMaxLength {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.AddressKey, consts.LimitExceedsKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.AddressKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.LimitExceedsKey},
			}

		}
	}

	//city field validation
	if !utils.IsEmpty(data.AddressDetails.City) {
		city := strings.TrimSpace(data.AddressDetails.City)
		if len(city) > consts.CityMaxLength {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.CityKey, consts.LimitExceedsKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.CityKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.LimitExceedsKey},
			}

		}
	}

	// street field validation
	if !utils.IsEmpty(data.AddressDetails.Street) {
		street := strings.TrimSpace(data.AddressDetails.Street)
		if len(street) > consts.StreetMaxLength {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.StreetKey, consts.LimitExceedsKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.StreetKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.LimitExceedsKey},
			}

		}
	}

	// Country field validation
	if utils.IsEmpty(data.AddressDetails.Country) {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.CountryKey, consts.RequiredKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.CountryKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.RequiredKey},
		}

	} else {
		exists, err := utilities.IsCountryExists(ctx, partner.cache, data.AddressDetails.Country, consts.UtilityServiceURL)
		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err.Error())
			return nil, err
		}
		if !exists {
			code, _ := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.CountryKey, consts.InvalidKey)
			errMap[consts.CountryKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}
		}
	}

	// country state validation
	if !utils.IsEmpty(data.AddressDetails.State) {

		exists, err := utilities.IsStateIsoExists(ctx, partner.cache, data.AddressDetails.Country, data.AddressDetails.State, consts.UtilityServiceURL)
		if !exists {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.StateKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.StateKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}
		}

		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err.Error())
			return nil, err

		}
	}

	//music language validation
	if !utils.IsEmpty(data.MusicLanguage) {
		exists, err := utilities.IsLanguageIsoExists(ctx, partner.cache, data.MusicLanguage, consts.UtilityServiceURL)
		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err.Error())
		}
		if !exists {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.MusicLanguageKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.MusicLanguageKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}
		}
	}

	// Payment_gateway field validation
	if data.PaymentGatewayDetails.PaymentGateways == nil {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PaymentGateWayIDKey, consts.RequiredKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.PaymentGateWayIDKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.RequiredKey},
		}

	}

	//background color validations
	if data.BackgroundColor != "" {
		if !utilities.IsValidHexColorCode(data.BackgroundColor) {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.BackgroundColorKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.BackgroundColorKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}

		}
	}

	//background image validations
	if !utils.IsEmpty(data.BackgroundImage) {
		if !utilities.IsValidURL(data.BackgroundImage) {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.BackgroundImageKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.BackgroundImageKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}

		}
	}

	// Free Plan Limit validation
	if data.FreePlanLimit > consts.MaxFreePlanLimit || data.FreePlanLimit < 0 {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.FreePlanLimitKey, consts.InvalidKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.FreePlanLimitKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}

	}
	// Mobile verify Interval validation
	if data.MobileVerifyInterval < 0 {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.MobileVerifyIntervalKey, consts.InvalidKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.MobileVerifyIntervalKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}

	}

	// Expiry Warning count validation
	if data.ExpiryWarningCount > consts.MaxExpiryWarningCount || data.ExpiryWarningCount < 0 {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.ExpiryWarningCountKey, consts.InvalidKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.ExpiryWarningCountKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}

	}

	// theme validation
	if data.ThemeID != 0 {
		theme, err := utilities.GetTheme(ctx, partner.cache, data.ThemeID, consts.UtilityServiceURL)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		if theme == "" {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.ThemeKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.ThemeKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}
		}
	}

	//login type validation
	if !utils.IsEmpty(data.LoginType) {
		data.LoginTypeID, err = utilities.GetLookupId(ctx, partner.cache, consts.LoginLookupTypeName, data.LoginType, consts.UtilityServiceURL)

		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err.Error())
			return nil, err
		}
		if data.LoginTypeID == 0 {

			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.LoginTypeKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.LoginTypeKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}
		}
	}

	//default price code  currency validation
	if !utils.IsEmpty(data.DefaultPriceCodeCurrency.Name) {
		data.DefaultPriceCodeCurrencyID, err = utilities.GetCurrencyId(ctx, partner.cache, data.DefaultPriceCodeCurrency.Name, consts.UtilityServiceURL)
		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err.Error())
			return nil, err
		}
		if data.DefaultPriceCodeCurrencyID == 0 {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.DefaultPriceCodeCurrencyKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.DefaultPriceCodeCurrencyKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}
		}
	}

	// MaxRemittancePerMonth validation
	if data.PaymentGatewayDetails.MaxRemittancePerMonth > consts.MaxRemittancePerMonth || data.PaymentGatewayDetails.MaxRemittancePerMonth < 0 {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.MaxRemittancePerMonthKey, consts.InvalidKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.MaxRemittancePerMonthKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}

	}

	// OutletsProcessingDuration validation
	if data.OutletsProcessingDuration < consts.MinOutletsProcessingDuration || data.OutletsProcessingDuration > consts.MaxOutletsProcessingDuration {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.OutletsProcessingDurationKey, consts.LimitExceedsKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.OutletsProcessingDurationKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.LimitExceedsKey},
		}

	}

	//BrowserTitle validation
	if !utils.IsEmpty(data.BrowserTitle) {
		browserTitle := strings.TrimSpace(data.BrowserTitle)
		if len(browserTitle) > consts.BrowserTitleMaxLength {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.BrowserTitleKey, consts.LimitExceedsKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.BrowserTitleKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.LimitExceedsKey},
			}

		}

	}
	//payoutmin limit validation
	if data.PaymentGatewayDetails.PayoutMinLimit < 0 {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PayoutMinLimitKey, consts.InvalidKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.PayoutMinLimitKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}
	}

	// site_info validation
	if !utils.IsEmpty(data.SiteInfo) {
		siteInfo := strings.TrimSpace(data.SiteInfo)
		if len(siteInfo) > consts.SiteInfoMaxLength {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.SiteInfoKey, consts.LimitExceedsKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.SiteInfoKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.LimitExceedsKey},
			}

		}

	}

	//Postal code validation
	if !utilities.IsValidPostalCode(data.AddressDetails.PostalCode) {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PostalCodeKey, consts.InvalidKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.PostalCodeKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}

	}

	//member default country validation
	if !utils.IsEmpty(data.MemberDefaultCountry) {
		exists, err := utilities.IsCountryExists(ctx, partner.cache, data.MemberDefaultCountry, consts.UtilityServiceURL)
		if !exists {
			code, _ := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.MemberDefaultCountryKey, consts.InvalidKey)
			errMap[consts.MemberDefaultCountryKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}
		}
		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err.Error())
			return nil, err
		}
	}

	//product review validation
	if !utils.IsEmpty(data.ProductReview) {
		data.ProductReviewID, err = utilities.GetLookupId(ctx, partner.cache, consts.ProductReviewLookupTypeName, data.ProductReview, consts.UtilityServiceURL)

		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err.Error())
			return nil, err
		}
		if data.ProductReviewID == 0 {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.ProductReviewKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.ProductReviewKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}
		}
	}

	//member grace period validation
	if !utils.IsEmpty(data.MemberGracePeriod) {

		data.MemberGracePeriodID, err = utilities.GetSubscriptionDurationId(ctx, partner.cache, data.MemberGracePeriod, consts.SubcriptionServiceApiURL)
		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err.Error())
			return nil, err
		}
		if data.MemberGracePeriodID == 0 {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.MemberGracePeriodKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.MemberGracePeriodKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}
		}
	}

	//default currency validation
	if !utils.IsEmpty(data.PaymentGatewayDetails.DefaultCurrency) {
		data.DefaultCurrencyID, err = utilities.GetCurrencyId(ctx, partner.cache, data.PaymentGatewayDetails.DefaultCurrency, consts.UtilityServiceURL)
		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err.Error())
			return nil, err
		}
		if data.DefaultCurrencyID == 0 {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PayoutCurrencyKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.PayoutCurrencyKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}
		}
	}
	var gateways []string

	for i := range data.PaymentGatewayDetails.PaymentGateways {
		gateways = append(gateways, data.PaymentGatewayDetails.PaymentGateways[i].Gateway)
		if data.PaymentGatewayDetails.PaymentGateways[i].Gateway == "" {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PaymentGatewayKey, consts.RequiredKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.PaymentGatewayKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.RequiredKey},
			}

			break

		} else {
			data.PaymentGatewayDetails.PaymentGateways[i].GatewayId, err = utilities.GetPaymentGatewayId(ctx, partner.cache, data.PaymentGatewayDetails.PaymentGateways[i].Gateway, consts.UtilityServiceURL)
			if err != nil {
				log.Errorf(consts.CreatePartnerErrMsg, err.Error())
				return nil, err
			}
			if data.PaymentGatewayDetails.PaymentGateways[i].GatewayId == "" {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PaymentGatewayKey, consts.InvalidKey)
				if err != nil {
					logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
				}
				errMap[consts.PaymentGatewayKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.InvalidKey},
				}
			}

		}

	}
	for i := range data.PaymentGatewayDetails.PaymentGateways {
		if data.PaymentGatewayDetails.PaymentGateways[i].ClientSecret == "" {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.ClientSecretKey, consts.RequiredKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.ClientSecretKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.RequiredKey},
			}

			break

		}
	}
	for i := range data.PaymentGatewayDetails.PaymentGateways {
		if data.PaymentGatewayDetails.PaymentGateways[i].Email == "" {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PaymentGatewayEmailKey, consts.RequiredKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.PaymentGatewayEmailKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.RequiredKey},
			}

			break

		}

		if !utils.IsValidEmail(data.PaymentGatewayDetails.PaymentGateways[i].Email) {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.PaymentGatewayEmailKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.PaymentGatewayEmailKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}

			break
		}
	}
	for i := range data.PaymentGatewayDetails.PaymentGateways {
		if data.PaymentGatewayDetails.PaymentGateways[i].ClientId == "" {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.ClientIdKey, consts.RequiredKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.ClientIdKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.RequiredKey},
			}

			break

		}
	}
	for i := range data.PaymentGatewayDetails.PaymentGateways {
		if data.PaymentGatewayDetails.PaymentGateways[i].DefaultPayoutCurrency == "" {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.DefaultPayoutCurrencyKey, consts.RequiredKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.DefaultPayoutCurrencyKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.RequiredKey},
			}

			break

		}

		if !slices.Contains(consts.SupportedCurrencies, data.PaymentGatewayDetails.PaymentGateways[i].DefaultPayoutCurrency) {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.DefaultPayoutCurrencyKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.DefaultPayoutCurrencyKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}

			break

		}
	}
	for i := range data.PaymentGatewayDetails.PaymentGateways {
		if data.PaymentGatewayDetails.PaymentGateways[i].DefaultPayinCurrency == "" {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.DefaultPayinCurrencyKey, consts.RequiredKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.DefaultPayinCurrencyKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.RequiredKey},
			}

			break

		}

		if !slices.Contains(consts.SupportedCurrencies, data.PaymentGatewayDetails.PaymentGateways[i].DefaultPayinCurrency) {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.DefaultPayinCurrencyKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.DefaultPayinCurrencyKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}

			break

		}
	}

	if data.PaymentGatewayDetails.DefaultPaymentGateway == "" {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.DefaultPaymentGatewayKey, consts.RequiredKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
		}
		errMap[consts.DefaultPaymentGatewayKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.RequiredKey},
		}

	} else {
		if !slices.Contains(gateways, data.PaymentGatewayDetails.DefaultPaymentGateway) {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.DefaultPaymentGatewayKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
			}
			errMap[consts.DefaultPaymentGatewayKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}

		} else {
			data.DefaultPaymentGatewayId, err = utilities.GetPaymentGatewayId(ctx, partner.cache, data.PaymentGatewayDetails.DefaultPaymentGateway, consts.UtilityServiceURL)
			if err != nil {
				log.Errorf(consts.CreatePartnerErrMsg, err.Error())
				return nil, err
			}
			if utils.IsEmpty(data.DefaultPaymentGatewayId) {
				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.DefaultPaymentGatewayKey, consts.NotFoundKey)
				if err != nil {
					logger.Log().WithContext(ctx).Error(consts.CreatePartnerErrMsg, err)
				}
				errMap[consts.DefaultPaymentGatewayKey] = models.ErrorResponse{
					Code:    code,
					Message: []string{consts.NotFoundKey},
				}
			}

		}
	}

	return errMap, nil
}

// function to set default values for empty values while creating  partner
func SetPartnerDefaultVal(data *entities.Partner) {

	// LoginPageLogo default value
	if utils.IsEmpty(data.LoginPageLogo) {
		data.LoginPageLogo = consts.LoginPageLogoDefaultVal
	}
	//Loader default value
	if utils.IsEmpty(data.Loader) {
		data.Loader = consts.LoaderDefaultVal
	}
	//BackgroundColor default value
	if utils.IsEmpty(data.BackgroundColor) {
		data.BackgroundColor = consts.BackgroundColorDefaultVal
	}
	//BackgroundImage default value
	if utils.IsEmpty(data.BackgroundImage) {
		data.BackgroundImage = consts.BackgroundImageDefaultVal
	}
	//Language default value
	if utils.IsEmpty(data.Language) {
		data.Language = consts.LanguageDefaultVal
	}
	//BrowserTitle default value
	if utils.IsEmpty(data.BrowserTitle) {
		data.BrowserTitle = consts.BrowserTitleDefaultVal
	}
	//MobileVerifyInterval default value
	if data.MobileVerifyInterval == 0 {
		data.MobileVerifyInterval = consts.MobileVerifyIntervalDefaultVal
	}
	//PayoutTargetCurrency default value
	if utils.IsEmpty(data.PayoutTargetCurrency) {
		data.PayoutTargetCurrency = consts.PayoutTargetCurrencyDefaultVal
	}
	//ThemeId fefault value
	if data.ThemeID == 0 {
		data.ThemeID = consts.ThemeIdDefaultVal
	}
	//LoginTypeId default value
	if data.LoginType == "" {
		data.LoginType = consts.LoginTypeIdDefaultVal
	}
	//PayoutMinLimit default value
	if data.PaymentGatewayDetails.PayoutMinLimit == 0 {
		data.PaymentGatewayDetails.PayoutMinLimit = consts.PayoutMinLimitDefaultVal
	}
	//MaxRemittancePerMonth default value
	if data.PaymentGatewayDetails.MaxRemittancePerMonth == 0 {
		data.PaymentGatewayDetails.MaxRemittancePerMonth = consts.MaxRemittancePerMonthDefaultVal
	}
	//PayoutCurrency default value
	if utils.IsEmpty(data.PaymentGatewayDetails.DefaultCurrency) {
		data.PaymentGatewayDetails.DefaultCurrency = consts.PayoutCurrencyDefaultVal
	}

	//MemberGracePeriod default value
	if data.MemberGracePeriod == "" {
		data.MemberGracePeriod = consts.MemberGracePeriodDefaultVal
	}
	//ExpiryWarningCount default value
	if data.ExpiryWarningCount == 0 {
		data.ExpiryWarningCount = consts.ExpiryWarningCountDefaultVal
	}
	//DefaultPriceCodeCurrency default value
	if utils.IsEmpty(data.DefaultPriceCodeCurrency.Name) {
		data.DefaultPriceCodeCurrency.Name = consts.DefaultPriceCodeCurrencyDefaultVal
	}
	//MusicLanguage default value
	if utils.IsEmpty(data.MusicLanguage) {
		data.MusicLanguage = consts.MusicLanguageDefaultVal
	}
	//MemberDefaultCountry default value
	if utils.IsEmpty(data.MemberDefaultCountry) {
		data.MemberDefaultCountry = consts.MemberDefaultCountryDefaultVal
	}
	//OutletsProcessingDuration default value
	if data.OutletsProcessingDuration == 0 {
		data.OutletsProcessingDuration = consts.OutletsProcessingDurationDefaultVal
	}
	//FreePlanLimit default value
	if data.FreePlanLimit == 0 {
		data.FreePlanLimit = consts.FreePlanLimitDefaultVal
	}
	//ProductReview default value
	if data.ProductReview == "" {
		data.ProductReview = consts.ProductReviewDefaultVal
	}
}

// function to delete  a Partner based on partner id
func (partner *PartnerUseCases) DeletePartner(ctx *gin.Context, partnerId string, endpoint string, method string, errMap map[string]models.ErrorResponse) error {
	return partner.repo.DeletePartner(ctx, partnerId)
}

// function to delete  Partner genre language
func (partner *PartnerUseCases) DeletePartnerGenreLanguage(ctx *gin.Context, genreID string, endpoint string, method string, errMap map[string]models.ErrorResponse) error {
	return partner.repo.DeletePartnerGenreLanguage(ctx, genreID, endpoint, method, errMap)
}

// function to delete Partner artist role language
func (partner *PartnerUseCases) DeletePartnerArtistRoleLanguage(ctx *gin.Context, roleID string, endpoint string, method string, errMap map[string]models.ErrorResponse) error {
	return partner.repo.DeletePartnerArtistRoleLanguage(ctx, roleID, endpoint, method, errMap)
}

// function to create partner stores
func (partner *PartnerUseCases) CreatePartnerStores(ctx context.Context, newPartnerStores entities.PartnerStores, partnerID string, endpoint string, method string, errMap map[string]models.ErrorResponse) (map[string]models.ErrorResponse, []string, error) {
	storeIds, err := partner.repo.CreatePartnerStores(ctx, newPartnerStores, partnerID, endpoint, method, errMap)
	return nil, storeIds, err
}
