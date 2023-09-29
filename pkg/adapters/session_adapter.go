package adapters

import (
	"context"
	"time"

	//	json "github.com/bytedance/sonic"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/mrparano1d/noxite/pkg/core/entities"
	"github.com/mrparano1d/noxite/pkg/core/fields"
	"github.com/mrparano1d/noxite/pkg/core/ports"
	"github.com/redis/go-redis/v9"
)

type SessionAdapter struct {
	client *redis.Client
	prefix string
}

var _ ports.SessionPort = &SessionAdapter{}

func NewSessionAdapter(client *redis.Client) *SessionAdapter {
	return &SessionAdapter{
		client: client,
		prefix: "session:",
	}
}

// CreateSession creates a new session.
// It returns an error if the session could not be created.
func (s *SessionAdapter) CreateSession(ctx context.Context) (*entities.Session, error) {

	tokenUuid, err := uuid.NewRandom()
	if err != nil {
		return nil, &ports.CreateSessionFailedError{Err: err}
	}
	token, err := fields.SessionTokenFromString(tokenUuid.String())
	if err != nil {
		return nil, &ports.CreateSessionFailedError{Err: err}
	}

	expiry := time.Now().Add(24 * time.Hour)

	session := entities.NewSession(token, expiry)

	sessionB, err := s.Serialize(session)
	if err != nil {
		return nil, &ports.CreateSessionFailedError{Err: err}
	}

	cmd := s.client.HSet(ctx, s.prefix+token.String(), "session", sessionB)
	if cmd.Err() != nil {
		return nil, &ports.CreateSessionFailedError{Err: cmd.Err()}
	}

	cmdBool := s.client.ExpireAt(ctx, s.prefix+token.String(), expiry)
	if cmdBool.Err() != nil {
		return nil, &ports.CreateSessionFailedError{Err: cmdBool.Err()}
	}

	return session, nil
}

// LinkSessionToUser links a session to a user.
// It returns an error if the session could not be linked to the user.
func (s *SessionAdapter) LinkSessionToUser(ctx context.Context, token fields.SessionToken, userID fields.EntityID) error {
	cmd := s.client.HSet(ctx, s.prefix+userID.String(), token.String(), true)
	if cmd.Err() != nil {
		return &ports.SessionAdapterLinkSessionToUserFailedError{Err: cmd.Err()}
	}

	return nil
}

// InvalidateSession invalidates a session.
// It returns an error if the session could not be invalidated.
func (s *SessionAdapter) InvalidateSession(ctx context.Context, token fields.SessionToken) error {
	cmd := s.client.Del(ctx, s.prefix+token.String())
	if cmd.Err() != nil {
		return &ports.SessionAdapterInvalidateSessionFailedError{Err: cmd.Err()}
	}

	return nil
}

// GetLinkedSessions returns all sessions linked to the given user.
// It returns an error if the sessions could not be retrieved.
func (s *SessionAdapter) GetLinkedSessions(ctx context.Context, userID fields.EntityID) ([]fields.SessionToken, error) {
	tokens, err := s.client.HKeys(ctx, s.prefix+userID.String()).Result()

	if err != nil {
		return nil, &ports.SessionAdapterGetLinkedSessionsFailedError{Err: err}
	}

	sessionTokens := make([]fields.SessionToken, len(tokens))

	for i, token := range tokens {
		sessionToken, err := fields.SessionTokenFromString(token)
		if err != nil {
			continue
		}
		sessionTokens[i] = sessionToken
	}

	return sessionTokens, nil
}

// ValidateToken validates a session token.
// It returns an error if the token is invalid or expired.
func (s *SessionAdapter) ValidateToken(ctx context.Context, token fields.SessionToken) error {
	cmd := s.client.HExists(ctx, s.prefix+token.String(), "session")
	if cmd.Err() != nil {
		return &ports.ValidateTokenFailedError{Err: cmd.Err()}
	}

	if cmd.Val() == false {
		return &ports.SessionNotFoundError{Token: token}
	}

	return nil
}

// GetSession returns the session associated with the given token.
// It returns an error if the token is invalid or expired.
func (s *SessionAdapter) GetSession(ctx context.Context, token fields.SessionToken) (*entities.Session, error) {
	sessionB, err := s.client.Get(ctx, s.prefix+token.String()).Bytes()
	if err == redis.Nil {
		return nil, &ports.SessionNotFoundError{Token: token}
	} else if err != nil {
		return nil, &ports.GetSessionFailedError{Err: err}
	}

	var session entities.Session
	err = s.Deserialize(sessionB, &session)
	if err != nil {
		return nil, &ports.GetSessionFailedError{Err: err}
	}

	return &session, nil
}

// SetValue sets a value for the given key in the session associated with the given token.
// It returns an error if the token is invalid or expired.
func (s *SessionAdapter) SetValue(ctx context.Context, token fields.SessionToken, key fields.RequiredString, value any) error {
	valueB, err := s.Serialize(value)
	if err != nil {
		return &ports.SetValueFailedError{
			Token: token,
			Key:   key,
			Value: value,
			Err:   err,
		}
	}

	cmd := s.client.HSet(ctx, s.prefix+token.String(), key.String(), valueB)
	if cmd.Err() != nil {
		return &ports.SetValueFailedError{
			Token: token,
			Key:   key,
			Value: value,
			Err:   cmd.Err(),
		}
	}

	return nil
}

// GetValue returns the value associated with the given key in the session associated with the given token.
// It returns an error if the token is invalid or expired or if the key does not exist.
func (s *SessionAdapter) GetValue(ctx context.Context, token fields.SessionToken, key fields.RequiredString) ([]byte, error) {
	valueB, err := s.client.HGet(ctx, s.prefix+token.String(), key.String()).Bytes()
	if err == redis.Nil {
		return nil, &ports.KeyNotFoundError{Token: token, Key: key}
	} else if err != nil {
		return nil, &ports.GetValueFailedError{Err: err}
	}

	return valueB, nil
}

func (s *SessionAdapter) Serialize(value any) ([]byte, error) {
	return json.Marshal(value)
}

func (s *SessionAdapter) Deserialize(value []byte, target any) error {
	return json.Unmarshal(value, target)
}
