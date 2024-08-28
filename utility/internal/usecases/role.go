package usecases

import (
	"context"
	"utility/internal/consts"
	"utility/internal/entities"
	"utility/internal/repo"
	"utility/utilities"

	"github.com/google/uuid"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/utils"
)

// RoleUseCases represents use cases for handling role-related operations.
type RoleUseCases struct {
	repo repo.RoleRepoImply
}

// RoleUseCaseImply is an interface defining the methods for working with role use cases.
type RoleUseCaseImply interface {
	GetRoles(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errMap map[string]models.ErrorResponse) (*entities.Response, map[string]models.ErrorResponse, error)
	GetRoleByID(ctx context.Context, validation entities.Validation, errMap map[string]models.ErrorResponse) (entities.Role, map[string]models.ErrorResponse, error)
	DeleteRoles(ctx context.Context, validation entities.Validation, errMap map[string]models.ErrorResponse) (map[string]models.ErrorResponse, error)
	CreateRole(ctx context.Context, req entities.Role, validation entities.Validation, errMap map[string]models.ErrorResponse) (map[string]models.ErrorResponse, error)
	UpdateRole(ctx context.Context, roleInfo entities.Role, validation entities.Validation, errMap map[string]models.ErrorResponse) (map[string]models.ErrorResponse, error)
}

// NewRoleUseCases creates a new RoleUseCases instance.
func NewRoleUseCases(RoleRepo repo.RoleRepoImply) RoleUseCaseImply {
	return &RoleUseCases{
		repo: RoleRepo,
	}
}

// GetRoleByID to retreive role details
func (role *RoleUseCases) GetRoleByID(ctx context.Context, validation entities.Validation, errMap map[string]models.ErrorResponse) (entities.Role, map[string]models.ErrorResponse, error) {

	log := logger.Log().WithContext(ctx)

	_, err := uuid.Parse(validation.ID)

	if err != nil {

		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.RoleIdField, consts.InvalidKey)
		if err != nil {
			log.Errorf("[GetRoleByID],Error while loading service code, Error : %s", err)
		}

		errMap[consts.RoleIdField] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}

		return entities.Role{}, errMap, nil
	}

	if len(errMap) > 0 {
		return entities.Role{}, errMap, nil
	}

	result, err := role.repo.GetRoleByID(ctx, validation.ID)

	if err != nil {
		log.Errorf("[RoleUseCases][GetRoleByID] Error : %s", err.Error())
		return entities.Role{}, nil, err
	}

	return result, nil, nil

}

// CreateRole handles the creation of a role and does necessary validations.
func (role *RoleUseCases) CreateRole(ctx context.Context, req entities.Role, validation entities.Validation, errMap map[string]models.ErrorResponse) (map[string]models.ErrorResponse, error) {

	log := logger.Log().WithContext(ctx)

	if utils.IsEmpty(req.Name) {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.NameKey, consts.Required)
		if err != nil {
			log.Errorf("[RoleUseCases][CreateRole] Error while loading service code, Error : %s", err.Error())
		}

		errMap[consts.NameKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.Required},
		}

	} else if err := utils.ValidateStringLength(req.Name, consts.MinimumLength, consts.MaximumLength); err != nil {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.NameKey, "length")
		if err != nil {
			log.Errorf("[RoleUseCases][CreateRole][ValidateStringLength] Error while loading service code, Error : %s", err.Error())
		}

		errMap[consts.NameKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{"length"},
		}

		return errMap, nil

	} else {
		isExist, err := role.repo.RoleNameExists(ctx, req.Name)

		if err != nil {
			log.Errorf("[RoleUseCases][CreateRole][RoleNameExists], Error : %s", err.Error())
			return nil, err
		}
		if isExist {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.NameKey, consts.RoleExists)
			if err != nil {
				log.Errorf("[RoleUseCases][CreateRole][RoleNameExists] Error while loading service code, Error : %s", err.Error())
			}

			errMap[consts.NameKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.RoleExists},
			}

		}
	}

	if len(errMap) > 0 {
		log.Errorf("[RoleUseCases][CreateRole] ValidationError")
		return errMap, nil
	}

	langIdentifier, err := utilities.GenerateLangIdentifier(req.Name, consts.Role)

	if err != nil {
		log.Errorf("[RoleUseCases][CreateRole][RoleNameExists], GenerateLangIdentifier err=%s", err.Error())
		return nil, err
	}

	if utils.IsEmpty(req.CustomName) {
		req.CustomName = req.Name
	}

	err = role.repo.CreateRole(ctx, req, langIdentifier)
	if err != nil {
		log.Errorf("[CreateRole], Error : %s", err.Error())
		return nil, err
	}
	return nil, nil
}

// GetRoles retrieves role data based on specified parameters.
func (role *RoleUseCases) GetRoles(ctx context.Context, params entities.Params, pagination entities.Pagination, validation entities.Validation, errMap map[string]models.ErrorResponse) (*entities.Response, map[string]models.ErrorResponse, error) {

	log := logger.Log().WithContext(ctx)

	roles, totalRecords, err := role.repo.GetRoles(ctx, params, pagination, validation, errMap)
	if len(errMap) > 0 {
		log.Errorf("[RoleUseCases][GetRoles] ValidationError")
		return nil, errMap, nil
	}

	if err != nil {
		log.Errorf("[GetRoles], Error : %s", err.Error())
		return nil, nil, err
	}

	metaData := utils.MetaDataInfo(&models.MetaData{
		Total:       totalRecords,
		PerPage:     pagination.Limit,
		CurrentPage: pagination.Page,
	})

	resp := &entities.Response{
		MetaData: metaData,
		Data:     roles,
	}
	return resp, nil, nil

}

