package fields

import (
	"fmt"
	"net/mail"
)

type Email string

func (e Email) String() string {
	return string(e)
}

func EmailFromString(s string) (Email, error) {

	if _, err := mail.ParseAddress(s); err != nil {
		return "", &InvalidEmailError{Email: s}
	}
	return Email(s), nil
}

// errors

type InvalidEmailError struct {
	Email string
}

func (e *InvalidEmailError) Error() string {
	return fmt.Sprintf("invalid email %s", e.Email)
}
