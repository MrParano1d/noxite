package services

import (
	"context"
	"fmt"

	"github.com/mrparano1d/noxite/pkg/core/coreerrors"
	"github.com/mrparano1d/noxite/pkg/core/entities"
	"github.com/mrparano1d/noxite/pkg/core/fields"
	"github.com/mrparano1d/noxite/pkg/core/ports"
)

type RoleService struct {
	adapter ports.RolePort
}

func NewRoleService(adapter ports.RolePort) *RoleService {
	return &RoleService{adapter: adapter}
}

// use cases

func (s *RoleService) CreateRole(ctx context.Context, user *entities.User, req CreateRoleRequest) (fields.EntityID, error) {

	if user == nil || user.Role.Permissions.CreateRole == false {
		return fields.EntityID(0), &coreerrors.NotAllowedToCreateRoleError{}
	}

	input, err := CreateRoleRequestToInput(req)
	if err != nil {
		return fields.EntityID(0), err
	}

	roleID, err := s.adapter.CreateRole(ctx, input)
	if err != nil {
		return fields.EntityID(0), handleRoleServiceErrors(err)
	}

	return roleID, nil
}

func (s *RoleService) GetAllRoles(ctx context.Context, user *entities.User) ([]*entities.Role, error) {

	if user == nil || user.Role.Permissions.GetRole == false {
		return nil, &coreerrors.NotAllowedToGetRoleError{}
	}

	roles, err := s.adapter.GetAllRoles(ctx)
	if err != nil {
		return nil, handleRoleServiceErrors(err)
	}
	return roles, nil
}

func (s *RoleService) GetRoleByID(ctx context.Context, user *entities.User, roleID string) (*entities.Role, error) {

	if user == nil || user.Role.Permissions.GetRole == false {
		return nil, &coreerrors.NotAllowedToGetRoleError{}
	}

	id, err := fields.EntityIDFromString(roleID)
	if err != nil {
		return nil, &RoleServiceFieldValidationError{Field: "roleID", Reason: err.Error()}
	}

	role, err := s.adapter.GetRoleByID(ctx, id)
	if err != nil {
		return nil, handleRoleServiceErrors(err)
	}
	return role, nil
}

func (s *RoleService) UpdateRole(ctx context.Context, user *entities.User, roleID string, req UpdateRoleRequest) error {

	if user == nil || user.Role.Permissions.UpdateRole == false {
		return &coreerrors.NotAllowedToUpdateRoleError{}
	}

	id, err := fields.EntityIDFromString(roleID)
	if err != nil {
		return &RoleServiceFieldValidationError{Field: "roleID", Reason: err.Error()}
	}

	input, err := UpdateRoleRequestToInput(req)
	if err != nil {
		return err
	}

	err = s.adapter.UpdateRole(ctx, id, input)
	if err != nil {
		return handleRoleServiceErrors(err)
	}

	return nil
}

func (s *RoleService) DeleteRole(ctx context.Context, user *entities.User, roleID string) error {

	if user == nil || user.Role.Permissions.DeleteRole == false {
		return &coreerrors.NotAllowedToDeleteRoleError{}
	}

	id, err := fields.EntityIDFromString(roleID)
	if err != nil {
		return &RoleServiceFieldValidationError{Field: "roleID", Reason: err.Error()}
	}

	err = s.adapter.DeleteRole(ctx, id)
	if err != nil {
		return handleRoleServiceErrors(err)
	}

	return nil
}

// requests

type CreateRoleRequest struct {
	Name        string
	Description string
	Permissions []string
}

func CreateRoleRequestToInput(req CreateRoleRequest) (ports.CreateRoleInput, error) {
	permissions := entities.Permissions{}
	permissions.FromSlice(req.Permissions)

	return ports.CreateRoleInput{
		Name:        fields.RequiredString(req.Name),
		Description: req.Description,
		Permissions: permissions,
	}, nil
}

type UpdateRoleRequest struct {
	Name        *string
	Description *string
	Permissions *[]string
}

