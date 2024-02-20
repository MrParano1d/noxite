package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/mrparano1d/noxite/pkg/graphql"
)

type RepoPackage struct {
	ent.Schema
}

// Annotations of the RepoPackage.
func (RepoPackage) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField().Directives(graphql.AuthDirective(graphql.RoleRestricted)),
		entgql.RelayConnection(),
		entgql.MultiOrder(),
	}
}

func (p RepoPackage) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.String("name").Unique().NotEmpty(),
		field.Int("creator_id"),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Optional().Nillable(),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

func (p RepoPackage) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("versions", Version.Type).Annotations(entgql.MultiOrder(), entgql.RelayConnection()),
		edge.From("creator", User.Type).Ref("packages").Unique().Required().Field("creator_id"),
	}
}
