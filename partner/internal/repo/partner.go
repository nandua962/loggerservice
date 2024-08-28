// Package repo provides a repository implementation for retrieving partner data from a database.
package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"partner/internal/consts"
	"partner/internal/entities"
	"partner/utilities"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	cacheConf "gitlab.com/tuneverse/toolkit/core/cache"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils"
	"gitlab.com/tuneverse/toolkit/utils/crypto"
)

// PartnerRepo is a repository for partner data operations.
type PartnerRepo struct {
	db    *sql.DB
	cache cacheConf.Cache
}

// PartnerRepoImply is the is the implementation of PartnerRepo
type PartnerRepoImply interface {
	DeletePartner(*gin.Context, string) error
	GetAllPartners(*gin.Context, entities.Params, string, string, map[string]models.ErrorResponse) ([]entities.ListAllPartners, int64, error)
	CreatePartner(context.Context, string, string, map[string]models.ErrorResponse, entities.Partner, entities.PartnerOauthCredential) (string, error)
	GetPartnerById(context.Context, string) (entities.GetPartner, error)
	GetPartnerOauthCredential(context.Context, string, entities.PartnerOAuthHeader, string) (entities.GetPartnerOauthCredential, error)
	GetAllTermsAndConditions(context.Context, string) (entities.TermsAndConditions, error)
	UpdatePartner(context.Context, string, uuid.UUID, *entities.Partner, *entities.PartnerProperties, string, string, map[string]models.ErrorResponse) (map[string]models.ErrorResponse, error)
	IsPartnerExists(context.Context, string) (bool, error)
	IsExists(context.Context, string, string, string) (bool, error)
	EncryptPaymentData(context.Context, string) (string, error)
	UpdateTermsAndConditions(context.Context, string, uuid.UUID, entities.UpdateTermsAndConditions, string, string, map[string]models.ErrorResponse) error
	GetPartnerPaymentGateways(*gin.Context, string) (entities.GetPaymentGateways, error)
	GetPartnerStores(context.Context, string) (entities.GetPartnerStores, error)
	DeletePartnerGenreLanguage(*gin.Context, string, string, string, map[string]models.ErrorResponse) error
	DeletePartnerArtistRoleLanguage(*gin.Context, string, string, string, map[string]models.ErrorResponse) error
	CreatePartnerStores(context.Context, entities.PartnerStores, string, string, string, map[string]models.ErrorResponse) ([]string, error)
	GetID(context.Context, string, string, interface{}, string) (string, error)
	IsFieldValueUnique(context.Context, string, string, string) (bool, error)
	GetPartnerName(*gin.Context, string) (string, error)
	UpdatePartnerStatus(*gin.Context, string, entities.UpdatePartnerStatus) error
	GetPartnerProductTypes(context.Context, string, entities.QueryParams) (int64, []entities.GetPartnerProdTypesAndTrackQuality, error)
	GetPartnerTrackFileQuality(context.Context, string, entities.QueryParams) (int64, []entities.GetPartnerProdTypesAndTrackQuality, error)
}

// NewPartnerRepo creates a new PartnerRepo instance with a database connection.

func NewPartnerRepo(db *sql.DB, cache cacheConf.Cache) PartnerRepoImply {
	return &PartnerRepo{
		db:    db,
		cache: cache,
	}
}

func (partner *PartnerRepo) UpdatePartnerStatus(ctx *gin.Context, partnerID string, partnerStatus entities.UpdatePartnerStatus) error {
	query := `UPDATE partner SET is_active =$2 WHERE id=$1 `
	_, err := partner.db.ExecContext(ctx, query, partnerID, partnerStatus.Active)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf(consts.UpdatePartnerStatusErrMsg, err.Error())
		return err
	}
	return nil

}
func (partner *PartnerRepo) GetPartnerName(ctx *gin.Context, partnerID string) (string, error) {
	var name string
	query := `SELECT name FROM partner WHERE id=$1 `

	row := partner.db.QueryRowContext(ctx, query, partnerID)
	err := row.Scan(&name)
	if err != nil {
		logger.Log().WithContext(ctx).Errorf(consts.GetPartnerNameErrMsg, err.Error())
		return "", err
	}
	return name, nil
}

func (repo *PartnerRepo) GetPartnerStores(ctx context.Context, PartnerID string) (entities.GetPartnerStores, error) {
	var (
		store entities.GetPartnerStores
		log   = logger.Log().WithContext(ctx)
	)
	query := `SELECT id, store_id FROM partner_store WHERE partner_id=$1 AND is_active =$2`
	rows, err := repo.db.QueryContext(ctx, query, PartnerID, true)
	if err != nil {
		log.Errorf(consts.GetPartnerStoresErrMsg, err.Error())
		return entities.GetPartnerStores{}, err
	}
	for rows.Next() {
		var data entities.Store
		err := rows.Scan(&data.Id, &data.StoreId)
		if err != nil {
			log.Errorf(consts.GetPartnerStoresErrMsg, err.Error())
			return entities.GetPartnerStores{}, err
		}
		store.Stores = append(store.Stores, data)
	}
	return store, nil
}

