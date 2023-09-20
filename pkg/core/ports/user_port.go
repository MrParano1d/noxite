package ports

import (
	"context"

	"github.com/mrparano1d/getregd/pkg/core/entities"
	"github.com/mrparano1d/getregd/pkg/core/fields"
)

type CreateUserInput struct {
	Username fields.Username
	Email    fields.Email
	Password fields.Password
}

type UpdateUserInput struct {
	Username *fields.Username
	Email    *fields.Email
	Password *fields.Password
}

type UserPort interface {
	CreateUser(ctx context.Context, createUser CreateUserInput) (fields.EntityID, error)
	GetUserByID(ctx context.Context, userID fields.EntityID) (*entities.User, error)
	UpdateUser(ctx context.Context, id fields.EntityID, updateUser UpdateUserInput) error
	DeleteUser(ctx context.Context, userID fields.EntityID) error
}
