package services

import (
	"context"
	"fmt"

	"github.com/mrparano1d/getregd/pkg/core/entities"
	"github.com/mrparano1d/getregd/pkg/core/fields"
	"github.com/mrparano1d/getregd/pkg/core/ports"
)

type AuthService struct {
	adapter ports.AuthPort

	sessionService *SessionService
}

func NewAuthService(
	adapter ports.AuthPort,
	sessionService *SessionService,
) *AuthService {
	return &AuthService{
		adapter:        adapter,
		sessionService: sessionService,
	}
}

// usecases

func (s *AuthService) Login(ctx context.Context, username string, email string, password string) (*entities.Session, error) {

	name, err := fields.UsernameFromString(username)
	if err != nil {
		return nil, handleErrors(err)
	}

	pw, err := fields.PasswordFromString(password)
	if err != nil {
		return nil, handleErrors(err)
	}

	user, err := s.adapter.Login(ctx, name, pw)
	if err != nil {
		return nil, handleErrors(err)
	}

	sess, err := s.sessionService.CreateSessionForUser(ctx, user)
	if err != nil {
		return nil, handleErrors(err)
	}

	return sess, nil
}


// service errors

func handleErrors(err error) error {
	switch e := err.(type) {
	case *ports.AuthAdapterLoginFailedError:
		return &AuthServiceLoginFailedError{
			Username: e.Username,
			Err:      e.Err,
		}
	case *ports.AuthAdapterUserNotFoundError:
		return &AuthServiceLoginFailedError{
			Username: e.Username,
			Err:      e,
		}
	case *ports.AuthAdapterInvalidCredentialsError:
		return &AuthServiceLoginFailedError{
			Username: e.Username,
			Err:      e,
		}
	default:
		return &AuthServiceUnknownError{
			Err: e,
		}
	}
}

type AuthServiceLoginFailedError struct {
	Username fields.Username
	Err      error
}

func (e *AuthServiceLoginFailedError) Error() string {
	return fmt.Sprintf("login failed for user %s: %s", e.Username, e.Err)
}

type AuthServiceRegistrationFailedError struct {
	Username fields.Username
	Email    fields.Email
	Err      error
}

func (e *AuthServiceRegistrationFailedError) Error() string {
	return fmt.Sprintf("registration failed for user %s (%s): %s", e.Username, e.Email, e.Err)
}

type AuthServiceUnknownError struct {
	Err error
}

func (e *AuthServiceUnknownError) Error() string {
	return fmt.Sprintf("unknown error %T: %s", e.Err, e.Err.Error())
}
