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

// TestGetLookupByTypeName to check GetLookupByTypeName function
func TestGetLookupByTypeName(t *testing.T) {

	validation := entities.Validation{
		Endpoint: "lookup",
		Method:   "get",
		ID:       "business_model",
	}

	lookup := []entities.Lookup{
		{
			ID:           463,
			Name:         "Both",
			Value:        "both",
			LookupTypeId: 1,
		},
	}

	filter := map[string]string{"value": "Subscription"}
	testCases := []struct {
		name          string
		lookup        entities.Lookup
		validation    entities.Validation
		buildStubs    func(store *mock.MockLookupRepoImply)
		checkresponse func(t *testing.T, resp []entities.Lookup, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name:       "success",
			validation: validation,
			buildStubs: func(store *mock.MockLookupRepoImply) {
				store.EXPECT().
					GetLookupByTypeName(gomock.Any(), validation.ID, filter).
					Times(1).
					Return(lookup, nil)
			},
			checkresponse: func(t *testing.T, resp []entities.Lookup, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)
				require.Nil(t, fieldsMap)
			},
		},
		{
			name:       "internal server error",
			validation: validation,
			buildStubs: func(store *mock.MockLookupRepoImply) {
				store.EXPECT().
					GetLookupByTypeName(gomock.Any(), validation.ID, filter).
					Times(1).
					Return(lookup, errors.New("error occured"))
			},
			checkresponse: func(t *testing.T, resp []entities.Lookup, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, err)
				require.Nil(t, fieldsMap)
			},
		},
		{
			name: "validation error due to invalid type name",
			validation: entities.Validation{
				Endpoint: "lookup",
				Method:   "get",
				ID:       "123456",
			},
			buildStubs: func(store *mock.MockLookupRepoImply) {
				store.EXPECT().
					GetLookupByTypeName(gomock.Any(), validation.ID, filter).
					Times(0)
			},
			checkresponse: func(t *testing.T, resp []entities.Lookup, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)
				require.NotNil(t, fieldsMap)
			},
		},
		{
			name:       "invalid endpoint",
			validation: entities.Validation{},
			buildStubs: func(store *mock.MockLookupRepoImply) {
				store.EXPECT().
					GetLookupByTypeName(gomock.Any(), validation.ID, filter).
					Times(0)
			},
			checkresponse: func(t *testing.T, resp []entities.Lookup, fieldsMap map[string]models.ErrorResponse, err error) {
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

			store := mock.NewMockLookupRepoImply(ctrl)
			tc.buildStubs(store)
			lookupUsecase := NewLookupUseCases(store)
			resp, fieldsMap, err := lookupUsecase.GetLookupByTypeName(context.Background(), filter, tc.validation, errMap)
			tc.checkresponse(t, resp, fieldsMap, err)
		})
	}
}

// TestGetLookupByIdList to check GetLookupById function
func TestGetLookupByIdList(t *testing.T) {

	idList := entities.LookupIDs{
		ID: []int64{1, 2, 3, 4, 5, 6},
	}

	response := entities.LookupData{
		Lookup: []entities.Lookup{
			{
				ID:                 1,
				Name:               "Member",
				Value:              "member",
				Description:        "Test Description",
				Position:           1,
				LookupTypeId:       1,
				LanguageIdentifier: "en",
			},
			{
				ID:                 2,
				Name:               "Partner",
				Value:              "partner",
				Description:        "Test Description",
				Position:           2,
				LookupTypeId:       2,
				LanguageIdentifier: "en",
			},
			{
				ID:                 3,
				Name:               "Admin",
				Value:              "admin",
				Description:        "Test Description",
				Position:           3,
				LookupTypeId:       3,
				LanguageIdentifier: "en",
			},
		},
		InvalidLookupIds: []int64{4, 5, 6},
	}

	testCases := []struct {
		name          string
		lookuIDs      entities.LookupIDs
		buildStubs    func(store *mock.MockLookupRepoImply)
		checkresponse func(t *testing.T, resp entities.LookupData, err error)
	}{
		{
			name:     "success",
			lookuIDs: idList,
			buildStubs: func(store *mock.MockLookupRepoImply) {
				store.EXPECT().
					GetLookupByIdList(gomock.Any(), idList).
					Times(1).
					Return(response, nil)
			},
			checkresponse: func(t *testing.T, resp entities.LookupData, err error) {
				require.Nil(t, err)
			},
		},
		{
			name:     "failed",
			lookuIDs: idList,
			buildStubs: func(store *mock.MockLookupRepoImply) {
				store.EXPECT().
					GetLookupByIdList(gomock.Any(), idList).
					Times(1).
					Return(entities.LookupData{}, errors.New("Error occured"))
			},
			checkresponse: func(t *testing.T, resp entities.LookupData, err error) {
				require.NotNil(t, err)
				require.Empty(t, resp)

			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock.NewMockLookupRepoImply(ctrl)
			tc.buildStubs(store)
			lookupUsecase := NewLookupUseCases(store)
			resp, err := lookupUsecase.GetLookupByIdList(context.Background(), tc.lookuIDs)
			tc.checkresponse(t, resp, err)
		})
	}
}
