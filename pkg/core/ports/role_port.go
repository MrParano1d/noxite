package ports

import (
	"context"
	"fmt"

	"github.com/mrparano1d/noxite/pkg/core/entities"
	"github.com/mrparano1d/noxite/pkg/core/fields"
)

type RolePort interface {
	// CreateRole creates a new role.
	// Returns the ID of the created role.
	// Returns RoleAdapterCreateRoleFailedError if failed to create role.
	// Returns RoleAdapterRoleAlreadyExistsError if role with the same name already exists.
	CreateRole(ctx context.Context, createRole CreateRoleInput) (fields.EntityID, error)
	// GetAllRoles returns all roles.
	// Returns RoleAdapterGetAllRolesFailedError if failed to get all roles.
	GetAllRoles(ctx context.Context) ([]*entities.Role, error)
	// GetRoleByID returns the role with the given ID.
	// Returns RoleAdapterGetRoleByIDFailedError if failed to get role.
	// Returns RoleAdapterRoleNotFoundError if role with the given ID does not exist.
	GetRoleByID(ctx context.Context, roleID fields.EntityID) (*entities.Role, error)
	// UpdateRole updates the role with the given ID.
	// Returns RoleAdapterUpdateRoleFailedError if failed to update role.
	// Returns RoleAdapterRoleNotFoundError if role with the given ID does not exist.
	UpdateRole(ctx context.Context, id fields.EntityID, updateRole UpdateRoleInput) error
	// DeleteRole deletes the role with the given ID.
	// Returns RoleAdapterDeleteRoleFailedError if failed to delete role.
	// Returns RoleAdapterRoleNotFoundError if role with the given ID does not exist.
	DeleteRole(ctx context.Context, roleID fields.EntityID) error
}

// inputs

type CreateRoleInput struct {
	Name        fields.RequiredString
	Description string
	Permissions entities.Permissions
}

type UpdateRoleInput struct {
	Name        *fields.RequiredString
	Permissions *entities.Permissions
	Description *string
}

// errors

type RoleAdapterCreateRoleFailedError struct {
	Err error
}

func (e RoleAdapterCreateRoleFailedError) Error() string {
	return fmt.Sprintf("failed to create role: %v", e.Err)
}

type RoleAdapterRoleAlreadyExistsError struct {
	Name fields.RequiredString
}

func (e RoleAdapterRoleAlreadyExistsError) Error() string {
	return fmt.Sprintf("role with name %q already exists", e.Name)
}

type RoleAdapterGetAllRolesFailedError struct {
	Err error
}

func (e RoleAdapterGetAllRolesFailedError) Error() string {
	return fmt.Sprintf("failed to get all roles: %v", e.Err)
}

type RoleAdapterGetRoleByIDFailedError struct {
	ID  fields.EntityID
	Err error
}

func (e RoleAdapterGetRoleByIDFailedError) Error() string {
	return fmt.Sprintf("failed to get role by id %q: %v", e.ID, e.Err)
}

type RoleAdapterRoleNotFoundError struct {
	ID fields.EntityID
}

func (e RoleAdapterRoleNotFoundError) Error() string {
	return fmt.Sprintf("role with id %q not found", e.ID)
}

type RoleAdapterUpdateRoleFailedError struct {
	ID  fields.EntityID
	Err error
}

func (e RoleAdapterUpdateRoleFailedError) Error() string {
	return fmt.Sprintf("failed to update role with id %q: %v", e.ID, e.Err)
}

type RoleAdapterDeleteRoleFailedError struct {
	ID  fields.EntityID
	Err error
}

func (e RoleAdapterDeleteRoleFailedError) Error() string {
	return fmt.Sprintf("failed to delete role with id %q: %v", e.ID, e.Err)
}
