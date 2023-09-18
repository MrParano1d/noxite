package adapters

import (
	"io"

	json "github.com/bytedance/sonic"

	"github.com/mrparano1d/getregd/pkg/core/entities"
	"github.com/mrparano1d/getregd/pkg/core/ports"
)

type PackageAdapter struct {
}

var _ ports.PackagePort = (*PackageAdapter)(nil)

func NewPackageAdapter() *PackageAdapter {
	return &PackageAdapter{}
}

func (a *PackageAdapter) ParseManifest(r io.Reader) (*entities.Manifest, error) {
	var manifest entities.Manifest
	if err := json.ConfigDefault.NewDecoder(r).Decode(&manifest); err != nil {
		return nil, &ports.PackageAdapterManifestParseError{Err: err}
	}
	return &manifest, nil
}
