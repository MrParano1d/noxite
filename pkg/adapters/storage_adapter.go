package adapters

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	json "github.com/bytedance/sonic"
	"github.com/mrparano1d/getregd/pkg/core/entities"
	"github.com/mrparano1d/getregd/pkg/core/ports"
)

type FSStorageAdapter struct {
	StoragePath string
}

var _ ports.StoragePort = (*FSStorageAdapter)(nil)

func NewFSStorageAdapter(storagePath string) *FSStorageAdapter {

	if err := os.MkdirAll(storagePath, 0755); err != nil {
		panic(fmt.Errorf("failed to create storage directory: %w", err))
	}

	return &FSStorageAdapter{
		StoragePath: storagePath,
	}
}

func (a *FSStorageAdapter) PublishPackage(ctx context.Context, manifest *entities.Manifest) error {

	var version string
	for v := range manifest.Versions {
		version = v
		break
	}

	fileName := url.QueryEscape(manifest.Name) + "_" + version + ".json"

	filePath := filepath.Join(a.StoragePath, fileName)

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return &ports.StorageAdapterPublishPackageError{Err: err}
	}
	defer f.Close()

	if err := json.ConfigDefault.NewEncoder(f).Encode(manifest); err != nil {
		return &ports.StorageAdapterPublishPackageError{Err: err}
	}

	return nil
}
