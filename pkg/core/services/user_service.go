package services

import (
	"context"
	"fmt"

	"github.com/mrparano1d/noxite/pkg/core/coreerrors"
	"github.com/mrparano1d/noxite/pkg/core/entities"
	"github.com/mrparano1d/noxite/pkg/core/fields"
	"github.com/mrparano1d/noxite/pkg/core/ports"
)

type UserService struct {
	adapter ports.UserPort
}

func NewUserService(adapter ports.UserPort) *UserService {
	return &UserService{adapter: adapter}
}

func (s *UserService) CreateUser(ctx context.Context, user *entities.User, req CreateUserRequest) (fields.EntityID, error) {

	if user == nil || user.Role.Permissions.CreateUser == false {
		return fields.EntityID(0), &coreerrors.NotAllowedToCreateUserError{}
	}

	input, err := CreateUserRequestToInput(req)
	if err != nil {
		return fields.EntityID(0), err
	}

	userID, err := s.adapter.CreateUser(ctx, input)
	if err != nil {
		return fields.EntityID(0), handleUserServiceErrors(err)
	}

	return userID, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, user *entities.User) ([]*entities.User, error) {

	if user == nil || user.Role.Permissions.GetUser == false {
		return nil, &coreerrors.NotAllowedToGetUserError{}
	}

	users, err := s.adapter.GetAllUsers(ctx)
	if err != nil {
		return nil, handleUserServiceErrors(err)
	}
	return users, nil
}

