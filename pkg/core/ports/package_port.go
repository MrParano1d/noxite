package ports

import (
	"fmt"
	"io"

	"github.com/mrparano1d/getregd/pkg/core/entities"
)

type PackagePort interface {
	ParseManifest(r io.Reader) (*entities.Manifest, error)
}

// errors

type PackageAdapterManifestParseError struct {
	Err error
}

func (e *PackageAdapterManifestParseError) Error() string {
	return fmt.Sprintf("failed to parse manifest: %s", e.Err)
}
