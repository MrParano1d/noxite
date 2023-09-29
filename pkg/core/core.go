package core

import (
	"github.com/mrparano1d/noxite/pkg/core/ports"
	"github.com/mrparano1d/noxite/pkg/core/services"
)

type ApplicationCore struct {
	authService    *services.AuthService
	packageService *services.PackageService
	sessionService *services.SessionService
	userService    *services.UserService
	roleService    *services.RoleService
}

func NewCoreApp(
	sessionAdapter ports.SessionPort,
	authAdapter ports.AuthPort,
	packageAdapter ports.PackagePort,
	storageAdapter ports.StoragePort,
	userAdapter ports.UserPort,
	roleAdapter ports.RolePort,
) *ApplicationCore {

	sessService := services.NewSessionService(sessionAdapter)

	return &ApplicationCore{
		authService:    services.NewAuthService(authAdapter, sessService),
		packageService: services.NewPackageService(packageAdapter, storageAdapter),
		sessionService: sessService,
		userService:    services.NewUserService(userAdapter),
		roleService:    services.NewRoleService(roleAdapter),
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

func (a *ApplicationCore) UserService() *services.UserService {
	return a.userService
}

func (a *ApplicationCore) RoleService() *services.RoleService {
	return a.roleService
}
