package fields

// RequiredString is a string that must not be empty.
type RequiredString string

func (s RequiredString) String() string {
	return string(s)
}

// RequiredStringFromString validates the given string and returns a RequiredString.
// If the string is invalid, an error is returned.
func RequiredStringFromString(s string) (RequiredString, error) {
	if len(s) == 0 {
		return RequiredString(s), &RequiredStringEmptyError{}
	}
	return RequiredString(s), nil
}

// errors

// RequiredStringEmptyError is returned when the string is empty.
type RequiredStringEmptyError struct{}

func (e RequiredStringEmptyError) Error() string {
	return "string is empty"
}