func (repo *PartnerRepo) GetPartnerPaymentGateways(ctx *gin.Context, PartnerID string) (entities.GetPaymentGateways, error) {

	var (
		data                 entities.GetPaymentGateways
		gatewayData          entities.PaymentGatewayDetails
		gateways             []entities.PaymentGateways
		details              entities.PaymentGateways
		encryptedPaymentData []string
		encryptedData        string
		currencyId           int
		gatewayId            string
		gatewayIds           []string
		log                  = logger.Log().WithContext(ctx)
	)
	query := ` SELECT payment_gateway_id,payment_details FROM partner_payment_gateway WHERE partner_id=$1`
	rows, err := repo.db.QueryContext(ctx, query, PartnerID)
	if err != nil {
		return entities.GetPaymentGateways{}, err
	}

	for rows.Next() {
		err := rows.Scan(
			&gatewayId,
			&encryptedData,
		)

		if err != nil {

			return entities.GetPaymentGateways{}, err
		}
		encryptedPaymentData = append(encryptedPaymentData, encryptedData)
		gatewayIds = append(gatewayIds, gatewayId)
	}
	for i, encryptData := range encryptedPaymentData {
		decryptedData, err := repo.DecryptPaymentData(ctx, encryptData)
		if err != nil {
			log.Errorf(consts.GetPartnerPaymentGatewaysErrMsg, err.Error())
			return entities.GetPaymentGateways{}, err
		}
		err = json.Unmarshal([]byte(decryptedData), &details)
		if err != nil {
			log.Errorf(consts.GetPartnerPaymentGatewaysErrMsg, err.Error())
			return entities.GetPaymentGateways{}, err
		}
		gateways = append(gateways, details)
		gateways[i].GatewayId = gatewayIds[i]

	}

	gatewayData.PaymentGateways = gateways
	query =
		` SELECT
			COALESCE(partner.max_remittance_per_month,0) AS max_remittance_per_month,
			COALESCE(partner.payout_min_limit,0) AS payout_min_limit,
			COALESCE(partner.default_currency_id,0) AS default_currency_id ,
			COALESCE(partner.default_payment_gateway_id::uuid::text, '') AS default_payment_gateway_id
		FROM
			partner
		WHERE
			partner.id = $1`

	rows, err = repo.db.QueryContext(ctx, query, PartnerID)
	if err != nil {
		log.Errorf(consts.GetPartnerPaymentGatewaysErrMsg, err.Error())
		return entities.GetPaymentGateways{}, err
	}
	for rows.Next() {
		err := rows.Scan(
			&gatewayData.MaxRemittancePerMonth,
			&gatewayData.PayoutMinLimit,
			&currencyId,
			&gatewayId,
		)
		if err != nil {
			log.Errorf(consts.GetPartnerPaymentGatewaysErrMsg, err.Error())
			return entities.GetPaymentGateways{}, err
		}
	}
	currency, err := utilities.GetCurrencyName(ctx, repo.cache, currencyId, consts.UtilityServiceURL)
	if err != nil {
		log.Errorf(consts.GetPartnerPaymentGatewaysErrMsg, err.Error())
		return entities.GetPaymentGateways{}, err
	}
	gatewayData.DefaultPaymentGateway, err = utilities.GetPaymentGatewayName(ctx, repo.cache, gatewayId, consts.UtilityServiceURL)
	if err != nil {
		log.Errorf(consts.GetPartnerPaymentGatewaysErrMsg, err.Error())
		return entities.GetPaymentGateways{}, err
	}
	gatewayData.DefaultCurrency = currency
	data.PaymentGatewayDetails = gatewayData

	return data, nil
}

// function to decrypt payment details
func (repo *PartnerRepo) DecryptPaymentData(ctx *gin.Context, data string) (string, error) {
	encryptionKey := consts.PartnerPaymentGatewayEncryptionKey
	key := []byte(encryptionKey)
	decryptedData, err := crypto.Decrypt(data, key)
	if err != nil {
		return "", err
	}
	return decryptedData, nil

}

// GetAllPartners
func (repo *PartnerRepo) GetAllPartners(ctx *gin.Context, params entities.Params, endpoint string, method string, errMap map[string]models.ErrorResponse) ([]entities.ListAllPartners, int64, error) {

	var (
		data                   entities.ListAllPartners
		partners               []entities.ListAllPartners
		loginType              int
		businessModelId        int
		payoutTargetCurrencyId int
		themeID                int
		log                    = logger.Log().WithContext(ctx)
		recordCount            int64
		args                   []interface{}
		count                  = 1
	)

	query :=
		`SELECT
		COUNT(id) OVER() AS total_records,
		p.id,
		p.name,
		p.url,
		COALESCE(p.logo,'')AS logo,
		COALESCE(p.contact_person,'')AS contact_person,
		p.email,
		COALESCE(p.feedback_email, '') AS feedback_email,
		COALESCE(p.noreply_email, '') AS noreply_email,
		COALESCE(p.support_email, '') AS support_email,
		COALESCE(p.language_code, '') AS language_name,
		COALESCE(p.browser_title, '') AS browser_title,
		COALESCE(p.payment_url, '') AS payment_url,
		COALESCE(p.profile_url, '') AS profile_url,
		COALESCE(p.landing_page, '') AS landing_page,
		COALESCE(p.mobile_verify_interval, 0) AS mobile_verify_interval,
		COALESCE(p.payout_target_currency_id, 0) AS payout_target_currency,
		COALESCE(p.theme_id, 0) AS theme_id,
		p.enable_mail,
		COALESCE(p.login_type_id, 0) AS login_type,
		COALESCE(p.business_model_id, 0) AS business_model,
		p.member_pay_to_partner,
		p.is_active
	FROM
    	partner AS p`

	if strings.ToLower(params.Status) == consts.StatusActive || params.Status == "" {
		query += ` WHERE p.is_active = true AND p.is_deleted =false `
	} else if strings.ToLower(params.Status) == consts.StatusInActive {
		query += ` WHERE p.is_active = false AND p.is_deleted =false `
	} else if strings.ToLower(params.Status) == consts.StatusAll {
		query += ` WHERE  p.is_deleted =false `
	}

	if strings.HasPrefix(params.Name, "~") {
		name := fmt.Sprintf("%%%s%%", params.Name[1:])
		query = AddCondition(query, fmt.Sprintf("p.name ILIKE $%d", count))
		args = append(args, name)
		count++
	} else if params.Name != "" {
		query = AddCondition(query, fmt.Sprintf(" LOWER(p.name) = $%d", count))
		args = append(args, strings.ToLower(params.Name))
		count++
	}

	if params.Country != "" {
		exists, err := utilities.IsCountryExists(ctx, repo.cache, strings.ToUpper(params.Country), consts.UtilityServiceURL)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.GetAllPartnersErrMsg, err)
			return nil, 0, err
		}
		if !exists {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.CountryKey, consts.InvalidKey)
			if err != nil {
				logger.Log().WithContext(ctx).Error(consts.GetAllPartnersErrMsg, err)
			}
			errMap[consts.CountryKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}
		}
		query = AddCondition(query, fmt.Sprintf(" p.country_code = $%d", count))
		args = append(args, strings.ToUpper(params.Country))
	}

	if params.Sort != "" {
		sortFields := strings.Split(params.Sort, consts.Seperator)
		sortOrders := strings.Split(params.Order, consts.Seperator)
		query += " ORDER BY"
		for i := range sortFields {
			if i >= len(sortOrders) {
				sortOrders = append(sortOrders, "ASC")
			}
			query += " " + "p." + sortFields[i] + " " + sortOrders[i] + consts.Seperator

		}
		query = strings.TrimSuffix(query, consts.Seperator)
	}

	if params.Page != 0 && params.Limit != 0 {
		offset := (params.Page - 1) * params.Limit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", params.Limit, offset)
	}
	if len(errMap) != 0 {
		return nil, 0, nil

	}
	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Errorf(consts.GetAllPartnersErrMsg, err)
		return nil, 0, err
	}
	for rows.Next() {
		err := rows.Scan(
			&recordCount,
			&data.UUID,
			&data.Name,
			&data.URL,
			&data.Logo,
			&data.ContactDetails.ContactPerson,
			&data.ContactDetails.Email,
			&data.ContactDetails.FeedbackEmail,
			&data.ContactDetails.NoReplyEmail,
			&data.ContactDetails.SupportEmail,
			&data.Language,
			&data.BrowserTitle,
			&data.PaymentURL,
			&data.ProfileURL,
			&data.LandingPage,
			&data.MobileVerifyInterval,
			&payoutTargetCurrencyId,
			&themeID,
			&data.EnableMail,
			&loginType,
			&businessModelId,
			&data.MemberPayToPartner,
			&data.Active,
		)

		if err != nil {
			log.Errorf(consts.GetAllPartnersErrMsg, err)
			return nil, 0, err
		}
		data.Users, err = utilities.GetPartnerUserCount(ctx, data.UUID, consts.MemberServiceURL)
		if err != nil {
			log.Errorf(consts.GetAllPartnersErrMsg, err)
			return nil, 0, err
		}
		data.LoginType, err = utilities.GetLookupName(ctx, repo.cache, loginType, consts.UtilityServiceURL)
		if err != nil {
			log.Errorf(consts.GetAllPartnersErrMsg, err)
			return nil, 0, err
		}
		data.BusinessModel, err = utilities.GetLookupName(ctx, repo.cache, businessModelId, consts.UtilityServiceURL)
		if err != nil {
			log.Errorf(consts.GetAllPartnersErrMsg, err)
			return nil, 0, err
		}
		data.PayoutTargetCurrency, err = utilities.GetCurrencyName(ctx, repo.cache, payoutTargetCurrencyId, consts.UtilityServiceURL)
		if err != nil {
			log.Errorf(consts.GetAllPartnersErrMsg, err)
			return nil, 0, err
		}
		data.Theme, err = utilities.GetTheme(ctx, repo.cache, themeID, consts.UtilityServiceURL)
		if err != nil {
			log.Errorf(consts.GetAllPartnersErrMsg, err)
			return nil, 0, err
		}
		partners = append(partners, data)
	}

	return partners, recordCount, nil

}
func AddCondition(query string, condition string) string {
	if strings.Contains(query, "WHERE") {
		return query + " AND " + condition
	}
	return query + " WHERE " + condition
}

