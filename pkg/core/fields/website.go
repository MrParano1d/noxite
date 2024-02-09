package fields

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Website url.URL

func (w Website) String() string {
	uri := url.URL(w)
	u := &uri
	return u.String()
}

func (w *Website) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	website, err := WebsiteFromString(s)
	if err != nil {
		return err
	}
	*w = website
	return nil
}

func (w Website) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.String())
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
