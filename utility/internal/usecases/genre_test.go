// // nolint
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

// TestGetGenresByID to check GetGenresByID function for test
func TestGetGenresByID(t *testing.T) {

	validation := entities.Validation{
		Endpoint: "genre",
		Method:   "get",
		ID:       "5173291f-c2ba-41f8-a0ba-1c76f87c06e8",
	}

	response := entities.GenreDetails{
		Name:                    "classic pop spioicxppcz",
		LanguageLabelIdentifier: "classicPopSpioicxppczGenreLabel",
		IsDeleted:               false,
	}

	testCases := []struct {
		name          string
		validation    entities.Validation
		buildStubs    func(store *mock.MockGenreRepoImply)
		checkresponse func(t *testing.T, resp entities.GenreDetails, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name:       "success",
			validation: validation,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					GetGenresByID(gomock.Any(), validation.ID).
					Times(1).
					Return(response, nil)
			},
			checkresponse: func(t *testing.T, resp entities.GenreDetails, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)

			},
		},
		{
			name:       "Failed while returning error from GetGenresByID",
			validation: validation,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					GetGenresByID(gomock.Any(), validation.ID).
					Times(1).
					Return(entities.GenreDetails{}, errors.New("error occured"))
			},
			checkresponse: func(t *testing.T, resp entities.GenreDetails, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, err)
				require.Equal(t, entities.GenreDetails{}, resp)

			},
		},
		{
			name: "invalid id",
			validation: entities.Validation{
				Endpoint: "genre",
				Method:   "get",
				ID:       "5173291f6e8",
			},
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					GetGenresByID(gomock.Any(), validation.ID).
					Times(0)
			},
			checkresponse: func(t *testing.T, resp entities.GenreDetails, fieldsMap map[string]models.ErrorResponse, err error) {
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

			store := mock.NewMockGenreRepoImply(ctrl)
			tc.buildStubs(store)
			genreUsecase := NewGenreUseCases(store)
			resp, fieldsMap, err := genreUsecase.GetGenresByID(context.Background(), tc.validation, errMap)
			tc.checkresponse(t, resp, fieldsMap, err)
		})
	}
}

// TestCreateGenre to check CreateGenre function for test
func TestCreateGenre(t *testing.T) {

	validation := entities.Validation{
		Endpoint: "genre",
		Method:   "post",
	}

	//Uncomment this if you need to get random genre name and generate its LanguageIdentifier
	// arg, err := utils.RandomString(4)
	// require.Nil(t, err)
	// langIdentifier, err := utilities.GenerateLangIdentifier(arg, consts.Genre)
	// require.Nil(t, err)
	langIdentifier := "sampleGenreLabel"
	genre := entities.Genre{
		Name: "sample",
	}

	testCases := []struct {
		name          string
		genre         entities.Genre
		buildStubs    func(store *mock.MockGenreRepoImply)
		checkResponse func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name:  "Valid genre",
			genre: genre,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					GenreNameExists(gomock.Any(), genre.Name).
					Times(1).
					Return(false, nil)
				store.EXPECT().
					CreateGenre(gomock.Any(), genre.Name, langIdentifier).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, fieldsMap)
				require.Nil(t, err)
			},
		},
		{
			name:  "Invalid genre", // return error from GenreNameExists
			genre: genre,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					GenreNameExists(gomock.Any(), genre.Name).
					Times(1).
					Return(true, errors.New("error occured"))
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, err)
			},
		},
		{
			name:  "DB error", //Return error from
			genre: genre,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					GenreNameExists(gomock.Any(), genre.Name).
					Times(1).
					Return(false, nil)
				store.EXPECT().
					CreateGenre(gomock.Any(), genre.Name, langIdentifier).
					Times(1).
					Return(errors.New("error ocuured"))
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, err)
			},
		},
		{
			name:  "Invalid genre", // When name is already exists
			genre: genre,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					GenreNameExists(gomock.Any(), genre.Name).
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
			genre: entities.Genre{
				Name: "",
			},
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					GenreNameExists(gomock.Any(), genre.Name).
					Times(0)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, fieldsMap)
				require.Nil(t, err)
			},
		},
		{
			name: "Invalid name length",
			genre: entities.Genre{
				Name: "a",
			},
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					GenreNameExists(gomock.Any(), genre.Name).
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

			store := mock.NewMockGenreRepoImply(ctrl)
			tc.buildStubs(store)
			genreUsecase := NewGenreUseCases(store)

			fieldsMap, err := genreUsecase.CreateGenre(context.Background(), tc.genre, validation, errMap)
			tc.checkResponse(t, fieldsMap, err)
		})
	}
}

