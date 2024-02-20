package fields

// Password must be atleast 8 characters long and must contain atleast one uppercase letter, one lowercase letter and one number.
type Password []byte

func (p Password) String() string {
	return string(p)
}

func (p Password) Bytes() []byte {
	return []byte(p)
}

// Compare compares the given password with the current password.
// If the passwords are equal, true is returned.
func (p Password) Compare(other Password) bool {
	return string(p) == string(other)
}

// PasswordFromString validates the given string and returns a Password.
// If the string is invalid, an error is returned.
func PasswordFromString(s string) (Password, error) {

	if len(s) < 8 {
		return Password(s), &PasswordTooShortError{}
	}

	var (
		hasUpper       bool
		hasLower       bool
		hasNum         bool
		hasSpecialChar bool
	)

	for _, c := range s {
		switch {
		case c == ' ':
			return Password(s), &PasswordInvalidCharacter{}
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasNum = true
		case c == '!' || c == '@' || c == '#' || c == '$' || c == '%' || c == '^' || c == '&' || c == '*' || c == '(' || c == ')' || c == '-' || c == '+' || c == '=' || c == '{' || c == '}' || c == '[' || c == ']' || c == ':' || c == ';' || c == '<' || c == '>' || c == ',' || c == '.' || c == '?' || c == '/' || c == '|' || c == '\\' || c == '`' || c == '~':
			hasSpecialChar = true
		}
	}

	if !hasUpper {
		return Password(s), &PasswordNoUpperError{}
	}

	if !hasLower {
		return Password(s), &PasswordNoLowerError{}
	}

	if !hasNum {
		return Password(s), &PasswordNoNumberError{}
	}

	if !hasSpecialChar {
		return Password(s), &PasswordNoSpecialCharError{}
	}

	return Password(s), nil
}

func PasswordFromBytes(b []byte) (Password, error) {
	return PasswordFromString(string(b))
}

// errors

// PasswordTooShortError is returned when the password is too short.
type PasswordTooShortError struct {
}

func (e PasswordTooShortError) Error() string {
	return "Password is too short."
}

// PasswordNoUpperError is returned when the password does not contain an uppercase letter.
type PasswordNoUpperError struct {
}

func (e PasswordNoUpperError) Error() string {
	return "Password does not contain an uppercase letter."
}

// PasswordNoLowerError is returned when the password does not contain a lowercase letter.
type PasswordNoLowerError struct {
}

func (e PasswordNoLowerError) Error() string {
	return "Password does not contain a lowercase letter."
}

// PasswordNoNumberError is returned when the password does not contain a number.
type PasswordNoNumberError struct {
}

func (e PasswordNoNumberError) Error() string {
	return "Password does not contain a number."
}

// PasswordNoSpecialCharError is returned when the password does not contain a special character.
type PasswordNoSpecialCharError struct {
}

func (e PasswordNoSpecialCharError) Error() string {
	return "Password does not contain a special character."
}

// PasswordInvalidCharacter is returned when the password contains an invalid character.
type PasswordInvalidCharacter struct {
}

func (e PasswordInvalidCharacter) Error() string {
	return "Password contains an invalid character."
}
