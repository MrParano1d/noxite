package fields

import (
	"fmt"
	"strings"
)

type UrlType struct {
	URL  *Website
	Type *RequiredString
}

func (u UrlType) String() string {
	var buf strings.Builder

	if u.Type != nil {
		buf.WriteString(u.Type.String())
		buf.WriteRune(':')
	}

	if u.URL != nil {
		if u.Type != nil {
			buf.WriteRune(' ')
		}
		buf.WriteString(u.URL.String())
	}

	return buf.String()
}

// converters

type urlTypeBuilder struct {
	url *string
	typ *string
}

func UrlTypeBuilder() *urlTypeBuilder {
	return &urlTypeBuilder{}
}

func (u *urlTypeBuilder) URL(url string) *urlTypeBuilder {
	u.url = &url
	return u
}

func (u *urlTypeBuilder) Type(typ string) *urlTypeBuilder {
	u.typ = &typ
	return u
}

func (u *urlTypeBuilder) Build() (UrlType, error) {
	if u.url == nil && u.typ == nil {
		return UrlType{}, &EmptyUrlTypeError{}
	}

	urlType := UrlType{}

	if u.url != nil {
		url, err := WebsiteFromString(*u.url)
		if err != nil {
			return UrlType{}, err
		}
		urlType.URL = &url
	}

	if u.typ != nil {
		typ, err := RequiredStringFromString(*u.typ)
		if err != nil {
			return UrlType{}, err
		}
		urlType.Type = &typ
	}

	return urlType, nil
}

// errors

type EmptyUrlTypeError struct{}

func (e *EmptyUrlTypeError) Error() string {
	return "url type cannot be empty"
}

type UrlTypeError struct {
	Reason string
}

func (e *UrlTypeError) Error() string {
	return fmt.Sprintf("invalid url type: %s", e.Reason)
}
