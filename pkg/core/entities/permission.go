package entities

var knownPermissions = [11]string{
	"create_user",
	"get_user",
	"update_user",
	"delete_user",
	"create_role",
	"get_role",
	"update_role",
	"delete_role",
	"get_package",
	"publish_package",
	"unpublish_package",
}

type Permissions struct {
	CreateUser bool
	GetUser    bool
	UpdateUser bool
	DeleteUser bool

	CreateRole bool
	GetRole    bool
	UpdateRole bool
	DeleteRole bool

	PublishPackage   bool
	GetPackage       bool
	UpdatePackage    bool
	UnpublishPackage bool
}

func (p *Permissions) FromSlice(permissions []string) {

	for _, permission := range permissions {
		switch permission {
		case "create_user":
			p.CreateUser = true
		case "get_user":
			p.GetUser = true
		case "update_user":
			p.UpdateUser = true
		case "delete_user":
			p.DeleteUser = true
		case "create_role":
			p.CreateRole = true
		case "get_role":
			p.GetRole = true
		case "update_role":
			p.UpdateRole = true
		case "delete_role":
			p.DeleteRole = true
		case "publish_package":
			p.PublishPackage = true
		case "get_package":
			p.GetPackage = true
		case "update_package":
			p.UpdatePackage = true
		case "unpublish_package":
			p.UnpublishPackage = true
		}
	}

}
