package entities

import (
	"time"

	"github.com/mrparano1d/noxite/pkg/core/fields"
)

type User struct {
	ID        fields.EntityID
	Role      *Role
	Username  fields.Username
	Email     fields.Email
	Password  fields.Password
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
