// nolint
package usecases

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"partner/internal/consts"
	"partner/internal/entities"
	"partner/internal/repo/mock"
	"testing"

	cacheConf "gitlab.com/tuneverse/toolkit/core/cache"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils/crypto"
)

var (
	logLevel   = "info"
	Currency   = "INR"
	Language   = "EN"
	Country    = "IN"
	PostalCode = 122454
)

func TestGetPartnerOauthCredential(t *testing.T) {

	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            logLevel,
		IncludeRequestDump:  false,
		IncludeResponseDump: true,
	}
	_ = logger.InitLogger(clientOpt)

	clientId, err := crypto.Encrypt("hello", []byte("tuneverse-esreve"))
	require.Nil(t, err)
	clientSecret, err := crypto.Encrypt("hello@123", []byte("tuneverse-esreve"))
	require.Nil(t, err)
	values := entities.GetPartnerOauthCredential{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
	InvalidClientIdValues := entities.GetPartnerOauthCredential{
		ClientSecret: clientSecret,
		ClientId:     "clientId",
	}
	InvalidClientSecretValues := entities.GetPartnerOauthCredential{
		ClientSecret: "clientsecret",
		ClientId:     "clientId",
	}
	testCases := []struct {
		name          string
		partnerId     string
		buildStubs    func(partner *mock.MockPartnerRepoImply)
		checkResponse func(t *testing.T, OauthCrdential entities.GetPartnerOauthCredential, err error)
	}{
		{
			name:      "ok",
			partnerId: "a2a00db1-3c9a-4b75-9907-6ef55b63d379",
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					GetPartnerOauthCredential(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return(values, nil)
			},
			checkResponse: func(t *testing.T, OauthCredential entities.GetPartnerOauthCredential, err error) {
				require.NotNil(t, OauthCredential)
			},
		},
		{
			name:      "no data",
			partnerId: "",
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					GetPartnerOauthCredential(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return(entities.GetPartnerOauthCredential{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, OauthCrdential entities.GetPartnerOauthCredential, err error) {
				require.Empty(t, OauthCrdential)
				require.Error(t, err)
			},
		},
		{
			name:      "invalid client_id",
			partnerId: "a2a00db1-3c9a-4b75-9907-6ef55b63d379",
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					GetPartnerOauthCredential(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return(InvalidClientIdValues, nil)
			},
			checkResponse: func(t *testing.T, OauthCredential entities.GetPartnerOauthCredential, err error) {
				require.NotNil(t, err)
			},
		},
		{
			name:      "invalid client_secret",
			partnerId: "a2a00db1-3c9a-4b75-9907-6ef55b63d379",
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					GetPartnerOauthCredential(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return(InvalidClientSecretValues, nil)
			},
			checkResponse: func(t *testing.T, OauthCredential entities.GetPartnerOauthCredential, err error) {
				require.NotNil(t, err)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
				Host:     "",
				UserName: "",
				Password: "",
				DB:       0,
			})

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			partner := mock.NewMockPartnerRepoImply(ctrl)
			tc.buildStubs(partner)
			partnerUsecase := NewPartnerUseCases(partner, redisClient)

			data, err := partnerUsecase.GetPartnerOauthCredential(context.Background(), tc.partnerId, entities.PartnerOAuthHeader{}, "", "", map[string]models.ErrorResponse{})
			tc.checkResponse(t, data, err)

		})
	}
}

func TestGetPartnerById(t *testing.T) {
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName, // Service name.
		LogLevel:            "info",         // Log level.
		IncludeRequestDump:  false,          // Include request data in logs.
		IncludeResponseDump: false,          // Include response data in logs.
	}
	_ = logger.InitLogger(clientOpt)
	// Create a mock HTTP request
	req, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP response recorder
	resp := httptest.NewRecorder()

	// Create a Gin context from the request and response
	c, _ := gin.CreateTestContext(resp)
	c.Request = req
	// Create a new controller for the mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock instance of your repository
	mockRepo := mock.NewMockPartnerRepoImply(ctrl)
	redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
		Host:     "",
		UserName: "",
		Password: "",
		DB:       0,
	})

	// Create an instance of the PartnerUseCases
	partnerUseCases := NewPartnerUseCases(mockRepo, redisClient) // Replace with your actual constructor function

	// Prepare the expected data and context
	expectedPartner := entities.GetPartner{} // Define your expected partner data
	ctx := context.TODO()                    // Create a context

	// Set expectations on the mock
	mockRepo.EXPECT().
		GetPartnerById(ctx, "partnerID").
		Return(expectedPartner, nil)

	// Call the function you want to test (ViewPartnerByID)
	actualPartner, _, _, err := partnerUseCases.GetPartnerById(c, entities.QueryParams{}, "partnerID", "", "", map[string]models.ErrorResponse{})

	// Assert the expected outcomes
	assert.NoError(t, err)                          // Check if no error occurred
	assert.Equal(t, expectedPartner, actualPartner) // Check if the returned partner matches the expected data
}

func TestGetPartnerByIdError(t *testing.T) {
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName, // Service name.
		LogLevel:            "info",         // Log level.
		IncludeRequestDump:  false,          // Include request data in logs.
		IncludeResponseDump: false,          // Include response data in logs.
	}
	_ = logger.InitLogger(clientOpt)
	// Create a mock HTTP request
	req, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP response recorder
	resp := httptest.NewRecorder()

	// Create a Gin context from the request and response
	c, _ := gin.CreateTestContext(resp)
	c.Request = req
	// Create a new instance of the gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock database repository
	mockRepo := mock.NewMockPartnerRepoImply(ctrl)
	redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
		Host:     "",
		UserName: "",
		Password: "",
		DB:       0,
	})
	// Create a use case instance that uses the mock repository
	partnerUseCase := NewPartnerUseCases(mockRepo, redisClient)

	// Define a test partner ID
	testPartnerID := "12345"

	// Define an error that you expect to be returned by the mock
	expectedError := errors.New("database error")

	// Set up the expected behavior of the mock repository to return an error
	mockRepo.EXPECT().
		GetPartnerById(gomock.Any(), testPartnerID).
		Return(entities.GetPartner{}, expectedError).
		Times(1)

	// Call the function you want to test (ViewPartnerById)
	viewedPartner, _, _, err := partnerUseCase.GetPartnerById(c, entities.QueryParams{}, testPartnerID, "", "", map[string]models.ErrorResponse{})

	// Assert the expected outcomes
	assert.Error(t, err) // An error is expected
	assert.Equal(t, expectedError, err)
	assert.Equal(t, entities.GetPartner{}, viewedPartner) // Partner should be empty
}

func TestGetTAllTermsAndConditions(t *testing.T) {
	// Initialize a new controller
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            logLevel,
		IncludeRequestDump:  false,
		IncludeResponseDump: false,
	}
	logger.InitLogger(clientOpt)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockPartnerRepoImply(ctrl)
	redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
		Host:     "",
		UserName: "",
		Password: "",
		DB:       0,
	})
	useCase := NewPartnerUseCases(mockRepo, redisClient)

	expectedTermsAndConditions := entities.TermsAndConditions{
		Name:        "Test Name",
		Description: "Test Description",
		Language:    "en",
	}

	mockRepo.EXPECT().
		GetAllTermsAndConditions(gomock.Any(), "614608f2-6538-4733-aded-96f902007254").Times(1).Return(expectedTermsAndConditions, nil)
	result, err := useCase.GetAllTermsAndConditions(context.Background(), "614608f2-6538-4733-aded-96f902007254", "", "", map[string]models.ErrorResponse{})
	require.Equal(t, expectedTermsAndConditions, result)
	require.NoError(t, err)

	mockRepo.EXPECT().
		GetAllTermsAndConditions(gomock.Any(), "123").Times(1)
	result, _ = useCase.GetAllTermsAndConditions(context.Background(), "123", "", "", map[string]models.ErrorResponse{})
	require.Equal(t, entities.TermsAndConditions{}, result)
}

func TestUpdatePartnerSuccess(t *testing.T) {
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            "info",
		IncludeRequestDump:  false,
		IncludeResponseDump: true,
	}

	// Check if the application is in debug mode.

	_ = logger.InitLogger(clientOpt)
	id := uuid.New()
	testCases := []struct {
		name          string
		ID            uuid.UUID
		partner       entities.PartnerProperties
		buildStubs    func(store *mock.MockPartnerRepoImply)
		checkResponse func(t *testing.T, validationErr map[string]models.ErrorResponse, err error)
	}{
		{
			name: "ok",
			ID:   id,
			partner: map[string]interface{}{
				"name": "Mame2345",
				"url":  "https://www.aathi987.com",
				"contact_details": map[string]interface{}{
					"contact_person": "contact person",
					"email":          "aathi987@email.com",
					"noreply_email":  "aathi987@email.com",
					"feedback_email": "aathi987@email.com",
					"support_email":  "aathi987@email.com",
				},
				"language":               "en",
				"browser_title":          "aathi987",
				"profile_url":            "https://www.aathi987.com",
				"payment_url":            "https://www.partnerpayment.com",
				"landing_page":           "https://www.aathi987.com",
				"mobile_verify_interval": 2,
				"payout_target_currency": "USD",
				"theme_id":               1,
				"enable_mail":            true,
				"login_type":             "Normal",
				"business_model":         "Revenue Share",
				"member_pay_to_partner":  true,
				"address_details": map[string]interface{}{
					"address":     "sample partner addressyy",
					"street":      "partner streetyy",
					"country":     "IN",
					"state":       "KL",
					"city":        "partner cityy",
					"postal_code": "695583",
				},
				"subscription_details": map[string]interface{}{
					"plan_id":          "f45d6151-04c1-434d-8bc9-bd5d0bedbcb9",
					"plan_start_date":  "2023-05-22",
					"plan_launch_date": "2023-05-01",
				},
				"payment": map[string]interface{}{
					"payment_gateways": []interface{}{
						map[string]interface{}{
							"gateway":                 "Paypal",
							"email":                   "paypal@gmail.com",
							"client_id":               "01596d6c1205c95fe15e6cf1e3170c34df6e53b6a693a169abf766c1a1c",
							"client_secret":           "33d9462ab77b99fcbcc0fd4db59627d30d100d658bb906d8a60df9bd3f5a7d06",
							"payin":                   true,
							"payout":                  false,
							"default_payin_currency":  "USD",
							"default_payout_currency": "USD",
						},
						map[string]interface{}{
							"gateway":                 "PayTm",
							"email":                   "stripe@gmail.com",
							"client_id":               "01596d6c1acfc5c95fe15e6cf1e3170c34df6e53b6a693a169abf766c1a1c",
							"client_secret":           "33d9462ab77b99fcbcc0fd4db59627d30d100d658bb906d8a60df9bd3f5a7d06",
							"payin":                   true,
							"payout":                  true,
							"default_payin_currency":  "USD",
							"default_payout_currency": "USD",
						},
					},
					"default_payment_gateway":         "Paypal",
					"payout_min_limit":                2,
					"default_currency":                "USD",
					"payout_max_remittance_per_month": 1,
				},
				"member_grace_period":  "Quarterly",
				"expiry_warning_count": 5,
				"album_review_email":   "aathi987@email.com",
				"site_info":            "sample site info",
				"default_price_code_currency": map[string]interface{}{
					"name": "USD",
				},
				"music_language":              "en",
				"member_default_country":      "IN",
				"outlets_processing_duration": 10,
				"free_plan_limit":             4,
				"product_review":              "Both",
				"website_url":                 "https://www.aathi987.com",
				"background_image":            "https://background.com/image",
				"background_color":            "#ff0000",
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					UpdatePartner(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
				partner.EXPECT().IsFieldValueUnique(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(9).Return(true, nil)
				partner.EXPECT().GetID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
			},
			checkResponse: func(t *testing.T, validationErr map[string]models.ErrorResponse, err error) {

				require.Equal(t, 0, len(validationErr))
				require.Nil(t, err)
			},
		},
	}

	for i := range testCases {

		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
				Host:     "",
				UserName: "",
				Password: "",
				DB:       0,
			})
			partner := mock.NewMockPartnerRepoImply(ctrl)
			tc.buildStubs(partner)
			partnerUsecase := NewPartnerUseCases(partner, redisClient)
			validationErr, err := partnerUsecase.UpdatePartner(context.Background(), tc.ID.String(), tc.ID, tc.partner, "", "", map[string]models.ErrorResponse{})

			tc.checkResponse(t, validationErr, err)
		})
	}

}

