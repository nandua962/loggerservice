// nolint
package usecases

import (
	"context"
	"errors"
	"testing"
	"utility/internal/consts"
	"utility/internal/entities"
	"utility/internal/repo/mock"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils"
)

// TestGetRoleByID to check GetRoleByID function for test
func TestGetRoleByID(t *testing.T) {

	validation := entities.Validation{
		Endpoint: "roles",
		Method:   "get",
		ID:       "4f081dd9-a818-420c-9a40-563955cf93df",
	}

	response := entities.Role{
		ID:                      "4f081dd9-a818-420c-9a40-563955cf93df",
		Name:                    "Music Engineer",
		CustomName:              "Music Engineer",
		LanguageLabelIdentifier: "musicEngineerRoleLabel",
		IsDefault:               false,
	}

	testCases := []struct {
		name          string
		validation    entities.Validation
		buildStubs    func(store *mock.MockRoleRepoImply)
		checkresponse func(t *testing.T, resp entities.Role, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name:       "success",
			validation: validation,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					GetRoleByID(gomock.Any(), validation.ID).
					Times(1).
					Return(response, nil)
			},
			checkresponse: func(t *testing.T, resp entities.Role, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)

			},
		},
		{
			name:       "Failed while returning error from GetRoleByID",
			validation: validation,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					GetRoleByID(gomock.Any(), validation.ID).
					Times(1).
					Return(entities.Role{}, errors.New("error occured"))
			},
			checkresponse: func(t *testing.T, resp entities.Role, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, err)
				require.Equal(t, entities.Role{}, resp)

			},
		},
		{
			name: "invalid id",
			validation: entities.Validation{
				Endpoint: "roles",
				Method:   "get",
				ID:       "4f081d55cf93df",
			},
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					GetRoleByID(gomock.Any(), validation.ID).
					Times(0)
			},
			checkresponse: func(t *testing.T, resp entities.Role, fieldsMap map[string]models.ErrorResponse, err error) {
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

			store := mock.NewMockRoleRepoImply(ctrl)
			tc.buildStubs(store)
			roleUsecase := NewRoleUseCases(store)
			resp, fieldsMap, err := roleUsecase.GetRoleByID(context.Background(), tc.validation, errMap)
			tc.checkresponse(t, resp, fieldsMap, err)
		})
	}
}

// TestCreateRole to check CreateRole function for test
func TestCreateRole(t *testing.T) {

	validation := entities.Validation{
		Endpoint: "roles",
		Method:   "post",
	}

	langIdentifier := "sampleRoleLabel"
	role := entities.Role{
		Name:       "sample",
		CustomName: "sample",
	}

	testCases := []struct {
		name          string
		role          entities.Role
		buildStubs    func(store *mock.MockRoleRepoImply)
		checkResponse func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name: "Valid role",
			role: role,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					RoleNameExists(gomock.Any(), role.Name).
					Times(1).
					Return(false, nil)
				store.EXPECT().
					CreateRole(gomock.Any(), role, langIdentifier).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, fieldsMap)
				require.Nil(t, err)
			},
		},
		{
			name: "Invalid role", // return error from RoleNameExists
			role: role,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					RoleNameExists(gomock.Any(), role.Name).
					Times(1).
					Return(true, errors.New("error occured"))
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, err)
			},
		},
		{
			name: "DB error", //Return error from
			role: role,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					RoleNameExists(gomock.Any(), role.Name).
					Times(1).
					Return(false, nil)
				store.EXPECT().
					CreateRole(gomock.Any(), role, langIdentifier).
					Times(1).
					Return(errors.New("error ocuured"))
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, err)
			},
		},
		{
			name: "Invalid role", // When name is already exists
			role: role,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					RoleNameExists(gomock.Any(), role.Name).
					Times(1).
					Return(true, nil)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)
				require.NotNil(t, fieldsMap)

			},
		},
		{
			name: "Empty name",
			role: entities.Role{
				Name: "",
			},
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					RoleNameExists(gomock.Any(), role.Name).
					Times(0)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, fieldsMap)
				require.Nil(t, err)
			},
		},
		{
			name: "Invalid name length",
			role: entities.Role{
				Name: "a",
			},
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					RoleNameExists(gomock.Any(), role.Name).
					Times(0)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
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

			store := mock.NewMockRoleRepoImply(ctrl)
			tc.buildStubs(store)
			roleUsecase := NewRoleUseCases(store)

			fieldsMap, err := roleUsecase.CreateRole(context.Background(), tc.role, validation, errMap)

			tc.checkResponse(t, fieldsMap, err)
		})
	}
}

