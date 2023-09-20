package ports

import (
	"context"
	"fmt"

	"github.com/mrparano1d/getregd/pkg/core/entities"
	"github.com/mrparano1d/getregd/pkg/core/fields"
)

type AuthPort interface {
	Login(ctx context.Context, username fields.Username, password fields.Password) (*entities.User, error)
}

// errors

type AuthAdapterLoginFailedError struct {
	Username fields.Username
	Err      error
}

func (e *AuthAdapterLoginFailedError) Error() string {
	return fmt.Sprintf("login failed for user %s: %s", e.Username, e.Err)
}

type AuthAdapterUserNotFoundError struct {
	Username fields.Username
}

func (e *AuthAdapterUserNotFoundError) Error() string {
	return fmt.Sprintf("user %s not found", e.Username)
}

type AuthAdapterInvalidCredentialsError struct {
	Username fields.Username
}

func (e *AuthAdapterInvalidCredentialsError) Error() string {
	return fmt.Sprintf("invalid credentials for user %s", e.Username)
}