func TestUpdatePartnerMinimalFields(t *testing.T) {
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            "info",
		IncludeRequestDump:  false,
		IncludeResponseDump: true,
	}

	// Check if the application is in debug mode.

	_ = logger.InitLogger(clientOpt)
	id := uuid.New()
	testCases := []struct {
		name          string
		ID            uuid.UUID
		partner       entities.PartnerProperties
		buildStubs    func(store *mock.MockPartnerRepoImply)
		checkResponse func(t *testing.T, validationErr map[string]models.ErrorResponse, err error)
	}{
		{
			name: "update with minimal fields",
			ID:   id,
			partner: map[string]interface{}{
				"name":                   "Mame2345",
				"language":               "en",
				"browser_title":          "fjahghdddddd",
				"profile_url":            "https://www.partnerpayment.com",
				"payment_url":            "https://www.partnerpayment.com",
				"landing_page":           "https://www.ccccccccc.com",
				"mobile_verify_interval": 2,
				"payout_target_currency": "USD",
				"theme_id":               1,
				"enable_mail":            true,
				"login_type":             "normal",
				"business_model":         "Revenue Share",
				"member_pay_to_partner":  true,
				"payment": map[string]interface{}{
					"payment_gateways": []interface{}{
						map[string]interface{}{
							"gateway":                 "Paypal",
							"email":                   "paypal@gmail.com",
							"client_id":               "01596d6c1205c95fe15e6cf1e3170c34df6e53b6a693a169abf766c1a1c",
							"client_secret":           "33d9462ab77b99fcbcc0fd4db59627d30d100d658bb906d8a60df9bd3f5a7d06",
							"payin":                   true,
							"payout":                  false,
							"default_payin_currency":  "USD",
							"default_payout_currency": "USD",
						},
						map[string]interface{}{
							"gateway":                 "PayTm",
							"email":                   "stripe@gmail.com",
							"client_id":               "01596d6c1acfc5c95fe15e6cf1e3170c34df6e53b6a693a169abf766c1a1c",
							"client_secret":           "33d9462ab77b99fcbcc0fd4db59627d30d100d658bb906d8a60df9bd3f5a7d06",
							"payin":                   true,
							"payout":                  true,
							"default_payin_currency":  "USD",
							"default_payout_currency": "USD",
						},
					},
					"default_payment_gateway":         "Paypal",
					"payout_min_limit":                201,
					"default_currency":                "USD",
					"payout_max_remittance_per_month": 1,
				},
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					UpdatePartner(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, nil)
				partner.EXPECT().IsFieldValueUnique(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(2).Return(true, nil)
			},
			checkResponse: func(t *testing.T, validationErr map[string]models.ErrorResponse, err error) {

				require.Equal(t, 0, len(validationErr))
				require.Nil(t, err)
			},
		},
	}

	for i := range testCases {

		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
				Host:     "",
				UserName: "",
				Password: "",
				DB:       0,
			})
			store := mock.NewMockPartnerRepoImply(ctrl)
			tc.buildStubs(store)
			partnerUsecase := NewPartnerUseCases(store, redisClient)
			validationErr, err := partnerUsecase.UpdatePartner(context.Background(), tc.ID.String(), tc.ID, tc.partner, "", "", map[string]models.ErrorResponse{})

			tc.checkResponse(t, validationErr, err)
		})
	}

}
func TestUpdatePartnerEmptyFields(t *testing.T) {
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            "info",
		IncludeRequestDump:  false,
		IncludeResponseDump: true,
	}

	// Check if the application is in debug mode.

	_ = logger.InitLogger(clientOpt)
	id := uuid.New()
	testCases := []struct {
		name          string
		ID            uuid.UUID
		partner       entities.PartnerProperties
		buildStubs    func(store *mock.MockPartnerRepoImply)
		checkResponse func(t *testing.T, validationErr map[string]models.ErrorResponse, err error)
	}{
		{
			name: "validation error for empty fields",
			ID:   id,
			partner: map[string]interface{}{
				"name": "",
				"url":  "",
				"contact_details": map[string]interface{}{
					"contact_person": "",
					"email":          "",
					"noreply_email":  "",
					"feedback_email": "",
					"support_email":  "",
				},
				"language":               "ennnn",
				"browser_title":          "fjahghdddddd",
				"profile_url":            "partnerpayme",
				"payment_url":            "partnerpayme",
				"landing_page":           "ccccccc",
				"mobile_verify_interval": 2,
				"payout_target_currency": "USD",
				"theme_id":               1,
				"enable_mail":            true,
				"login_type":             "normal",
				"business_model":         "revenue_share",
				"member_pay_to_partner":  true,
				"address_details": map[string]interface{}{
					"address":     "",
					"street":      "partner streetyy",
					"country":     "",
					"state":       "",
					"city":        "partner cityy",
					"postal_code": "695583",
				},
				"subscription_details": map[string]interface{}{
					"plan_id":          "f45d6151-04c1-434d-8bc9-bd5d0bedbcb9",
					"plan_start_date":  "2023-05-22",
					"plan_launch_date": "2023-05-01",
				},
				"payment": map[string]interface{}{
					"payment_gateways": []interface{}{
						map[string]interface{}{
							"gateway":                 "",
							"email":                   "",
							"client_id":               "",
							"client_secret":           "",
							"payin":                   true,
							"payout":                  false,
							"default_payin_currency":  "",
							"default_payout_currency": "",
						},
						map[string]interface{}{
							"gateway":                 "",
							"email":                   "",
							"client_id":               "",
							"client_secret":           "",
							"payin":                   true,
							"payout":                  true,
							"default_payin_currency":  "",
							"default_payout_currency": "",
						},
					},
					"default_payment_gateway":         "",
					"payout_min_limit":                0,
					"default_currency":                "",
					"payout_max_remittance_per_month": 1,
				},
				"member_grace_period":  "Quarterly",
				"expiry_warning_count": 5,
				"album_review_email":   "",
				"site_info":            "sample site info",
				"default_price_code_currency": map[string]interface{}{
					"name": "USD",
				},
				"music_language":              "en",
				"member_default_country":      "IN",
				"outlets_processing_duration": 7,
				"free_plan_limit":             9,
				"product_review":              "both",
				"website_url":                 "",
				"background_image":            "https://background.com/image",
				"background_color":            "#0",
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {
				partner.EXPECT().GetID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
			},
			checkResponse: func(t *testing.T, validationErr map[string]models.ErrorResponse, err error) {
				require.NotNil(t, validationErr)
			},
		},
	}

	for i := range testCases {

		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
				Host:     "",
				UserName: "",
				Password: "",
				DB:       0,
			})
			partner := mock.NewMockPartnerRepoImply(ctrl)
			tc.buildStubs(partner)
			partnerUsecase := NewPartnerUseCases(partner, redisClient)
			validationErr, err := partnerUsecase.UpdatePartner(context.Background(), tc.ID.String(), tc.ID, tc.partner, "", "", map[string]models.ErrorResponse{})

			tc.checkResponse(t, validationErr, err)
		})
	}

}

func TestUpdatePartnerInvalidFields(t *testing.T) {
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            "info",
		IncludeRequestDump:  false,
		IncludeResponseDump: true,
	}

	// Check if the application is in debug mode.

	_ = logger.InitLogger(clientOpt)
	id := uuid.New()
	testCases := []struct {
		name          string
		ID            uuid.UUID
		partner       entities.PartnerProperties
		buildStubs    func(store *mock.MockPartnerRepoImply)
		checkResponse func(t *testing.T, validationErr map[string]models.ErrorResponse, err error)
	}{
		{
			name: "validation error for invalid fields",
			ID:   id,
			partner: map[string]interface{}{
				"name": "Mame2345",
				"url":  "ccccc.com",
				"contact_details": map[string]interface{}{
					"contact_person": "contact person",
					"email":          "ccccccccccc",
					"noreply_email":  "ccccccccc",
					"feedback_email": "cccccccc",
					"support_email":  "ccccccccc",
				},
				"language":               "ennnn",
				"browser_title":          "fjahghdddddd",
				"profile_url":            "partnerpayme",
				"payment_url":            "partnerpayme",
				"landing_page":           "ccccccc",
				"mobile_verify_interval": 2,
				"payout_target_currency": "USD",
				"theme_id":               1,
				"enable_mail":            true,
				"login_type":             "normal",
				"business_model":         "Revenue Share",
				"member_pay_to_partner":  true,
				"address_details": map[string]interface{}{
					"address":     "",
					"street":      "partner streetyy",
					"country":     "",
					"state":       "",
					"city":        "partner cityy",
					"postal_code": "695583",
				},
				"subscription_details": map[string]interface{}{
					"plan_id":          "f45d6151-04c1-434d-8bc9-bd5d0bedbcb9",
					"plan_start_date":  "2023-05-22",
					"plan_launch_date": "2023-05-01",
				},
				"payment": map[string]interface{}{
					"payment_gateways": []interface{}{
						map[string]interface{}{
							"gateway":                 "",
							"email":                   "",
							"client_id":               "",
							"client_secret":           "",
							"payin":                   true,
							"payout":                  false,
							"default_payin_currency":  "",
							"default_payout_currency": "",
						},
						map[string]interface{}{
							"gateway":                 "",
							"email":                   "",
							"client_id":               "",
							"client_secret":           "",
							"payin":                   true,
							"payout":                  true,
							"default_payin_currency":  "",
							"default_payout_currency": "",
						},
					},
					"default_payment_gateway":         "",
					"payout_min_limit":                0,
					"default_currency":                "",
					"payout_max_remittance_per_month": 1,
				},
				"member_grace_period":  "Quarterly",
				"expiry_warning_count": 45,
				"album_review_email":   "cccccccc",
				"site_info":            "sample site info",
				"default_price_code_currency": map[string]interface{}{
					"name": "USD",
				},
				"music_language":              "en",
				"member_default_country":      "IN",
				"outlets_processing_duration": 7,
				"free_plan_limit":             9,
				"product_review":              "Both",
				"website_url":                 "cccccc",
				"background_image":            "https://background.com/image",
				"background_color":            "#0",
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().IsFieldValueUnique(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(2).Return(true, nil)
				partner.EXPECT().GetID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
			},
			checkResponse: func(t *testing.T, validationErr map[string]models.ErrorResponse, err error) {
				require.NotNil(t, validationErr)
			},
		},
	}

	for i := range testCases {

		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
				Host:     "",
				UserName: "",
				Password: "",
				DB:       0,
			})
			partner := mock.NewMockPartnerRepoImply(ctrl)
			tc.buildStubs(partner)
			partnerUsecase := NewPartnerUseCases(partner, redisClient)
			validationErr, err := partnerUsecase.UpdatePartner(context.Background(), tc.ID.String(), tc.ID, tc.partner, "", "", map[string]models.ErrorResponse{})

			tc.checkResponse(t, validationErr, err)
		})
	}

}

