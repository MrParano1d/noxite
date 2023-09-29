package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/mrparano1d/noxite/ent"
	"github.com/mrparano1d/noxite/ent/user"
	"github.com/mrparano1d/noxite/pkg/core/entities"
	"github.com/mrparano1d/noxite/pkg/core/fields"
	"github.com/mrparano1d/noxite/pkg/core/ports"
)

type UserAdapter struct {
	entClient *ent.Client
}

var _ ports.UserPort = (*UserAdapter)(nil)

func NewUserAdapter(entClient *ent.Client) *UserAdapter {
	return &UserAdapter{entClient: entClient}
}

func (u *UserAdapter) CreateUser(ctx context.Context, createUser ports.CreateUserInput) (fields.EntityID, error) {

	user, err := u.entClient.User.Create().
		SetName(createUser.Username.String()).
		SetEmail(createUser.Email.String()).
		SetPassword(createUser.Password).
		SetRoleID(createUser.RoleID.Int()).
		Save(ctx)

	if err != nil {
		return fields.EntityID(0), &ports.UserAdapterCreateUserFailedError{
			Err: err,
		}
	}

	id, err := fields.EntityIDFromInt(user.ID)
	if err != nil {
		return fields.EntityID(0), &ports.UserAdapterCreateUserFailedError{
			Err: fmt.Errorf("failed to convert ent.User.ID to fields.EntityID: %w", err),
		}
	}

	return id, nil
}

func (u *UserAdapter) GetUserByID(ctx context.Context, userID fields.EntityID) (*entities.User, error) {

	user, err := u.entClient.User.Query().WithRole().Where(user.DeletedAtIsNil()).Only(ctx)
	if err != nil {
		return nil, &ports.UserAdapterGetUserByIDFailedError{
			Err: err,
		}
	}

	return UserFromEntUser(user)
}

func (u *UserAdapter) GetAllUsers(ctx context.Context) ([]*entities.User, error) {

	users, err := u.entClient.User.Query().WithRole().Where(user.DeletedAtIsNil()).All(ctx)
	if err != nil {
		return nil, &ports.UserAdapterGetAllUsersFailedError{
			Err: err,
		}
	}

	result := make([]*entities.User, 0, len(users))
	for _, user := range users {
		user, err := UserFromEntUser(user)
		if err != nil {
			return nil, &ports.UserAdapterGetAllUsersFailedError{
				Err: err,
			}
		}
		result = append(result, user)
	}

	return result, nil
}

func (u *UserAdapter) UpdateUser(ctx context.Context, id fields.EntityID, updateUser ports.UpdateUserInput) error {

	query := u.entClient.User.UpdateOneID(id.Int())

	query = query.Where(user.DeletedAtIsNil())

	if updateUser.Username != nil {
		query = query.SetName(updateUser.Username.String())
	}

	if updateUser.Email != nil {
		query = query.SetEmail(updateUser.Email.String())
	}

	if updateUser.Password != nil {
		query = query.SetPassword(updateUser.Password.Bytes())
	}

	if updateUser.RoleID != nil {
		query = query.SetRoleID(updateUser.RoleID.Int())
	}

	query = query.SetUpdatedAt(time.Now())

	_, err := query.Save(ctx)

	if err != nil {
		return &ports.UserAdapterUpdateUserFailedError{
			ID:  id,
			Err: err,
		}
	}

	return nil
}

func (u *UserAdapter) DeleteUser(ctx context.Context, userID fields.EntityID) error {
	err := u.entClient.User.UpdateOneID(userID.Int()).SetDeletedAt(time.Now()).Exec(ctx)
	if err != nil {
		return &ports.UserAdapterDeleteUserFailedError{
			ID:  userID,
			Err: err,
		}
	}

	return nil
}

func (u *UserAdapter) FindUsersByEmailAddress(ctx context.Context, emails []fields.Email) ([]*entities.User, error) {
	emailStrings := make([]string, 0, len(emails))
	for _, email := range emails {
		emailStrings = append(emailStrings, email.String())
	}
	users, err := u.entClient.User.Query().WithRole().Where(user.EmailIn(emailStrings...)).All(ctx)
	if err != nil {
		return nil, &ports.UserAdapterGetAllUsersFailedError{
			Err: err,
		}
	}

	result := make([]*entities.User, 0, len(users))
	for _, user := range users {
		user, err := UserFromEntUser(user)
		if err != nil {
			return nil, &ports.UserAdapterGetAllUsersFailedError{
				Err: err,
			}
		}
		result = append(result, user)
	}

	return result, nil
}