func UpdateRoleRequestToInput(req UpdateRoleRequest) (ports.UpdateRoleInput, error) {
	var permissions *entities.Permissions
	var name *fields.RequiredString
	if req.Permissions != nil {
		p := entities.Permissions{}
		p.FromSlice(*req.Permissions)
		permissions = &p
	}

	if req.Name != nil {
		n, err := fields.RequiredStringFromString(*req.Name)
		if err != nil {
			return ports.UpdateRoleInput{}, &RoleServiceFieldValidationError{Field: "name", Reason: err.Error()}
		}
		name = &n
	}

	return ports.UpdateRoleInput{
		Name:        name,
		Description: req.Description,
		Permissions: permissions,
	}, nil
}

// errors

func handleRoleServiceErrors(err error) error {
	switch e := err.(type) {
	case *ports.RoleAdapterRoleNotFoundError:
		return RoleServiceRoleNotFoundError{ID: e.ID}
	case *ports.RoleAdapterRoleAlreadyExistsError:
		return RoleServiceRoleAlreadyExistsError{Name: e.Name}
	case *ports.RoleAdapterCreateRoleFailedError:
		return RoleServiceCreateRoleFailedError{Err: e.Err}
	case *ports.RoleAdapterGetRoleByIDFailedError:
		return RoleServiceGetRoleByIDFailedError{ID: e.ID, Err: e.Err}
	case *ports.RoleAdapterUpdateRoleFailedError:
		return RoleServiceUpdateRoleFailedError{ID: e.ID, Err: e.Err}
	case *ports.RoleAdapterDeleteRoleFailedError:
		return RoleServiceDeleteRoleFailedError{ID: e.ID, Err: e.Err}
	case *ports.RoleAdapterGetAllRolesFailedError:
		return RoleServiceGetAllRolesFailedError{Err: e.Err}
	default:
		return RoleServiceUnknownError{Err: err}
	}
}

type RoleServiceUnknownError struct {
	Err error
}

func (e RoleServiceUnknownError) Error() string {
	return fmt.Sprintf("role service unknown error: %v", e.Err)
}

type RoleServiceFieldValidationError struct {
	Field  string
	Reason string
}

func (e RoleServiceFieldValidationError) Error() string {
	return fmt.Sprintf("invalid role service field: %v: %v", e.Field, e.Reason)
}

type RoleServiceRoleNotFoundError struct {
	ID fields.EntityID
}

func (e RoleServiceRoleNotFoundError) Error() string {
	return fmt.Sprintf("role with ID %v not found", e.ID)
}

type RoleServiceRoleAlreadyExistsError struct {
	Name fields.RequiredString
}

func (e RoleServiceRoleAlreadyExistsError) Error() string {
	return fmt.Sprintf("role with name %q already exists", e.Name)
}

type RoleServiceCreateRoleFailedError struct {
	Err error
}

func (e RoleServiceCreateRoleFailedError) Error() string {
	return fmt.Sprintf("role service failed to create role: %v", e.Err)
}

type RoleServiceGetRoleByIDFailedError struct {
	ID  fields.EntityID
	Err error
}

func (e RoleServiceGetRoleByIDFailedError) Error() string {
	return fmt.Sprintf("role service failed to get role with ID %q: %v", e.ID, e.Err)
}

type RoleServiceUpdateRoleFailedError struct {
	ID  fields.EntityID
	Err error
}

func (e RoleServiceUpdateRoleFailedError) Error() string {
	return fmt.Sprintf("role service failed to update role with ID %q: %v", e.ID, e.Err)
}

type RoleServiceDeleteRoleFailedError struct {
	ID  fields.EntityID
	Err error
}

func (e RoleServiceDeleteRoleFailedError) Error() string {
	return fmt.Sprintf("role service failed to delete role with ID %q: %v", e.ID, e.Err)
}

type RoleServiceGetAllRolesFailedError struct {
	Err error
}

func (e RoleServiceGetAllRolesFailedError) Error() string {
	return fmt.Sprintf("role service failed to get all roles: %v", e.Err)
}
