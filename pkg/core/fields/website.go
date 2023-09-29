package fields

import (
	"fmt"
	"net/url"
)

type Website url.URL

func (w Website) String() string {
	uri := url.URL(w)
	u := &uri
	return u.String()
}

// converters

func WebsiteFromString(s string) (Website, error) {
	u, err := url.Parse(s)
	if err != nil {
		return Website{}, &WebsiteValidationError{URL: s, Err: err}
	}
	return Website(*u), nil
}

// errors

type WebsiteValidationError struct {
	URL string
	Err error
}

func (e *WebsiteValidationError) Error() string {
	return fmt.Sprintf("invalid website: %s: %s", e.URL, e.Err.Error())
}
