package usecases

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"utility/internal/consts"
	"utility/internal/entities"
	mockdb "utility/internal/repo/mock"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils"
)

// TestGetPaymentGatewayByID to check GetPaymentGatewayByID function
func TestGetPaymentGatewayByID(t *testing.T) {
	// Create a sample member ID.

	response := entities.PaymentGatewayName{
		Name: "Paypal",
	}
	validation := entities.Validation{
		ID:       "47db82c7-452d-4b83-9ec4-90400c9c1eaf",
		Endpoint: "gateway",
		Method:   "get",
	}
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName, // Service name.
		LogLevel:            "info",         // Log level.
		IncludeRequestDump:  false,          // Include request data in logs.
		IncludeResponseDump: false,          // Include response data in logs.
	}
	logger.InitLogger(clientOpt)

	// Define test cases.
	testCases := []struct {
		name string
		// paymentID     int
		validation    entities.Validation
		response      entities.PaymentGatewayName
		buildStubs    func(mockPaymentRepo *mockdb.MockPaymentGatewayRepoImply)
		checkResponse func(t *testing.T, response entities.PaymentGatewayName, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name:       "Valid payment id",
			validation: validation,
			response:   response,
			buildStubs: func(mockPaymentRepo *mockdb.MockPaymentGatewayRepoImply) {

				// Expect the GetPaymentGatewayByID function to be called with the given payment ID.
				mockPaymentRepo.EXPECT().GetPaymentGatewayByID(gomock.Any(), validation.ID).
					Times(1).Return(response, nil)
			},
			checkResponse: func(t *testing.T, response entities.PaymentGatewayName, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)
				require.Empty(t, fieldsMap)
				require.NotNil(t, response)
			},
		},
		{
			name:       "Invalid payment",
			validation: validation,
			response:   response,
			buildStubs: func(mockPaymentRepo *mockdb.MockPaymentGatewayRepoImply) {

				// Expect the GetPaymentGatewayByID function to be called with the given payment ID.
				mockPaymentRepo.EXPECT().GetPaymentGatewayByID(gomock.Any(), validation.ID).
					Times(1).Return(response, errors.New("Error occured"))
			},
			checkResponse: func(t *testing.T, response entities.PaymentGatewayName, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, err)
				require.Nil(t, fieldsMap)
				require.Equal(t, entities.PaymentGatewayName{}, response)
			},
		},
		{
			name: "Wrong endpoint in errormap",
			validation: entities.Validation{
				ID:       "47db82c7-452d-4b83-9ec4-90400c9c1eaf",
				Endpoint: "",
				Method:   "get",
			},
			response: response,
			buildStubs: func(mockPaymentRepo *mockdb.MockPaymentGatewayRepoImply) {

				// Expect the GetPaymentGatewayByID function to be called with the given payment ID.
				mockPaymentRepo.EXPECT().GetPaymentGatewayByID(gomock.Any(), validation.ID).
					Times(1).Return(response, errors.New("Error occured"))
			},
			checkResponse: func(t *testing.T, response entities.PaymentGatewayName, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, err)
				require.Nil(t, fieldsMap)
			},
		},
	}
	// Convert context.Background() to *gin.Context for testing or specific use cases.
	ginCtx := createTestGinContext()

	for i := range testCases {

		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPaymentRepo := mockdb.NewMockPaymentGatewayRepoImply(ctrl)
			tc.buildStubs(mockPaymentRepo)
			errMap := make(map[string]models.ErrorResponse)

			paymentUseCase := NewPaymentGatewayUseCases(mockPaymentRepo)
			response, fieldsMap, err := paymentUseCase.GetPaymentGatewayByID(ginCtx, tc.validation, errMap)

			tc.checkResponse(t, response, fieldsMap, err)
		})
	}
}

// createTestGinContext creates a dummy context
func createTestGinContext() *gin.Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("headerKey", "headerValue")

	return c
}

func TestGetAllPaymentGateway(t *testing.T) {

	validation := entities.Validation{
		Endpoint: "gateway",
		Method:   "get",
	}

	var payment []entities.PaymentGateway
	for i := 0; i < 3; i++ {
		arg, err := utils.RandomString(4)
		require.Nil(t, err)
		payment = append(payment, entities.PaymentGateway{
			ID:   "02827988-dade-45ba-9312-83caf5541f46",
			Name: arg,
		})
	}

	paginationInfo := entities.Pagination{
		Page:  10,
		Limit: 10,
	}

	params := entities.Params{
		Name:  "Paypal",
		Sort:  "asc",
		Order: "name",
	}

	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName, // Service name.
		LogLevel:            "info",         // Log level.
		IncludeRequestDump:  false,          // Include request data in logs.
		IncludeResponseDump: false,          // Include response data in logs.
	}
	logger.InitLogger(clientOpt)

	// Define test cases.
	testCases := []struct {
		name          string
		validation    entities.Validation
		buildStubs    func(mockPaymentRepo *mockdb.MockPaymentGatewayRepoImply)
		checkResponse func(t *testing.T, response *entities.Response, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name:       "Valid payment",
			validation: validation,
			buildStubs: func(mockPaymentRepo *mockdb.MockPaymentGatewayRepoImply) {

				// Expect the GetPaymentGatewayByID function to be called with the given payment ID.
				mockPaymentRepo.EXPECT().GetAllPaymentGateway(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(payment, int64(len(payment)), nil)
			},
			checkResponse: func(t *testing.T, response *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)
				require.Empty(t, fieldsMap)
				require.NotNil(t, response)
			},
		},
		{
			name:       "Internal server error",
			validation: validation,
			buildStubs: func(mockPaymentRepo *mockdb.MockPaymentGatewayRepoImply) {

				// Expect the GetPaymentGatewayByID function to be called with the given payment ID.
				mockPaymentRepo.EXPECT().GetAllPaymentGateway(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, int64(0), errors.New("error occured"))
			},
			checkResponse: func(t *testing.T, response *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, err)
				require.Empty(t, fieldsMap)
				require.Nil(t, response)
			},
		},
	}
	// Convert context.Background() to *gin.Context for testing or specific use cases.
	ginCtx := createTestGinContext()

	for i := range testCases {

		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPaymentRepo := mockdb.NewMockPaymentGatewayRepoImply(ctrl)
			tc.buildStubs(mockPaymentRepo)
			errMap := make(map[string]models.ErrorResponse)

			paymentUseCase := NewPaymentGatewayUseCases(mockPaymentRepo)
			response, fieldsMap, err := paymentUseCase.GetAllPaymentGateway(ginCtx, params, paginationInfo, tc.validation, errMap)

			tc.checkResponse(t, response, fieldsMap, err)
		})
	}
}