func TestUpdatePartnerMaxLengthErr(t *testing.T) {
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            "info",
		IncludeRequestDump:  false,
		IncludeResponseDump: true,
	}

	// Check if the application is in debug mode.

	_ = logger.InitLogger(clientOpt)
	id := uuid.New()
	testCases := []struct {
		name          string
		ID            uuid.UUID
		partner       entities.PartnerProperties
		buildStubs    func(store *mock.MockPartnerRepoImply)
		checkResponse func(t *testing.T, validationErr map[string]models.ErrorResponse, err error)
	}{
		{
			name: "validation error for maximum length",
			ID:   id,
			partner: map[string]interface{}{
				"name": "Lorem ipsum dolor sit amet, consectetur adipidfggafcing elit. Fusce vitae dolor vel urna placerat dignissim. Vivamus dictum ante nec arcu lobortis, eu tempor quam vehicula. Proin vitae quam vitae eros tincidunt auctor ac in risus. Maecenas non purus libero. Suspendisse potenti. Sed sed odio sed justo vulputate iaculis a sed odio. Quisque ut urna ac turpis aliquam bibendum. Duis euismod erat nec elit malesuada, et malesuada nisi rhoncus. Pellentesque sit amet felis eget eros sodales auctor a quis metus. Ut placerat dolor vitae libero vehicula, non faucibus est vehicula. Mauris sit amet nisi ipsum. Integer eleifend tellus eu augue dictum, non ultricies eros tincidunt. Aliquam tincidunt nulla id leo sagittis posuere. Vestibulum pharetra scelerisque ex, in semper odio volutpat vel. Cras lobortis elit vel velit fermentum, a venenatis nisi dapibus",
				"url":  "ccccc.com",
				"contact_details": map[string]interface{}{
					"contact_person": "Lorem ipsum dolor sit amet, conseuictetur adipiscing elit. Fusce vitae dolor vel urna placerat dignissim. Vivamus dictum ante nec arcu lobortis, eu tempor quam vehicula. Proin vitae quam vitae eros tincidunt auctor ac in risus. Maecenas non purus libero. Suspendisse potenti. Sed sed odio sed justo vulputate iaculis a sed odio. Quisque ut urna ac turpis aliquam bibendum. Duis euismod erat nec elit malesuada, et malesuada nisi rhoncus. Pellentesque sit amet felis eget eros sodales auctor a quis metus. Ut placerat dolor vitae libero vehicula, non faucibus est vehicula. Mauris sit amet nisi ipsum. Integer eleifend tellus eu augue dictum, non ultricies eros tincidunt. Aliquam tincidunt nulla id leo sagittis posuere. Vestibulum pharetra scelerisque ex, in semper odio volutpat vel. Cras lobortis elit vel velit fermentum, a venenatis nisi dapibus",
					"email":          "ccccccccccc",
					"noreply_email":  "ccccccccc",
					"feedback_email": "cccccccc",
					"support_email":  "ccccccccc",
				},
				"language":               "ennnn",
				"browser_title":          "Lorem ipsum dolor sit amet, conseiictetur adipiscing elit. Fusce vitae dolor vel urna placerat dignissim. Vivamus dictum ante nec arcu lobortis, eu tempor quam vehicula. Proin vitae quam vitae eros tincidunt auctor ac in risus. Maecenas non purus libero. Suspendisse potenti. Sed sed odio sed justo vulputate iaculis a sed odio. Quisque ut urna ac turpis aliquam bibendum. Duis euismod erat nec elit malesuada, et malesuada nisi rhoncus. Pellentesque sit amet felis eget eros sodales auctor a quis metus. Ut placerat dolor vitae libero vehicula, non faucibus est vehicula. Mauris sit amet nisi ipsum. Integer eleifend tellus eu augue dictum, non ultricies eros tincidunt. Aliquam tincidunt nulla id leo sagittis posuere. Vestibulum pharetra scelerisque ex, in semper odio volutpat vel. Cras lobortis elit vel velit fermentum, a venenatis nisi dapibus",
				"profile_url":            "partnerpayme",
				"payment_url":            "partnerpayme",
				"landing_page":           "ccccccc",
				"mobile_verify_interval": 2,
				"payout_target_currency": "USD",
				"theme_id":               1,
				"enable_mail":            true,
				"login_type":             "normal",
				"business_model":         "revenue_share",
				"member_pay_to_partner":  true,
				"address_details": map[string]interface{}{
					"address":     "Lorem ipsum dolor sit amet, consectetur adipiscyyying elit. Fusce vitae dolor vel urna placerat dignissim. Vivamus dictum ante nec arcu lobortis, eu tempor quam vehicula. Proin vitae quam vitae eros tincidunt auctor ac in risus. Maecenas non purus libero. Suspendisse potenti. Sed sed odio sed justo vulputate iaculis a sed odio. Quisque ut urna ac turpis aliquam bibendum. Duis euismod erat nec elit malesuada, et malesuada nisi rhoncus. Pellentesque sit amet felis eget eros sodales auctor a quis metus. Ut placerat dolor vitae libero vehicula, non faucibus est vehicula. Mauris sit amet nisi ipsum. Integer eleifend tellus eu augue dictum, non ultricies eros tincidunt. Aliquam tincidunt nulla id leo sagittis posuere. Vestibulum pharetra scelerisque ex, in semper odio volutpat vel. Cras lobortis elit vel velit fermentum, a venenatis nisi dapibus",
					"street":      "Lorem ipsum dolor sit amet, consectetur adipiseecing elit. Fusce vitae dolor vel urna placerat dignissim. Vivamus dictum ante nec arcu lobortis, eu tempor quam vehicula. Proin vitae quam vitae eros tincidunt auctor ac in risus. Maecenas non purus libero. Suspendisse potenti. Sed sed odio sed justo vulputate iaculis a sed odio. Quisque ut urna ac turpis aliquam bibendum. Duis euismod erat nec elit malesuada, et malesuada nisi rhoncus. Pellentesque sit amet felis eget eros sodales auctor a quis metus. Ut placerat dolor vitae libero vehicula, non faucibus est vehicula. Mauris sit amet nisi ipsum. Integer eleifend tellus eu augue dictum, non ultricies eros tincidunt. Aliquam tincidunt nulla id leo sagittis posuere. Vestibulum pharetra scelerisque ex, in semper odio volutpat vel. Cras lobortis elit vel velit fermentum, a venenatis nisi dapibus",
					"country":     "",
					"state":       "",
					"city":        "Lorem ipsum dolor sit amet, consectetur adipiscffing elit. Fusce vitae dolor vel urna placerat dignissim. Vivamus dictum ante nec arcu lobortis, eu tempor quam vehicula. Proin vitae quam vitae eros tincidunt auctor ac in risus. Maecenas non purus libero. Suspendisse potenti. Sed sed odio sed justo vulputate iaculis a sed odio. Quisque ut urna ac turpis aliquam bibendum. Duis euismod erat nec elit malesuada, et malesuada nisi rhoncus. Pellentesque sit amet felis eget eros sodales auctor a quis metus. Ut placerat dolor vitae libero vehicula, non faucibus est vehicula. Mauris sit amet nisi ipsum. Integer eleifend tellus eu augue dictum, non ultricies eros tincidunt. Aliquam tincidunt nulla id leo sagittis posuere. Vestibulum pharetra scelerisque ex, in semper odio volutpat vel. Cras lobortis elit vel velit fermentum, a venenatis nisi dapibus",
					"postal_code": "695583898000000",
				},
				"subscription_details": map[string]interface{}{
					"plan_id":          "f45d6151-04c1-434d-8bc9-bd5d0bedbcb9",
					"plan_start_date":  "2023-05-22",
					"plan_launch_date": "2023-05-01",
				},
				"payment": map[string]interface{}{
					"payment_gateways": []interface{}{
						map[string]interface{}{
							"gateway":                 "Paypal",
							"email":                   "paypal@gmail.com",
							"client_id":               "01596d6c1205c95fe15e6cf1e3170c34df6e53b6a693a169abf766c1a1c",
							"client_secret":           "33d9462ab77b99fcbcc0fd4db59627d30d100d658bb906d8a60df9bd3f5a7d06",
							"payin":                   true,
							"payout":                  false,
							"default_payin_currency":  "UUU",
							"default_payout_currency": "UUU",
						},
						map[string]interface{}{
							"gateway":                 "PayTm",
							"email":                   "stripe@gmail.com",
							"client_id":               "01596d6c1acfc5c95fe15e6cf1e3170c34df6e53b6a693a169abf766c1a1c",
							"client_secret":           "33d9462ab77b99fcbcc0fd4db59627d30d100d658bb906d8a60df9bd3f5a7d06",
							"payin":                   true,
							"payout":                  true,
							"default_payin_currency":  "UUU",
							"default_payout_currency": "UUU",
						},
					},
					"default_payment_gateway":         "ppppp",
					"payout_min_limit":                0,
					"default_currency":                "",
					"payout_max_remittance_per_month": 1,
				},
				"member_grace_period":  "Quarterly",
				"expiry_warning_count": 5,
				"album_review_email":   "cccccccc",
				"site_info":            "sample site info",
				"default_price_code_currency": map[string]interface{}{
					"name": "USD",
				},
				"music_language":              "en",
				"member_default_country":      "IN",
				"outlets_processing_duration": 40,
				"free_plan_limit":             100,
				"product_review":              "both",
				"website_url":                 "",
				"background_image":            "https://background.com/image",
				"background_color":            "#0",
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {
				partner.EXPECT().IsFieldValueUnique(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(true, nil)
				partner.EXPECT().GetID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
			},
			checkResponse: func(t *testing.T, validationErr map[string]models.ErrorResponse, err error) {
				require.NotNil(t, validationErr)
			},
		},
	}

	for i := range testCases {

		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
				Host:     "",
				UserName: "",
				Password: "",
				DB:       0,
			})
			partner := mock.NewMockPartnerRepoImply(ctrl)
			tc.buildStubs(partner)
			partnerUsecase := NewPartnerUseCases(partner, redisClient)
			validationErr, err := partnerUsecase.UpdatePartner(context.Background(), tc.ID.String(), tc.ID, tc.partner, "", "", map[string]models.ErrorResponse{})

			tc.checkResponse(t, validationErr, err)
		})
	}

}
func TestUpdatePartnerServerError(t *testing.T) {

	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            "info",
		IncludeRequestDump:  false,
		IncludeResponseDump: true,
	}

	// Check if the application is in debug mode.

	_ = logger.InitLogger(clientOpt)
	id := uuid.New()
	testCases := []struct {
		name          string
		ID            uuid.UUID
		partner       entities.PartnerProperties
		buildStubs    func(store *mock.MockPartnerRepoImply)
		checkResponse func(t *testing.T, validationErr map[string]models.ErrorResponse, err error)
	}{

		{
			name: "internal server error",
			ID:   id,
			partner: map[string]interface{}{
				"name": "Mame2345",
				"url":  "https://www.ccccc.com",
				"contact_details": map[string]interface{}{
					"contact_person": "contact person",
					"email":          "cccccccccccc@email.com",
					"noreply_email":  "ccccccccccc@email.com",
					"feedback_email": "ccccccccccc@email.com",
					"support_email":  "ccccccccccc@email.com",
				},
				"language":               "en",
				"browser_title":          "fjahghdddddd",
				"profile_url":            "https://www.partnerpayment.com",
				"payment_url":            "https://www.partnerpayment.com",
				"landing_page":           "https://www.ccccccccc.com",
				"mobile_verify_interval": 2,
				"payout_target_currency": "USD",
				"theme_id":               1,
				"enable_mail":            true,
				"login_type":             "normal",
				"business_model":         "Revenue Share",
				"member_pay_to_partner":  true,
				"address_details": map[string]interface{}{
					"address":     "sample partner addressyy",
					"street":      "partner streetyy",
					"country":     "IN",
					"state":       "KL",
					"city":        "partner cityy",
					"postal_code": "695583",
				},
				"subscription_details": map[string]interface{}{
					"plan_id":          "f45d6151-04c1-434d-8bc9-bd5d0bedbcb9",
					"plan_start_date":  "2023-05-22",
					"plan_launch_date": "2023-05-01",
				},
				"payment": map[string]interface{}{
					"payment_gateways": []interface{}{
						map[string]interface{}{
							"gateway":                 "Paypal",
							"email":                   "paypal@gmail.com",
							"client_id":               "01596d6c1205c95fe15e6cf1e3170c34df6e53b6a693a169abf766c1a1c",
							"client_secret":           "33d9462ab77b99fcbcc0fd4db59627d30d100d658bb906d8a60df9bd3f5a7d06",
							"payin":                   true,
							"payout":                  false,
							"default_payin_currency":  "USD",
							"default_payout_currency": "USD",
						},
						map[string]interface{}{
							"gateway":                 "PayTm",
							"email":                   "stripe@gmail.com",
							"client_id":               "01596d6c1acfc5c95fe15e6cf1e3170c34df6e53b6a693a169abf766c1a1c",
							"client_secret":           "33d9462ab77b99fcbcc0fd4db59627d30d100d658bb906d8a60df9bd3f5a7d06",
							"payin":                   true,
							"payout":                  true,
							"default_payin_currency":  "USD",
							"default_payout_currency": "USD",
						},
					},
					"default_payment_gateway":         "Paypal",
					"payout_min_limit":                207,
					"default_currency":                "USD",
					"payout_max_remittance_per_month": 1,
				},
				"member_grace_period":  "Quarterly",
				"expiry_warning_count": 5,
				"album_review_email":   "cccccccccc@email.com",
				"site_info":            "sample site info",
				"default_price_code_currency": map[string]interface{}{
					"name": "USD",
				},
				"music_language":              "en",
				"member_default_country":      "IN",
				"outlets_processing_duration": 7,
				"free_plan_limit":             9,
				"product_review":              "Both",
				"website_url":                 "https://www.cccccccccc.com",
				"background_image":            "https://background.com/image",
				"background_color":            "#ff0000",
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					UpdatePartner(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, sql.ErrConnDone)
				partner.EXPECT().IsFieldValueUnique(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(9).Return(true, nil)
				partner.EXPECT().GetID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
			},
			checkResponse: func(t *testing.T, validationErr map[string]models.ErrorResponse, err error) {

				require.Nil(t, validationErr)
				require.Error(t, err)
			},
		},
	}

	for i := range testCases {

		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
				Host:     "",
				UserName: "",
				Password: "",
				DB:       0,
			})
			partner := mock.NewMockPartnerRepoImply(ctrl)
			tc.buildStubs(partner)
			partnerUsecase := NewPartnerUseCases(partner, redisClient)
			validationErr, err := partnerUsecase.UpdatePartner(context.Background(), tc.ID.String(), tc.ID, tc.partner, "", "", map[string]models.ErrorResponse{})

			tc.checkResponse(t, validationErr, err)
		})
	}
}

