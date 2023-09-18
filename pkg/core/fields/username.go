package fields

import "fmt"

// Username needs to be atleast 2 characters but must not be more than 20 characters.
// The Username can contain alphanumeric values and "_" / "-" as special characters.
type Username string

func (u Username) String() string {
	return string(u)
}

// UsernameFromString validates the given string and returns a Username.
// If the string is invalid, an error is returned.
func UsernameFromString(s string) (Username, error) {
	if len(s) < 2 {
		return Username(s), &UsernameTooShortError{s}
	}
	if len(s) > 20 {
		return Username(s), &UsernameTooLongError{s}
	}
	for _, c := range s {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '_' || c == '-') {
			return Username(s), &UsernameInvalidCharacterError{s}
		}
	}
	return Username(s), nil
}

// errors

// UsernameTooShortError is returned when the username is too short.
type UsernameTooShortError struct {
	Username string
}

func (e UsernameTooShortError) Error() string {
	return fmt.Sprintf("Username %s is too short.", e.Username)
}

// UsernameTooLongError is returned when the username is too long.
type UsernameTooLongError struct {
	Username string
}

func (e UsernameTooLongError) Error() string {
	return fmt.Sprintf("Username %s is too long.", e.Username)
}

// UsernameInvalidCharacterError is returned when the username contains invalid characters.
type UsernameInvalidCharacterError struct {
	Username string
}

func (e UsernameInvalidCharacterError) Error() string {
	return fmt.Sprintf("Username %s contains invalid characters.", e.Username)
}

// UsernameInvalidError is returned when the username is invalid.
type UsernameInvalidError struct {
	Username string
}

func (e UsernameInvalidError) Error() string {
	return fmt.Sprintf("Username %s is invalid.", e.Username)
}
