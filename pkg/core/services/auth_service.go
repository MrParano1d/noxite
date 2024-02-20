package services

import (
	"context"
	"fmt"

	"github.com/mrparano1d/noxite/pkg/core/entities"
	"github.com/mrparano1d/noxite/pkg/core/fields"
	"github.com/mrparano1d/noxite/pkg/core/ports"
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

func (s *AuthService) Login(ctx context.Context, usernameOrEmail string, password string) (*entities.Session, error) {

	var err error
	var userEmail *fields.Email
	var userName *fields.Username

	possibleEmail, err := fields.EmailFromString(usernameOrEmail)
	if err != nil {
		name, err := fields.UsernameFromString(usernameOrEmail)
		if err != nil {
			return nil, handleErrors(err)
		}
		userName = &name
	} else {
		userEmail = &possibleEmail
	}

	pw, err := fields.PasswordFromString(password)
	if err != nil {
		return nil, handleErrors(err)
	}

	var user *entities.User

	if userName != nil {
		user, err = s.adapter.Login(ctx, *userName, pw)
		if err != nil {
			return nil, handleErrors(err)
		}
	}

	if userEmail != nil {
		user, err = s.adapter.LoginByEmail(ctx, *userEmail, pw)
		if err != nil {
			return nil, handleErrors(err)
		}
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