// TestUpdateGenre to check UpdateGenre function for test
func TestUpdateGenre(t *testing.T) {

	validation := entities.Validation{
		Endpoint: "genre",
		Method:   "patch",
		ID:       "01050762-1488-4bbb-8a9b-443be77ad395",
	}
	uuid, err := uuid.Parse(validation.ID)
	require.Nil(t, err)
	//Uncomment this if you need to get random genre name and generate its LanguageIdentifier
	// arg, err := utils.RandomString(4)
	// require.Nil(t, err)
	// langIdentifier, err := utilities.GenerateLangIdentifier(arg, consts.Genre)
	// require.Nil(t, err)
	langIdentifier := "samplesGenreLabel"
	genre := entities.Genre{
		Name: "samples",
		ID:   uuid,
	}

	testCases := []struct {
		name          string
		genre         entities.Genre
		validation    entities.Validation
		buildStubs    func(store *mock.MockGenreRepoImply)
		checkResponse func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name:       "Valid genre",
			genre:      genre,
			validation: validation,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					IsGenreExists(gomock.Any(), uuid).
					Times(1).
					Return(nil)
				store.EXPECT().
					IsDuplicateExists(gomock.Any(), "name", genre.Name, validation.ID).
					Return(false, nil)
				store.EXPECT().
					UpdateGenre(gomock.Any(), genre, langIdentifier).
					Return(nil)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, fieldsMap)
				require.Nil(t, err)
			},
		},
		{
			name:       "Return error from IsDuplicateExists",
			genre:      genre,
			validation: validation,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					IsGenreExists(gomock.Any(), uuid).
					Times(1).
					Return(nil)
				store.EXPECT().
					IsDuplicateExists(gomock.Any(), "name", genre.Name, validation.ID).
					Return(false, errors.New("error occured"))

			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, err)
			},
		},
		{
			name:       "If name is already exists",
			genre:      genre,
			validation: validation,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					IsGenreExists(gomock.Any(), uuid).
					Times(1).
					Return(nil)
				store.EXPECT().
					IsDuplicateExists(gomock.Any(), "name", genre.Name, validation.ID).
					Return(true, nil)

			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)
				require.NotNil(t, fieldsMap)

			},
		},
		{
			name:       "Error from database",
			genre:      genre,
			validation: validation,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					IsGenreExists(gomock.Any(), uuid).
					Times(1).
					Return(nil)
				store.EXPECT().
					IsDuplicateExists(gomock.Any(), "name", genre.Name, validation.ID).
					Return(false, nil)
				store.EXPECT().
					UpdateGenre(gomock.Any(), genre, langIdentifier).
					Return(errors.New(" error from database,"))
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, fieldsMap)
				require.NotNil(t, err)
			},
		},
		{
			name:       "Non exists genre",
			genre:      genre,
			validation: validation,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					IsGenreExists(gomock.Any(), uuid).
					Times(1).
					Return(errors.New("genre already deleted"))
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, fieldsMap)
				require.Equal(t, consts.ErrNotExist, err)
			},
		},
		{
			name:       "Return error from IsGenreExists",
			genre:      genre,
			validation: validation,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					IsGenreExists(gomock.Any(), uuid).
					Times(1).
					Return(errors.New("error occured"))
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, fieldsMap)
				require.NotNil(t, err)
			},
		},
		{
			name:  "Invalid genre id",
			genre: genre,
			validation: entities.Validation{
				Endpoint: "genre",
				Method:   "patch",
				ID:       "0105076395",
			},
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					IsGenreExists(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)
				require.NotNil(t, fieldsMap)
			},
		},
		{
			name: "Empty genre name",
			genre: entities.Genre{
				Name: "",
				ID:   uuid,
			},
			validation: validation,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					IsGenreExists(gomock.Any(), uuid).
					Times(1).
					Return(nil)
				store.EXPECT().
					IsDuplicateExists(gomock.Any(), "name", genre.Name, validation.ID).
					Times(0)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.NotNil(t, fieldsMap)
				require.Nil(t, err)
			},
		},
		{
			name: "Invalid genre name length",
			genre: entities.Genre{
				Name: "s",
				ID:   uuid,
			},
			validation: validation,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					IsGenreExists(gomock.Any(), uuid).
					Times(1).
					Return(nil)
				store.EXPECT().
					IsDuplicateExists(gomock.Any(), "name", genre.Name, validation.ID).
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

			store := mock.NewMockGenreRepoImply(ctrl)
			tc.buildStubs(store)
			genreUsecase := NewGenreUseCases(store)

			fieldsMap, err := genreUsecase.UpdateGenre(context.Background(), tc.genre, tc.validation, errMap)
			tc.checkResponse(t, fieldsMap, err)
		})
	}
}

