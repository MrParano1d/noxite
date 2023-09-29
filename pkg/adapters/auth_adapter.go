package adapters

import (
	"context"

	"github.com/mrparano1d/noxite/ent"
	"github.com/mrparano1d/noxite/ent/user"
	"github.com/mrparano1d/noxite/pkg/core/entities"
	"github.com/mrparano1d/noxite/pkg/core/fields"
	"github.com/mrparano1d/noxite/pkg/core/ports"
)

type AuthAdapter struct {
	entClient *ent.Client
}

var _ ports.AuthPort = (*AuthAdapter)(nil)

func NewAuthAdapter(entClient *ent.Client) *AuthAdapter {
	return &AuthAdapter{
		entClient: entClient,
	}
}

func (a *AuthAdapter) Login(ctx context.Context, username fields.Username, password fields.Password) (*entities.User, error) {
	user, err := a.entClient.User.Query().
		WithRole().
		Where(user.DeletedAtIsNil()).
		Where(user.Name(username.String())).
		Where(user.Password(password.Bytes())).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, &ports.AuthAdapterUserNotFoundError{
				Username: username,
			}
		}
		return nil, &ports.AuthAdapterLoginFailedError{
			Err: err,
		}
	}

	return UserFromEntUser(user)
}
