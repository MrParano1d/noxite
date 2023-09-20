package adapters

import (
	"context"
	"time"

	"github.com/mrparano1d/getregd/pkg/core/entities"
	"github.com/mrparano1d/getregd/pkg/core/fields"
	"github.com/mrparano1d/getregd/pkg/core/ports"
)

type AuthAdapter struct {
}

var _ ports.AuthPort = (*AuthAdapter)(nil)

func NewAuthAdapter() *AuthAdapter {
	return &AuthAdapter{}
}

func (a *AuthAdapter) Login(ctx context.Context, username fields.Username, password fields.Password) (*entities.User, error) {
	// TODO implement
	return &entities.User{
		ID: fields.EntityID(123),
		Role: &entities.Role{
			ID:   fields.EntityID(1),
			Name: "user",
			Permissions: entities.Permissions{
				GetPackage:       true,
				PublishPackage:   false,
				UpdatePackage:    false,
				UnpublishPackage: false,

				CreateUser: false,
				GetUser:    false,
				UpdateUser: false,
				DeleteUser: false,

				GetRole:    false,
				CreateRole: false,
				UpdateRole: false,
				DeleteRole: false,
			},
		},
		Username:  username,
		Email:     fields.Email("test@test.de"),
		Password:  password,
		CreatedAt: time.Now(),
	}, nil
}
