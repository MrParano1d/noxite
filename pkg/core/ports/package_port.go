package ports

import (
	"context"
	"fmt"
	"io"

	"github.com/mrparano1d/noxite/pkg/core/entities"
)

type PackagePort interface {
	ParseManifest(ctx context.Context, r io.Reader) (*entities.PackageVersion, error)
	SerializeManifest(ctx context.Context, manifest *entities.PackageVersion) ([]byte, error)
}

// errors

type PackageAdapterManifestParseError struct {
	Err error
}

func (e *PackageAdapterManifestParseError) Error() string {
	return fmt.Sprintf("package adapter failed to parse manifest: %s", e.Err)
}

type PackageAdapterManifestConvertError struct {
	Err error
}

func (e *PackageAdapterManifestConvertError) Error() string {
	return fmt.Sprintf("package adapter failed to convert manifest: %s", e.Err)
}
