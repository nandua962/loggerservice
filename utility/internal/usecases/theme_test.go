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
)

// TestGetThemeByID to check GetThemeByID function for test
func TestGetThemeByID(t *testing.T) {

	validation := entities.Validation{
		Endpoint: "theme",
		Method:   "get",
		ID:       "1",
	}

	response := entities.Theme{
		ID:       1,
		Name:     "Example Theme",
		Value:    "theme_value",
		LayoutID: 3,
	}

	testCases := []struct {
		name          string
		validation    entities.Validation
		buildStubs    func(store *mock.MockThemeRepoImply)
		checkresponse func(t *testing.T, resp entities.Theme, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name:       "success",
			validation: validation,
			buildStubs: func(store *mock.MockThemeRepoImply) {
				store.EXPECT().
					GetThemeByID(gomock.Any(), validation.ID).
					Times(1).
					Return(response, nil)
			},
			checkresponse: func(t *testing.T, resp entities.Theme, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)

			},
		},
		{
			name:       "Failed while returning error from GetGenresByID",
			validation: validation,
			buildStubs: func(store *mock.MockThemeRepoImply) {
				store.EXPECT().
					GetThemeByID(gomock.Any(), validation.ID).
					Times(1).
					Return(entities.Theme{}, errors.New("error occured"))
			},
			checkresponse: func(t *testing.T, resp entities.Theme, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, err)
				require.Equal(t, entities.Theme{}, resp)

			},
		},
		{
			name: "invalid theme id",
			validation: entities.Validation{
				Endpoint: "theme",
				Method:   "get",
				ID:       "av",
			},
			buildStubs: func(store *mock.MockThemeRepoImply) {
				store.EXPECT().
					GetThemeByID(gomock.Any(), validation.ID).
					Times(0)
			},
			checkresponse: func(t *testing.T, resp entities.Theme, fieldsMap map[string]models.ErrorResponse, err error) {
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

			store := mock.NewMockThemeRepoImply(ctrl)
			tc.buildStubs(store)
			themeUsecase := NewThemeUseCases(store)
			resp, fieldsMap, err := themeUsecase.GetThemeByID(context.Background(), tc.validation, errMap)
			tc.checkresponse(t, resp, fieldsMap, err)
		})
	}
}
