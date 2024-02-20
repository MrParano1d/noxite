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

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Annotations of the RepoPackage.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField().Directives(graphql.AuthDirective(graphql.RoleRestricted)),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		entgql.MultiOrder(),
		entgql.RelayConnection(),
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.String("name").NotEmpty().Unique(),
		field.String("email").NotEmpty().Unique(),
		field.String("about_me").Optional().Nillable(),
		field.String("website").Optional().Nillable(),
		field.Bytes("avatar").Optional().Nillable().Annotations(entgql.Type("Bytes")),
		field.Bytes("password").NotEmpty().Sensitive().Annotations(entgql.Skip()),
		field.Int("role_id").Positive(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Optional().Nillable(),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("role", Role.Type).Ref("user_role").Unique().Required().Field("role_id"),
		edge.To("packages", RepoPackage.Type).Annotations(entgql.RelayConnection(), entgql.MultiOrder()),
		edge.To("publishes", Version.Type).Annotations(entgql.RelayConnection(), entgql.MultiOrder()),
	}
}
