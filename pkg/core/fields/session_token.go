package fields

// SessionToken is a string that represents a session token.
// It must not be empty.
type SessionToken string

func (t SessionToken) String() string {
	return string(t)
}

func SessionTokenFromString(token string) (SessionToken, error) {
	if token == "" {
		return "", &EmptySessionTokenError{}
	}
	return SessionToken(token), nil
}

// errors

type EmptySessionTokenError struct{}

func (e EmptySessionTokenError) Error() string {
	return "SessionToken is empty"
}