func TestCreatePartnerFailure(t *testing.T) {
	// logger.InitLogger(clientOpt)

	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName, // Service name.
		LogLevel:            logLevel,       // Log level.
		IncludeRequestDump:  false,          // Include request data in logs.
		IncludeResponseDump: true,           // Include response data in logs.
	}
	// Check if the application is in debug mode.

	_ = logger.InitLogger(clientOpt)
	testCases := []struct {
		name          string
		partner       entities.Partner
		buildStubs    func(partner *mock.MockPartnerRepoImply)
		checkResponse func(t *testing.T, validationErr map[string]models.ErrorResponse, id string, err error)
	}{

		{
			name: "internal server error",
			partner: entities.Partner{
				Name:                 "name B",
				URL:                  "https://www.website.com",
				Logo:                 "https://www.website.com/logo.png",
				Favicon:              "https://www.website.com/favicon.ico",
				Loader:               "https://www.website.com/loader.ico",
				LoginPageLogo:        "https://www.website.com/pagelogo.ico",
				BackgroundColor:      "#F0F0F0",
				BackgroundImage:      "https://www.website.com/background.jpg",
				WebsiteURL:           "https://www.pqr.com",
				Language:             "en",
				BrowserTitle:         "Welcome to My Website",
				ProfileURL:           "https://www.website.com/pqr",
				PaymentURL:           "https://www.website.com/payment",
				LandingPage:          "https://www.website.com/landing",
				MobileVerifyInterval: 0,
				PayoutTargetCurrency: "USD",
				ThemeID:              0,
				EnableMail:           true,
				LoginType:            "Normal",
				BusinessModel:        "Subscription",
				MemberPayToPartner:   true,
				MemberGracePeriod:    "Monthly",
				ExpiryWarningCount:   1,
				AlbumReviewEmail:     "reviews@mywebsite.com",
				SiteInfo:             "Welcome to My Website!",
				DefaultPriceCode:     1,
				DefaultPriceCodeCurrency: entities.Currency{
					Name: Currency,
				},
				MusicLanguage:             Language,
				MemberDefaultCountry:      Country,
				OutletsProcessingDuration: 8,
				FreePlanLimit:             1,
				ProductReview:             "Both",
				ContactDetails: entities.ContactDetails{
					ContactPerson: "contact person1",
					Email:         "person@gmail.com",
					NoReplyEmail:  "noreplyemail@gmail.com",
					FeedbackEmail: "feedbackemail@gmail.com",
					SupportEmail:  "supportemail@gmail.com",
				},
				AddressDetails: entities.AddressDetails{
					Address:    "thiruvathira",
					Country:    Country,
					State:      "UP",
					City:       "patna",
					Street:     "torus",
					PostalCode: "PostalCode",
				},
				PaymentGatewayDetails: entities.PaymentGatewayDetails{
					PayoutMinLimit:        1,
					DefaultCurrency:       "payout_currency",
					MaxRemittancePerMonth: 1,
					DefaultPaymentGateway: "paypal",
					PaymentGateways: []entities.PaymentGateways{
						{Gateway: "paypal",
							Email:                 "aathi@gmail.com",
							ClientId:              "1223435",
							ClientSecret:          "aathi@12344",
							Payin:                 true,
							Payout:                true,
							DefaultPayinCurrency:  "USD",
							DefaultPayoutCurrency: "USD"},
					},
				},
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {
				partner.EXPECT().
					IsExists(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(9)

			},
			checkResponse: func(t *testing.T, validationErr map[string]models.ErrorResponse, id string, err error) {
				require.NotNil(t, validationErr)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
				Host:     "",
				UserName: "",
				Password: "",
				DB:       0,
			})
			partner := mock.NewMockPartnerRepoImply(ctrl)
			tc.buildStubs(partner)
			partnerUsecase := NewPartnerUseCases(partner, redisClient)

			validationErr, id, err := partnerUsecase.CreatePartner(context.Background(), tc.partner, "", "", map[string]models.ErrorResponse{})
			tc.checkResponse(t, validationErr, id, err)

		})
	}

}
func TestCreatePartnerDefaultField(t *testing.T) {
	// logger.InitLogger(clientOpt)

	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName, // Service name.
		LogLevel:            logLevel,       // Log level.
		IncludeRequestDump:  false,          // Include request data in logs.
		IncludeResponseDump: true,           // Include response data in logs.
	}
	// Check if the application is in debug mode.

	_ = logger.InitLogger(clientOpt)
	testCases2 := []struct {
		name          string
		partner       entities.Partner
		buildStubs    func(partner *mock.MockPartnerRepoImply)
		checkResponse func(t *testing.T, validationErr map[string]models.ErrorResponse, id string, err error)
	}{
		{
			name: " set default values for empty fields",
			partner: entities.Partner{
				Name:                      "name A",
				URL:                       "https://www.mywebsite.com",
				Logo:                      "https://www.mywebsite.com/logo.png",
				Favicon:                   "https://www.mywebsite.com/favicon.ico",
				Loader:                    "",
				LoginPageLogo:             "",
				BackgroundColor:           "",
				BackgroundImage:           "",
				WebsiteURL:                "https://www.xyz.com",
				Language:                  "",
				BrowserTitle:              "",
				ProfileURL:                "https://www.mywebsite.com/pqt",
				PaymentURL:                "https://www.mywebsite.com/payment",
				LandingPage:               "https://www.mywebsite.com/landing",
				MobileVerifyInterval:      0,
				PayoutTargetCurrency:      "",
				ThemeID:                   0,
				EnableMail:                false,
				LoginType:                 "normal",
				BusinessModel:             "Subscription",
				MemberPayToPartner:        false,
				MemberGracePeriod:         "Monthly",
				ExpiryWarningCount:        0,
				AlbumReviewEmail:          "reviews@mywebsite.com",
				SiteInfo:                  "Welcome to My Website!",
				DefaultPriceCode:          1,
				DefaultPriceCodeCurrency:  entities.Currency{},
				MusicLanguage:             "",
				MemberDefaultCountry:      "",
				OutletsProcessingDuration: 8,
				FreePlanLimit:             0,
				ProductReview:             "Both",
				ContactDetails: entities.ContactDetails{
					ContactPerson: "contact person1",
					Email:         "person@gmail.com",
					NoReplyEmail:  "noreplyemail@gmail.com",
					FeedbackEmail: "feedbackemail@gmail.com",
					SupportEmail:  "supportemail@gmail.com",
				},
				AddressDetails: entities.AddressDetails{
					Address:    "thiruvathira",
					Country:    Country,
					State:      "KL",
					City:       "patna",
					Street:     "torus",
					PostalCode: "PostalCode",
				},
				PaymentGatewayDetails: entities.PaymentGatewayDetails{
					PayoutMinLimit:        0,
					DefaultCurrency:       "",
					MaxRemittancePerMonth: 0,
					DefaultPaymentGateway: "Paypal",
					PaymentGateways: []entities.PaymentGateways{
						{Gateway: "Paypal",
							Email:                 "aathi@gmail.com",
							ClientId:              "1223435",
							ClientSecret:          "aathi@12344",
							Payin:                 true,
							Payout:                true,
							DefaultPayinCurrency:  "USD",
							DefaultPayoutCurrency: "USD"},
					},
				},
			},
			buildStubs: func(mockUserRepo *mock.MockPartnerRepoImply) {
				mockUserRepo.EXPECT().
					IsExists(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(9)

				mockUserRepo.EXPECT().CreatePartner(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

			},
			checkResponse: func(t *testing.T, validationErr map[string]models.ErrorResponse, id string, err error) {
				require.Nil(t, validationErr)
			},
		},
	}

	for i := range testCases2 {
		tc := testCases2[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
				Host:     "",
				UserName: "",
				Password: "",
				DB:       0,
			})
			partner := mock.NewMockPartnerRepoImply(ctrl)
			tc.buildStubs(partner)
			partnerUsecase := NewPartnerUseCases(partner, redisClient)

			validationErr, id, err := partnerUsecase.CreatePartner(context.Background(), tc.partner, "", "", map[string]models.ErrorResponse{})
			tc.checkResponse(t, validationErr, id, err)

		})
	}
}
func TestCreatePartnerMaxLengthErr(t *testing.T) {
	// logger.InitLogger(clientOpt)

	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName, // Service name.
		LogLevel:            logLevel,       // Log level.
		IncludeRequestDump:  false,          // Include request data in logs.
		IncludeResponseDump: true,           // Include response data in logs.
	}
	// Check if the application is in debug mode.

	_ = logger.InitLogger(clientOpt)
	testCases2 := []struct {
		name          string
		partner       entities.Partner
		buildStubs    func(partner *mock.MockPartnerRepoImply)
		checkResponse func(t *testing.T, validationErr map[string]models.ErrorResponse, id string, err error)
	}{
		{
			name: " validation errors for maximum length",
			partner: entities.Partner{
				Name:                 "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce vitae dolor vel urna placerat dignissim. Vivamus dictum ante nec arcu lobortis, eu tempor quam vehicula. Proin vitae quam vitae eros tincidunt auctor ac in risus. Maecenas non purus libero. Suspendisse potenti. Sed sed odio sed justo vulputate iaculis a sed odio. Quisque ut urna ac turpis aliquam bibendum. Duis euismod erat nec elit malesuada, et malesuada nisi rhoncus. Pellentesque sit amet felis eget eros sodales auctor a quis metus. Ut placerat dolor vitae libero vehicula, non faucibus est vehicula. Mauris sit amet nisi ipsum. Integer eleifend tellus eu augue dictum, non ultricies eros tincidunt. Aliquam tincidunt nulla id leo sagittis posuere. Vestibulum pharetra scelerisque ex, in semper odio volutpat vel. Cras lobortis elit vel velit fermentum, a venenatis nisi dapibus",
				URL:                  "website",
				Logo:                 "logo",
				Favicon:              "https://www.website.com/favicon.ico",
				Loader:               "https://www.website.com/loader.ico",
				LoginPageLogo:        "https://www.website.com/pagelogo.ico",
				BackgroundColor:      "#9",
				BackgroundImage:      "https://www.website.com/background.jpg",
				WebsiteURL:           "url",
				Language:             Language,
				BrowserTitle:         "Welcome to My Website",
				ProfileURL:           "url",
				PaymentURL:           "payment",
				LandingPage:          "landing",
				MobileVerifyInterval: 0,
				PayoutTargetCurrency: "USD",
				ThemeID:              0,
				EnableMail:           true,
				LoginType:            "normal",
				BusinessModel:        "business model",
				MemberPayToPartner:   true,
				MemberGracePeriod:    "test plan",
				ExpiryWarningCount:   1000,
				AlbumReviewEmail:     "email",
				SiteInfo:             "Welcome to My Website!",
				DefaultPriceCode:     1,
				DefaultPriceCodeCurrency: entities.Currency{
					Name: Currency,
				},
				MusicLanguage:             Language,
				MemberDefaultCountry:      Country,
				OutletsProcessingDuration: 5,
				FreePlanLimit:             200,
				ProductReview:             "both",
				ContactDetails: entities.ContactDetails{
					ContactPerson: "Lorem ipsum dolor sit amyyuiet, consectetur adipiscing elit. Fusce vitae dolor vel urna placerat dignissim. Vivamus dictum ante nec arcu lobortis, eu tempor quam vehicula. Proin vitae quam vitae eros tincidunt auctor ac in risus. Maecenas non purus libero. Suspendisse potenti. Sed sed odio sed justo vulputate iaculis a sed odio. Quisque ut urna ac turpis aliquam bibendum. Duis euismod erat nec elit malesuada, et malesuada nisi rhoncus. Pellentesque sit amet felis eget eros sodales auctor a quis metus. Ut placerat dolor vitae libero vehicula, non faucibus est vehicula. Mauris sit amet nisi ipsum. Integer eleifend tellus eu augue dictum, non ultricies eros tincidunt. Aliquam tincidunt nulla id leo sagittis posuere. Vestibulum pharetra scelerisque ex, in semper odio volutpat vel. Cras lobortis elit vel velit fermentum, a venenatis nisi dapibus",
					Email:         "email",
					NoReplyEmail:  "email",
					FeedbackEmail: "email",
					SupportEmail:  "email",
				},
				AddressDetails: entities.AddressDetails{
					Address:    "Lorem ipsum dolor sit amuiet, consectetur adipiscing elit. Fusce vitae dolor vel urna placerat dignissim. Vivamus dictum ante nec arcu lobortis, eu tempor quam vehicula. Proin vitae quam vitae eros tincidunt auctor ac in risus. Maecenas non purus libero. Suspendisse potenti. Sed sed odio sed justo vulputate iaculis a sed odio. Quisque ut urna ac turpis aliquam bibendum. Duis euismod erat nec elit malesuada, et malesuada nisi rhoncus. Pellentesque sit amet felis eget eros sodales auctor a quis metus. Ut placerat dolor vitae libero vehicula, non faucibus est vehicula. Mauris sit amet nisi ipsum. Integer eleifend tellus eu augue dictum, non ultricies eros tincidunt. Aliquam tincidunt nulla id leo sagittis posuere. Vestibulum pharetra scelerisque ex, in semper odio volutpat vel. Cras lobortis elit vel velit fermentum, a venenatis nisi dapibus",
					Country:    Country,
					State:      "UP",
					City:       "Lorem ipsum dolor sit amttttet, consectetur adipiscing elit. Fusce vitae dolor vel urna placerat dignissim. Vivamus dictum ante nec arcu lobortis, eu tempor quam vehicula. Proin vitae quam vitae eros tincidunt auctor ac in risus. Maecenas non purus libero. Suspendisse potenti. Sed sed odio sed justo vulputate iaculis a sed odio. Quisque ut urna ac turpis aliquam bibendum. Duis euismod erat nec elit malesuada, et malesuada nisi rhoncus. Pellentesque sit amet felis eget eros sodales auctor a quis metus. Ut placerat dolor vitae libero vehicula, non faucibus est vehicula. Mauris sit amet nisi ipsum. Integer eleifend tellus eu augue dictum, non ultricies eros tincidunt. Aliquam tincidunt nulla id leo sagittis posuere. Vestibulum pharetra scelerisque ex, in semper odio volutpat vel. Cras lobortis elit vel velit fermentum, a venenatis nisi dapibus",
					Street:     "Lorem ipsum dolor sit amgget, consectetur adipiscing elit. Fusce vitae dolor vel urna placerat dignissim. Vivamus dictum ante nec arcu lobortis, eu tempor quam vehicula. Proin vitae quam vitae eros tincidunt auctor ac in risus. Maecenas non purus libero. Suspendisse potenti. Sed sed odio sed justo vulputate iaculis a sed odio. Quisque ut urna ac turpis aliquam bibendum. Duis euismod erat nec elit malesuada, et malesuada nisi rhoncus. Pellentesque sit amet felis eget eros sodales auctor a quis metus. Ut placerat dolor vitae libero vehicula, non faucibus est vehicula. Mauris sit amet nisi ipsum. Integer eleifend tellus eu augue dictum, non ultricies eros tincidunt. Aliquam tincidunt nulla id leo sagittis posuere. Vestibulum pharetra scelerisque ex, in semper odio volutpat vel. Cras lobortis elit vel velit fermentum, a venenatis nisi dapibus",
					PostalCode: "1234567898888",
				},
				PaymentGatewayDetails: entities.PaymentGatewayDetails{
					PayoutMinLimit:        201,
					DefaultCurrency:       "payout_currency",
					MaxRemittancePerMonth: 202,
					DefaultPaymentGateway: "pppp",
					PaymentGateways: []entities.PaymentGateways{
						{Gateway: "Paypal",
							Email:                 "aathi@gmail.com",
							ClientId:              "1223435",
							ClientSecret:          "aathi@12344",
							Payin:                 true,
							Payout:                true,
							DefaultPayinCurrency:  "UUU",
							DefaultPayoutCurrency: "UUU"},
					},
				},
			},
			buildStubs: func(mockUserRepo *mock.MockPartnerRepoImply) {
				mockUserRepo.EXPECT().
					IsExists(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(9)

			},
			checkResponse: func(t *testing.T, validationErr map[string]models.ErrorResponse, id string, err error) {
				require.NotNil(t, validationErr)
			},
		},
	}

	for i := range testCases2 {
		tc := testCases2[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
				Host:     "",
				UserName: "",
				Password: "",
				DB:       0,
			})
			partner := mock.NewMockPartnerRepoImply(ctrl)
			tc.buildStubs(partner)
			partnerUsecase := NewPartnerUseCases(partner, redisClient)

			validationErr, id, err := partnerUsecase.CreatePartner(context.Background(), tc.partner, "", "", map[string]models.ErrorResponse{})
			tc.checkResponse(t, validationErr, id, err)

		})
	}
}
func TestCreatePartnerInvalidFields(t *testing.T) {
	// logger.InitLogger(clientOpt)

	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName, // Service name.
		LogLevel:            logLevel,       // Log level.
		IncludeRequestDump:  false,          // Include request data in logs.
		IncludeResponseDump: true,           // Include response data in logs.
	}
	// Check if the application is in debug mode.

	_ = logger.InitLogger(clientOpt)
	testCases2 := []struct {
		name          string
		partner       entities.Partner
		buildStubs    func(partner *mock.MockPartnerRepoImply)
		checkResponse func(t *testing.T, validationErr map[string]models.ErrorResponse, id string, err error)
	}{
		{
			name: " validation errors for invalid fields",
			partner: entities.Partner{
				Name:                       "KANNAN@123",
				URL:                        "https://www.partccccdccccner.com",
				Logo:                       "https://www.fdfdfdf.com",
				Favicon:                    "",
				Loader:                     "",
				LoginPageLogo:              "",
				BackgroundColor:            "#ffffff",
				BackgroundImage:            "background.com/image",
				WebsiteURL:                 "cccccdccccccccdd.com",
				Language:                   "en",
				BrowserTitle:               "partner xyz",
				ProfileURL:                 "https://www.abcccccdcccccccdd.com",
				PaymentURL:                 "httacccccdccccyment.com",
				LandingPage:                "bcdccccccccccccdd.com",
				MobileVerifyInterval:       2,
				PayoutTargetCurrency:       "USD",
				ThemeID:                    1,
				EnableMail:                 true,
				LoginType:                  "oauth 2.0 - client_credentails",
				BusinessModel:              "subscription",
				MemberPayToPartner:         false,
				MemberGracePeriod:          "",
				ExpiryWarningCount:         5,
				AlbumReviewEmail:           "abccccccdccccccccdd",
				SiteInfo:                   "site info",
				DefaultPriceCode:           0,
				MusicLanguage:              "en",
				MemberDefaultCountry:       "IN",
				OutletsProcessingDuration:  7,
				FreePlanLimit:              1,
				ProductReview:              "admin",
				BusinessModelID:            1,
				PayoutTargetCurrencyID:     2,
				DefaultPriceCodeCurrencyID: 2,
				ContactDetails: entities.ContactDetails{
					ContactPerson: "contact person22",
					Email:         "abcccccdccccccccdd",
					NoReplyEmail:  "abcccccdccccccccdd",
					FeedbackEmail: "abcccccdccccccccc",
					SupportEmail:  "abccccccdcccccd",
				},
				AddressDetails: entities.AddressDetails{
					Address:    "address",
					Street:     "torus",
					Country:    "IN",
					State:      "KL",
					City:       "patna",
					PostalCode: "141414",
				},
				DefaultPriceCodeCurrency: entities.Currency{
					Name: "USD",
				},
				PaymentGatewayDetails: entities.PaymentGatewayDetails{
					PayoutMinLimit:        100,
					DefaultCurrency:       "USD",
					MaxRemittancePerMonth: 1,
					DefaultPaymentGateway: "Paypal",
					PaymentGateways: []entities.PaymentGateways{
						{
							Gateway:               "1f81655a-7486-4750-ae66-879b117b0bf6",
							GatewayId:             "123456789",
							Email:                 "paypal@gmail.com",
							ClientId:              "01596d6c1acfcf205c95fe15e6cf1e3170c34df6e53b6a693a169abf766c1a1c",
							ClientSecret:          "33d9462ab77b99fcbcc0fd4db59627d30d100d658bb906d8a60df9bd3f5a7d06",
							Payin:                 true,
							Payout:                false,
							DefaultPayinCurrency:  "USD",
							DefaultPayoutCurrency: "USD",
						},
						{
							Gateway:               "PayTm",
							GatewayId:             "1f81655a-7486-4750-ae66-879b117b0bf6",
							Email:                 "stripe@gmail.com",
							ClientId:              "01596d6c1acfcf205c95fe15e6cf1e3170c34df6e53b6a693a169abf766c1a1c",
							ClientSecret:          "33d9462ab77b99fcbcc0fd4db59627d30d100d658bb906d8a60df9bd3f5a7d06",
							Payin:                 true,
							Payout:                true,
							DefaultPayinCurrency:  "USD",
							DefaultPayoutCurrency: "USD",
						},
					},
				},
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {
				partner.EXPECT().
					IsExists(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(9)

			},
			checkResponse: func(t *testing.T, validationErr map[string]models.ErrorResponse, id string, err error) {
				require.NotEqual(t, 0, len(validationErr))
			},
		},
	}

	for i := range testCases2 {
		tc := testCases2[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
				Host:     "",
				UserName: "",
				Password: "",
				DB:       0,
			})
			partner := mock.NewMockPartnerRepoImply(ctrl)
			tc.buildStubs(partner)
			partnerUsecase := NewPartnerUseCases(partner, redisClient)

			validationErr, id, err := partnerUsecase.CreatePartner(context.Background(), tc.partner, "", "", map[string]models.ErrorResponse{})
			tc.checkResponse(t, validationErr, id, err)

		})
	}
}
func TestCreatePartnerEmptyFields(t *testing.T) {
	// logger.InitLogger(clientOpt)

	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName, // Service name.
		LogLevel:            logLevel,       // Log level.
		IncludeRequestDump:  false,          // Include request data in logs.
		IncludeResponseDump: true,           // Include response data in logs.
	}
	// Check if the application is in debug mode.

	_ = logger.InitLogger(clientOpt)
	testCases2 := []struct {
		name          string
		partner       entities.Partner
		buildStubs    func(partner *mock.MockPartnerRepoImply)
		checkResponse func(t *testing.T, validationErr map[string]models.ErrorResponse, id string, err error)
	}{
		{
			name: " validation error for empty fields",
			partner: entities.Partner{
				Name:                 "",
				URL:                  "",
				Logo:                 "",
				Favicon:              "",
				Loader:               "",
				LoginPageLogo:        "https://www.website.com/pagelogo.ico",
				BackgroundColor:      "#F0F0F0",
				BackgroundImage:      "https://www.website.com/background.jpg",
				WebsiteURL:           "https://www.pqr.com",
				Language:             "en",
				BrowserTitle:         "Welcome to My Website",
				ProfileURL:           "https://www.website.com/pqr",
				PaymentURL:           "https://www.website.com/payment",
				LandingPage:          "https://www.website.com/landing",
				MobileVerifyInterval: 0,
				PayoutTargetCurrency: "",
				ThemeID:              0,
				EnableMail:           true,
				LoginType:            "",
				BusinessModel:        "",
				MemberPayToPartner:   true,
				MemberGracePeriod:    "",
				ExpiryWarningCount:   1,
				AlbumReviewEmail:     "",
				SiteInfo:             "Welcome to My Website!",
				DefaultPriceCode:     1,
				DefaultPriceCodeCurrency: entities.Currency{
					Name: Currency,
				},
				MusicLanguage:             "en",
				MemberDefaultCountry:      "IN",
				OutletsProcessingDuration: 8,
				FreePlanLimit:             1,
				ProductReview:             "",
				ContactDetails: entities.ContactDetails{
					ContactPerson: "",
					Email:         "",
					NoReplyEmail:  "",
					FeedbackEmail: "",
					SupportEmail:  "",
				},
				AddressDetails: entities.AddressDetails{
					Address:    "",
					Country:    "",
					State:      "",
					City:       "patna",
					Street:     "torus",
					PostalCode: "PostalCode",
				},
				PaymentGatewayDetails: entities.PaymentGatewayDetails{
					PayoutMinLimit:        1,
					DefaultCurrency:       "payout_currency",
					MaxRemittancePerMonth: 1,
					DefaultPaymentGateway: "",
					PaymentGateways: []entities.PaymentGateways{
						{Gateway: "",
							Email:                 "",
							ClientId:              "",
							ClientSecret:          "",
							Payin:                 true,
							Payout:                true,
							DefaultPayinCurrency:  "",
							DefaultPayoutCurrency: ""},
					},
				},
			},
			buildStubs: func(mockUserRepo *mock.MockPartnerRepoImply) {
				mockUserRepo.EXPECT().
					IsExists(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(2)

			},
			checkResponse: func(t *testing.T, validationErr map[string]models.ErrorResponse, id string, err error) {
				require.NotEmpty(t, validationErr)
			},
		},
	}

	for i := range testCases2 {
		tc := testCases2[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
				Host:     "",
				UserName: "",
				Password: "",
				DB:       0,
			})
			partner := mock.NewMockPartnerRepoImply(ctrl)
			tc.buildStubs(partner)
			partnerUsecase := NewPartnerUseCases(partner, redisClient)

			validationErr, id, err := partnerUsecase.CreatePartner(context.Background(), tc.partner, "", "", map[string]models.ErrorResponse{})
			tc.checkResponse(t, validationErr, id, err)

		})
	}
}

