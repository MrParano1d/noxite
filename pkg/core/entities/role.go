package entities

import (
	"time"

	"github.com/mrparano1d/getregd/pkg/core/fields"
)

type Role struct {
	ID          fields.EntityID
	Name        fields.RequiredString
	Description string
	Permissions Permissions
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}