// GetPartnerByID retrieves partner details using a partner_id from the database.
func (repo *PartnerRepo) GetPartnerById(ctx context.Context, partnerId string) (entities.GetPartner, error) {
	var (
		defaultCurrencyId      int
		payoutTargetCurrencyId int
		priceCodecurrencyId    int
		memberGracePeriod      int
		loginType              int
		businessModelId        int
		productReview          int
		themeID                int
		log                    = logger.Log().WithContext(ctx)
		gatewayId              string
		wg                     sync.WaitGroup
		mu                     sync.Mutex
		partnerData            entities.GetPartner
	)

	query := `SELECT 
        COALESCE(partner.name, '') AS name,
        COALESCE(partner.url, '') AS url,
        COALESCE(partner.logo, '') AS logo,
        COALESCE(partner.contact_person, '') AS contact_person,
        COALESCE(partner.email, '') AS email,
        COALESCE(partner.noreply_email, '') AS noreply_email,
        COALESCE(partner.feedback_email, '') AS feedback_email,
        COALESCE(partner.support_email, '') AS support_email,
        COALESCE(partner.language_code::text, '') AS language_code,
        COALESCE(partner.browser_title, '') AS browser_title,
        COALESCE(partner.profile_url, '') AS profile_url,
        COALESCE(partner.payment_url, '') AS payment_url,
        COALESCE(partner.landing_page, '') AS landing_page,
        COALESCE(partner.terms_and_conditions_id, 0) AS terms_and_conditions_id,
        COALESCE(partner.mobile_verify_interval,0) AS mobile_verify_interval,
        COALESCE(partner.payout_target_currency_id,0) AS payout_target_currency_id,
        COALESCE(partner.theme_id,0) AS theme_id,
        COALESCE(partner.enable_mail,false) AS enable_mail,
        COALESCE(partner.login_type_id,0) AS login_type_id,
        COALESCE(partner.business_model_id,0) AS business_model_id,
        COALESCE(partner.member_pay_to_partner,false) AS member_pay_to_partner,
        COALESCE(partner.country_code,'') AS country_code,
        COALESCE(partner.state_code,'') AS state_code,
        COALESCE(partner.city, '') AS city,
        COALESCE(partner.street, '') AS street,
        COALESCE(partner.postal_code, '') AS postal_code,
        COALESCE(partner.partner_plan_id::text,'') AS partner_plan_id,
        COALESCE(partner.plan_start_date::text, '') AS plan_start_date, 
        COALESCE(partner.plan_launch_date::text, '') AS plan_launch_date,
        COALESCE(partner.member_grace_period,0) AS member_grace_period,
        COALESCE(partner.expiry_warning_count,0 ) AS expiry_warning_count,
        COALESCE(partner.album_review_email, '') AS album_review_email,
        COALESCE(partner.site_info, '') AS site_info,
        COALESCE(partner.default_price_code_currency_id,0) AS default_price_code_currency_id,
        COALESCE(partner.music_language_code,'') AS music_language_code,
        COALESCE(partner.member_default_country_code,'') AS member_default_country_code,
        COALESCE(partner.outlets_processing_duration,0) AS outlets_processing_duration,
        COALESCE(partner.free_plan_limit,0) AS free_plan_limit,
        COALESCE(partner.default_currency_id,0) AS default_currency_id,
        COALESCE(product_review, 0) AS product_review,
        COALESCE(partner.language_code,'') AS language,
        COALESCE(partner.music_language_code,'') AS music_language,
        COALESCE(partner.state_code,'') AS state,
        COALESCE(partner.country_code,'') AS country,
        COALESCE(partner.is_active,false) AS status,
        COALESCE(partner.default_payment_gateway_id::uuid::text,'') AS default_payment_gateway_id
    FROM partner
    WHERE  partner.id = $1 AND  partner.is_deleted = $2 AND  partner.is_active =$3;`

	if err := repo.db.Ping(); err != nil {
		log.Errorf(consts.GetPartnerByIdErrMsg, err.Error())
	}

	row := repo.db.QueryRow(query, partnerId, false, true)
	err := row.Scan(&partnerData.Name,
		&partnerData.URL,
		&partnerData.Logo,
		&partnerData.ContactDetails.ContactPerson,
		&partnerData.ContactDetails.Email,
		&partnerData.ContactDetails.NoReplyEmail,
		&partnerData.ContactDetails.FeedbackEmail,
		&partnerData.ContactDetails.SupportEmail,
		&partnerData.Language,
		&partnerData.BrowserTitle,
		&partnerData.ProfileURL,
		&partnerData.PaymentURL,
		&partnerData.LandingPage,
		&partnerData.TermsAndConditionsVersionID,
		&partnerData.MobileVerifyInterval,
		&payoutTargetCurrencyId,
		&themeID,
		&partnerData.EnableMail,
		&loginType,
		&businessModelId,
		&partnerData.MemberPayToPartner,
		&partnerData.AddressDetails.Country,
		&partnerData.AddressDetails.State,
		&partnerData.AddressDetails.City,
		&partnerData.AddressDetails.Street,
		&partnerData.AddressDetails.PostalCode,
		&partnerData.SubscriptionPlanDetails.ID,
		&partnerData.SubscriptionPlanDetails.StartDate,
		&partnerData.SubscriptionPlanDetails.LaunchDate,
		&memberGracePeriod,
		&partnerData.ExpiryWarningCount,
		&partnerData.AlbumReviewEmail,
		&partnerData.SiteInfo,
		&priceCodecurrencyId,
		&partnerData.MusicLanguage,
		&partnerData.MemberDefaultCountry,
		&partnerData.OutletsProcessingDuration,
		&partnerData.FreePlanLimit,
		&defaultCurrencyId,
		&productReview,
		&partnerData.Language,
		&partnerData.MusicLanguage,
		&partnerData.AddressDetails.State,
		&partnerData.AddressDetails.Country,
		&partnerData.Active,
		&gatewayId)

	if err != nil {
		log.Errorf(consts.GetPartnerByIdErrMsg, err.Error())
		return entities.GetPartner{}, err
	}

	results := make(map[string]entities.GetPartnerByIdMapResult)

	wg.Add(consts.GetPartnerByIDGoroutinesNum)

	go func() {
		defer wg.Done()
		name, err := utilities.GetPaymentGatewayName(ctx, repo.cache, gatewayId, consts.UtilityServiceURL)
		mu.Lock()
		results[consts.DefaultPaymentGatewayKey] = entities.GetPartnerByIdMapResult{Value: name, Err: err}
		mu.Unlock()
	}()
	go func() {
		defer wg.Done()
		name, err := utilities.GetCurrencyName(ctx, repo.cache, defaultCurrencyId, consts.UtilityServiceURL)
		mu.Lock()
		results[consts.DefaultCurrencyKey] = entities.GetPartnerByIdMapResult{Value: name, Err: err}
		mu.Unlock()
	}()
	go func() {
		defer wg.Done()
		name, err := utilities.GetCurrencyName(ctx, repo.cache, payoutTargetCurrencyId, consts.UtilityServiceURL)
		mu.Lock()
		results[consts.PayoutTargetCurrencyKey] = entities.GetPartnerByIdMapResult{Value: name, Err: err}
		mu.Unlock()
	}()
	go func() {
		defer wg.Done()
		name, err := utilities.GetCurrencyName(ctx, repo.cache, priceCodecurrencyId, consts.UtilityServiceURL)
		mu.Lock()
		results[consts.DefaultPriceCodeCurrencyKey] = entities.GetPartnerByIdMapResult{Value: name, Err: err}
		mu.Unlock()
	}()
	go func() {
		defer wg.Done()
		name, err := utilities.GetSubcriptionDurationName(ctx, repo.cache, memberGracePeriod, consts.SubcriptionServiceApiURL)
		mu.Lock()
		results[consts.MemberGracePeriodKey] = entities.GetPartnerByIdMapResult{Value: name, Err: err}
		mu.Unlock()
	}()
	go func() {
		defer wg.Done()
		name, err := utilities.GetLookupName(ctx, repo.cache, loginType, consts.UtilityServiceURL)
		mu.Lock()
		results[consts.LoginTypeKey] = entities.GetPartnerByIdMapResult{Value: name, Err: err}
		mu.Unlock()
	}()
	go func() {
		defer wg.Done()
		name, err := utilities.GetLookupName(ctx, repo.cache, businessModelId, consts.UtilityServiceURL)
		mu.Lock()
		results[consts.BusinessModelKey] = entities.GetPartnerByIdMapResult{Value: name, Err: err}
		mu.Unlock()
	}()
	go func() {
		defer wg.Done()
		name, err := utilities.GetLookupName(ctx, repo.cache, productReview, consts.UtilityServiceURL)
		mu.Lock()
		results[consts.ProductReviewKey] = entities.GetPartnerByIdMapResult{Value: name, Err: err}
		mu.Unlock()
	}()
	go func() {
		defer wg.Done()
		name, err := utilities.GetTheme(ctx, repo.cache, themeID, consts.UtilityServiceURL)
		mu.Lock()
		results[consts.ThemeKey] = entities.GetPartnerByIdMapResult{Value: name, Err: err}
		mu.Unlock()
	}()

	wg.Wait()

	for key, res := range results {
		if res.Err != nil {
			log.Errorf(consts.GetPartnerByIdErrMsg, res.Err.Error())
			return entities.GetPartner{}, res.Err
		}
		switch key {
		case consts.DefaultPaymentGatewayKey:
			partnerData.DefaultPaymentGateway = res.Value.(string)
		case consts.DefaultCurrencyKey:
			partnerData.DefaultCurrency = res.Value.(string)
		case consts.PayoutTargetCurrencyKey:
			partnerData.PayoutTargetCurrency = res.Value.(string)
		case consts.DefaultPriceCodeCurrencyKey:
			partnerData.DefaultPriceCodeCurrency.Name = res.Value.(string)
		case consts.MemberGracePeriodKey:
			partnerData.MemberGracePeriod = res.Value.(string)
		case consts.LoginTypeKey:
			partnerData.LoginType = res.Value.(string)
		case consts.BusinessModelKey:
			partnerData.BusinessModel = res.Value.(string)
		case consts.ProductReviewKey:
			partnerData.ProductReview = res.Value.(string)
		case consts.ThemeKey:
			partnerData.Theme = res.Value.(string)
		}
	}

	return partnerData, nil
}

