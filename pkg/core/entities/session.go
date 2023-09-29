package entities

import (
	"time"

	"github.com/mrparano1d/noxite/pkg/core/fields"
)

type Session struct {
	Token     fields.SessionToken `json:"token"`
	ExpiresAt time.Time           `json:"expires_at"`
}

func NewSession(token fields.SessionToken, expiresAt time.Time) *Session {
	return &Session{
		Token:     token,
		ExpiresAt: expiresAt,
	}
}
