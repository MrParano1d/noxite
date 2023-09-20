package ports

import (
	"context"
	"fmt"

	"github.com/mrparano1d/getregd/pkg/core/entities"
	"github.com/mrparano1d/getregd/pkg/core/fields"
)

type SessionPort interface {
	// CreateSession creates a new session.
	// It returns an error if the session could not be created.
	CreateSession(ctx context.Context) (*entities.Session, error)
	// LinkSessionToUser links a session to a user.
	// It returns an error if the session could not be linked to the user.
	LinkSessionToUser(ctx context.Context, token fields.SessionToken, userID fields.EntityID) error
	// InvalidateSession invalidates a session.
	// It returns an error if the session could not be invalidated.
	InvalidateSession(ctx context.Context, token fields.SessionToken) error
	// GetLinkedSessions returns all sessions linked to the given user.
	// It returns an error if the sessions could not be retrieved.
	GetLinkedSessions(ctx context.Context, userID fields.EntityID) ([]fields.SessionToken, error)
	// ValidateToken validates a session token.
	// It returns an error if the token is invalid or expired.
	ValidateToken(ctx context.Context, token fields.SessionToken) error
	// GetSession returns the session associated with the given token.
	// It returns an error if the token is invalid or expired.
	GetSession(ctx context.Context, token fields.SessionToken) (*entities.Session, error)
	// SetValue sets a value for the given key in the session associated with the given token.
	// It returns an error if the token is invalid or expired.
	SetValue(ctx context.Context, token fields.SessionToken, key fields.RequiredString, value any) error
	// GetValue returns the value associated with the given key in the session associated with the given token.
	// It returns an error if the token is invalid or expired or if the key does not exist.
	GetValue(ctx context.Context, token fields.SessionToken, key fields.RequiredString) ([]byte, error)

	Serialize(value any) ([]byte, error)
	Deserialize(value []byte, target any) error
}

func SessionValueFromAdapter[V any](adapter SessionPort, ctx context.Context, token fields.SessionToken, key fields.RequiredString, defaultValue ...V) (V, error) {
	var zero V
	value, err := adapter.GetValue(ctx, token, key)
	if err != nil {
		if _, ok := err.(KeyNotFoundError); ok && len(defaultValue) > 0 {
			return defaultValue[0], nil
		} else {
			return zero, err
		}
	}

	err = adapter.Deserialize(value, &zero)
	if err != nil {
		return zero, &InvalidValueTypeError{
			Token:  token,
			Key:    key,
			Reason: err.Error(),
		}
	}
	return zero, nil
}

// errors

type CreateSessionFailedError struct {
	Err error
}

func (e CreateSessionFailedError) Error() string {
	return "CreateSession failed: " + e.Err.Error()
}

type ValidateTokenFailedError struct {
	Token fields.SessionToken
	Err   error
}

func (e ValidateTokenFailedError) Error() string {
	return "ValidateToken failed for token " + e.Token.String() + ": " + e.Err.Error()
}

type SessionNotFoundError struct {
	Token fields.SessionToken
}

func (e SessionNotFoundError) Error() string {
	return "Session not found for token " + e.Token.String()
}

type InvalidTokenError struct {
	Token fields.SessionToken
}

func (e InvalidTokenError) Error() string {
	return "Token " + e.Token.String() + " is invalid"
}

type ExpiredTokenError struct {
	Token fields.SessionToken
}

func (e ExpiredTokenError) Error() string {
	return "Token " + e.Token.String() + " is expired"
}

type GetSessionFailedError struct {
	Token fields.SessionToken
	Err   error
}

func (e GetSessionFailedError) Error() string {
	return "GetSession failed for token " + e.Token.String() + ": " + e.Err.Error()
}

type SetValueFailedError struct {
	Token fields.SessionToken
	Key   fields.RequiredString
	Value any
	Err   error
}

func (e SetValueFailedError) Error() string {
	return "SetValue failed for token " + e.Token.String() + " and key " + e.Key.String() + ": " + e.Err.Error()
}

type GetValueFailedError struct {
	Token fields.SessionToken
	Key   fields.RequiredString
	Err   error
}

func (e GetValueFailedError) Error() string {
	return "GetValue failed for token " + e.Token.String() + " and key " + e.Key.String() + ": " + e.Err.Error()
}

type KeyNotFoundError struct {
	Token fields.SessionToken
	Key   fields.RequiredString
}

func (e KeyNotFoundError) Error() string {
	return "Key " + e.Key.String() + " not found in session for token " + e.Token.String()
}

type InvalidValueTypeError struct {
	Token  fields.SessionToken
	Key    fields.RequiredString
	Reason string
}

func (e InvalidValueTypeError) Error() string {
	return "Invalid value type for key " + e.Key.String() + " in session for token " + e.Token.String() + ": " + e.Reason
}

type SessionAdapterLinkSessionToUserFailedError struct {
	Token  fields.SessionToken
	UserID fields.EntityID
	Err    error
}

func (e SessionAdapterLinkSessionToUserFailedError) Error() string {
	return "LinkSessionToUser failed for token " + e.Token.String() + " and user " + e.UserID.String() + ": " + e.Err.Error()
}

type SessionAdapterInvalidateSessionFailedError struct {
	Token fields.SessionToken
	Err   error
}

func (e SessionAdapterInvalidateSessionFailedError) Error() string {
	return "InvalidateSession failed for token " + e.Token.String() + ": " + e.Err.Error()
}

type SessionAdapterGetLinkedSessionsFailedError struct {
	UserID fields.EntityID
	Err    error
}

func (e SessionAdapterGetLinkedSessionsFailedError) Error() string {
	return "GetLinkedSessions failed for user " + e.UserID.String() + ": " + e.Err.Error()
}

type SessionAdapterDeserializeFailedError struct {
	Value []byte
	Err   error
}

func (e SessionAdapterDeserializeFailedError) Error() string {
	return "deserialize failed for value " + string(e.Value) + ": " + e.Err.Error()
}

type SessionAdapterSerializeFailedError struct {
	Value any
	Err   error
}

func (e SessionAdapterSerializeFailedError) Error() string {
	return fmt.Sprintf("Serialize failed for value %T: %s", e.Value, e.Err.Error())
}
