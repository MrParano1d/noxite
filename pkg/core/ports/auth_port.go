package ports

import (
	"fmt"

	"github.com/mrparano1d/getregd/pkg/core/entities"
	"github.com/mrparano1d/getregd/pkg/core/fields"
)

type AuthPort interface {
	Login(username fields.Username, password fields.Password) (*entities.User, error)
	Register(username fields.Username, email fields.Email, password fields.Password) (*entities.User, error)
	GenerateToken(user *entities.User) (fields.RequiredString, error)
}

// errors

type AuthAdapterLoginFailedError struct {
	Username fields.Username
	Err      error
}

func (e *AuthAdapterLoginFailedError) Error() string {
	return fmt.Sprintf("login failed for user %s: %s", e.Username, e.Err)
}

type AuthAdapterRegistrationFailedError struct {
	Username fields.Username
	Email    fields.Email
	Err      error
}

func (e *AuthAdapterRegistrationFailedError) Error() string {
	return fmt.Sprintf("registration failed for user %s (%s): %s", e.Username, e.Email, e.Err)
}

type AuthAdapterUserAlreadyExistsError struct {
	Username fields.Username
	Email    fields.Email
}

func (e *AuthAdapterUserAlreadyExistsError) Error() string {
	return fmt.Sprintf("user %s (%s) already exists", e.Username, e.Email)
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
