package ports

import (
	"context"
	"fmt"

	"github.com/mrparano1d/noxite/pkg/core/entities"
	"github.com/mrparano1d/noxite/pkg/core/fields"
)

// CreateUserInput contains the fields required to create a new user.
type CreateUserInput struct {
	Username fields.Username
	Email    fields.Email
	Password fields.Password
	RoleID   fields.EntityID
}

// UpdateUserInput contains the fields that can be updated.
// If a field is nil, it will not be updated.
type UpdateUserInput struct {
	Username *fields.Username
	Email    *fields.Email
	Password *fields.Password
	RoleID   *fields.EntityID
}

// UserPort is the interface that must be implemented by the user adapter.
// The user adapter is responsible for managing users.
type UserPort interface {
	// CreateUser creates a new user.
	// Returns the ID of the created user.
	// Returns UserAdapterCreateUserFailedError if failed to create user.
	// Returns UserAdapterUserAlreadyExistsError if user with the same username or email already exists.
	CreateUser(ctx context.Context, createUser CreateUserInput) (fields.EntityID, error)
	// GetUserByID returns the user with the given ID.
	// Returns UserAdapterGetUserByIDFailedError if failed to get user.
	// Returns UserAdapterUserNotFoundError if user with the given ID does not exist.
	GetUserByID(ctx context.Context, userID fields.EntityID) (*entities.User, error)
	// GetAllUsers returns all users.
	// Returns UserAdapterGetAllUsersFailedError if failed to get all users.
	GetAllUsers(ctx context.Context) ([]*entities.User, error)
	// UpdateUser updates the user with the given ID.
	// Returns UserAdapterUpdateUserFailedError if failed to update user.
	// Returns UserAdapterUserNotFoundError if user with the given ID does not exist.
	UpdateUser(ctx context.Context, id fields.EntityID, updateUser UpdateUserInput) error
	// DeleteUser deletes the user with the given ID.
	// Returns UserAdapterDeleteUserFailedError if failed to delete user.
	// Returns UserAdapterUserNotFoundError if user with the given ID does not exist.
	DeleteUser(ctx context.Context, userID fields.EntityID) error
	// FindUsersByEmailAddress returns all users with the given email address.
	// Returns UserAdapterGetAllUsersFailedError if failed to get all users.
	FindUsersByEmailAddress(ctx context.Context, emails []fields.Email) ([]*entities.User, error)
}

// errors

type UserAdapterCreateUserFailedError struct {
	Err error
}

func (e UserAdapterCreateUserFailedError) Error() string {
	return fmt.Sprintf("failed to create user: %v", e.Err)
}

type UserAdapterUserAlreadyExistsError struct {
	Username fields.Username
	Email    fields.Email
}

func (e UserAdapterUserAlreadyExistsError) Error() string {
	return fmt.Sprintf("user with username %v or email %v already exists", e.Username, e.Email)
}

type UserAdapterGetAllUsersFailedError struct {
	Err error
}

func (e UserAdapterGetAllUsersFailedError) Error() string {
	return fmt.Sprintf("failed to get all users: %v", e.Err)
}

type UserAdapterGetUserByIDFailedError struct {
	ID  fields.EntityID
	Err error
}

func (e UserAdapterGetUserByIDFailedError) Error() string {
	return fmt.Sprintf("failed to get user by id %v: %v", e.ID, e.Err)
}

type UserAdapterUserNotFoundError struct {
	ID fields.EntityID
}

func (e UserAdapterUserNotFoundError) Error() string {
	return fmt.Sprintf("user with id %v not found", e.ID)
}

type UserAdapterUpdateUserFailedError struct {
	ID  fields.EntityID
	Err error
}

func (e UserAdapterUpdateUserFailedError) Error() string {
	return fmt.Sprintf("failed to update user with id %v: %v", e.ID, e.Err)
}

type UserAdapterDeleteUserFailedError struct {
	ID  fields.EntityID
	Err error
}

func (e UserAdapterDeleteUserFailedError) Error() string {
	return fmt.Sprintf("failed to delete user with id %v: %v", e.ID, e.Err)
}
