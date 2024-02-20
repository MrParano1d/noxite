package graphql

import "fmt"

type Role string

func (r Role) String() string {
	return string(r)
}

const (
	RoleRestricted Role = "RESTRICTED"
	RolePublic     Role = "PUBLIC"
)

var AllRole = []Role{
	RoleRestricted,
	RolePublic,
}

func RoleFromString(roleStr string) (Role, error) {
	switch Role(roleStr) {
	case RolePublic, RoleRestricted:
		return Role(roleStr), nil
	}
	return "", fmt.Errorf("invalid role: %s", roleStr)
}
