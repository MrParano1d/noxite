package fields

import (
	"fmt"
	"net/url"
	"strings"
)

type PackageName string

func (n PackageName) String() string {
	return string(n)
}

// converters

func PackageNameFromString(s string) (PackageName, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return PackageName(""), &InvalidPackageNameError{
			Name:   s,
			Reason: "package name cannot be empty",
		}
	}

	name, err := url.QueryUnescape(s)
	if err != nil {
		return PackageName(""), &InvalidPackageNameError{
			Name:   s,
			Reason: err.Error(),
		}
	}

	return PackageName(name), nil
}

// errors

type InvalidPackageNameError struct {
	Name   string
	Reason string
}

func (e *InvalidPackageNameError) Error() string {
	return fmt.Sprintf("invalid package name %s: %s", e.Name, e.Reason)
}