// TestGetGenres to check GetGenres function for test
func TestGetGenres(t *testing.T) {

	validation := entities.Validation{
		Endpoint: "genre",
		Method:   "get",
		ID:       "5173291f-c2ba-41f8-a0ba-1c76f87c06e8",
	}
	paginationInfo := entities.Pagination{
		Page:  10,
		Limit: 10,
	}
	params := entities.Params{
		Name:   "John",
		Sort:   "asc",
		Order:  "name",
		IdList: "01050762-1488-4bbb-8a9b-443be77ad395,0e252f0a-ab6d-4b94-acd9-198eb3f01e58,12dfcb93-8d8a-4895-b8ff-f392a5084ccf",
	}
	n := 5
	var genres []*entities.Genre
	for i := 0; i < n; i++ {
		arg, err := utils.RandomString(4)
		require.Nil(t, err)
		genres = append(genres, &entities.Genre{
			ID:   uuid.New(),
			Name: arg,
		})
	}

	testCases := []struct {
		name          string
		pagination    entities.Pagination
		params        entities.Params
		buildStubs    func(store *mock.MockGenreRepoImply)
		checkresponse func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name:   "success",
			params: params,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					GetGenres(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(genres, int64(len(genres)), nil)
			},
			checkresponse: func(t *testing.T, resp *entities.Response, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, err)
				require.Equal(t, int64(len(genres)), resp.MetaData.Total)
			},
		},
		{
			name:   "internal server error",
			params: params,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					GetGenres(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					GetGenres(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					GetGenres(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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

			store := mock.NewMockGenreRepoImply(ctrl)
			tc.buildStubs(store)
			genreUsecase := NewGenreUseCases(store)
			resp, fieldsMap, err := genreUsecase.GetGenres(context.Background(), tc.params, paginationInfo, validation, errMap)
			tc.checkresponse(t, resp, fieldsMap, err)
		})
	}
}

// TestDeleteGenre to check DeleteGenre function for test
func TestDeleteGenre(t *testing.T) {

	validation := entities.Validation{
		Endpoint: "genre",
		Method:   "delete",
		ID:       "01050762-1488-4bbb-8a9b-443be77ad395",
	}
	uuid, err := uuid.Parse(validation.ID)
	require.Nil(t, err)

	genre := entities.Genre{
		ID: uuid,
	}

	testCases := []struct {
		name          string
		genre         entities.Genre
		validation    entities.Validation
		buildStubs    func(store *mock.MockGenreRepoImply)
		checkResponse func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error)
	}{
		{
			name:       "Valid genre",
			genre:      genre,
			validation: validation,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					IsGenreExists(gomock.Any(), uuid).
					Times(1).
					Return(nil)
				store.EXPECT().
					DeleteGenre(gomock.Any(), uuid).
					Return(int64(1), nil)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, fieldsMap)
				require.Nil(t, err)
			},
		},
		{
			name:       "Doesn't affect any rows",
			genre:      genre,
			validation: validation,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					IsGenreExists(gomock.Any(), uuid).
					Times(1).
					Return(nil)
				store.EXPECT().
					DeleteGenre(gomock.Any(), uuid).
					Return(int64(0), nil)
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, fieldsMap)
				require.Equal(t, consts.ErrNotFound, err)
			},
		},
		{
			name:       "Error from database",
			genre:      genre,
			validation: validation,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					IsGenreExists(gomock.Any(), uuid).
					Times(1).
					Return(nil)
				store.EXPECT().
					DeleteGenre(gomock.Any(), uuid).
					Return(int64(0), errors.New(" error from database,"))
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, fieldsMap)
				require.NotNil(t, err)
			},
		},
		{
			name:       "Non exists genre",
			genre:      genre,
			validation: validation,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					IsGenreExists(gomock.Any(), uuid).
					Times(1).
					Return(errors.New("genre already deleted"))
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, fieldsMap)
				require.Equal(t, consts.ErrNotExist, err)
			},
		},
		{
			name:       "Return error from IsGenreExists",
			genre:      genre,
			validation: validation,
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					IsGenreExists(gomock.Any(), uuid).
					Times(1).
					Return(errors.New("error occured"))
			},
			checkResponse: func(t *testing.T, fieldsMap map[string]models.ErrorResponse, err error) {
				require.Nil(t, fieldsMap)
				require.NotNil(t, err)
			},
		},
		{
			name:  "Invalid genre id",
			genre: genre,
			validation: entities.Validation{
				Endpoint: "genre",
				Method:   "delete",
				ID:       "0105076395",
			},
			buildStubs: func(store *mock.MockGenreRepoImply) {
				store.EXPECT().
					IsGenreExists(gomock.Any(), gomock.Any()).
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

			store := mock.NewMockGenreRepoImply(ctrl)
			tc.buildStubs(store)
			genreUsecase := NewGenreUseCases(store)

			fieldsMap, err := genreUsecase.DeleteGenre(context.Background(), tc.validation, errMap)
			tc.checkResponse(t, fieldsMap, err)
		})
	}
}
