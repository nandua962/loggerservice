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

// RoleRepo is responsible for handling role-related data.
type RoleRepo struct {
	db *sql.DB
}

// RoleRepoImply is the interface defining the methods for working with role data.
type RoleRepoImply interface {
	GetRoles(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errs map[string]models.ErrorResponse) ([]*entities.Role, int64, error)
	GetRoleByID(ctx context.Context, id string) (entities.Role, error)
	DeleteRoles(ctx context.Context, id uuid.UUID) (int64, error)
	CreateRole(ctx context.Context, role entities.Role, langIdentifier string) error
	UpdateRole(ctx context.Context, role entities.Role, langIdentifier string) (int64, error)
	RoleNameExists(ctx context.Context, role string) (bool, error)
	IsDuplicateExists(ctx context.Context, fieldName string, value string, roleID string) (bool, error)
	IsRoleExists(ctx context.Context, id uuid.UUID) error
}

// NewRoleRepo creates a new RoleRepo instance.
func NewRoleRepo(db *sql.DB) RoleRepoImply {
	return &RoleRepo{db: db}
}

func (roleRepo *RoleRepo) GetRoleByID(ctx context.Context, id string) (entities.Role, error) {

	var (
		res entities.Role
		log = logger.Log().WithContext(ctx)
	)

	query := `
		SELECT 
			id,
			name,
			language_label_identifier,
			custom_name,
			is_default
		FROM role 
		WHERE id= $1
		`

	row := roleRepo.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&res.ID, &res.Name, &res.LanguageLabelIdentifier, &res.CustomName, &res.IsDefault)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Errorf("[RoleRepo][GetRoleByID], No Content, Error : %s", err.Error())
			return entities.Role{}, err
		}
		log.Errorf("[RoleRepo][GetRoleByID], Error : %s", err.Error())
		return entities.Role{}, err
	}
	return res, nil

}

// CreateRole creates a new role with the provided name and language identifier.
func (roleRepo *RoleRepo) CreateRole(ctx context.Context, role entities.Role, langIdentifier string) error {

	log := logger.Log().WithContext(ctx)
	createRole := `
		INSERT INTO role(name, language_label_identifier,custom_name)
		VALUES ($1, $2, $3)
		RETURNING id;
	`
	var newRoleId uuid.UUID
	err := roleRepo.db.QueryRowContext(ctx, createRole, role.Name, langIdentifier, role.CustomName).Scan(&newRoleId)
	if err != nil {
		log.Errorf("[RoleRepo][CreateRole], Error : %s", err.Error())
		return err // Error while creating the role or fetching its new ID
	}

	// Fetch all store_ingestion_spec IDs
	fetchStores := `
        SELECT id FROM store_ingestion_spec;
    `
	rows, err := roleRepo.db.QueryContext(ctx, fetchStores)
	if err != nil {
		log.Errorf("[RoleRepo][CreateRole], Error : %s", err.Error())
		return err // Error while fetching store ingestion spec IDs
	}
	defer rows.Close()

	// Prepare the insert statement for mapping role to stores
	insertMapping := `
        INSERT INTO store_role(role_id, store_ingestion_spec_id)
        VALUES ($1, $2);
    `
	for rows.Next() {
		var storeIngestionSpecId int64
		if err := rows.Scan(&storeIngestionSpecId); err != nil {
			return err
		}

		// Map the new role ID with each store_ingestion_spec_id
		if _, err := roleRepo.db.ExecContext(ctx, insertMapping, newRoleId, storeIngestionSpecId); err != nil {
			log.Errorf("[RoleRepo][CreateRole], Error : %s", err.Error())
			return err
		}
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return err
	}

	return err
}

// GetRoles retrieves role data based on specified parameters.
func (roleRepo *RoleRepo) GetRoles(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errs map[string]models.ErrorResponse) ([]*entities.Role, int64, error) {

	var (
		roles        []*entities.Role
		totalRecords int64
		log          = logger.Log().WithContext(ctx)
	)
	getRoles := `
		SELECT
			id,
			name ,
			COALESCE(custom_name, '') AS custom_name,
			language_label_identifier ,
			is_default ,
			COUNT(id) OVER() 
		FROM
			role
	`

	if !utils.IsEmpty(params.Name) {
		getRoles = utilities.Search(getRoles, params.Name, "name")
	}
	getRoles = utilities.GroupBy(getRoles, "id", "name")
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
		log.Errorf("[RoleRepo][GetRoles] Error : %s", err.Error())
		return nil, 0, err
	}
	getRoles = fmt.Sprintf("%s %s", getRoles, sortQ)
	getRoles = fmt.Sprintf("%s %s", getRoles, utils.CalculateOffset(pagination.Page, pagination.Limit))
	if len(errs) > 0 {
		return nil, 0, nil
	}
	res, err := roleRepo.db.QueryContext(ctx, getRoles)
	if err != nil {
		log.Errorf("[RoleRepo][GetRoles], Error : %s", err.Error())
		return nil, 0, err
	}

	defer res.Close()
	for res.Next() {
		var role entities.Role
		err := res.Scan(
			&role.ID,
			&role.Name,
			&role.CustomName,
			&role.LanguageLabelIdentifier,
			&role.IsDefault,
			&totalRecords,
		)

		if err != nil {
			log.Errorf("[RoleRepo][GetRoles], Error : %s", err.Error())
			return nil, 0, err
		}

		roles = append(roles, &role)
	}

	return roles, totalRecords, nil
}

