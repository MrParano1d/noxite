package services

import (
	"context"
	"fmt"

	"github.com/mrparano1d/noxite/pkg/core/entities"
	"github.com/mrparano1d/noxite/pkg/core/fields"
	"github.com/mrparano1d/noxite/pkg/core/ports"
)

type SessionService struct {
	adapter ports.SessionPort
}

func NewSessionService(adapter ports.SessionPort) *SessionService {
	return &SessionService{
		adapter: adapter,
	}
}

func (s *SessionService) CreateSession(ctx context.Context) (*entities.Session, error) {
	session, err := s.adapter.CreateSession(ctx)
	if err != nil {
		return nil, handleSessionErrors(err)
	}
	return session, nil
}

func (s *SessionService) CreateSessionForUser(ctx context.Context, user *entities.User) (*entities.Session, error) {
	session, err := s.adapter.CreateSession(ctx)
	if err != nil {
		return nil, handleSessionErrors(err)
	}

	if err := s.adapter.LinkSessionToUser(ctx, session.Token, user.ID); err != nil {
		return nil, handleSessionErrors(err)
	}

	err = s.adapter.SetValue(ctx, session.Token, "user", user)
	if err != nil {
		return nil, handleSessionErrors(err)
	}

	return session, nil
}

func (s *SessionService) GetLinkedSessions(ctx context.Context, userID string) ([]fields.SessionToken, error) {
	uID, err := fields.EntityIDFromString(userID)
	if err != nil {
		return nil, handleSessionErrors(err)
	}

	tokens, err := s.adapter.GetLinkedSessions(ctx, uID)
	if err != nil {
		return nil, handleSessionErrors(err)
	}
	return tokens, nil
}

func (s *SessionService) InvalidateSession(ctx context.Context, token string) error {
	sessionToken, err := fields.SessionTokenFromString(token)
	if err != nil {
		return handleSessionErrors(err)
	}

	err = s.adapter.InvalidateSession(ctx, sessionToken)
	if err != nil {
		return handleSessionErrors(err)
	}
	return nil
}

func (s *SessionService) ValidateToken(ctx context.Context, token string) error {
	sessionToken, err := fields.SessionTokenFromString(token)
	if err != nil {
		return handleSessionErrors(err)
	}

	err = s.adapter.ValidateToken(ctx, sessionToken)
	if err != nil {
		return handleSessionErrors(err)
	}
	return nil
}

func (s *SessionService) GetSession(ctx context.Context, token string) (*entities.Session, error) {
	sessionToken, err := fields.SessionTokenFromString(token)
	if err != nil {
		return nil, handleSessionErrors(err)
	}

	session, err := s.adapter.GetSession(ctx, sessionToken)
	if err != nil {
		return nil, handleSessionErrors(err)
	}
	return session, nil
}

func (s *SessionService) SetValue(ctx context.Context, token string, key string, value any) error {
	sessionToken, err := fields.SessionTokenFromString(token)
	if err != nil {
		return handleSessionErrors(err)
	}

	keyNZ, err := fields.RequiredStringFromString(key)
	if err != nil {
		return handleSessionErrors(err)
	}

	err = s.adapter.SetValue(ctx, sessionToken, keyNZ, value)
	if err != nil {
		return handleSessionErrors(err)
	}
	return nil
}

func (s *SessionService) GetValue(ctx context.Context, token string, key string) (any, error) {
	sessionToken, err := fields.SessionTokenFromString(token)
	if err != nil {
		return nil, handleSessionErrors(err)
	}

	keyNZ, err := fields.RequiredStringFromString(key)
	if err != nil {
		return nil, handleSessionErrors(err)
	}

	value, err := s.adapter.GetValue(ctx, sessionToken, keyNZ)
	if err != nil {
		return nil, handleSessionErrors(err)
	}
	return value, nil
}

func SessionValueFromService[V any](service *SessionService, ctx context.Context, token string, key string, defaultValue ...*V) (*V, error) {
	var v V
	value, err := service.GetValue(ctx, token, key)
	if err != nil {
		if _, ok := err.(*KeyNotFoundError); ok && len(defaultValue) > 0 {
			return defaultValue[0], nil
		} else {
			return nil, err
		}
	}

	bs, ok := value.([]byte)
	if ok {
		err := service.adapter.Deserialize(bs, &v)
		if err != nil {
			return nil, &ports.SessionAdapterDeserializeFailedError{
				Value: bs,
				Err:   err,
			}
		}
	} else {
		v, ok = value.(V)
		if !ok {
			return nil, &InvalidValueTypeError{
				Token: token,
				Key:   key,
				Value: value,
			}
		}
	}
	return &v, nil
}