func TestCreatePartnerSuccess(t *testing.T) {
	// logger.InitLogger(clientOpt)

	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName, // Service name.
		LogLevel:            logLevel,       // Log level.
		IncludeRequestDump:  false,          // Include request data in logs.
		IncludeResponseDump: true,           // Include response data in logs.
	}
	// Check if the application is in debug mode.

	_ = logger.InitLogger(clientOpt)
	testCases2 := []struct {
		name          string
		partner       entities.Partner
		buildStubs    func(partner *mock.MockPartnerRepoImply)
		checkResponse func(t *testing.T, validationErr map[string]models.ErrorResponse, id string, err error)
	}{
		{
			name: "ok",
			partner: entities.Partner{
				Name:                       "KANNAN1234",
				URL:                        "https://www.partccccdccccner.com",
				Logo:                       "https://www.fdfdfdf.com",
				Favicon:                    "",
				Loader:                     "",
				LoginPageLogo:              "",
				BackgroundColor:            "#ffffff",
				BackgroundImage:            "https://background.com/image",
				WebsiteURL:                 "https://www.abcccccdccccccccdd.com",
				Language:                   "en",
				BrowserTitle:               "partner xyz",
				ProfileURL:                 "https://www.abcccccdcccccccdd.com",
				PaymentURL:                 "httacccccdccccyment.com",
				LandingPage:                "https://www.abcdccccccccccccdd.com",
				MobileVerifyInterval:       2,
				PayoutTargetCurrency:       "USD",
				ThemeID:                    1,
				EnableMail:                 true,
				LoginType:                  "oauth 2.0 - client_credentails",
				BusinessModel:              "subscription",
				MemberPayToPartner:         false,
				MemberGracePeriod:          "",
				ExpiryWarningCount:         5,
				AlbumReviewEmail:           "abccccccdccccccccdd@part.com",
				SiteInfo:                   "site info",
				DefaultPriceCode:           0,
				MusicLanguage:              "en",
				MemberDefaultCountry:       "IN",
				OutletsProcessingDuration:  7,
				FreePlanLimit:              1,
				ProductReview:              "Admin",
				BusinessModelID:            1,
				PayoutTargetCurrencyID:     2,
				DefaultPriceCodeCurrencyID: 2,
				ContactDetails: entities.ContactDetails{
					ContactPerson: "contact person22",
					Email:         "abcccccdccccccccdd@abssss.com",
					NoReplyEmail:  "abcccccdccccccccdd@part.com",
					FeedbackEmail: "abcccccdcccccccccdd@part.com",
					SupportEmail:  "abccccccdcccccdd@part.com",
				},
				AddressDetails: entities.AddressDetails{
					Address:    "address",
					Street:     "torus",
					Country:    "IN",
					State:      "KL",
					City:       "patna",
					PostalCode: "141414",
				},
				DefaultPriceCodeCurrency: entities.Currency{
					Name: "USD",
				},
				PaymentGatewayDetails: entities.PaymentGatewayDetails{
					PayoutMinLimit:        100,
					DefaultCurrency:       "USD",
					MaxRemittancePerMonth: 1,
					DefaultPaymentGateway: "Paypal",
					PaymentGateways: []entities.PaymentGateways{
						{
							Gateway:               "Paypal",
							GatewayId:             "1f81655a-7486-4750-ae66-879b117b0bf6",
							Email:                 "paypal@gmail.com",
							ClientId:              "01596d6c1acfcf205c95fe15e6cf1e3170c34df6e53b6a693a169abf766c1a1c",
							ClientSecret:          "33d9462ab77b99fcbcc0fd4db59627d30d100d658bb906d8a60df9bd3f5a7d06",
							Payin:                 true,
							Payout:                false,
							DefaultPayinCurrency:  "USD",
							DefaultPayoutCurrency: "USD",
						},
						{
							Gateway:               "PayTm",
							GatewayId:             "1f81655a-7486-4750-ae66-879b117b0bf6",
							Email:                 "stripe@gmail.com",
							ClientId:              "01596d6c1acfcf205c95fe15e6cf1e3170c34df6e53b6a693a169abf766c1a1c",
							ClientSecret:          "33d9462ab77b99fcbcc0fd4db59627d30d100d658bb906d8a60df9bd3f5a7d06",
							Payin:                 true,
							Payout:                true,
							DefaultPayinCurrency:  "USD",
							DefaultPayoutCurrency: "USD",
						},
					},
				},
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {
				partner.EXPECT().
					IsExists(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(9)

				partner.EXPECT().
					CreatePartner(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

			},

			checkResponse: func(t *testing.T, validationErr map[string]models.ErrorResponse, id string, err error) {
				require.Equal(t, 0, len(validationErr))
			},
		},
	}

	for i := range testCases2 {
		tc := testCases2[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
				Host:     "",
				UserName: "",
				Password: "",
				DB:       0,
			})
			partner := mock.NewMockPartnerRepoImply(ctrl)
			tc.buildStubs(partner)
			partnerUsecase := NewPartnerUseCases(partner, redisClient)

			validationErr, id, err := partnerUsecase.CreatePartner(context.Background(), tc.partner, "", "", map[string]models.ErrorResponse{})
			tc.checkResponse(t, validationErr, id, err)

		})
	}
}