// function to get all terms and conditions of a partner

func (repo *PartnerRepo) GetAllTermsAndConditions(ctx context.Context, partnerId string) (entities.TermsAndConditions, error) {

	var (
		TermsAndConditions entities.TermsAndConditions
		log                = logger.Log().WithContext(ctx)
	)

	query := `
	SELECT
		tc.heading,
		tc.content,
		tc.language_code AS language
	FROM
		terms_and_conditions AS tc
	WHERE
		tc.partner_id = $1
		AND tc.is_active = true;
	`

	rows, err := repo.db.QueryContext(ctx, query, partnerId)
	if err != nil {
		log.Errorf(consts.GetAllTermsAndConditionErrMsg, err)
		return entities.TermsAndConditions{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&TermsAndConditions.Name,
			&TermsAndConditions.Description,
			&TermsAndConditions.Language,
		)

		if err != nil {
			log.Errorf(consts.GetAllTermsAndConditionErrMsg, err)
			return entities.TermsAndConditions{}, err
		}
	}

	return TermsAndConditions, nil
}

func (repo *PartnerRepo) UpdateTermsAndConditions(ctx context.Context, partnerID string, memberID uuid.UUID, termsAndConditionsData entities.UpdateTermsAndConditions, endpoint string, method string, errMap map[string]models.ErrorResponse) error {
	var (
		log      = logger.Log().WithContext(ctx)
		err      error
		id       int
		language interface{}
	)
	lang := termsAndConditionsData.GetData(consts.LanguageKey)
	isExists, err := utilities.IsLanguageIsoExists(ctx, repo.cache, lang, consts.UtilityServiceURL)
	if !isExists {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.LanguageKey, consts.InvalidKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.UpdateTermsAndConditionsErrMsg, err)
		}
		errMap[consts.LanguageKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}

	}
	if err != nil {
		return err
	}

	if len(errMap) != 0 {
		return nil
	}

	query := `UPDATE terms_and_conditions
	SET is_active = $3
	WHERE partner_id = $1
	AND EXISTS (SELECT 1 FROM terms_and_conditions WHERE partner_id = $2)`

	_, err = repo.db.Exec(query, partnerID, partnerID, false)
	if err != nil {
		log.Errorf(consts.UpdateTermsAndConditionsErrMsg, err.Error())
		return err
	}

	data := map[string]interface{}{
		consts.HeadingKey:      termsAndConditionsData.GetData(consts.NameKey),
		consts.ContentKey:      termsAndConditionsData.GetData(consts.DescriptionKey),
		consts.LanguageCodeKey: language,
		consts.PartnerIDKey:    partnerID,
		consts.IsActiveKey:     true,
		consts.UpdatedByKey:    memberID,
		consts.UpdatedOnKey:    time.Now(),
	}

	query, value := InsertQueryBuilder("terms_and_conditions", data)
	row := repo.db.QueryRowContext(ctx, query, value...)
	err = row.Scan(&id)
	if err != nil {
		log.Errorf(consts.UpdateTermsAndConditionsErrMsg, err.Error())
		return err
	}
	// Start a transaction.
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
		}
	}()

	query = `UPDATE partner
	SET terms_and_conditions_id = terms_and_conditions.id
	FROM terms_and_conditions
	WHERE partner.id = terms_and_conditions.partner_id
	AND terms_and_conditions.is_active = $2
	AND terms_and_conditions.partner_id = $1;`

	_, err = tx.Exec(query, partnerID, true)
	if err != nil {

		log.Errorf(consts.UpdateTermsAndConditionsErrMsg, err.Error())
		return err
	}

	isUpdated, err := utilities.UpdateMemberTermsAndConditions(tx, id, false, partnerID, consts.MemberServiceURL)
	if err != nil {
		log.Errorf(consts.UpdateTermsAndConditionsErrMsg, err.Error())
		return err
	}
	if !isUpdated {
		return errors.New("member service updation failed")
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// UpdatePartner function
func (repo *PartnerRepo) UpdatePartner(ctx context.Context, partnerID string, memberID uuid.UUID, data *entities.Partner, partner *entities.PartnerProperties, endpoint string, method string, errMap map[string]models.ErrorResponse) (map[string]models.ErrorResponse, error) {

	var (
		log = logger.Log().WithContext(ctx)
		// productReviewID, loginTypeID interface{}
		err            error
		paymentDetails = make([]map[string]interface{}, 0)
		// businessModelID              int
		query string
	)

	if len(errMap) != 0 {
		return nil, nil
	}
	// Start a transaction.
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
		}
	}()

	var (
		updates       []string
		values        []interface{}
		encryptedData string
	)

	counter := 1
	paymentGateways, err := partner.GetPaymentGateways()
	if err != nil {
		return nil, err
	}
	if paymentGateways != nil {
		query = `DELETE FROM partner_payment_gateway WHERE partner_id =$1`
		_, err = tx.Exec(query, partnerID)
		if err != nil {
			log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
			return nil, err
		}
	}

	for _, gateway := range data.PaymentGatewayDetails.PaymentGateways {

		// Create an anonymous struct for serialisation without PayoutGateway
		serializedData, _ := json.Marshal(struct {
			Gateway               string   `json:"gateway"`
			Email                 string   `json:"email"`
			ClientId              string   `json:"client_id"`
			ClientSecret          string   `json:"client_secret"`
			Payin                 bool     `json:"payin"`
			Payout                bool     `json:"payout"`
			DefaultPayinCurrency  string   `json:"default_payin_currency"`
			DefaultPayoutCurrency string   `json:"default_payout_currency"`
			SupportedCurrencies   []string `json:"supported_currency"`
		}{
			Gateway:               gateway.Gateway,
			Email:                 gateway.Email,
			ClientId:              gateway.ClientId,
			ClientSecret:          gateway.ClientSecret,
			Payin:                 gateway.Payin,
			Payout:                gateway.Payout,
			SupportedCurrencies:   consts.SupportedCurrencies,
			DefaultPayinCurrency:  gateway.DefaultPayinCurrency,
			DefaultPayoutCurrency: gateway.DefaultPayoutCurrency,
		})

		encryptedData, err = repo.EncryptPaymentData(ctx, string(serializedData))
		if err != nil {
			log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
			return nil, err
		}
		data := make(map[string]interface{})
		data[consts.PartnerIDKey] = partnerID
		data[consts.PaymentGateWayIDKey] = gateway.GatewayId
		data[consts.PaymentDetailsKey] = encryptedData
		paymentDetails = append(paymentDetails, data)

		if err != nil {
			log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
			return nil, err
		}

	}
	err = InsertData(tx, consts.PartnerPaymentGatewayTable, paymentDetails)

	if err != nil {
		log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
		return nil, err
	}

	updatePartnerData := utilities.UpdatePartnerData(*data, partnerID, memberID)

	for field, value := range updatePartnerData {
		if !utilities.IsValueEmpty(value) {
			updates = append(updates, fmt.Sprintf("%s=$%d", field, counter))
			values = append(values, value)
			counter++
		}
	}

	query = fmt.Sprintf(`
	 	UPDATE partner
	 	SET %s
	 	WHERE id= '%s' `, strings.Join(updates, consts.Seperator), partnerID)

	_, err = tx.Exec(query, values...)
	if err != nil {
		log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
		return nil, err
	}
	// Commit the transaction if everything succeeds.
	if err = tx.Commit(); err != nil {
		log.Errorf(consts.UpdatePartnerErrMsg, err.Error())
		return nil, err
	}
	return nil, nil
}