// TestUpdateRole to check UpdateRole function for test
func TestUpdateRole(t *testing.T) {

	validation := entities.Validation{
		Endpoint: "roles",
		Method:   "patch",
		ID:       "01050762-1488-4bbb-8a9b-443be77ad395",
	}
	uuid, err := uuid.Parse(validation.ID)
	require.Nil(t, err)

	langIdentifier := "samplesRoleLabel"
	role := entities.Role{
		Name:       "samples",
		ID:         "01050762-1488-4bbb-8a9b-443be77ad395",
		CustomName: "samples",
	}

	testCases := []struct {
		name          string
		role          entities.Role
		validation    entities.Validation
		buildStubs    func(store *mock.MockRoleRepoImply)
		checkResponse func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name:       "Valid role",
			role:       role,
			validation: validation,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					IsRoleExists(gomock.Any(), uuid).
					Times(1).
					Return(nil)
				store.EXPECT().
					IsDuplicateExists(gomock.Any(), "name", role.Name, validation.ID).
					Return(false, nil)
				store.EXPECT().
					UpdateRole(gomock.Any(), role, langIdentifier).
					Return(int64(1), nil)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Empty(t, fieldsMap)
				require.Nil(t, err)
			},
		},
		{
			name:       "Return error from IsDuplicateExists",
			role:       role,
			validation: validation,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					IsRoleExists(gomock.Any(), uuid).
					Times(1).
					Return(nil)
				store.EXPECT().
					IsDuplicateExists(gomock.Any(), "name", role.Name, validation.ID).
					Return(false, errors.New("error occured"))

			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, err)
			},
		},
		{
			name:       "If name is already exists",
			role:       role,
			validation: validation,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					IsRoleExists(gomock.Any(), uuid).
					Times(1).
					Return(nil)
				store.EXPECT().
					IsDuplicateExists(gomock.Any(), "name", role.Name, validation.ID).
					Return(true, nil)

			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)
				require.NotNil(t, fieldsMap)

			},
		},
		{
			name:       "Error from database",
			role:       role,
			validation: validation,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					IsRoleExists(gomock.Any(), uuid).
					Times(1).
					Return(nil)
				store.EXPECT().
					IsDuplicateExists(gomock.Any(), "name", role.Name, validation.ID).
					Return(false, nil)
				store.EXPECT().
					UpdateRole(gomock.Any(), role, langIdentifier).
					Return(int64(0), errors.New(" error from database,"))
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, fieldsMap)
				require.NotNil(t, err)
			},
		},
		{
			name:       "Non exists role",
			role:       role,
			validation: validation,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					IsRoleExists(gomock.Any(), uuid).
					Times(1).
					Return(errors.New("role already deleted"))
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, fieldsMap)
				require.Equal(t, consts.ErrNotExist, err)
			},
		},
		{
			name:       "Return error from IsRoleExists",
			role:       role,
			validation: validation,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					IsRoleExists(gomock.Any(), uuid).
					Times(1).
					Return(errors.New("error occured"))
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, fieldsMap)
				require.NotNil(t, err)
			},
		},
		{
			name: "Invalid role id",
			role: role,
			validation: entities.Validation{
				Endpoint: "roles",
				Method:   "patch",
				ID:       "0105076395",
			},
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					IsRoleExists(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)
				require.NotNil(t, fieldsMap)
			},
		},
		{
			name: "Empty role name",
			role: entities.Role{
				Name:       "",
				ID:         "01050762-1488-4bbb-8a9b-443be77ad395",
				CustomName: "",
			},
			validation: validation,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					IsRoleExists(gomock.Any(), uuid).
					Times(1).
					Return(nil)
				store.EXPECT().
					IsDuplicateExists(gomock.Any(), "name", role.Name, validation.ID).
					Times(0)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, fieldsMap)
				require.Nil(t, err)
			},
		},
		{
			name: "Invalid role name length",
			role: entities.Role{
				Name:       "sa",
				ID:         "01050762-1488-4bbb-8a9b-443be77ad395",
				CustomName: "s",
			},
			validation: validation,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					IsRoleExists(gomock.Any(), uuid).
					Times(1).
					Return(nil)
				store.EXPECT().
					IsDuplicateExists(gomock.Any(), "name", role.Name, validation.ID).
					Times(0)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
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

			store := mock.NewMockRoleRepoImply(ctrl)
			tc.buildStubs(store)
			roleUsecase := NewRoleUseCases(store)

			fieldsMap, err := roleUsecase.UpdateRole(context.Background(), tc.role, tc.validation, errMap)
			tc.checkResponse(t, fieldsMap, err)
		})
	}
}

