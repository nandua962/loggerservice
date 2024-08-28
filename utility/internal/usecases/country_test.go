package usecases

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"utility/internal/entities"
	"utility/internal/repo/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils"
)

// TestGetCountries to check GetCountries function
func TestGetCountries(t *testing.T) {

	validation := entities.Validation{
		Endpoint: "countries",
		Method:   "get",
	}
	paginationInfo := entities.Pagination{
		//nolint
		Page:  10,
		Limit: 10,
	}
	params := entities.Params{
		Sort:  "asc",
		Order: "name",
		Iso:   "an",
	}
	validation.ID = params.Iso
	var (
		countries []*entities.GeographicInfo
		n         = 5
	)
	for i := 0; i < n; i++ {
		arg, err := utils.RandomString(4)
		require.Nil(t, err)
		countries = append(countries, &entities.GeographicInfo{
			ID:   1,
			Name: arg,
			ISO:  "USD",
		})
	}

	testCases := []struct {
		name          string
		funcArguments struct {
			params     entities.Params
			validation entities.Validation
		}
		buildStubs    func(store *mock.MockCountryRepoImply)
		checkresponse func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name: "success",
			funcArguments: struct {
				params     entities.Params
				validation entities.Validation
			}{params, validation},
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					GetCountries(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(countries, int64(len(countries)), nil)

			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Equal(t, resp.MetaData.Total, int64(len(countries)))
				require.Nil(t, err)
			},
		},
		{
			name: "internal server error",
			funcArguments: struct {
				params     entities.Params
				validation entities.Validation
			}{params, validation},
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					GetCountries(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, int64(0), errors.New("Error occured"))

			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, err)
			},
		},
		{
			name: "no data found",
			funcArguments: struct {
				params     entities.Params
				validation entities.Validation
			}{params, validation},
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					GetCountries(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, int64(0), nil)

			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)
				require.Nil(t, resp.MetaData)
			},
		},
		{
			name: "invalid iso",
			funcArguments: struct {
				params     entities.Params
				validation entities.Validation
			}{entities.Params{
				Sort:  "asc",
				Order: "name",
				Iso:   "angs",
			}, validation},
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					GetCountries(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, resp)
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
			store := mock.NewMockCountryRepoImply(ctrl)
			tc.buildStubs(store)
			countryUsecase := NewCountryUseCases(store)
			resp, fieldsMap, err := countryUsecase.GetCountries(context.Background(), tc.funcArguments.params, paginationInfo, tc.funcArguments.validation, errMap)
			tc.checkresponse(t, resp, fieldsMap, err)
		})
	}
}

// TestGetStatesOfState to check GetStatesOfState function
func TestGetStatesOfState(t *testing.T) {

	validation := entities.Validation{
		Endpoint: "states",
		Method:   "get",
		ID:       "194",
	}
	var (
		states []*entities.GeographicInfo
		n      = 5
	)
	id := int64(194)
	for i := 0; i < n; i++ {
		arg, err := utils.RandomString(4)
		require.Nil(t, err)
		states = append(states, &entities.GeographicInfo{
			ID:   1,
			Name: arg,
			ISO:  "CL",
		})
	}

	paginationInfo := entities.Pagination{
		Page:  10,
		Limit: 10,
	}
	params := entities.Params{
		Sort:  "asc",
		Order: "name",
	}
	testCases := []struct {
		name          string
		funcArguments struct {
			params     entities.Params
			validation entities.Validation
		}
		buildStubs    func(store *mock.MockCountryRepoImply)
		checkresponse func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name: "success",
			funcArguments: struct {
				params     entities.Params
				validation entities.Validation
			}{params, validation},
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					GetStatesOfCountry(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), id, gomock.Any()).
					Times(1).
					Return(states, int64(len(states)), nil)

			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Equal(t, resp.MetaData.Total, int64(len(states)))
				require.Nil(t, err)
			},
		},
		{
			name: "internal server error",
			funcArguments: struct {
				params     entities.Params
				validation entities.Validation
			}{params, validation},
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					GetStatesOfCountry(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), id, gomock.Any()).
					Times(1).
					Return(nil, int64(0), errors.New("error occured"))

			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "no data found",
			funcArguments: struct {
				params     entities.Params
				validation entities.Validation
			}{params, validation},
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					GetStatesOfCountry(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), id, gomock.Any()).
					Times(1).
					Return(nil, int64(0), nil)

			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Empty(t, resp.MetaData)
				require.Nil(t, err)
			},
		},

		{
			name: "Failed due to negative Id",
			funcArguments: struct {
				params     entities.Params
				validation entities.Validation
			}{params, entities.Validation{
				Endpoint: "states",
				Method:   "get",
				ID:       "-19",
			}},
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					GetStatesOfCountry(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), id, gomock.Any()).
					Times(0)

			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, fieldsMap)
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

			store := mock.NewMockCountryRepoImply(ctrl)
			tc.buildStubs(store)
			StateUsecase := NewCountryUseCases(store)
			resp, fieldsMap, err := StateUsecase.GetStatesOfCountry(context.Background(), tc.funcArguments.params, paginationInfo, tc.funcArguments.validation, errMap)
			tc.checkresponse(t, resp, fieldsMap, err)

		})
	}
}