// UpdateRole updates an existing role with the provided data.
func (roleRepo *RoleRepo) UpdateRole(ctx context.Context, role entities.Role, langIdentifier string) (int64, error) {

	log := logger.Log().WithContext(ctx)

	updtRole := `
		UPDATE role
		SET  name=$1, language_label_identifier=$2, custom_name=$3
		WHERE id=$4;
	`
	res, err := roleRepo.db.ExecContext(ctx, updtRole, role.Name, langIdentifier, role.CustomName, role.ID)
	if err != nil {
		log.Error("[RoleRepo][UpdateRole], Error : %s", err)
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	return rowsAffected, err
}

// DeleteRoles deletes a role with the specified ID.
func (roleRepo *RoleRepo) DeleteRoles(ctx context.Context, id uuid.UUID) (int64, error) {
	var (
		rowsAffected int64
		log          = logger.Log().WithContext(ctx)
		isExist      bool
	)

	err := utilities.ExecuteTx(ctx, roleRepo.db, func(tx *sql.Tx) error {

		sqlst := `
		SELECT EXISTS(SELECT 1 FROM track_artist_role WHERE role_id = $1)`

		err := roleRepo.db.QueryRowContext(ctx, sqlst, id).Scan(&isExist)
		if err != nil {
			log.Errorf("[RoleRepo][DeleteRoles], Error : %s", err.Error())
			return err
		}

		if isExist {
			//update from role
			updateRole := `
				UPDATE role 
				SET 
					is_deleted = true  
				WHERE id = $1
	`
			rows, err := roleRepo.db.ExecContext(ctx, updateRole, id)
			if err != nil {
				log.Errorf("[RoleRepo][DeleteRoles], Error : %s", err.Error())
				return err
			}
			rowsAffected, err = rows.RowsAffected()
			if err != nil {
				log.Errorf("[RoleRepo][DeleteRoles], Error : %s", err.Error())
				return err
			}

			return nil
		}

		//delete from partner_artist_role_language roles
		_, err = roleRepo.db.ExecContext(ctx,
			`	DELETE FROM partner_artist_role_language
					WHERE role_id = $1`,
			id,
		)
		if err != nil {
			log.Errorf("[RoleRepo][DeleteRoles], Error : %s", err.Error())
			return err
		}

		//delete from store_role roles
		_, err = roleRepo.db.ExecContext(ctx,
			`	DELETE FROM store_role
					WHERE role_id = $1`,
			id,
		)
		if err != nil {
			log.Errorf("[RoleRepo][DeleteRoles], Error : %s", err.Error())
			return err
		}

		//delete from role
		delRole := `
			DELETE FROM role
			WHERE id = $1
		`
		rows, err := roleRepo.db.ExecContext(ctx, delRole, id)
		if err != nil {
			log.Errorf("[RoleRepo][DeleteRoles], Error : %s", err.Error())
			return err
		}

		rowsAffected, err = rows.RowsAffected()
		if err != nil {
			log.Errorf("[RoleRepo][DeleteRoles], Error : %s", err.Error())
			return err
		}

		return nil
	})

	return rowsAffected, err
}

func (roleRepo *RoleRepo) IsRoleExists(ctx context.Context, id uuid.UUID) error {
	var (
		log         = logger.Log().WithContext(ctx)
		isRoleExist bool
	)
	query := `SELECT EXISTS(SELECT 1 FROM role WHERE id = $1 AND is_deleted = false)`
	err := roleRepo.db.QueryRowContext(ctx, query, id).Scan(&isRoleExist)

	if err != nil {
		log.Errorf("[RoleRepo][UpdateRole], Error : %s", err.Error())
		return err
	}
	if !isRoleExist {
		log.Errorf("[RoleRepo][UpdateRole], Error :  role already deleted.")
		return fmt.Errorf("role already deleted")
	}
	return nil

}

// RoleExists retrieves the count of roles with the specified value in the given field from the database.
func (roleRepo *RoleRepo) RoleNameExists(ctx context.Context, role string) (bool, error) {

	var isExist bool
	// Construct the SQL query with the dynamic field name.
	sqlst := `
	SELECT EXISTS(SELECT * FROM role WHERE LOWER(name) = LOWER($1))
    `
	// Execute the query with the specified role value.
	err := roleRepo.db.QueryRowContext(ctx, sqlst, role).Scan(&isExist)

	return isExist, err
}

func (roleRepo *RoleRepo) IsDuplicateExists(ctx context.Context, fieldName string, value string, roleID string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM  role  WHERE LOWER(` + fieldName + `) = LOWER($1) AND  id != $2)`

	err := roleRepo.db.QueryRowContext(ctx, query, value, roleID).Scan(&exists)

	return exists, err
}
