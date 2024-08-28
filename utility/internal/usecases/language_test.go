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

// TestGetLanguages to check GetLanguages function
func TestGetLanguages(t *testing.T) {

	n := 5
	languages := make([]entities.Language, n)
	for i := 0; i < n; i++ {
		arg, err := utils.RandomString(4)
		require.Nil(t, err)
		languages[i] = entities.Language{
			ID:   int64(i + 1),
			Name: arg,
			Code: "en",
		}
	}

	validation := entities.Validation{
		Endpoint: "languages",
		Method:   "get",
	}

	paginationInfo := entities.Pagination{
		Page:  10,
		Limit: 10,
	}
	params := entities.LangParams{
		Params: entities.Params{
			Name:  "John",
			Sort:  "asc",
			Order: "name",
		},
		Status: "active",
	}
	validation.ID = params.Params.Code
	testCases := []struct {
		name          string
		pagination    entities.Pagination
		params        entities.LangParams
		buildStubs    func(store *mock.MockLanguageRepoImply)
		checkresponse func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name:   "success",
			params: params,
			buildStubs: func(store *mock.MockLanguageRepoImply) {
				store.EXPECT().
					GetLanguages(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return([]*entities.Language{}, int64(len(languages)), nil)
			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)
				require.Equal(t, resp.MetaData.Total, int64(len(languages)))
			},
		},
		{
			name:   "Return error from GetLanguages",
			params: params,
			buildStubs: func(store *mock.MockLanguageRepoImply) {
				store.EXPECT().
					GetLanguages(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return([]*entities.Language{}, int64(0), errors.New("error occured"))
			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, err)
				require.Nil(t, fieldsMap)
			},
		},
		{
			name: "success",
			params: entities.LangParams{
				Params: entities.Params{
					Name:  "John",
					Sort:  "asc",
					Order: "name",
					Code:  "5",
				},
				Status: "active",
			},
			buildStubs: func(store *mock.MockLanguageRepoImply) {
				store.EXPECT().
					GetLanguages(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
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

			store := mock.NewMockLanguageRepoImply(ctrl)
			tc.buildStubs(store)
			languageUsecase := NewLanguageUseCases(store)
			resp, fieldsMap, err := languageUsecase.GetLanguages(context.Background(), tc.params, paginationInfo, validation, errMap)
			tc.checkresponse(t, resp, fieldsMap, err)
		})
	}
}

// TestGetLanguageCodeExists to check GetLanguageCodeExists function
func TestGetLanguageCodeExists(t *testing.T) {

	validation := entities.Validation{
		Endpoint: "languages",
		Method:   "get",
		ID:       "ro",
	}

	testCases := []struct {
		name          string
		validation    entities.Validation
		buildStubs    func(store *mock.MockLanguageRepoImply)
		checkresponse func(t *testing.T, resp bool, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name:       "success",
			validation: validation,
			buildStubs: func(store *mock.MockLanguageRepoImply) {
				store.EXPECT().
					GetLanguageCodeExists(gomock.Any(), validation.ID).
					Times(1).
					Return(true, nil)
			},
			checkresponse: func(t *testing.T, resp bool, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)
				require.Equal(t, true, resp)

			},
		},
		{
			name:       "return error from GetLanguageCodeExists",
			validation: validation,
			buildStubs: func(store *mock.MockLanguageRepoImply) {
				store.EXPECT().
					GetLanguageCodeExists(gomock.Any(), validation.ID).
					Times(1).
					Return(false, errors.New("error occured"))
			},
			checkresponse: func(t *testing.T, resp bool, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, err)
				require.Equal(t, false, resp)

			},
		},
		{
			name:       "Return false when code not found",
			validation: validation,
			buildStubs: func(store *mock.MockLanguageRepoImply) {
				store.EXPECT().
					GetLanguageCodeExists(gomock.Any(), validation.ID).
					Times(1).
					Return(false, nil)
			},
			checkresponse: func(t *testing.T, resp bool, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)
				require.Equal(t, false, resp)

			},
		},
		{
			name: "success",
			validation: entities.Validation{
				Endpoint: "languages",
				Method:   "get",
				ID:       "ros",
			},
			buildStubs: func(store *mock.MockLanguageRepoImply) {
				store.EXPECT().
					GetLanguageCodeExists(gomock.Any(), validation.ID).
					Times(0)
			},
			checkresponse: func(t *testing.T, resp bool, fieldsMap map[string]models.ErrorResponse, err error) {
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

			store := mock.NewMockLanguageRepoImply(ctrl)
			tc.buildStubs(store)
			languageUsecase := NewLanguageUseCases(store)
			isExists, fieldsMap, err := languageUsecase.GetLanguageCodeExists(context.Background(), tc.validation, errMap)
			tc.checkresponse(t, isExists, fieldsMap, err)
		})
	}
}