func (s *UserService) GetUserByID(ctx context.Context, user *entities.User, userID string) (*entities.User, error) {

	if user == nil || user.Role.Permissions.GetUser == false {
		return nil, &coreerrors.NotAllowedToGetUserError{}
	}

	id, err := fields.EntityIDFromString(userID)
	if err != nil {
		return nil, handleUserServiceRequestValidationError("id", err.Error())
	}

	user, err = s.adapter.GetUserByID(ctx, id)
	if err != nil {
		return nil, handleUserServiceErrors(err)
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *entities.User, userID string, req UpdateUserRequest) error {

	if user == nil || user.Role.Permissions.UpdateUser == false {
		return &coreerrors.NotAllowedToUpdateUserError{}
	}

	id, err := fields.EntityIDFromString(userID)
	if err != nil {
		return handleUserServiceRequestValidationError("id", err.Error())
	}

	input, err := UpdateUserRequestToInput(req)
	if err != nil {
		return err
	}

	err = s.adapter.UpdateUser(ctx, id, input)
	if err != nil {
		return handleUserServiceErrors(err)
	}

	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, user *entities.User, userID string) error {

	if user == nil || user.Role.Permissions.DeleteUser == false {
		return &coreerrors.NotAllowedToDeleteUserError{}
	}

	id, err := fields.EntityIDFromString(userID)
	if err != nil {
		return handleUserServiceRequestValidationError("id", err.Error())
	}

	err = s.adapter.DeleteUser(ctx, id)
	if err != nil {
		return handleUserServiceErrors(err)
	}

	return nil
}

// requests

type CreateUserRequest struct {
	Username string
	Email    string
	Password string
}

func CreateUserRequestToInput(req CreateUserRequest) (ports.CreateUserInput, error) {
	var input ports.CreateUserInput
	var err error

	if input.Username, err = fields.UsernameFromString(req.Username); err != nil {
		return input, handleUserServiceRequestValidationError("username", err.Error())
	}

	if input.Email, err = fields.EmailFromString(req.Email); err != nil {
		return input, handleUserServiceRequestValidationError("email", err.Error())
	}

	if input.Password, err = fields.PasswordFromString(req.Password); err != nil {
		return input, handleUserServiceRequestValidationError("password", err.Error())
	}

	return input, nil
}

type UpdateUserRequest struct {
	Username *string
	Email    *string
	Password *string
}

func UpdateUserRequestToInput(req UpdateUserRequest) (ports.UpdateUserInput, error) {
	var input ports.UpdateUserInput

	if req.Username != nil {
		if username, err := fields.UsernameFromString(*req.Username); err != nil {
			return input, handleUserServiceRequestValidationError("username", err.Error())
		} else {
			input.Username = &username
		}
	}

	if req.Email != nil {
		if email, err := fields.EmailFromString(*req.Email); err != nil {
			return input, handleUserServiceRequestValidationError("email", err.Error())
		} else {
			input.Email = &email
		}
	}

	if req.Password != nil {
		if password, err := fields.PasswordFromString(*req.Password); err != nil {
			return input, handleUserServiceRequestValidationError("password", err.Error())
		} else {
			input.Password = &password
		}
	}

	return input, nil
}

// errors

func handleUserServiceRequestValidationError(field string, reason string) error {
	return &UserServiceRequestValidationError{Field: field, Reason: reason}
}

func handleUserServiceErrors(err error) error {
	switch e := err.(type) {
	case *ports.UserAdapterGetUserByIDFailedError:
		return &UserServiceGetUserByIDFailedError{ID: e.ID, Err: e.Err}
	case *ports.UserAdapterUserNotFoundError:
		return &UserServiceUserNotFoundError{ID: e.ID}
	case *ports.UserAdapterUpdateUserFailedError:
		return &UserServiceUpdateUserFailedError{ID: e.ID, Err: e.Err}
	case *ports.UserAdapterDeleteUserFailedError:
		return &UserServiceDeleteUserFailedError{ID: e.ID, Err: e.Err}
	case *ports.UserAdapterCreateUserFailedError:
		return &UserServiceCreateUserFailedError{Err: e.Err}
	case *ports.UserAdapterUserAlreadyExistsError:
		return &UserServiceUserAlreadyExistsError{Username: e.Username, Email: e.Email}
	case *ports.UserAdapterGetAllUsersFailedError:
		return &UserServiceGetAllUsersFailedError{Err: e.Err}
	default:
		return &UserServiceUnknownError{Err: err}
	}
}

type UserServiceUnknownError struct {
	Err error
}

func (e UserServiceUnknownError) Error() string {
	return fmt.Sprintf("unknown user service error: %v", e.Err)
}

type UserServiceRequestValidationError struct {
	Field  string
	Reason string
}

func (e UserServiceRequestValidationError) Error() string {
	return fmt.Sprintf("invalid user service request: %s: %s", e.Field, e.Reason)
}

type UserServiceCreateUserFailedError struct {
	Err error
}

func (e UserServiceCreateUserFailedError) Error() string {
	return fmt.Sprintf("failed to create user: %v", e.Err)
}

type UserServiceUserAlreadyExistsError struct {
	Username fields.Username
	Email    fields.Email
}

func (e UserServiceUserAlreadyExistsError) Error() string {
	return fmt.Sprintf("user with username %q or email %q already exists", e.Username, e.Email)
}

type UserServiceGetUserByIDFailedError struct {
	ID  fields.EntityID
	Err error
}

func (e UserServiceGetUserByIDFailedError) Error() string {
	return fmt.Sprintf("failed to get user with ID %q: %v", e.ID, e.Err)
}

type UserServiceUserNotFoundError struct {
	ID fields.EntityID
}

func (e UserServiceUserNotFoundError) Error() string {
	return fmt.Sprintf("user with ID %q not found", e.ID)
}

type UserServiceUpdateUserFailedError struct {
	ID  fields.EntityID
	Err error
}

func (e UserServiceUpdateUserFailedError) Error() string {
	return fmt.Sprintf("failed to update user with ID %q: %v", e.ID, e.Err)
}

type UserServiceDeleteUserFailedError struct {
	ID  fields.EntityID
	Err error
}

func (e UserServiceDeleteUserFailedError) Error() string {
	return fmt.Sprintf("failed to delete user with ID %q: %v", e.ID, e.Err)
}

type UserServiceGetAllUsersFailedError struct {
	Err error
}

func (e UserServiceGetAllUsersFailedError) Error() string {
	return fmt.Sprintf("failed to get all users: %v", e.Err)
}