// func TestGenerateClientIdSecret(t *testing.T) {

// 	clientOpt := &logger.ClientOptions{
// 		Service:             consts.AppName,
// 		LogLevel:            logLevel,
// 		IncludeRequestDump:  false,
// 		IncludeResponseDump: true,
// 	}
// 	_ = logger.InitLogger(clientOpt)
// 	cfg := &entities.EnvConfig{
// 		Encryption: entities.Encryption{
// 			Key: "tuneverse-esrevenuttuneverse-tue",
// 		},
// 	}
// 	testCases := []struct {
// 		name                   string
// 		partnerOauthCredential entities.PartnerOauthCredential
// 		buildStubs             func(partner *mock.MockPartnerRepoImply)
// 		checkResponse          func(t *testing.T, err error)
// 	}{
// 		{
// 			name: "ok",
// 			partnerOauthCredential: entities.PartnerOauthCredential{
// 				PartnerId: "28822555-2467-4022-bae9-c7bf8a0e0bc7",
// 			},
// 			buildStubs: func(partner *mock.MockPartnerRepoImply) {

// 				partner.EXPECT().
// 					GenerateClientIdSecret(gomock.Any(), gomock.Any()).
// 					Times(1)

// 			},
// 			checkResponse: func(t *testing.T, err error) {
// 				require.Nil(t, err)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			partner := mock.NewMockPartnerRepoImply(ctrl)
// 			tc.buildStubs(partner)
// 			partnerUsecase := NewPartnerUseCases(partner, cfg)

