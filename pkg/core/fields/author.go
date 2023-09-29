package fields

import (
	"fmt"
	"strings"
)

type ForeignAuthor struct {
	Name    RequiredString
	Email   *Email
	Website *Website
}

type MixedAuthor Author[any]
type MixedAuthors []MixedAuthor
type Authors []Author[EntityID]
type ForeignAuthors []Author[ForeignAuthor]

type Author[T ForeignAuthor | EntityID | any] struct {
	value T
}

func (a Author[T]) TryForeignAuthor() (ForeignAuthor, bool) {
	v, ok := any(a.value).(ForeignAuthor)
	return v, ok
}

func (a Author[T]) TryEntityID() (EntityID, bool) {
	v, ok := any(a.value).(EntityID)
	return v, ok
}

func (a Author[T]) ToMixedAuthor() MixedAuthor {
	return MixedAuthor(MixedAuthor{
		value: any(a.value),
	})
}

func (a Author[T]) Value() T {
	return a.value
}

// converters

// author string format: "name <email> (website)"
// <email> and (website) are optional

type foreignAuthorBuilder struct {
	name    string
	email   *string
	website *string
}

func ForeignAuthorBuilder() *foreignAuthorBuilder {
	return &foreignAuthorBuilder{}
}

func (b *foreignAuthorBuilder) Name(name string) *foreignAuthorBuilder {
	b.name = name
	return b
}

func (b *foreignAuthorBuilder) Email(email string) *foreignAuthorBuilder {
	b.email = &email
	return b
}

func (b *foreignAuthorBuilder) Website(website string) *foreignAuthorBuilder {
	b.website = &website
	return b
}

func (b *foreignAuthorBuilder) Build() (Author[ForeignAuthor], error) {

	foreignAuthor := ForeignAuthor{}

	name, err := RequiredStringFromString(b.name)
	if err != nil {
		return Author[ForeignAuthor]{}, err
	}
	foreignAuthor.Name = name

	if b.email != nil {
		email, err := EmailFromString(*b.email)
		if err != nil {
			return Author[ForeignAuthor]{}, err
		}

		foreignAuthor.Email = &email
	}

	if b.website != nil {
		website, err := WebsiteFromString(*b.website)
		if err != nil {
			return Author[ForeignAuthor]{}, err
		}
		foreignAuthor.Website = &website
	}

	return Author[ForeignAuthor]{value: foreignAuthor}, nil
}

func AuthorFromString(s string) (Author[ForeignAuthor], error) {
	parts := strings.Split(s, " ")
	if len(parts) < 1 || len(parts) > 3 {
		return Author[ForeignAuthor]{}, &AuthorConversionError{Author: s, Err: fmt.Errorf("invalid author format")}
	}

	name := parts[0]
	email := ""
	website := ""

	if len(parts) > 1 {
		email = strings.Trim(parts[1], "<>")

		if email == "" {
			website = strings.Trim(parts[1], "()")
		}
	}

	if len(parts) > 2 {
		website = strings.Trim(parts[2], "()")
	}

	authorName, err := RequiredStringFromString(name)
	if err != nil {
		return Author[ForeignAuthor]{}, &AuthorConversionError{Author: s, Err: err}
	}

	var authorEmail *Email

	if email != "" {
		email, err := EmailFromString(email)
		if err != nil {
			return Author[ForeignAuthor]{}, &AuthorConversionError{Author: s, Err: err}
		}
		authorEmail = &email
	}

	var authorWebsite *Website

	if website != "" {
		website, err := WebsiteFromString(website)
		if err != nil {
			return Author[ForeignAuthor]{}, &AuthorConversionError{Author: s, Err: err}
		}
		authorWebsite = &website
	}

	return Author[ForeignAuthor]{value: ForeignAuthor{Name: authorName, Email: authorEmail, Website: authorWebsite}}, nil
}

func AuthorFromEntityID(id EntityID) Author[EntityID] {
	return Author[EntityID]{value: id}
}

func AuthorFromIDString(id string) (Author[EntityID], error) {
	entityID, err := EntityIDFromString(id)
	if err != nil {
		return Author[EntityID]{}, &AuthorConversionError{Author: id, Err: err}
	}

	return AuthorFromEntityID(entityID), nil
}

// errors

type AuthorConversionError struct {
	Author string
	Err    error
}

func (e *AuthorConversionError) Error() string {
	return fmt.Sprintf("invalid author: %s: %s", e.Author, e.Err.Error())
}
