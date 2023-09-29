package services

import (
	"context"
	"fmt"
	"io"

	"github.com/mrparano1d/noxite/pkg/core/coreerrors"
	"github.com/mrparano1d/noxite/pkg/core/entities"
	"github.com/mrparano1d/noxite/pkg/core/fields"
	"github.com/mrparano1d/noxite/pkg/core/ports"
)

type PackageService struct {
	packageAdapter ports.PackagePort
	storageAdapter ports.StoragePort
}

func NewPackageService(
	packageAdapter ports.PackagePort,
	storageAdapter ports.StoragePort,
) *PackageService {
	return &PackageService{
		packageAdapter: packageAdapter,
		storageAdapter: storageAdapter,
	}
}

// usecases

func (s *PackageService) ParseManifest(ctx context.Context, user *entities.User, r io.Reader) (*entities.PackageVersion, error) {

	if user.Role.Permissions.GetPackage == false {
		return nil, &coreerrors.NotAllowedToGetPackageError{}
	}

	manifest, err := s.packageAdapter.ParseManifest(ctx, r)
	if err != nil {
		return nil, handlePackageErrors(err)
	}
	return manifest, nil
}

func (s *PackageService) PublishPackage(ctx context.Context, user *entities.User, manifest *entities.PackageVersion) error {

	if user.Role.Permissions.PublishPackage == false {
		return &coreerrors.NotAllowedToPublishPackageError{}
	}

	if err := s.storageAdapter.PublishPackage(ctx, user.ID, manifest); err != nil {
		return handlePackageErrors(err)
	}
	return nil
}

func (s *PackageService) GetPackage(ctx context.Context, user *entities.User, name string, version string) ([]byte, error) {
	if user.Role.Permissions.GetPackage == false {
		return nil, &coreerrors.NotAllowedToGetPackageError{}
	}

	packageName, err := fields.PackageNameFromString(name)
	if err != nil {
		return nil, &InvalidGetPackageFieldError{
			Field:  "name",
			Reason: err.Error(),
		}
	}

	packageVersion, err := fields.RequiredStringFromString(version)
	if err != nil {
		return nil, &InvalidGetPackageFieldError{
			Field:  "version",
			Reason: err.Error(),
		}
	}

	data, err := s.storageAdapter.GetPackage(ctx, packageName, packageVersion)
	if err != nil {
		return nil, handlePackageErrors(err)
	}

	return data, nil
}

// errors

type PackageServiceManifestParseError struct {
	Err error
}

func (e *PackageServiceManifestParseError) Error() string {
	return "failed to parse manifest: " + e.Err.Error()
}

type PackageServicePublishPackageError struct {
	Err error
}

func (e *PackageServicePublishPackageError) Error() string {
	return "failed to publish package: " + e.Err.Error()
}

type PackageServiceUnknownError struct {
	Err error
}

func (e *PackageServiceUnknownError) Error() string {
	return fmt.Sprintf("unknown package service error: %s", e.Err)
}

type PackageServicePackageNotFoundError struct {
	Name    string
	Version string
}

func (e *PackageServicePackageNotFoundError) Error() string {
	return fmt.Sprintf("package %s@%s not found", e.Name, e.Version)
}

type PackageServiceGetPackageError struct {
	Name    string
	Version string
	Err     error
}

func (e *PackageServiceGetPackageError) Error() string {
	return fmt.Sprintf("failed to get package %s@%s: %s", e.Name, e.Version, e.Err)
}

type InvalidGetPackageFieldError struct {
	Field  string
	Reason string
}

func (e *InvalidGetPackageFieldError) Error() string {
	return fmt.Sprintf("invalid get package field %s: %s", e.Field, e.Reason)
}

// service errors

func handlePackageErrors(err error) error {
	switch e := err.(type) {
	case *ports.PackageAdapterManifestParseError:
		return &PackageServiceManifestParseError{
			Err: e.Err,
		}
	case *ports.StorageAdapterPublishPackageError:
		return &PackageServicePublishPackageError{
			Err: e.Err,
		}
	case *ports.StorageAdapterPackageNotFoundError:
		return &PackageServicePackageNotFoundError{
			Name:    e.Name.String(),
			Version: e.Version.String(),
		}
	case *ports.StorageAdapterGetPackageError:
		return &PackageServiceGetPackageError{
			Name:    e.Name.String(),
			Version: e.Version.String(),
			Err:     e.Err,
		}
	default:
		return &PackageServiceUnknownError{
			Err: e,
		}
	}
}
