package repo

import (
	"context"
	"database/sql"
	"fmt"
	"utility/internal/consts"
	"utility/internal/entities"
	"utility/utilities"

	"github.com/google/uuid"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils"
)

// GenreRepo is responsible for handling genre-related data.
type GenreRepo struct {
	db *sql.DB
}

// GenreRepoImply is the interface defining the methods for working with genre data.
type GenreRepoImply interface {
	GetGenres(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errs map[string]models.ErrorResponse) ([]*entities.Genre, int64, error)
	CreateGenre(ctx context.Context, genreName, langIdentifier string) error
	UpdateGenre(ctx context.Context, genre entities.Genre, langIdentifier string) error
	DeleteGenre(ctx context.Context, id uuid.UUID) (int64, error)
	GenreNameExists(ctx context.Context, role string) (bool, error)
	IsDuplicateExists(ctx context.Context, fieldName string, value string, roleID string) (bool, error)
	IsGenreExists(ctx context.Context, id uuid.UUID) error
	GetGenresByID(ctx context.Context, id string) (entities.GenreDetails, error)
}

// NewGenreRepo creates a new GenreRepo instance.
func NewGenreRepo(db *sql.DB) GenreRepoImply {
	return &GenreRepo{db: db}
}

// GetGenres retrieves genre data based on specified parameters.
func (genreRepo *GenreRepo) GetGenres(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errs map[string]models.ErrorResponse) ([]*entities.Genre, int64, error) {

	var (
		genres       []*entities.Genre
		totalRecords int64
		log          = logger.Log().WithContext(ctx)
	)
	getGenres := `
		SELECT
			id AS genre_id,
			name AS name,
			SUM(COUNT(DISTINCT id)) OVER() AS total_records
		FROM
			genre
	`

	if !utils.IsEmpty(params.Name) {
		getGenres = utilities.Search(getGenres, params.Name, "name")
	} else {
		if len(params.IdList) != 0 {
			getGenres = utilities.SearchIdList(getGenres, params.IdList, "id")
		}
	}

	getGenres = utilities.GroupBy(getGenres, "id, name")
	sortQ, err := utilities.OrderBy(
		ctx,
		params.Sort,
		params.Order,
		consts.SortOptns,
		validation.Endpoint,
		validation.Method,
		errs,
	)
	if err != nil {
		log.Errorf("[GenreRepo][GetGenres] Error : %s", err.Error())
		return nil, 0, err
	}
	getGenres = fmt.Sprintf("%s %s", getGenres, sortQ)
	if len(params.IdList) == 0 {
		getGenres = fmt.Sprintf("%s %s", getGenres, utils.CalculateOffset(pagination.Page, pagination.Limit))
	}

	if len(errs) > 0 {
		return nil, 0, nil
	}
	res, err := genreRepo.db.QueryContext(ctx, getGenres)

	if err != nil {
		log.Errorf("[GenreRepo][GetGenres], Error : %s", err.Error())
		return nil, 0, err
	}
	defer res.Close()
	for res.Next() {
		var genre entities.Genre
		err := res.Scan(
			&genre.ID,
			&genre.Name,
			&totalRecords,
		)

		if err != nil {
			log.Errorf("[GenreRepo][GetGenres], Error : %s", err.Error())
			return nil, 0, err
		}

		genres = append(genres, &genre)
	}

	return genres, totalRecords, nil
}

// CreateGenre creates a new genre with the provided name and language identifier.
func (genreRepo *GenreRepo) CreateGenre(ctx context.Context, genreName, langIdentifier string) error {

	log := logger.Log().WithContext(ctx)

	createGenre := `
		INSERT INTO genre(name, language_label_identifier)
		VALUES ($1, $2) RETURNING id;
	`
	var newRoleId uuid.UUID
	err := genreRepo.db.QueryRowContext(ctx, createGenre, genreName, langIdentifier).Scan(&newRoleId)
	if err != nil {
		log.Errorf("[GenreRepo][CreateGenre], Error : %s", err.Error())
		return err // Error while creating the role or fetching its new ID
	}

	// Fetch all store_ingestion_spec IDs
	fetchStores := `
			SELECT si.id
			FROM store_ingestion_spec si
			JOIN store s ON si.store_id = s.id
			WHERE s.is_active = true;
    `
	rows, err := genreRepo.db.QueryContext(ctx, fetchStores)
	if err != nil {
		return err // Error while fetching store ingestion spec IDs
	}
	defer rows.Close()

	// Prepare the insert statement for mapping role to stores
	insertMapping := `
        INSERT INTO store_genre(genre_id, store_ingestion_spec_id)
        VALUES ($1, $2);
    `
	for rows.Next() {
		var storeIngestionSpecId int64
		if err := rows.Scan(&storeIngestionSpecId); err != nil {
			return err
		}

		// Map the new role ID with each store_ingestion_spec_id
		if _, err := genreRepo.db.ExecContext(ctx, insertMapping, newRoleId, storeIngestionSpecId); err != nil {
			return err
		}
	}

	return nil
}