// TestCheckStateExists to check CheckStateExists function
func TestCheckStateExists(t *testing.T) {

	validations := entities.Validation{
		Endpoint: "states",
		Method:   "head",
	}

	testCases := []struct {
		name          string
		funcArguments struct {
			iso         string
			countryCode string
			validation  entities.Validation
		}
		buildStubs    func(store *mock.MockCountryRepoImply)
		checkresponse func(t *testing.T, isExists bool, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name: "success",
			funcArguments: struct {
				iso         string
				countryCode string
				validation  entities.Validation
			}{"kl", "in", validations},
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					CheckStateExists(gomock.Any(), "kl", "in").
					Times(1).
					Return(true, nil)
			},
			checkresponse: func(t *testing.T, isExists bool, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Equal(t, true, isExists)
				require.Nil(t, err)
			},
		},
		{
			name: "invalid code length",
			funcArguments: struct {
				iso         string
				countryCode string
				validation  entities.Validation
			}{"kl", "inr", validations},
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					CheckStateExists(gomock.Any(), "kl", "inr").
					Times(0)
			},
			checkresponse: func(t *testing.T, isExists bool, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, fieldsMap)
			},
		},
		{
			name: "If state is not exists",
			funcArguments: struct {
				iso         string
				countryCode string
				validation  entities.Validation
			}{"kl", "in", validations},
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					CheckStateExists(gomock.Any(), "kl", "in").
					Times(1).
					Return(false, nil)
			},
			checkresponse: func(t *testing.T, isExists bool, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Equal(t, false, isExists)
				require.Nil(t, err)
			},
		},
		{
			name: "Internal server error",
			funcArguments: struct {
				iso         string
				countryCode string
				validation  entities.Validation
			}{"kl", "in", validations},
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					CheckStateExists(gomock.Any(), "kl", "in").
					Times(1).
					Return(false, errors.New("error occured"))
			},
			checkresponse: func(t *testing.T, isExists bool, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Equal(t, false, isExists)
				require.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			errMap := make(map[string]models.ErrorResponse)
			store := mock.NewMockCountryRepoImply(ctrl)
			tc.buildStubs(store)
			countryUsecase := NewCountryUseCases(store)
			isExists, fieldsMap, err := countryUsecase.CheckStateExists(context.Background(), tc.funcArguments.iso, tc.funcArguments.countryCode, tc.funcArguments.validation, errMap)
			tc.checkresponse(t, isExists, fieldsMap, err)
		})
	}
}

// TestGetStatesOfCountry to check GetStatesOfCountry function
func TestGetStatesOfCountry(t *testing.T) {

	validation := entities.Validation{
		Endpoint: "states",
		Method:   "get",
		ID:       "194",
	}
	id, err := strconv.ParseInt(validation.ID, 10, 64)
	require.Nil(t, err)
	paginationInfo := entities.Pagination{
		Page:  10,
		Limit: 10,
	}
	params := entities.Params{
		Sort:  "asc",
		Order: "name",
	}
	var (
		states []*entities.GeographicInfo
		n      = 3
	)
	for i := 0; i < n; i++ {
		arg, err := utils.RandomString(4)
		require.Nil(t, err)
		states = append(states, &entities.GeographicInfo{
			ID:   1,
			Name: arg,
			ISO:  "USD",
		})
	}

	testCases := []struct {
		name          string
		params        entities.Params
		buildStubs    func(store *mock.MockCountryRepoImply)
		checkresponse func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name:   "success",
			params: params,
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					GetStatesOfCountry(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), id, gomock.Any()).
					Times(1).
					Return(states, int64(len(states)), nil)

			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Equal(t, resp.MetaData.Total, int64(len(states)))
				require.Nil(t, err)
			},
		},
		{
			name:   "internal server error",
			params: params,
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					GetStatesOfCountry(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), id, gomock.Any()).
					Times(1).
					Return(nil, int64(0), errors.New("error occured"))

			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Error(t, err)
			},
		},
		{
			name:   "no data found",
			params: params,
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					GetStatesOfCountry(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), id, gomock.Any()).
					Times(1).
					Return(nil, int64(0), nil)

			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, resp.MetaData)
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
			store := mock.NewMockCountryRepoImply(ctrl)
			tc.buildStubs(store)
			countryUsecase := NewCountryUseCases(store)
			resp, fieldsMap, err := countryUsecase.GetStatesOfCountry(context.Background(), params, paginationInfo, validation, errMap)
			tc.checkresponse(t, resp, fieldsMap, err)
		})
	}
}