// UpdateRole handles the updating of a role and does necessary validattions for the input data.
func (role *RoleUseCases) UpdateRole(ctx context.Context, roleInfo entities.Role, validation entities.Validation, errMap map[string]models.ErrorResponse) (map[string]models.ErrorResponse, error) {

	log := logger.Log().WithContext(ctx)
	id := validation.ID

	uuid, err := uuid.Parse(id)

	if err != nil {
		log.Errorf("[UpdateRole]: Parsing uuid failed,  Error : %s", err.Error())
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.CommonMsg, consts.InvalidRole)
		if err != nil {
			log.Errorf("[RoleUseCases][UpdateRole] Error while loading service code, Error : %s", err.Error())
		}

		errMap[consts.CommonMsg] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidRole},
		}

		return errMap, nil
	}

	err = role.repo.IsRoleExists(ctx, uuid)
	if err != nil {
		if err.Error() == "role already deleted" {
			log.Errorf("[RoleUseCases.UpdateRole] UpdateRole failed, validation error,  Error : %s", err)
			return nil, consts.ErrNotExist
		}
		log.Errorf("[UpdateRole], Error : %s", err.Error())
		return nil, err

	}

	if utils.IsEmpty(roleInfo.Name) {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.NameKey, consts.Required)
		if err != nil {
			log.Errorf("[RoleUseCases][UpdateRole] Error while loading service code, Error : %s", err.Error())
		}

		errMap[consts.NameKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.Required},
		}
	} else if err := utils.ValidateStringLength(roleInfo.Name, consts.MinimumLength, consts.MaximumLength); err != nil {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.NameKey, consts.LengthKey)
		if err != nil {
			log.Errorf("[RoleUseCases][UpdateRole] Error while loading service code, Error : %s", err.Error())
		}

		errMap[consts.NameKey] = models.ErrorResponse{
			Code:    code,
			Message: []string{"length"},
		}

		return errMap, nil
	} else {
		isNameExists, err := role.repo.IsDuplicateExists(ctx, "name", roleInfo.Name, id)

		if err != nil {
			log.Errorf("[RoleUseCases][UpdateRole][IsDuplicateExists], Error : %s", err.Error())
			return nil, err
		}

		if isNameExists {
			log.Errorf("[RoleUseCases][UpdateRole][IsDuplicateExists], Name already exists")
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.NameKey, consts.RoleExists)
			if err != nil {
				log.Errorf("[RoleUseCases][UpdateRole][IsDuplicateExists] Error while loading service code, Error : %s", err.Error())
			}

			errMap[consts.NameKey] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.RoleExists},
			}
		}
	}

	if len(errMap) > 0 {
		log.Errorf("[RoleUseCases][UpdateRole] ValidationError")
		return errMap, nil
	}

	if utils.IsEmpty(roleInfo.CustomName) {
		roleInfo.CustomName = roleInfo.Name
	}

	langIdentifier, err := utilities.GenerateLangIdentifier(roleInfo.Name, consts.Role)
	if err != nil {
		log.Errorf("[UpdateRole][GenerateLangIdentifier], Error : %s", err.Error())
		return nil, err
	}

	roleInfo.ID = id

	rowsAffected, err := role.repo.UpdateRole(ctx, roleInfo, langIdentifier)
	if err != nil {
		log.Errorf("[UpdateRole], Error : %s", err.Error())
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, consts.ErrNotFound
	}

	return errMap, nil
}

// DeleteRoles handles the deletion of a role.
func (role *RoleUseCases) DeleteRoles(ctx context.Context, validation entities.Validation, errMap map[string]models.ErrorResponse) (map[string]models.ErrorResponse, error) {

	log := logger.Log().WithContext(ctx)

	uuid, err := uuid.Parse(validation.ID)

	if err != nil {
		log.Errorf("[DeleteRoles], Parsing uuid failed, Error : %s", err.Error())
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, consts.CommonMsg, consts.InvalidRole)
		if err != nil {
			log.Errorf("[RoleUseCases][DeleteRoles] Error while loading service code, Error : %s", err.Error())
		}

		errMap[consts.CommonMsg] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidRole},
		}

		return errMap, nil
	}

	err = role.repo.IsRoleExists(ctx, uuid)
	if err != nil {

		if err.Error() == "role already deleted" {
			log.Errorf("[RoleUseCases.DeleteRoles] DeleteRoles failed, validation error, Error : %s", err.Error())
			return nil, consts.ErrNotExist
		}
		log.Errorf("[DeleteRoles][IsRoleExists], Error : %s", err.Error())
		return nil, err

	}
	rowsAffected, err := role.repo.DeleteRoles(ctx, uuid)

	if err != nil {
		log.Errorf("[DeleteRoles], Error : %s", err.Error())
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, consts.ErrNotFound
	}

	return nil, nil
}
