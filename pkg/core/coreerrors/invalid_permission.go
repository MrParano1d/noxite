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
