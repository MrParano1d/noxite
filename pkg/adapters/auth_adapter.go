package adapters

import (
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

func (a *AuthAdapter) Login(username fields.Username, password fields.Password) (*entities.User, error) {
	// TODO implement
	return &entities.User{
		ID:        fields.EntityID(123),
		Username:  username,
		Email:     fields.Email("test@test.de"),
		Password:  password,
		CreatedAt: time.Now(),
	}, nil
}

func (a *AuthAdapter) Register(username fields.Username, email fields.Email, password fields.Password) (*entities.User, error) {
	// TODO implement
	return &entities.User{
		ID:        fields.EntityID(123),
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
	}, nil
}

func (a *AuthAdapter) GenerateToken(user *entities.User) (fields.RequiredString, error) {
	// TODO implement
	return fields.RequiredString("token"), nil
}