func handleSessionErrors(err error) error {
	switch et := err.(type) {
	case *ports.SessionAdapterGetLinkedSessionsFailedError:
		return &GetLinkedSessionsFailedError{
			UserID: et.UserID.String(),
			Err:    et.Err,
		}
	case *ports.SessionAdapterInvalidateSessionFailedError:
		return &InvalidateSessionFailedError{
			Token: et.Token.String(),
			Err:   et.Err,
		}
	case *ports.SessionAdapterLinkSessionToUserFailedError:
		return &LinkSessionToUserFailedError{
			Token:  et.Token.String(),
			UserID: et.UserID.String(),
			Err:    et.Err,
		}
	case *fields.EmptySessionTokenError:
		return &InvalidTokenError{
			Token:  "",
			Reason: "empty token",
		}
	case *ports.CreateSessionFailedError:
		return &CreateSessionFailedError{
			Err: et.Err,
		}
	case *ports.ValidateTokenFailedError:
		return &ValidateTokenFailedError{
			Token: et.Token.String(),
			Err:   et.Err,
		}
	case *ports.InvalidTokenError:
		return &InvalidTokenError{
			Token:  et.Token.String(),
			Reason: et.Error(),
		}
	case *ports.ExpiredTokenError:
		return &ExpiredTokenError{
			Token: et.Token.String(),
		}
	case *ports.KeyNotFoundError:
		return &KeyNotFoundError{
			Token: et.Token.String(),
			Key:   et.Key.String(),
		}
	case *ports.InvalidValueTypeError:
		return &InvalidValueTypeError{
			Token: et.Token.String(),
			Key:   et.Key.String(),
		}
	case *ports.GetSessionFailedError:
		return &GetSessionFailedError{
			Token: et.Token.String(),
			Err:   et.Err,
		}
	case *ports.SetValueFailedError:
		return &SetValueFailedError{
			Token: et.Token.String(),
			Key:   et.Key.String(),
			Err:   et.Err,
		}
	case *ports.GetValueFailedError:
		return &GetValueFailedError{
			Token: et.Token.String(),
			Key:   et.Key.String(),
			Err:   et.Err,
		}
	case *fields.RequiredStringEmptyError:
		return &KeyNotFoundError{
			Token: "",
			Key:   "",
		}
	default:
		return err
	}

}

// errors

type CreateSessionFailedError struct {
	Err error
}

func (e CreateSessionFailedError) Error() string {
	return "CreateSession failed: " + e.Err.Error()
}

type ValidateTokenFailedError struct {
	Token string
	Err   error
}

func (e ValidateTokenFailedError) Error() string {
	return "ValidateToken failed for token " + e.Token + ": " + e.Err.Error()
}

type InvalidTokenError struct {
	Token  string
	Reason string
}

func (e InvalidTokenError) Error() string {
	return "Token " + e.Token + " is invalid: " + e.Reason
}

type ExpiredTokenError struct {
	Token string
}

func (e ExpiredTokenError) Error() string {
	return "Token " + e.Token + " is expired"
}

type KeyNotFoundError struct {
	Token string
	Key   string
}

func (e KeyNotFoundError) Error() string {
	return "Key " + e.Key + " not found for token " + e.Token
}

type InvalidValueTypeError struct {
	Token string
	Key   string
	Value any
}

func (e InvalidValueTypeError) Error() string {
	return fmt.Sprintf("Value %T for key %s is invalid for token %s", e.Value, e.Key, e.Token)
}

type GetSessionFailedError struct {
	Token string
	Err   error
}

func (e GetSessionFailedError) Error() string {
	return "GetSession failed for token " + e.Token + ": " + e.Err.Error()
}

type SetValueFailedError struct {
	Token string
	Key   string
	Err   error
}

func (e SetValueFailedError) Error() string {
	return "SetValue failed for token " + e.Token + " and key " + e.Key + ": " + e.Err.Error()
}

type GetValueFailedError struct {
	Token string
	Key   string
	Err   error
}

func (e GetValueFailedError) Error() string {
	return "GetValue failed for token " + e.Token + " and key " + e.Key + ": " + e.Err.Error()
}

type GetLinkedSessionsFailedError struct {
	UserID string
	Err    error
}

func (e GetLinkedSessionsFailedError) Error() string {
	return "GetLinkedSessions failed for user " + e.UserID + ": " + e.Err.Error()
}

type InvalidateSessionFailedError struct {
	Token string
	Err   error
}

func (e InvalidateSessionFailedError) Error() string {
	return "InvalidateSession failed for token " + e.Token + ": " + e.Err.Error()
}

type LinkSessionToUserFailedError struct {
	Token  string
	UserID string
	Err    error
}

func (e LinkSessionToUserFailedError) Error() string {
	return "LinkSessionToUser failed for token " + e.Token + " and user " + e.UserID + ": " + e.Err.Error()
}
