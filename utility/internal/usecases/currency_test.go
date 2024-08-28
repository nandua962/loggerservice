// // nolint
package usecases

import (
	"context"
	"errors"
	"testing"
	"utility/internal/entities"
	"utility/internal/repo/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils"
)

// TestGetCurrencies to check GetCurrencies function for test
func TestGetCurrencies(t *testing.T) {

	validation := entities.Validation{
		Endpoint: "currencies",
		Method:   "get",
	}
	paginationInfo := entities.Pagination{
		Page:  10,
		Limit: 10,
	}
	params := entities.Params{
		Name:  "John",
		Sort:  "asc",
		Order: "name",
	}

	var currencies []entities.Currency
	for i := 0; i < 5; i++ {
		arg, err := utils.RandomString(4)
		require.Nil(t, err)
		currencies = append(currencies, entities.Currency{
			ID:     int64(i + 1),
			Name:   arg,
			ISO:    "USD",
			Symbol: "$",
		})
	}

	testCases := []struct {
		name          string
		pagination    entities.Pagination
		params        entities.Params
		buildStubs    func(store *mock.MockCurrencyRepoImply)
		checkresponse func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name:   "success",
			params: params,
			buildStubs: func(store *mock.MockCurrencyRepoImply) {
				store.EXPECT().
					GetCurrencies(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(currencies, int64(len(currencies)), nil)

			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Equal(t, resp.MetaData.Total, int64(len(currencies)))
				require.Nil(t, err)
			},
		},
		{
			name:   "internal server error",
			params: params,
			buildStubs: func(store *mock.MockCurrencyRepoImply) {
				store.EXPECT().
					GetCurrencies(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, int64(0), errors.New("error occured"))

			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Error(t, err)
			},
		},
		{
			name:   "data not found",
			params: params,
			buildStubs: func(store *mock.MockCurrencyRepoImply) {
				store.EXPECT().
					GetCurrencies(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, int64(0), nil)

			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)
				require.Nil(t, resp.MetaData)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			errMap := make(map[string]models.ErrorResponse)

			store := mock.NewMockCurrencyRepoImply(ctrl)
			tc.buildStubs(store)
			currencyUsecase := NewCurrencyUseCases(store)
			resp, fieldsMap, err := currencyUsecase.GetCurrencies(context.Background(), params, paginationInfo, validation, errMap)
			tc.checkresponse(t, resp, fieldsMap, err)
		})
	}
}

// TestGetCurrencyByID to check GetCurrencyByID function for test
func TestGetCurrencyByID(t *testing.T) {

	validation := entities.Validation{
		Endpoint: "currencies",
		Method:   "get",
		ID:       "93",
	}

	response := entities.GeographicInfo{
		ID:   93,
		Name: "Albanian Lek",
		ISO:  "ALL",
	}

	testCases := []struct {
		name          string
		validation    entities.Validation
		buildStubs    func(store *mock.MockCurrencyRepoImply)
		checkresponse func(t *testing.T, resp entities.GeographicInfo, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name:       "success",
			validation: validation,
			buildStubs: func(store *mock.MockCurrencyRepoImply) {
				store.EXPECT().
					GetCurrencyByID(gomock.Any(), validation.ID).
					Times(1).
					Return(response, nil)
			},
			checkresponse: func(t *testing.T, resp entities.GeographicInfo, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)

			},
		},
		{
			name:       "Failed while returning error from GetCurrencyByID",
			validation: validation,
			buildStubs: func(store *mock.MockCurrencyRepoImply) {
				store.EXPECT().
					GetCurrencyByID(gomock.Any(), validation.ID).
					Times(1).
					Return(entities.GeographicInfo{}, errors.New("error occured"))
			},
			checkresponse: func(t *testing.T, resp entities.GeographicInfo, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, err)
				require.Equal(t, entities.GeographicInfo{}, resp)

			},
		},
		{
			name: "Failed due to invalid id",
			validation: entities.Validation{
				Endpoint: "currencies",
				Method:   "get",
				ID:       "id",
			},
			buildStubs: func(store *mock.MockCurrencyRepoImply) {
				store.EXPECT().
					GetCurrencyByID(gomock.Any(), validation.ID).
					Times(0)
			},
			checkresponse: func(t *testing.T, resp entities.GeographicInfo, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)
				require.NotNil(t, fieldsMap)

			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			errMap := make(map[string]models.ErrorResponse)

			store := mock.NewMockCurrencyRepoImply(ctrl)
			tc.buildStubs(store)
			currencyUsecase := NewCurrencyUseCases(store)
			resp, fieldsMap, err := currencyUsecase.GetCurrencyByID(context.Background(), tc.validation, errMap)
			tc.checkresponse(t, resp, fieldsMap, err)
		})
	}
}

// TestGetCurrencyByISO to check GetCurrencyByISO function for test
func TestGetCurrencyByISO(t *testing.T) {

	validation := entities.Validation{
		Endpoint: "currencies",
		Method:   "get",
		ID:       "all",
	}

	response := entities.Currency{
		ID:   93,
		Name: "Albanian Lek",
	}

	testCases := []struct {
		name          string
		validation    entities.Validation
		buildStubs    func(store *mock.MockCurrencyRepoImply)
		checkresponse func(t *testing.T, resp entities.Currency, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name:       "success",
			validation: validation,
			buildStubs: func(store *mock.MockCurrencyRepoImply) {
				store.EXPECT().
					GetCurrencyByISO(gomock.Any(), validation.ID).
					Times(1).
					Return(response, nil)
			},
			checkresponse: func(t *testing.T, resp entities.Currency, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)

			},
		},
		{
			name:       "Failed while returning error from GetCurrencyByID",
			validation: validation,
			buildStubs: func(store *mock.MockCurrencyRepoImply) {
				store.EXPECT().
					GetCurrencyByISO(gomock.Any(), validation.ID).
					Times(1).
					Return(entities.Currency{}, errors.New("error occured"))
			},
			checkresponse: func(t *testing.T, resp entities.Currency, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, err)
				require.Equal(t, entities.Currency{}, resp)

			},
		},
		{
			name: "Failed due to invalid id",
			validation: entities.Validation{
				Endpoint: "currencies",
				Method:   "get",
				ID:       "12",
			},
			buildStubs: func(store *mock.MockCurrencyRepoImply) {
				store.EXPECT().
					GetCurrencyByISO(gomock.Any(), validation.ID).
					Times(0)
			},
			checkresponse: func(t *testing.T, resp entities.Currency, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)

			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			errMap := make(map[string]models.ErrorResponse)

			store := mock.NewMockCurrencyRepoImply(ctrl)
			tc.buildStubs(store)
			currencyUsecase := NewCurrencyUseCases(store)
			resp, fieldsMap, err := currencyUsecase.GetCurrencyByISO(context.Background(), tc.validation, errMap)
			tc.checkresponse(t, resp, fieldsMap, err)
		})
	}
}
