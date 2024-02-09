package fields

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Bugs struct {
	URL   *Website
	Email *Email
}

func (b Bugs) MarshalJSON() ([]byte, error) {
	if b.URL != nil {
		return json.Marshal(b.URL)
	}

	if b.Email != nil {
		return json.Marshal(b.Email)
	}

	return []byte("null"), nil
}

func (b *Bugs) UnmarshalJSON(data []byte) error {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v := v.(type) {
	case string:
		email, err := EmailFromString(v)
		if err != nil {
			website, err := WebsiteFromString(v)
			if err != nil {
				return &InvalidBugsError{
					Reason: err.Error(),
				}
			}
			b.URL = &website
			return nil
		}
		b.Email = &email
		return nil
	case map[string]any:
		if len(v) == 0 {
			*b = Bugs{}
			return nil
		}

		builder := BugsBuilder()

		url, ok := v["url"].(string)
		if ok {
			builder.URL(url)
		}

		email, ok := v["email"].(string)
		if ok {
			builder.Email(email)
		}

		bugs, err := builder.Build()
		if err != nil {
			return err
		}

		*b = bugs
		return nil
	default:
		return &InvalidBugsError{
			Reason: fmt.Sprintf("%T is not a valid type for Bugs", v),
		}
	}
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
