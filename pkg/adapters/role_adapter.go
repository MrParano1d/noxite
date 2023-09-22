package adapters

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mrparano1d/getregd/ent"
	"github.com/mrparano1d/getregd/ent/role"
	"github.com/mrparano1d/getregd/pkg/core/entities"
	"github.com/mrparano1d/getregd/pkg/core/fields"
	"github.com/mrparano1d/getregd/pkg/core/ports"
)

type RoleAdapter struct {
	entClient *ent.Client
}

var _ ports.RolePort = (*RoleAdapter)(nil)

func NewRoleAdapter(entClient *ent.Client) *RoleAdapter {
	return &RoleAdapter{entClient: entClient}
}

func (r *RoleAdapter) CreateRole(ctx context.Context, createRole ports.CreateRoleInput) (fields.EntityID, error) {

	role, err := r.entClient.Role.Create().
		SetName(createRole.Name.String()).
		SetDescription(createRole.Description).
		SetPermissions(createRole.Permissions).
		Save(ctx)

	if err != nil {
		return fields.EntityID(0), &ports.RoleAdapterCreateRoleFailedError{
			Err: err,
		}
	}

	id, err := fields.EntityIDFromInt(role.ID)
	if err != nil {
		return fields.EntityID(0), &ports.RoleAdapterCreateRoleFailedError{
			Err: fmt.Errorf("failed to convert ent.Role.ID to fields.EntityID: %w", err),
		}
	}

	return id, nil

}

func (r *RoleAdapter) GetAllRoles(ctx context.Context) ([]*entities.Role, error) {

	roles, err := r.entClient.Role.Query().Where(role.DeletedAtIsNil()).All(ctx)
	if err != nil {
		return nil, &ports.RoleAdapterGetAllRolesFailedError{
			Err: err,
		}
	}

	result := make([]*entities.Role, 0, len(roles))
	for _, role := range roles {
		role, err := RoleFromEntRole(role)
		if err != nil {
			// TODO improve error handling
			log.Printf("failed to convert ent.Role to entities.Role: %v", err)
			continue
		}
		result = append(result, role)
	}
	return result, nil
}

func (r *RoleAdapter) GetRoleByID(ctx context.Context, id fields.EntityID) (*entities.Role, error) {

	role, err := r.entClient.Role.Query().Where(role.ID(id.Int()), role.DeletedAtIsNil()).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, &ports.RoleAdapterRoleNotFoundError{ID: id}
		}
		return nil, &ports.RoleAdapterGetRoleByIDFailedError{
			ID:  id,
			Err: err,
		}
	}

	return RoleFromEntRole(role)
}

func (r *RoleAdapter) UpdateRole(ctx context.Context, id fields.EntityID, updateRole ports.UpdateRoleInput) error {

	query := r.entClient.Role.UpdateOneID(id.Int()).Where(role.DeletedAtIsNil())

	if updateRole.Name != nil {
		query = query.SetName(updateRole.Name.String())
	}

	if updateRole.Description != nil {
		query = query.SetDescription(*updateRole.Description)
	}

	if updateRole.Permissions != nil {
		query = query.SetPermissions(*updateRole.Permissions)
	}

	query = query.SetUpdatedAt(time.Now())

	_, err := query.Save(ctx)
	if err != nil {
		return &ports.RoleAdapterUpdateRoleFailedError{
			ID:  id,
			Err: err,
		}
	}
	return nil
}

func (r *RoleAdapter) DeleteRole(ctx context.Context, id fields.EntityID) error {

	_, err := r.entClient.Role.UpdateOneID(id.Int()).SetDeletedAt(time.Now()).Save(ctx)
	if err != nil {
		return &ports.RoleAdapterDeleteRoleFailedError{
			ID:  id,
			Err: err,
		}
	}
	return nil
}
