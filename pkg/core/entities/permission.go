package entities

type Permissions struct {
	CreateUser bool
	GetUser    bool
	UpdateUser bool
	DeleteUser bool

	CreateRole bool
	GetRole    bool
	UpdateRole bool
	DeleteRole bool

	PublishPackage bool
	GetPackage     bool
	UpdatePackage  bool
	UnpublishPackage bool
}
