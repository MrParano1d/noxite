package services

import (
	"fmt"

	"github.com/mrparano1d/getregd/pkg/core/entities"
	"github.com/mrparano1d/getregd/pkg/core/fields"
	"github.com/mrparano1d/getregd/pkg/core/ports"
)

type AuthService struct {
	adapter ports.AuthPort
}

func NewAuthService(adapter ports.AuthPort) *AuthService {
	return &AuthService{adapter: adapter}
}

// usecases

func (s *AuthService) Login(username string, email string, password string) (*entities.User, error) {

	name, err := fields.UsernameFromString(username)
	if err != nil {
		return nil, handleErrors(err)
	}

	pw, err := fields.PasswordFromString(password)
	if err != nil {
		return nil, handleErrors(err)
	}

	user, err := s.adapter.Login(name, pw)
	if err != nil {
		if _, ok := err.(*ports.AuthAdapterUserNotFoundError); ok {
			return s.Register(username, email, password)
		}
		return nil, handleErrors(err)
	}
	return user, nil
}

func (s *AuthService) Register(username string, email string, password string) (*entities.User, error) {

	name, err := fields.UsernameFromString(username)
	if err != nil {
		return nil, handleErrors(err)
	}

	em, err := fields.EmailFromString(email)
	if err != nil {
		return nil, handleErrors(err)
	}

	pw, err := fields.PasswordFromString(password)
	if err != nil {
		return nil, handleErrors(err)
	}

	user, err := s.adapter.Register(name, em, pw)
	if err != nil {
		return nil, handleErrors(err)
	}
	return user, nil
}

func (s *AuthService) GenerateToken(user *entities.User) (fields.RequiredString, error) {
	return s.adapter.GenerateToken(user)
}

// service errors

func handleErrors(err error) error {
	switch e := err.(type) {
	case *ports.AuthAdapterLoginFailedError:
		return &AuthServiceLoginFailedError{
			Username: e.Username,
			Err:      e.Err,
		}
	case *ports.AuthAdapterRegistrationFailedError:
		return &AuthServiceRegistrationFailedError{
			Username: e.Username,
			Email:    e.Email,
			Err:      e.Err,
		}
	case *ports.AuthAdapterUserAlreadyExistsError:
		return &AuthServiceRegistrationFailedError{
			Username: e.Username,
			Email:    e.Email,
			Err:      e,
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
	case *fields.UsernameInvalidError:
		return &AuthServiceInvalidUsernameError{
			Username: e.Username,
			Reason:   e.Error(),
		}
	case *fields.UsernameTooLongError:
		return &AuthServiceInvalidUsernameError{
			Username: e.Username,
			Reason:   e.Error(),
		}
	case *fields.UsernameTooShortError:
		return &AuthServiceInvalidUsernameError{
			Username: e.Username,
			Reason:   e.Error(),
		}
	case *fields.UsernameInvalidCharacterError:
		return &AuthServiceInvalidUsernameError{
			Username: e.Username,
			Reason:   e.Error(),
		}
	case *fields.PasswordNoLowerError, *fields.PasswordNoNumberError, *fields.PasswordNoUpperError, *fields.PasswordTooShortError:
		return &AuthServiceInvalidPasswordError{
			Reason: e.Error(),
		}
	case *fields.InvalidEmailError:
		return &AuthServiceInvalidEmailError{
			Email:  e.Email,
			Reason: e.Error(),
		}
	default:
		return &AuthServiceUnknownError{
			Err: e,
		}
	}
}

type AuthServiceInvalidEmailError struct {
	Email  string
	Reason string
}

func (e *AuthServiceInvalidEmailError) Error() string {
	return fmt.Sprintf("invalid email %s: %s", e.Email, e.Reason)
}

type AuthServiceInvalidUsernameError struct {
	Username string
	Reason   string
}

func (e *AuthServiceInvalidUsernameError) Error() string {
	return fmt.Sprintf("invalid username %s: %s", e.Username, e.Reason)
}

type AuthServiceInvalidPasswordError struct {
	Reason string
}

func (e *AuthServiceInvalidPasswordError) Error() string {
	return fmt.Sprintf("invalid password: %s", e.Reason)
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
	return fmt.Sprintf("unknown error: %w", e.Err)
}
