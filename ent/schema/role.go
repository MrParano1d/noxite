package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/mrparano1d/noxite/pkg/core/entities"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.String("name").NotEmpty().Unique(),
		field.String("description").Optional(),
		field.JSON("permissions", entities.Permissions{}),
		field.Time("created_at").Immutable().Default(time.Now),
		field.Time("updated_at").Optional().Nillable(),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user_role", User.Type),
	}
}