// function to encrypt the payment data
func (repo *PartnerRepo) EncryptPaymentData(ctx context.Context, data string) (string, error) {
	var log = logger.Log().WithContext(ctx)
	encryptionKey := consts.PartnerPaymentGatewayEncryptionKey
	key := []byte(encryptionKey)
	string, err := crypto.Encrypt(data, key)
	if err != nil {
		log.Errorf("Encryption of payment details failed , error=%s", err)
		return "", err
	}
	return string, nil

}

func (partner *PartnerRepo) IsFieldValueUnique(ctx context.Context, fieldName, fieldValue string, partnerID string) (bool, error) {
	var count int
	query := fmt.Sprintf(`
        SELECT COUNT(*)
        FROM partner
        WHERE %s = $1 AND id != $2
    `, fieldName)

	err := partner.db.QueryRowContext(ctx, query, fieldValue, partnerID).Scan(&count)

	if err != nil || count > 0 {
		return false, err
	}

	return true, nil
}

// IsPartnerrExists function to check whether the member exists or not
func (repo *PartnerRepo) IsPartnerExists(ctx context.Context, partnerID string) (bool, error) {

	var exists int
	isPartnerExistsQ := `SELECT 1 FROM partner WHERE id = $1 AND is_deleted =$2 ANd is_active =$3`
	row := repo.db.QueryRowContext(ctx, isPartnerExistsQ, partnerID, false, true)
	err := row.Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// function to store partner details
func (repo *PartnerRepo) CreatePartner(ctx context.Context, endpoint string, method string, errMap map[string]models.ErrorResponse, partner entities.Partner, oauthCredentials entities.PartnerOauthCredential) (string, error) {
	var (
		partnerId      string
		err            error
		log            = logger.Log().WithContext(ctx)
		paymentDetails = make([]map[string]interface{}, 0)
	)

	// Start a transaction.
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Errorf(consts.CreatePartnerErrMsg, err.Error())
		return "", err
	}

	defer func() {
		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err.Error())
			err := tx.Rollback()
			if err != nil {
				log.Errorf(consts.CreatePartnerErrMsg, err.Error())
				return
			}
		}
	}()

	partnerData := utilities.GeneratePartnerData(partner)
	query, values := InsertQueryBuilder(consts.PartnerTable, partnerData)

	row := tx.QueryRowContext(ctx, query, values...)
	err = row.Scan(&partnerId)

	if err != nil {
		log.Errorf(consts.CreatePartnerErrMsg, err.Error())
		return "", err
	}
	for i, gateway := range partner.PaymentGatewayDetails.PaymentGateways {
		serializedData, _ := json.Marshal(struct {
			Gateway               string   `json:"gateway"`
			Email                 string   `json:"email"`
			ClientId              string   `json:"client_id"`
			ClientSecret          string   `json:"client_secret"`
			Payin                 bool     `json:"payin"`
			Payout                bool     `json:"payout"`
			SupportedCurrencies   []string `json:"supported_currency"`
			DefaultPayinCurrency  string   `json:"default_payin_currency"`
			DefaultPayoutCurrency string   `json:"default_payout_currency"`
		}{
			Gateway:               gateway.Gateway,
			Email:                 gateway.Email,
			ClientId:              gateway.ClientId,
			ClientSecret:          gateway.ClientSecret,
			Payin:                 gateway.Payin,
			Payout:                gateway.Payout,
			SupportedCurrencies:   consts.SupportedCurrencies,
			DefaultPayinCurrency:  gateway.DefaultPayinCurrency,
			DefaultPayoutCurrency: gateway.DefaultPayoutCurrency,
		})
		encryptedData, err := repo.EncryptPaymentData(ctx, string(serializedData))
		data := make(map[string]interface{})

		data[consts.PartnerIDKey] = partnerId
		data[consts.PaymentGateWayIDKey] = partner.PaymentGatewayDetails.PaymentGateways[i].GatewayId
		data[consts.PaymentDetailsKey] = encryptedData
		paymentDetails = append(paymentDetails, data)

		if err != nil {
			log.Errorf(consts.CreatePartnerErrMsg, err.Error())
			return "", err
		}
	}
	err = InsertData(tx, consts.PartnerPaymentGatewayTable, paymentDetails)
	if err != nil {
		log.Errorf(consts.CreatePartnerErrMsg, err.Error())
		return "", err
	}
	oauthCredentials.PartnerId = partnerId
	partnerCredentialData := utilities.OauthCredentialData(oauthCredentials)
	query, values = InsertQueryBuilder(consts.PartnerOauthCredentialTable, partnerCredentialData)
	_, err = tx.ExecContext(ctx, query, values...)

	if err != nil {
		log.Errorf(consts.PartnerOauthCredentialErrMsg, err.Error())
		return "", err
	}

	// Commit the transaction if everything succeeds.
	if err = tx.Commit(); err != nil {
		log.Errorf(consts.CreatePartnerErrMsg, err.Error())
		return "", err
	}

	return partnerId, nil
}
func InsertData(tx *sql.Tx, tableName string, data []map[string]interface{}) error {
	if len(data) == 0 {
		return nil
	}
	var columns []string
	for key := range data[0] {
		columns = append(columns, key)
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ", tableName, strings.Join(columns, ", "))
	var params []string
	var values []interface{}

	for i, entry := range data {
		var entryParams []string
		for j, col := range columns {
			param := fmt.Sprintf("$%d", i*len(columns)+j+1)
			entryParams = append(entryParams, param)
			values = append(values, entry[col])
		}
		params = append(params, fmt.Sprintf("(%s)", strings.Join(entryParams, ", ")))
	}
	_, err := tx.Exec(query+strings.Join(params, ", "), values...)
	return err
}

func (repo *PartnerRepo) IsExists(ctx context.Context, tableName, fieldName, fieldValue string) (bool, error) {
	var (
		log    = logger.Log().WithContext(ctx)
		result bool
	)

	query := fmt.Sprintf(`SELECT CASE WHEN EXISTS 
				(SELECT 1 FROM %s WHERE %s = $1 AND %s IS NOT NULL) 
				THEN true ELSE false END`, tableName, fieldName, fieldName)

	if err := repo.db.QueryRowContext(ctx, query, fieldValue).Scan(&result); err != nil {
		log.Errorf("Error checking existence: %s", err.Error())
		return false, err
	}

	return result, nil
}

// function to get partner Oauth credential
func (repo *PartnerRepo) GetPartnerOauthCredential(ctx context.Context, partnerId string, oauthHeader entities.PartnerOAuthHeader, oauthProviderID string) (entities.GetPartnerOauthCredential, error) {
	var (
		data     entities.GetPartnerOauthCredential
		log      = logger.Log().WithContext(ctx)
		err      error
		query    string
		scope    []string
		scopeVal string
		args     []interface{}
		count    = 1
	)

	query = `
	SELECT 
	COALESCE(client_id, '') AS client_id,
	COALESCE(client_secret, '') AS client_secret,
	COALESCE(redirect_uri, '') AS redirect_uri,
	COALESCE(
        CASE 
            WHEN jsonb_typeof(scope) = 'array' THEN scope
            ELSE '[]'::jsonb
        END,
        '[]'::jsonb
    ) AS scope,
	COALESCE(access_token_endpoint, '') AS access_token_endpoint
	FROM partner_oauth_credential
	WHERE partner_id = $1`
	args = append(args, partnerId)
	appendcondition := func(field string, value interface{}) {
		query = AddCondition(query, fmt.Sprintf(" %s = $%d", field, count+1))
		args = append(args, value)
		count++

	}

	if !utils.IsEmpty(oauthHeader.ProviderName) {
		appendcondition("oauth_provider_id", oauthProviderID)

	}
	if !utils.IsEmpty(oauthHeader.ClientID) {
		appendcondition("client_id", oauthHeader.ClientID)
	}

	if !utils.IsEmpty(oauthHeader.ClientSecret) {
		appendcondition("client_secret", oauthHeader.ClientSecret)
	}

	err = repo.db.QueryRowContext(ctx, query, args...).Scan(
		&data.ClientId,
		&data.ClientSecret,
		&data.RedirectUri,
		&scopeVal,
		&data.AccessTokenEndpoint,
	)

	data.ProviderName = oauthHeader.ProviderName
	if err != nil {
		log.Errorf(consts.GetPartnerOauthCredentialErrMsg, err.Error())
		return entities.GetPartnerOauthCredential{}, err
	}
	err = json.Unmarshal([]byte(scopeVal), &scope)
	if err != nil {
		log.Errorf(consts.GetPartnerOauthCredentialErrMsg, err.Error())
		return entities.GetPartnerOauthCredential{}, err
	}
	data.Scope = scope

	return data, nil
}

// DeletePartner
func (partner *PartnerRepo) DeletePartner(ctx *gin.Context, partnerId string) error {
	var log = logger.Log().WithContext(ctx)

	// delete the partner
	query := `
	UPDATE partner
	SET is_active = false, is_deleted = true
	WHERE id = $1;
	`
	_, err := partner.db.ExecContext(ctx, query, partnerId)

	if err != nil {
		log.Errorf(consts.DeletePartnerErrMsg, err)
		return err
	}

	return nil

}

// function to return id of a particular table for the given fields and value.
func (repo *PartnerRepo) GetID(ctx context.Context, tableName string, field string, fieldValue interface{}, returnField string) (string, error) {
	var result string

	if fieldValue == "" {
		return "", nil
	}
	query := `SELECT ` + returnField + ` FROM ` + tableName + `
	WHERE ` + field + ` = $1 `
	row := repo.db.QueryRowContext(ctx, query, fieldValue)

	err := row.Scan(&result)
	if err != nil {
		return "", err
	}
	return result, nil
}

// Delete Partner genre language
func (partner *PartnerRepo) DeletePartnerGenreLanguage(ctx *gin.Context, genreID string, endpoint string, method string, errMap map[string]models.ErrorResponse) error {
	var (
		log = logger.Log().WithContext(ctx)
	)

	// Check if the row with the given genre_id exists
	existsQuery := `
	SELECT EXISTS(SELECT 1 FROM partner_genre_language WHERE genre_id=$1);
	`
	var exists bool
	err := partner.db.QueryRowContext(ctx, existsQuery, genreID).Scan(&exists)
	if err != nil {
		log.Errorf(consts.DeletePartnerGenreLanguageErrMsg, err.Error())
		return err
	}

	if !exists {
		log.Errorf(consts.DeletePartnerGenreLanguageErrMsg, "genre not exists")
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.GenreIdKey, consts.NotFoundKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.DeletePartnerGenreLanguageErrMsg, err)
		}
		errMap[consts.GenreIdKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.NotFoundKey},
		}
		return fmt.Errorf("genre not exists")
	}

	// Delete the partner genre language
	deleteQuery := `
	DELETE FROM partner_genre_language WHERE genre_id=$1;
	`
	_, err = partner.db.ExecContext(ctx, deleteQuery, genreID)
	if err != nil {
		log.Errorf(consts.DeletePartnerGenreLanguageErrMsg, err.Error())
		return err
	}

	return nil
}