// TestDeleteRoles to check DeleteRoles function for test
func TestDeleteRoles(t *testing.T) {

	validation := entities.Validation{
		Endpoint: "roles",
		Method:   "delete",
		ID:       "01050762-1488-4bbb-8a9b-443be77ad395",
	}
	uuid, err := uuid.Parse(validation.ID)
	require.Nil(t, err)

	role := entities.Role{
		ID: "01050762-1488-4bbb-8a9b-443be77ad395",
	}

	testCases := []struct {
		name          string
		role          entities.Role
		validation    entities.Validation
		buildStubs    func(store *mock.MockRoleRepoImply)
		checkResponse func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name:       "Valid role",
			role:       role,
			validation: validation,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					IsRoleExists(gomock.Any(), uuid).
					Times(1).
					Return(nil)
				store.EXPECT().
					DeleteRoles(gomock.Any(), uuid).
					Return(int64(1), nil)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, fieldsMap)
				require.Nil(t, err)
			},
		},
		{
			name:       "Doesn't affect any rows",
			role:       role,
			validation: validation,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					IsRoleExists(gomock.Any(), uuid).
					Times(1).
					Return(nil)
				store.EXPECT().
					DeleteRoles(gomock.Any(), uuid).
					Return(int64(0), nil)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, fieldsMap)
				require.Equal(t, consts.ErrNotFound, err)
			},
		},
		{
			name:       "Error from database",
			role:       role,
			validation: validation,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					IsRoleExists(gomock.Any(), uuid).
					Times(1).
					Return(nil)
				store.EXPECT().
					DeleteRoles(gomock.Any(), uuid).
					Return(int64(0), errors.New(" error from database,"))
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, fieldsMap)
				require.NotNil(t, err)
			},
		},
		{
			name:       "Non exists role",
			role:       role,
			validation: validation,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					IsRoleExists(gomock.Any(), uuid).
					Times(1).
					Return(errors.New("role already deleted"))
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, fieldsMap)
				require.Equal(t, consts.ErrNotExist, err)
			},
		},
		{
			name:       "Return error from IsRoleExists",
			role:       role,
			validation: validation,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					IsRoleExists(gomock.Any(), uuid).
					Times(1).
					Return(errors.New("error occured"))
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, fieldsMap)
				require.NotNil(t, err)
			},
		},
		{
			name: "Invalid role id",
			role: role,
			validation: entities.Validation{
				Endpoint: "roles",
				Method:   "delete",
				ID:       "0105076395",
			},
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					IsRoleExists(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
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

			store := mock.NewMockRoleRepoImply(ctrl)
			tc.buildStubs(store)
			roleUsecase := NewRoleUseCases(store)

			fieldsMap, err := roleUsecase.DeleteRoles(context.Background(), tc.validation, errMap)
			tc.checkResponse(t, fieldsMap, err)
		})
	}
}

// TestGetRoles to check GetRoles function for test
func TestGetRoles(t *testing.T) {

	validation := entities.Validation{
		Endpoint: "roles",
		Method:   "get",
		ID:       "5173291f-c2ba-41f8-a0ba-1c76f87c06e8",
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
	n := 5
	var roles []*entities.Role
	for i := 0; i < n; i++ {
		arg, err := utils.RandomString(4)
		require.Nil(t, err)
		roles = append(roles, &entities.Role{
			ID:   "5173291f-c2ba-41f8-a0ba-1c76f87c06e8",
			Name: arg,
		})
	}

	testCases := []struct {
		name          string
		pagination    entities.Pagination
		params        entities.Params
		buildStubs    func(store *mock.MockRoleRepoImply)
		checkresponse func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name:   "success",
			params: params,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					GetRoles(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(roles, int64(len(roles)), nil)
			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)
				require.Equal(t, int64(len(roles)), resp.MetaData.Total)
			},
		},
		{
			name:   "internal server error",
			params: params,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					GetRoles(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, int64(0), errors.New("error occured"))
			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Error(t, err)
			},
		},
		{
			name:   "No records",
			params: params,
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					GetRoles(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, int64(0), nil)
			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)
			},
		},
		{
			name: "search with empty name and return no records",
			params: entities.Params{
				Name:  "",
				Sort:  "asc",
				Order: "date",
			},
			buildStubs: func(store *mock.MockRoleRepoImply) {
				store.EXPECT().
					GetRoles(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, int64(0), nil)
			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
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

			store := mock.NewMockRoleRepoImply(ctrl)
			tc.buildStubs(store)
			roleUsecase := NewRoleUseCases(store)
			resp, fieldsMap, err := roleUsecase.GetRoles(context.Background(), tc.params, paginationInfo, validation, errMap)
			tc.checkresponse(t, resp, fieldsMap, err)
		})
	}
}
