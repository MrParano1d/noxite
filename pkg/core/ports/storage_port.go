package ports

import (
	"context"

	"github.com/mrparano1d/getregd/pkg/core/entities"
)

type StoragePort interface {
	PublishPackage(ctx context.Context, manifest *entities.Manifest) error
}

// errors

type StorageAdapterPublishPackageError struct {
	Err error
}

func (e *StorageAdapterPublishPackageError) Error() string {
	return "failed to publish package: " + e.Err.Error()
}
