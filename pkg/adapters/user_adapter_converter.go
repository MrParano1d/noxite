package adapters

import (
	"github.com/mrparano1d/getregd/ent"
	"github.com/mrparano1d/getregd/pkg/core/entities"
	"github.com/mrparano1d/getregd/pkg/core/fields"
)

func UserFromEntUser(user *ent.User) (*entities.User, error) {
	username, err := fields.UsernameFromString(user.Name)
	if err != nil {
		return nil, err
	}

	email, err := fields.EmailFromString(user.Email)
	if err != nil {
		return nil, err
	}

	password, err := fields.PasswordFromBytes(user.Password)
	if err != nil {
		return nil, err
	}

	id, err := fields.EntityIDFromInt(user.ID)
	if err != nil {
		return nil, err
	}

	role, err := RoleFromEntRole(user.Edges.Role)
	if err != nil {
		return nil, err
	}

	return &entities.User{
		ID:        id,
		Username:  username,
		Email:     email,
		Password:  password,
		Role:      role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}, nil
}
