package fields

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type ForeignAuthor struct {
	Name    RequiredString `json:"name"`
	Email   *Email         `json:"email,omitempty"`
	Website *Website       `json:"website,omitempty"`
}

type MixedAuthor Author[any]

func (a MixedAuthor) MarshalJSON() ([]byte, error) {
	switch v := a.value.(type) {
	case Author[ForeignAuthor]:
		return json.Marshal(v)
	case ForeignAuthor:
		return json.Marshal(v)
	case EntityID:
		return json.Marshal(v)
	default:
		return nil, fmt.Errorf("invalid marshal MixedAuthor type")
	}
}

func (a *MixedAuthor) UnmarshalJSON(data []byte) error {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v := v.(type) {
	case map[string]any:

		if len(v) == 0 {
			*a = MixedAuthor{}
			return nil
		}

		builder := ForeignAuthorBuilder()

		name, ok := v["name"].(string)
		if ok {
			builder.Name(name)
		}

		email, ok := v["email"].(string)
		if ok {
			builder.Email(email)
		}

		website, ok := v["website"].(string)
		if ok {
			builder.Website(website)
		}

		author, err := builder.Build()
		if err != nil {
			return fmt.Errorf("invalid author from json: %v", err)
		}

		a.value = author
	case string:
		author, err := EntityIDFromString(v)
		if err != nil {
			return err
		}
		a.value = author
	case float64:
		author, err := EntityIDFromInt(int(v))
		if err != nil {
			return err
		}
		a.value = author
	default:
		return fmt.Errorf("invalid unmarshal MixedAuthor type: %T", v)
	}

	return nil
}

type MixedAuthors []MixedAuthor

func (a MixedAuthors) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, author := range a {
		if i > 0 {
			buf.WriteString(",")
		}
		authorJSON, err := json.Marshal(author)
		if err != nil {
			return nil, err
		}
		buf.Write(authorJSON)
	}
	buf.WriteString("]")
	return buf.Bytes(), nil
}

type Authors []Author[EntityID]
type ForeignAuthors []Author[ForeignAuthor]

type Author[T ForeignAuthor | EntityID | any] struct {
	value T
}

func (a Author[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.value)
}

func (a *Author[T]) UnmarshalJSON(data []byte) error {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v := v.(type) {
	case map[string]any:

		if len(v) == 0 {
			*a = Author[T]{}
			return nil
		}

		builder := ForeignAuthorBuilder()

		name, ok := v["name"].(string)
		if ok {
			builder.Name(name)
		}

		email, ok := v["email"].(string)
		if ok {
			builder.Email(email)
		}

		website, ok := v["website"].(string)
		if ok {
			builder.Website(website)
		}

		author, err := builder.Build()
		if err != nil {
			return fmt.Errorf("invalid author from json: %v", err)
		}

		a.value = any(author).(T)
	case string:
		author, err := EntityIDFromString(v)
		if err != nil {
			return err
		}
		a.value = any(author).(T)
	default:
		return fmt.Errorf("invalid unmarshal Author type: %T", v)
	}

	return nil
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
	return MixedAuthor{
		value: any(a.value),
	}
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
