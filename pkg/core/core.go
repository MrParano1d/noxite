package core

import (
	"github.com/mrparano1d/getregd/pkg/core/ports"
	"github.com/mrparano1d/getregd/pkg/core/services"
)

type ApplicationCore struct {
	authService    *services.AuthService
	packageService *services.PackageService
}

func NewCoreApp(
	authAdapter ports.AuthPort,
	packageAdapter ports.PackagePort,
	storageAdapter ports.StoragePort,
) *ApplicationCore {
	return &ApplicationCore{
		authService:    services.NewAuthService(authAdapter),
		packageService: services.NewPackageService(packageAdapter, storageAdapter),
	}
}

func (a *ApplicationCore) AuthService() *services.AuthService {
	return a.authService
}

func (a *ApplicationCore) PackageService() *services.PackageService {
	return a.packageService
}
