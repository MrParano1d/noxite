package ports

import (
	"context"
	"fmt"

	"github.com/mrparano1d/noxite/pkg/core/entities"
	"github.com/mrparano1d/noxite/pkg/core/fields"
)

type StoragePort interface {
	PublishPackage(ctx context.Context, creatorID fields.EntityID, manifest *entities.PackageVersion) error
	GetPackage(ctx context.Context, name fields.PackageName, version fields.RequiredString) ([]byte, error)
}

// errors

type StorageAdapterPublishPackageError struct {
	Err error
}

func (e *StorageAdapterPublishPackageError) Error() string {
	return "storage adapter failed to publish package: " + e.Err.Error()
}

type StorageAdapterPackageNotFoundError struct {
	Name    fields.PackageName
	Version fields.RequiredString
}

func (e *StorageAdapterPackageNotFoundError) Error() string {
	return fmt.Sprintf("storage adapter didn't find package: %s@%s", e.Name, e.Version)
}

type StorageAdapterGetPackageError struct {
	Name    fields.PackageName
	Version fields.RequiredString
	Err     error
}

func (e *StorageAdapterGetPackageError) Error() string {
	return fmt.Sprintf("storage adapter failed to get package: %s@%s: %s", e.Name, e.Version, e.Err)
}