// Delete Partner artist role language
func (partner *PartnerRepo) DeletePartnerArtistRoleLanguage(ctx *gin.Context, roleID string, endpoint string, method string, errMap map[string]models.ErrorResponse) error {
	var (
		log = logger.Log().WithContext(ctx)
	)

	// Check if the role exists
	existsQuery := `
	SELECT EXISTS(SELECT 1 FROM partner_artist_role_language WHERE role_id = $1)
	`
	var exists bool
	err := partner.db.QueryRowContext(ctx, existsQuery, roleID).Scan(&exists)
	if err != nil {
		log.Errorf(consts.DeletePartnerArtistRoleLanguageErrMsg, err.Error())
		return err
	}
	if !exists {
		log.Errorf(consts.DeletePartnerArtistRoleLanguageErrMsg, "role not exists")
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.RoleIdKey, consts.NotFoundKey)
		if err != nil {
			logger.Log().WithContext(ctx).Error(consts.DeletePartnerArtistRoleLanguageErrMsg, err)
		}
		errMap[consts.RoleIdKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.NotFoundKey},
		}
		return fmt.Errorf("role not exists")
	}

	// Delete the partner artist role language
	deleteQuery := `
	DELETE FROM partner_artist_role_language WHERE role_id = $1
	`
	_, err = partner.db.ExecContext(ctx, deleteQuery, roleID)
	if err != nil {
		log.Errorf(consts.DeletePartnerArtistRoleLanguageErrMsg, err.Error())
		return err
	}

	return nil
}