// DeleteGenre deletes a genre with the specified ID.
func (genreRepo *GenreRepo) DeleteGenre(ctx context.Context, id uuid.UUID) (int64, error) {
	var (
		rowsAffected int64
		log          = logger.Log().WithContext(ctx)
		isExist      bool
	)

	err := utilities.ExecuteTx(ctx, genreRepo.db, func(tx *sql.Tx) error {

		sqlst := `
		SELECT EXISTS(SELECT 1 FROM track_genre WHERE genre_id = $1)`

		err := genreRepo.db.QueryRowContext(ctx, sqlst, id).Scan(&isExist)
		if err != nil {
			log.Errorf("[GenreRepo][DeleteGenre], Error : %s", err.Error())
			return err
		}

		if isExist {

			//update from genre
			updateGenre := `
				UPDATE genre 
				SET 
					is_deleted = true  
				WHERE id = $1
	`
			rows, err := genreRepo.db.ExecContext(ctx, updateGenre, id)
			if err != nil {
				log.Errorf("[GenreRepo][DeleteGenre], Error : %s", err.Error())
				return err
			}
			rowsAffected, err = rows.RowsAffected()
			if err != nil {
				log.Errorf("[GenreRepo][DeleteGenre], Error : %s", err.Error())
				return err
			}

			return nil
		}

		//delete from partner genres
		_, err = genreRepo.db.ExecContext(ctx,
			`	DELETE FROM partner_genre_language
				WHERE genre_id = $1`,
			id,
		)
		if err != nil {
			log.Errorf("[GenreRepo][DeleteGenre]: error from database, Error : %s", err.Error())
			return err
		}

		//delete from store_genre genre
		_, err = genreRepo.db.ExecContext(ctx,
			`	DELETE FROM store_genre
					WHERE genre_id = $1`,
			id,
		)
		if err != nil {
			log.Errorf("[GenreRepo][DeleteGenre]: error from database, Error : %s", err.Error())
			return err
		}

		//delete from genre
		delGenre := `
			DELETE FROM genre
			WHERE id = $1
		`
		rows, err := genreRepo.db.ExecContext(ctx, delGenre, id)
		if err != nil {
			log.Errorf("[GenreRepo][DeleteGenre]: error from database, Error : %s", err.Error())
			return err
		}

		rowsAffected, err = rows.RowsAffected()
		if err != nil {
			log.Errorf("[GenreRepo][DeleteGenre] Error : %s", err.Error())
			return err
		}

		return nil
	})

	return rowsAffected, err
}

// UpdateGenre updates an existing genre with the provided data.
func (genreRepo *GenreRepo) UpdateGenre(ctx context.Context, genre entities.Genre, langIdentifier string) error {

	log := logger.Log().WithContext(ctx)

	updtGenre := `
		UPDATE genre
		SET  name=$1, language_label_identifier=$2
		WHERE id=$3;
	`
	_, err := genreRepo.db.ExecContext(ctx, updtGenre, genre.Name, langIdentifier, genre.ID)
	if err != nil {
		log.Error("[GenreRepo][UpdateGenre], Error : %s", err)
		return err
	}
	//rowsAffected, err := res.RowsAffected()
	return nil
}

func (genreRepo *GenreRepo) IsGenreExists(ctx context.Context, id uuid.UUID) error {
	var (
		log          = logger.Log().WithContext(ctx)
		isGenreExist bool
	)
	query := `SELECT EXISTS(SELECT 1 FROM genre WHERE id = $1 AND is_deleted = false)`
	err := genreRepo.db.QueryRowContext(ctx, query, id).Scan(&isGenreExist)

	if err != nil {
		log.Errorf("[GenreRepo][DeleteGenre], Error : %s", err.Error())
		return err
	}
	if !isGenreExist {
		log.Errorf("[GenreRepo][DeleteGenre], Error =  genre already deleted.")
		return fmt.Errorf("genre already deleted")
	}
	return nil

}

func (genreRepo *GenreRepo) GenreNameExists(ctx context.Context, role string) (bool, error) {

	var isExist bool
	// Construct the SQL query with the dynamic field name.
	sqlst := `
	SELECT EXISTS(SELECT * FROM genre WHERE LOWER(name) = LOWER($1))
    `

	// Execute the query with the specified role value.
	err := genreRepo.db.QueryRowContext(ctx, sqlst, role).Scan(&isExist)

	return isExist, err
}

func (genreRepo *GenreRepo) IsDuplicateExists(ctx context.Context, fieldName string, value string, genreID string) (bool, error) {

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM  genre  WHERE LOWER(` + fieldName + `) = LOWER($1) AND  id != $2)`

	err := genreRepo.db.QueryRowContext(ctx, query, value, genreID).Scan(&exists)

	return exists, err
}

func (genreRepo *GenreRepo) GetGenresByID(ctx context.Context, id string) (entities.GenreDetails, error) {

	var (
		res entities.GenreDetails
		log = logger.Log().WithContext(ctx)
	)

	query := `
		SELECT 
			name,
			language_label_identifier,
			is_deleted
		FROM genre 
		WHERE id= $1
		`

	row := genreRepo.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&res.Name, &res.LanguageLabelIdentifier, &res.IsDeleted)

	if err != nil {
		log.Errorf("[GenreRepo][GetGenresByID], Error : %s", err.Error())
		return entities.GenreDetails{}, err
	}
	return res, nil

}
