package core

import (
	"github.com/mrparano1d/getregd/pkg/core/ports"
	"github.com/mrparano1d/getregd/pkg/core/services"
)

type ApplicationCore struct {
	authService    *services.AuthService
	packageService *services.PackageService
	sessionService *services.SessionService
}

func NewCoreApp(
	sessionAdapter ports.SessionPort,
	authAdapter ports.AuthPort,
	packageAdapter ports.PackagePort,
	storageAdapter ports.StoragePort,
) *ApplicationCore {

	sessService := services.NewSessionService(sessionAdapter)

	return &ApplicationCore{
		authService:    services.NewAuthService(authAdapter, sessService),
		packageService: services.NewPackageService(packageAdapter, storageAdapter),
		sessionService: sessService,
	}
}

func (a *ApplicationCore) AuthService() *services.AuthService {
	return a.authService
}

func (a *ApplicationCore) PackageService() *services.PackageService {
	return a.packageService
}

func (a *ApplicationCore) SessionService() *services.SessionService {
	return a.sessionService
}
