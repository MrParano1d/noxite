package services

import (
	"fmt"
	"io"

	"github.com/mrparano1d/getregd/pkg/core/entities"
	"github.com/mrparano1d/getregd/pkg/core/ports"
)

type PackageService struct {
	packageAdapter ports.PackagePort
	storageAdapter ports.StoragePort
}

func NewPackageService(packageAdapter ports.PackagePort, storageAdapter ports.StoragePort) *PackageService {
	return &PackageService{
		packageAdapter: packageAdapter,
		storageAdapter: storageAdapter,
	}
}

// usecases

func (s *PackageService) ParseManifest(r io.Reader) (*entities.Manifest, error) {
	manifest, err := s.packageAdapter.ParseManifest(r)
	if err != nil {
		return nil, handlePackageErrors(err)
	}
	return manifest, nil
}

func (s *PackageService) PublishPackage(manifest *entities.Manifest) error {
	if err := s.storageAdapter.PublishPackage(manifest); err != nil {
		return handlePackageErrors(err)
	}
	return nil
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
	default:
		return &PackageServiceUnknownError{
			Err: e,
		}
	}
}
