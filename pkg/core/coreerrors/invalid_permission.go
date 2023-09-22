package coreerrors

type NotAllowedToPublishPackageError struct {
}

func (e *NotAllowedToPublishPackageError) Error() string {
	return "not allowed to publish package"
}

type NotAllowedToGetPackageError struct {
}

func (e *NotAllowedToGetPackageError) Error() string {
	return "not allowed to get package"
}

type NotAllowedToCreateUserError struct {
}

func (e *NotAllowedToCreateUserError) Error() string {
	return "not allowed to create user"
}

type NotAllowedToGetUserError struct {
}

func (e *NotAllowedToGetUserError) Error() string {
	return "not allowed to get user"
}

type NotAllowedToUpdateUserError struct {
}

func (e *NotAllowedToUpdateUserError) Error() string {
	return "not allowed to update user"
}

type NotAllowedToDeleteUserError struct {
}

func (e *NotAllowedToDeleteUserError) Error() string {
	return "not allowed to delete user"
}

type NotAllowedToCreateRoleError struct {
}

func (e *NotAllowedToCreateRoleError) Error() string {
	return "not allowed to create role"
}

type NotAllowedToGetRoleError struct {
}

func (e *NotAllowedToGetRoleError) Error() string {
	return "not allowed to get role"
}

type NotAllowedToUpdateRoleError struct {
}

func (e *NotAllowedToUpdateRoleError) Error() string {
	return "not allowed to update role"
}

type NotAllowedToDeleteRoleError struct {
}

func (e *NotAllowedToDeleteRoleError) Error() string {
	return "not allowed to delete role"
}