// 			err := partnerUsecase.GenerateClientIdSecret(context.Background(), tc.partnerOauthCredential.PartnerId)
// 			tc.checkResponse(t, err)

// 		})
// 	}
// }

func TestGetAllPartners(t *testing.T) {
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            "info",
		IncludeRequestDump:  false,
		IncludeResponseDump: false,
	}
	logger.InitLogger(clientOpt)
	// Create a mock HTTP request
	req, err := http.NewRequest("GET", "http://localhost:8031/api/v1/partners?page=1&limit=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP response recorder
	resp := httptest.NewRecorder()

	// Create a Gin context from the request and response
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	testCases := []struct {
		name          string
		params        entities.Params
		partner       entities.ListAllPartners
		buildStubs    func(partner *mock.MockPartnerRepoImply)
		checkResponse func(t *testing.T, count int64, data []entities.ListAllPartners, metadata models.MetaData, validationErr map[string]models.ErrorResponse, err error)
	}{
		{
			name: "ok",
			partner: entities.ListAllPartners{
				UUID: "26ad89f0-372a-4d42-8d17-04021f654e30",

				ContactDetails: entities.ContactDetails{
					ContactPerson: "contact person1",
					Email:         "person1@gmail.com",
					NoReplyEmail:  "noreplyemail@gmail.com",
					FeedbackEmail: "feedbackemail@gmail.com",
					SupportEmail:  "supportemail@gmail.com",
				},
				Name:                 "Partner 1",
				URL:                  "mywebsite1.com",
				Logo:                 "https://www.mywebsite.com/logo.png",
				Language:             "English",
				BrowserTitle:         "Welcome to My Website1",
				ProfileURL:           "https://www.mywebsite.com/profile",
				PaymentURL:           "https://www.mywebsite.com/payment",
				LandingPage:          "https://www.mywebsite.com/landing",
				MobileVerifyInterval: 1,
				PayoutTargetCurrency: "USD",
				EnableMail:           false,
				BusinessModel:        "buisness model1",
				MemberPayToPartner:   true,

				Theme:     "Example Theme1",
				LoginType: "partner",
				Users:     3,
				Active:    true,
			},
			params: entities.Params{
				Limit:   1,
				Page:    1,
				Status:  "active",
				Order:   "ASC",
				Sort:    "name",
				Country: "IN",
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					GetAllPartners(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

			},
			checkResponse: func(t *testing.T, count int64, data []entities.ListAllPartners, metadata models.MetaData, validationErr map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)
			},
		},

		{
			name: "no-content",
			partner: entities.ListAllPartners{
				UUID: "26ad89f0-372a-4d42-8d17-04021f654e30",

				ContactDetails: entities.ContactDetails{
					ContactPerson: "contact person1",
					Email:         "person1@gmail.com",
					NoReplyEmail:  "noreplyemail@gmail.com",
					FeedbackEmail: "feedbackemail@gmail.com",
					SupportEmail:  "supportemail@gmail.com",
				},
				Name:                 "Partner 1",
				URL:                  "mywebsite1.com",
				Logo:                 "https://www.mywebsite.com/logo.png",
				Language:             "English",
				BrowserTitle:         "Welcome to My Website1",
				ProfileURL:           "https://www.mywebsite.com/profile",
				PaymentURL:           "https://www.mywebsite.com/payment",
				LandingPage:          "https://www.mywebsite.com/landing",
				MobileVerifyInterval: 1,
				PayoutTargetCurrency: "USD",
				EnableMail:           false,
				BusinessModel:        "buisness model1",
				MemberPayToPartner:   true,

				Theme:     "Example Theme1",
				LoginType: "partner",
				Users:     3,
				Active:    true,
			},
			params: entities.Params{
				Limit:   10,
				Page:    4,
				Status:  "active",
				Order:   "ASC",
				Sort:    "name",
				Country: "IN",
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					GetAllPartners(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

			},
			checkResponse: func(t *testing.T, count int64, data []entities.ListAllPartners, metadata models.MetaData, validationErr map[string]models.ErrorResponse, err error) {
				require.Zero(t, count)
			},
		},
		{
			name: "invalid page",
			partner: entities.ListAllPartners{
				UUID: "26ad89f0-372a-4d42-8d17-04021f654e30",

				ContactDetails: entities.ContactDetails{
					ContactPerson: "contact person1",
					Email:         "person1@gmail.com",
					NoReplyEmail:  "noreplyemail@gmail.com",
					FeedbackEmail: "feedbackemail@gmail.com",
					SupportEmail:  "supportemail@gmail.com",
				},
				Name:                 "Partner 1",
				URL:                  "mywebsite1.com",
				Logo:                 "https://www.mywebsite.com/logo.png",
				Language:             "English",
				BrowserTitle:         "Welcome to My Website1",
				ProfileURL:           "https://www.mywebsite.com/profile",
				PaymentURL:           "https://www.mywebsite.com/payment",
				LandingPage:          "https://www.mywebsite.com/landing",
				MobileVerifyInterval: 1,
				PayoutTargetCurrency: "USD",
				EnableMail:           false,
				BusinessModel:        "buisness model1",
				MemberPayToPartner:   true,

				Theme:     "Example Theme1",
				LoginType: "partner",
				Users:     3,
				Active:    true,
			},
			params: entities.Params{
				Limit:   10,
				Page:    -4,
				Status:  "active",
				Order:   "ASC",
				Sort:    "name",
				Country: "IN",
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					GetAllPartners(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkResponse: func(t *testing.T, count int64, data []entities.ListAllPartners, metadata models.MetaData, validationErr map[string]models.ErrorResponse, err error) {
				require.NotEmpty(t, validationErr)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
				Host:     "",
				UserName: "",
				Password: "",
				DB:       0,
			})
			partner := mock.NewMockPartnerRepoImply(ctrl)
			tc.buildStubs(partner)
			partnerUsecase := NewPartnerUseCases(partner, redisClient)

			data, metadata, validationErr, err := partnerUsecase.GetAllPartners(c, tc.params, "", "", map[string]models.ErrorResponse{})
			tc.checkResponse(t, metadata.Total, data, metadata, validationErr, err)

		})
	}

}

func TestDeletePartner(t *testing.T) {
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            "info",
		IncludeRequestDump:  false,
		IncludeResponseDump: false,
	}
	logger.InitLogger(clientOpt)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockPartnerRepoImply(ctrl)

	useCases := &PartnerUseCases{
		repo: mockRepo,
	}

	mockCtx := &gin.Context{}

	testCases := []struct {
		partnerID                string
		expectedRowsAffected     int64
		expectedValidationErrors map[string][]string
		expectedError            error
		description              string
	}{
		{
			partnerID:     "614608f2-6538-4733-aded-96f902007251",
			expectedError: nil,
			description:   "success",
		},
		{
			partnerID:     "614608f2-6538-4733-aded-96f902007251",
			expectedError: errors.New("internal server error"),
			description:   "internal server error",
		},
	}

	for _, testCase := range testCases {

		if testCase.expectedValidationErrors == nil {
			mockRepo.EXPECT().DeletePartner(gomock.Any(), testCase.partnerID).Times(1).
				Return(testCase.expectedError)

		}
		err := useCases.DeletePartner(mockCtx, testCase.partnerID, "", "", map[string]models.ErrorResponse{})
		assert.Equal(t, testCase.expectedError, err)

		if testCase.expectedError != nil {
			assert.EqualError(t, err, testCase.expectedError.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestUpdateTermsAndConditions(t *testing.T) {
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            logLevel,
		IncludeRequestDump:  false,
		IncludeResponseDump: true,
	}
	_ = logger.InitLogger(clientOpt)

	testCases := []struct {
		name               string
		partnerId          uuid.UUID
		termsAndConditions entities.UpdateTermsAndConditions
		buildStubs         func(partner *mock.MockPartnerRepoImply)
		checkResponse      func(t *testing.T, validationErr map[string]models.ErrorResponse, err error)
	}{
		{
			name:      "ok",
			partnerId: uuid.New(),
			termsAndConditions: map[string]interface{}{
				"name":        "condition1",
				"description": "condition description1",
				"language":    "en",
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					UpdateTermsAndConditions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(t *testing.T, validationErr map[string]models.ErrorResponse, err error) {
				require.Empty(t, validationErr)
				require.Nil(t, err)
			},
		},
		{
			name:      "validation_error",
			partnerId: uuid.New(),
			termsAndConditions: map[string]interface{}{
				"name":        "",
				"description": " ",
				"language":    "",
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					UpdateTermsAndConditions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(t *testing.T, validationErr map[string]models.ErrorResponse, err error) {
				require.NotNil(t, validationErr)
			},
		},
		{
			name:      "max_length_error",
			partnerId: uuid.New(),
			termsAndConditions: map[string]interface{}{
				"name":        "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce vitae dolor vel urna placerat dignissim. Vivamus dictum ante nec arcu lobortis, eu tempor quam vehicula. Proin vitae quam vitae eros tincidunt auctor ac in risus. Maecenas non purus libero. Suspendisse potenti. Sed sed odio sed justo vulputate iaculis a sed odio. Quisque ut urna ac turpis aliquam bibendum. Duis euismod erat nec elit malesuada, et malesuada nisi rhoncus. Pellentesque sit amet felis eget eros sodales auctor a quis metus. Ut placerat dolor vitae libero vehicula, non faucibus est vehicula. Mauris sit amet nisi ipsum. Integer eleifend tellus eu augue dictum, non ultricies eros tincidunt. Aliquam tincidunt nulla id leo sagittis posuere. Vestibulum pharetra scelerisque ex, in semper odio volutpat vel. Cras lobortis elit vel velit fermentum, a venenatis nisi dapibus",
				"description": "condition description1",
				"language":    "en",
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					UpdateTermsAndConditions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(t *testing.T, validationErr map[string]models.ErrorResponse, err error) {
				require.NotNil(t, validationErr)
			},
		},
		{
			name:      "internal_server_error",
			partnerId: uuid.New(),
			termsAndConditions: map[string]interface{}{
				"name":        "condition1",
				"description": "condition description1",
				"language":    "en",
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					UpdateTermsAndConditions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, validationErr map[string]models.ErrorResponse, err error) {
				require.NotNil(t, err)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
				Host:     "",
				UserName: "",
				Password: "",
				DB:       0,
			})
			partner := mock.NewMockPartnerRepoImply(ctrl)
			tc.buildStubs(partner)
			partnerUsecase := NewPartnerUseCases(partner, redisClient)

			validationErr, err := partnerUsecase.UpdateTermsAndConditions(context.Background(), "f45d6151-04c1-434d-8bc9-bd5d0bedbcb9", tc.partnerId, tc.termsAndConditions, "", "", map[string]models.ErrorResponse{})
			tc.checkResponse(t, validationErr, err)

		})
	}
}

func TestGetPartnerPaymentGateways(t *testing.T) {
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            logLevel,
		IncludeRequestDump:  false,
		IncludeResponseDump: true,
	}
	_ = logger.InitLogger(clientOpt)

	// Create a mock HTTP request
	req, err := http.NewRequest("GET", "http://localhost:8031/api/v1/partners/61c1eb1c-4be4-44e2-80a9-6662cfed6dda/payment-gateways", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP response recorder
	resp := httptest.NewRecorder()

	// Create a Gin context from the request and response
	c, _ := gin.CreateTestContext(resp)
	c.Request = req
	testCases := []struct {
		name          string
		partnerId     uuid.UUID
		buildStubs    func(partner *mock.MockPartnerRepoImply)
		checkResponse func(t *testing.T, data entities.GetPaymentGateways, err error)
	}{
		{
			name:      "ok",
			partnerId: uuid.New(),
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					GetPartnerPaymentGateways(gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(t *testing.T, data entities.GetPaymentGateways, err error) {
				require.Nil(t, err)
			},
		},
		{
			name:      "internal_server_error",
			partnerId: uuid.New(),
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					GetPartnerPaymentGateways(gomock.Any(), gomock.Any()).
					Times(1).Return(entities.GetPaymentGateways{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, data entities.GetPaymentGateways, err error) {
				require.NotNil(t, err)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
				Host:     "",
				UserName: "",
				Password: "",
				DB:       0,
			})
			partner := mock.NewMockPartnerRepoImply(ctrl)
			tc.buildStubs(partner)
			partnerUsecase := NewPartnerUseCases(partner, redisClient)

			data, err := partnerUsecase.GetPartnerPaymentGateways(c, "f45d6151-04c1-434d-8bc9-bd5d0bedbcb9", "", "", map[string]models.ErrorResponse{})
			tc.checkResponse(t, data, err)

		})
	}
}

func TestGetPartnerStores(t *testing.T) {
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            logLevel,
		IncludeRequestDump:  false,
		IncludeResponseDump: true,
	}
	_ = logger.InitLogger(clientOpt)

	testCases := []struct {
		name          string
		partnerId     string
		buildStubs    func(partner *mock.MockPartnerRepoImply)
		checkResponse func(t *testing.T, data entities.GetPartnerStores, err error)
	}{
		{
			name:      "ok",
			partnerId: "f45d6151-04c1-434d-8bc9-bd5d0bedbcb9",
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					GetPartnerStores(gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(t *testing.T, data entities.GetPartnerStores, err error) {
				require.Nil(t, err)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
				Host:     "",
				UserName: "",
				Password: "",
				DB:       0,
			})
			partner := mock.NewMockPartnerRepoImply(ctrl)
			tc.buildStubs(partner)
			partnerUsecase := NewPartnerUseCases(partner, redisClient)

			data, err := partnerUsecase.GetPartnerStores(context.Background(), tc.partnerId, "", "", map[string]models.ErrorResponse{})
			tc.checkResponse(t, data, err)

		})
	}

}

func TestIsPartnerExists(t *testing.T) {
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            logLevel,
		IncludeRequestDump:  false,
		IncludeResponseDump: true,
	}
	_ = logger.InitLogger(clientOpt)
	testCases := []struct {
		name          string
		partnerId     string
		buildStubs    func(partner *mock.MockPartnerRepoImply)
		checkResponse func(t *testing.T, data bool, err error)
	}{
		{
			name:      "ok",
			partnerId: "f45d6151-04c1-434d-8bc9-bd5d0bedbcb9",
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					IsPartnerExists(gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(t *testing.T, data bool, err error) {
				require.Nil(t, err)
			},
		},
		{
			name:      "invalid id",
			partnerId: "f45d6151-04c1-434d-8bc9",
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					IsPartnerExists(gomock.Any(), gomock.Any()).
					Times(0).Return(false, nil)
			},
			checkResponse: func(t *testing.T, data bool, err error) {
				require.Equal(t, data, false)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
				Host:     "",
				UserName: "",
				Password: "",
				DB:       0,
			})
			partner := mock.NewMockPartnerRepoImply(ctrl)
			tc.buildStubs(partner)
			partnerUsecase := NewPartnerUseCases(partner, redisClient)

			data, err := partnerUsecase.IsPartnerExists(context.Background(), tc.partnerId, "", "", map[string]models.ErrorResponse{})
			tc.checkResponse(t, data, err)

		})
	}

}