// function to get product types of a partner
func (partner *PartnerRepo) GetPartnerProductTypes(ctx context.Context, partnerId string, params entities.QueryParams) (int64, []entities.GetPartnerProdTypesAndTrackQuality, error) {
	var (
		log            = logger.Log().WithContext(ctx)
		records        []entities.GetPartnerProdTypesAndTrackQuality
		data           entities.GetPartnerProdTypesAndTrackQuality
		productTypeMap = make(map[int]string)
		id             int
		recordCount    int64
	)
	productTypes, err := utilities.GetAllProductTypes(ctx, partner.cache, consts.PartnerServiceURL)

	if err != nil {
		log.Errorf(consts.GetPartnerProductTypesErrMsg, err.Error())
		return 0, nil, err
	}

	for _, val := range productTypes {
		productTypeMap[val.ID] = val.Name
	}

	query := `WITH total AS (
    SELECT COUNT(id) AS total_records
    FROM partner_product_type
    WHERE partner_id = $1 AND is_active = $2
)
SELECT 
    total.total_records,
    ppt.product_type_id
FROM 
    partner_product_type ppt,
    total
WHERE 
    ppt.partner_id = $1 AND ppt.is_active = $2 `

	if params.Page != 0 && params.Limit != 0 {
		offset := (params.Page - 1) * params.Limit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", params.Limit, offset)
	}
	rows, err := partner.db.QueryContext(ctx, query, partnerId, true)
	if err != nil {
		log.Errorf(consts.GetPartnerProductTypesErrMsg, err.Error())
		return 0, []entities.GetPartnerProdTypesAndTrackQuality{}, err
	}
	for rows.Next() {
		err := rows.Scan(
			&recordCount,
			&id,
		)
		if err != nil {
			log.Errorf(consts.GetPartnerProductTypesErrMsg, err.Error())
			return 0, []entities.GetPartnerProdTypesAndTrackQuality{}, err
		}
		data.Id = id
		data.Name = productTypeMap[id]
		records = append(records, data)

	}
	// to sort the name field either in ascending or descending
	if params.Order == consts.Descending {
		sort.Slice(records, func(i, j int) bool {
			return records[i].Name > records[j].Name
		})
	} else {
		sort.Slice(records, func(i, j int) bool {
			return records[i].Name < records[j].Name
		})
	}

	return recordCount, records, nil
}

