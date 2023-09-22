package adapters

import (
	"fmt"

	"github.com/mrparano1d/getregd/ent"
	"github.com/mrparano1d/getregd/pkg/core/entities"
	"github.com/mrparano1d/getregd/pkg/core/fields"
)

func RoleFromEntRole(role *ent.Role) (*entities.Role, error) {
	roleName, err := fields.RequiredStringFromString(role.Name)
	if err != nil {
		return nil, err
	}

	id, err := fields.EntityIDFromInt(role.ID)
	if err != nil {
		return nil, err
	}

	return &entities.Role{
		ID:          id,
		Name:        roleName,
		Description: role.Description,
		Permissions: role.Permissions,
	}, nil
}

// errors

type RoleConverterInvalidFieldError struct {
	Field  string
	Reason string
}

func (e RoleConverterInvalidFieldError) Error() string {
	return fmt.Sprintf("invalid role field while converting: %s: %s", e.Field, e.Reason)
}