func TestCreatePartnerStores(t *testing.T) {

	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            logLevel,
		IncludeRequestDump:  false,
		IncludeResponseDump: true,
	}
	_ = logger.InitLogger(clientOpt)

	testCases := []struct {
		name          string
		partnerId     string
		stores        entities.PartnerStores
		buildStubs    func(partner *mock.MockPartnerRepoImply)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name:      "ok",
			partnerId: "a2a00db1-3c9a-4b75-9907-6ef55b63d379",
			stores: entities.PartnerStores{
				Stores: []string{"a2a00db1-3c9a-4b75-9907-6ef55b63d379",
					"a2a00db1-3c9a-4b75-9907-6ef55b63d379"},
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					CreatePartnerStores(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Nil(t, err)
			},
		},
		{
			name:      "no data",
			partnerId: "",
			stores: entities.PartnerStores{
				Stores: []string{},
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					CreatePartnerStores(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Nil(t, err)
			},
		},
		{
			name:      "server error",
			partnerId: "a2a00db1-3c9a-4b75-9907-6ef55b63d379",
			stores: entities.PartnerStores{
				Stores: []string{"a2a00db1-3c9a-4b75-9907-6ef55b63d379",
					"a2a00db1-3c9a-4b75-9907-6ef55b63d379"},
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					CreatePartnerStores(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NotNil(t, err)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
				Host:     "",
				UserName: "",
				Password: "",
				DB:       0,
			})
			partner := mock.NewMockPartnerRepoImply(ctrl)
			tc.buildStubs(partner)
			partnerUsecase := NewPartnerUseCases(partner, redisClient)

			_, _, err := partnerUsecase.CreatePartnerStores(context.Background(), tc.stores, tc.partnerId, "", "", map[string]models.ErrorResponse{})
			tc.checkResponse(t, err)

		})
	}
}

func TestUpdatePartnerStatus(t *testing.T) {

	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            logLevel,
		IncludeRequestDump:  false,
		IncludeResponseDump: true,
	}
	_ = logger.InitLogger(clientOpt)
	// Create a mock HTTP request
	req, err := http.NewRequest("GET", "http://localhost:8031/api/v1/partners", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP response recorder
	resp := httptest.NewRecorder()

	// Create a Gin context from the request and response
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	testCases := []struct {
		name          string
		partnerId     string
		status        entities.UpdatePartnerStatus
		buildStubs    func(partner *mock.MockPartnerRepoImply)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name:      "ok",
			partnerId: "a2a00db1-3c9a-4b75-9907-6ef55b63d379",
			status: entities.UpdatePartnerStatus{
				Active: true,
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					UpdatePartnerStatus(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Nil(t, err)
			},
		},
		{
			name:      "server error",
			partnerId: "a2a00db1-3c9a-4b75-9907-6ef55b63d379",
			status: entities.UpdatePartnerStatus{
				Active: true,
			},
			buildStubs: func(partner *mock.MockPartnerRepoImply) {

				partner.EXPECT().
					UpdatePartnerStatus(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NotNil(t, err)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisClient, _ := cacheConf.New(&cacheConf.RedisCacheOptions{
				Host:     "",
				UserName: "",
				Password: "",
				DB:       0,
			})
			partner := mock.NewMockPartnerRepoImply(ctrl)
			tc.buildStubs(partner)
			partnerUsecase := NewPartnerUseCases(partner, redisClient)

			err := partnerUsecase.UpdatePartnerStatus(c, tc.partnerId, tc.status, "", "", map[string]models.ErrorResponse{})
			tc.checkResponse(t, err)

		})
	}
}