// function to get track file quality of a partner
func (partner *PartnerRepo) GetPartnerTrackFileQuality(ctx context.Context, partnerId string, params entities.QueryParams) (int64, []entities.GetPartnerProdTypesAndTrackQuality, error) {
	var (
		log              = logger.Log().WithContext(ctx)
		records          []entities.GetPartnerProdTypesAndTrackQuality
		data             entities.GetPartnerProdTypesAndTrackQuality
		trackQuality     []entities.ProductTypes
		trackFileQuality entities.ProductTypes
		trackQualityMap  = make(map[int]string)
		id               int
		recordCount      int64
	)
	query := `SELECT l.id,l.name FROM lookup l where l.lookup_type_id=19`
	rows, err := partner.db.QueryContext(ctx, query)
	if err != nil {
		log.Errorf(consts.GetPartnerTrackQualityErrMsg, err.Error())
		return 0, []entities.GetPartnerProdTypesAndTrackQuality{}, err
	}
	for rows.Next() {
		err := rows.Scan(
			&trackFileQuality.ID,
			&trackFileQuality.Name,
		)
		if err != nil {
			log.Errorf(consts.GetPartnerTrackQualityErrMsg, err.Error())
			return 0, []entities.GetPartnerProdTypesAndTrackQuality{}, err
		}
		trackQuality = append(trackQuality, trackFileQuality)

	}
	for _, val := range trackQuality {
		trackQualityMap[val.ID] = val.Name
	}

	query = `WITH total AS (
    SELECT COUNT(id) AS total_records
    FROM partner_track_file_quality
    WHERE partner_id = $1
)
SELECT 
    total.total_records,
    ptfq.track_file_quality_id
FROM 
    partner_track_file_quality ptfq,
    total
WHERE 
    ptfq.partner_id = $1
`
	// q += fmt.Sprintf("%s %s", params.Sort, params.Order)
	if params.Page != 0 && params.Limit != 0 {
		offset := (params.Page - 1) * params.Limit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", params.Limit, offset)
	}
	rows, err = partner.db.QueryContext(ctx, query, partnerId)
	if err != nil {
		log.Errorf(consts.GetPartnerTrackQualityErrMsg, err.Error())
		return 0, []entities.GetPartnerProdTypesAndTrackQuality{}, err
	}
	for rows.Next() {
		err := rows.Scan(
			&recordCount,
			&id,
		)
		if err != nil {
			log.Errorf(consts.GetPartnerTrackQualityErrMsg, err.Error())
			return 0, []entities.GetPartnerProdTypesAndTrackQuality{}, err
		}
		data.Id = id
		data.Name = trackQualityMap[id]
		records = append(records, data)

	}
	// to sort the name field either in ascending or descending
	if params.Order == consts.Descending {
		sort.Slice(records, func(i, j int) bool {
			return records[i].Name > records[j].Name
		})
	} else {
		sort.Slice(records, func(i, j int) bool {
			return records[i].Name < records[j].Name
		})
	}
	return recordCount, records, nil
}

// function to create partner stores
func (partner *PartnerRepo) CreatePartnerStores(ctx context.Context, newPartnerStores entities.PartnerStores, partnerId string, endpoint string, method string, errMap map[string]models.ErrorResponse) ([]string, error) {

	var (
		storeIds     []string
		foundStoreId = make([]string, 0)
		data         []entities.StoreData
		log          = logger.Log().WithContext(ctx)
	)
	query := `UPDATE partner_store
	SET is_active = $3
	WHERE partner_id = $1
	AND EXISTS (
		SELECT 1
		FROM partner_store
		WHERE partner_id = $2
	);`

	_, err := partner.db.ExecContext(ctx, query, partnerId, partnerId, false)
	if err != nil {
		log.Errorf(consts.CreatePartnerStoresErrMsg, err.Error())
		return nil, err
	}
	data, err = utilities.GetStores(ctx, partner.cache, consts.StoreServiceURL)
	if err != nil {
		log.Errorf(consts.CreatePartnerStoresErrMsg, err.Error())
		return nil, err
	}

	stringData := make([]string, len(data))
	for i, v := range data {
		stringData[i] = fmt.Sprintf("%v", v.Id)
	}

	storeIds = stringData

	if len(newPartnerStores.Stores) != 0 {
		for i, newStoreId := range newPartnerStores.Stores {
			if slices.Contains(storeIds, newStoreId) {
				foundStoreId = append(foundStoreId, newPartnerStores.Stores[i])
			}
		}
	} else {
		foundStoreId = storeIds
	}

	q := `
	INSERT INTO partner_store(
		partner_id,
		store_id,
		is_active, 
		processing_duration, 
		review_processing_duration
	) VALUES `

	valueArgs := []interface{}{}

	for i, v := range foundStoreId {
		q = fmt.Sprintf("%s($%d, $%d, $%d, $%d ,$%d)", q, 5*i+1, 5*i+2, 5*i+3, 5*i+4, 5*i+5)

		if i < len(foundStoreId)-1 {
			q += ","
		}

		valueArgs = append(valueArgs, partnerId)
		valueArgs = append(valueArgs, v)
		valueArgs = append(valueArgs, true)
		valueArgs = append(valueArgs, consts.ProcessingDuration)
		valueArgs = append(valueArgs, consts.ReviewProcessingDuration)

	}

	_, err = partner.db.ExecContext(ctx, q, valueArgs...)
	if err != nil {
		log.Errorf(consts.CreatePartnerStoresErrMsg, err.Error())
		return nil, err
	}
	return foundStoreId, nil

}

// function to create dynamic insert query based on the input
func InsertQueryBuilder(tableName string, data map[string]interface{}) (string, []interface{}) {
	// Prepare the SQL query
	query := fmt.Sprintf("INSERT INTO %s ", tableName)

	var (
		columns []string
		values  []interface{}
	)

	for column, value := range data {
		if !utilities.IsValueEmpty(value) {
			columns = append(columns, column)
			values = append(values, value)
		}
	}
	joinedString := strings.Join(columns, ",")
	query = fmt.Sprintf("%s %s %s %s", query, "(", joinedString, ")")
	query = fmt.Sprintf("%s %s", query, "VALUES (")
	for i := range columns {
		query += fmt.Sprintf("$%d", i+1)
		if i < len(values)-1 {
			query += ","
		}
	}
	query = fmt.Sprintf("%s %s", query, ") RETURNING id")
	return query, values

}
