package fields

import (
	"fmt"
	"strings"
)

type Bugs struct {
	URL   *Website
	Email *Email
}

func (b Bugs) String() string {
	var buf strings.Builder

	if b.Email != nil {
		buf.WriteRune('<')
		buf.WriteString(b.Email.String())
		buf.WriteRune('>')
	}

	if b.URL != nil {
		if b.Email != nil {
			buf.WriteRune(' ')
		}
		buf.WriteRune('(')
		buf.WriteString(b.URL.String())
		buf.WriteRune(')')
	}

	return buf.String()
}

// converters

type bugsBuilder struct {
	url   *string
	email *string
}

func BugsBuilder() *bugsBuilder {
	return &bugsBuilder{}
}

func (b *bugsBuilder) URL(url string) *bugsBuilder {
	b.url = &url
	return b
}

func (b *bugsBuilder) Email(email string) *bugsBuilder {
	b.email = &email
	return b
}

func (b *bugsBuilder) Build() (Bugs, error) {
	if b.url == nil && b.email == nil {
		return Bugs{}, &EmptyBugsError{}
	}

	bugs := Bugs{}

	if b.url != nil {
		url, err := WebsiteFromString(*b.url)
		if err != nil {
			return Bugs{}, &InvalidBugsError{
				Reason: err.Error(),
			}
		}
		bugs.URL = &url
	}

	if b.email != nil {
		email, err := EmailFromString(*b.email)
		if err != nil {
			return Bugs{}, &InvalidBugsError{
				Reason: err.Error(),
			}
		}
		bugs.Email = &email
	}

	return bugs, nil
}

// errors

type InvalidBugsError struct {
	Reason string
}

func (e *InvalidBugsError) Error() string {
	return fmt.Sprintf("invalid bugs: %s", e.Reason)
}

type EmptyBugsError struct {
}

func (e *EmptyBugsError) Error() string {
	return "empty bugs"
}