// TestGetAllCountryCodes to check GetAllCountryCodes function
func TestGetAllCountryCodes(t *testing.T) {

	iso := entities.IsoList{
		Iso: []string{"AW",
			"BS",
			"CU"},
	}

	testCases := []struct {
		name          string
		buildStubs    func(store *mock.MockCountryRepoImply)
		checkresponse func(t *testing.T, resp entities.IsoList, err error)
	}{
		{
			name: "success",
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					GetAllCountryCodes(gomock.Any()).
					Times(1).
					Return(iso, nil)

			},
			checkresponse: func(t *testing.T, resp entities.IsoList, err error) {
				require.Equal(t, resp, iso)
				require.Nil(t, err)
			},
		},
		{
			name: "internal server error",
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					GetAllCountryCodes(gomock.Any()).
					Times(1).
					Return(entities.IsoList{}, errors.New("error occured"))

			},
			checkresponse: func(t *testing.T, resp entities.IsoList, err error) {
				require.Equal(t, resp, entities.IsoList{})
				require.Error(t, err)
			},
		},
		{
			name: "no data found",
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					GetAllCountryCodes(gomock.Any()).
					Times(1).
					Return(entities.IsoList{}, nil)

			},
			checkresponse: func(t *testing.T, resp entities.IsoList, err error) {
				require.Equal(t, resp, entities.IsoList{})
				require.Nil(t, err)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mock.NewMockCountryRepoImply(ctrl)
			tc.buildStubs(store)
			countryUsecase := NewCountryUseCases(store)
			resp, err := countryUsecase.GetAllCountryCodes(context.Background())
			tc.checkresponse(t, resp, err)
		})
	}
}

// TestCheckCountryExists to check CheckCountryExists function
func TestCheckCountryExists(t *testing.T) {

	validations := entities.Validation{
		Endpoint: "countries",
		Method:   "get",
	}

	params := entities.IsoParam{
		Iso: "AL",
	}
	countryExists := entities.CountryExists{
		Exists: true,
	}

	testCases := []struct {
		name          string
		funcArguments struct {
			iso        entities.IsoParam
			validation entities.Validation
		}
		buildStubs    func(store *mock.MockCountryRepoImply)
		checkresponse func(t *testing.T, resp entities.CountryExists, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name: "success",
			funcArguments: struct {
				iso        entities.IsoParam
				validation entities.Validation
			}{params, validations},
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					CheckCountryExists(gomock.Any(), params).
					Times(1).
					Return(countryExists, nil)
			},
			checkresponse: func(t *testing.T, resp entities.CountryExists, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Equal(t, resp, countryExists)
				require.Nil(t, err)
			},
		},
		{
			name: "return error from CheckCountryExists",
			funcArguments: struct {
				iso        entities.IsoParam
				validation entities.Validation
			}{params, validations},
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					CheckCountryExists(gomock.Any(), params).
					Times(1).
					Return(entities.CountryExists{}, errors.New("error occured"))
			},
			checkresponse: func(t *testing.T, resp entities.CountryExists, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Equal(t, resp, entities.CountryExists{})
				require.NotNil(t, err)
			},
		},
		{
			name: "empty iso value",
			funcArguments: struct {
				iso        entities.IsoParam
				validation entities.Validation
			}{entities.IsoParam{
				Iso: "",
			}, validations},
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					CheckCountryExists(gomock.Any(), params).
					Times(0)
			},
			checkresponse: func(t *testing.T, resp entities.CountryExists, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, fieldsMap)
				require.Nil(t, err)
			},
		},
		{
			name: "invalid iso value",
			funcArguments: struct {
				iso        entities.IsoParam
				validation entities.Validation
			}{entities.IsoParam{
				Iso: "12",
			}, validations},
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					CheckCountryExists(gomock.Any(), params).
					Times(0)
			},
			checkresponse: func(t *testing.T, resp entities.CountryExists, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, fieldsMap)
				require.Nil(t, err)
			},
		},
		{
			name: "invalid iso value length",
			funcArguments: struct {
				iso        entities.IsoParam
				validation entities.Validation
			}{entities.IsoParam{
				Iso: "asd",
			}, validations},
			buildStubs: func(store *mock.MockCountryRepoImply) {
				store.EXPECT().
					CheckCountryExists(gomock.Any(), params).
					Times(0)
			},
			checkresponse: func(t *testing.T, resp entities.CountryExists, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, fieldsMap)
				require.Nil(t, err)
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			errMap := make(map[string]models.ErrorResponse)
			store := mock.NewMockCountryRepoImply(ctrl)
			tc.buildStubs(store)
			countryUsecase := NewCountryUseCases(store)
			resp, fieldsMap, err := countryUsecase.CheckCountryExists(context.Background(), tc.funcArguments.iso, tc.funcArguments.validation, errMap)
			tc.checkresponse(t, resp, fieldsMap, err)
		})
	}
}
