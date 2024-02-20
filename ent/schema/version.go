package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/mrparano1d/noxite/pkg/core/fields"
	"github.com/mrparano1d/noxite/pkg/graphql"
)

type Version struct {
	ent.Schema
}

// Annotations of the Role.
func (Version) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField().Directives(graphql.AuthDirective(graphql.RoleRestricted)),
		entgql.MultiOrder(),
		entgql.RelayConnection(),
	}
}

func (v Version) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.String("version").NotEmpty(),
		field.String("description").Optional(),
		field.Strings("keywords").Optional().Default([]string{}),
		field.String("homepage").Optional(),
		field.JSON("bugs", &fields.Bugs{}).Optional(),
		field.String("license").Optional(),
		field.JSON("author", &fields.MixedAuthor{}).Optional(),
		field.JSON("contributors", fields.MixedAuthors{}).Optional(),
		field.JSON("funding", []fields.UrlType{}).Optional(),
		field.Strings("files").Optional().Default([]string{}),
		field.String("main").Optional(),
		field.String("browser").Optional(),
		field.JSON("bin", map[fields.RequiredString]fields.RequiredString{}).Optional().Annotations(entgql.Type("RequiredMap")),
		field.JSON("man", []fields.RequiredString{}).Optional().Annotations(entgql.Type("RequiredList")),
		field.JSON("directories", &fields.Directories{}).Optional(),
		field.JSON("repository", &fields.Repository{}).Optional(),
		field.JSON("scripts", map[fields.RequiredString]fields.RequiredString{}).Optional().Annotations(entgql.Type("RequiredMap")),
		field.JSON("config", map[fields.RequiredString]fields.RequiredString{}).Optional().Annotations(entgql.Type("RequiredMap")),
		field.JSON("dependencies", map[fields.RequiredString]fields.RequiredString{}).Optional().Annotations(entgql.Type("RequiredMap")),
		field.JSON("dev_dependencies", map[fields.RequiredString]fields.RequiredString{}).Optional().Annotations(entgql.Type("RequiredMap")),
		field.JSON("peer_dependencies", map[fields.RequiredString]fields.RequiredString{}).Optional().Annotations(entgql.Type("RequiredMap")),
		field.JSON("peer_dependencies_meta", map[fields.RequiredString]map[fields.RequiredString]interface{}{}).Optional().Annotations(entgql.Type("RequiredMapRequiredKeyMap")),
		field.Strings("bundled_dependencies").Optional().Default([]string{}),
		field.JSON("optional_dependencies", map[fields.RequiredString]fields.RequiredString{}).Optional().Annotations(entgql.Type("RequiredMap")),
		field.JSON("overrides", map[fields.RequiredString]fields.RequiredString{}).Optional().Annotations(entgql.Type("RequiredMap")),
		field.JSON("engines", map[fields.RequiredString]fields.RequiredString{}).Optional().Annotations(entgql.Type("RequiredMap")),
		field.Strings("os").Optional().Default([]string{}),
		field.Strings("cpu").Optional().Default([]string{}),
		field.Bool("private").Optional().Default(false),
		field.JSON("publish_config", map[fields.RequiredString]interface{}{}).Optional().Annotations(entgql.Type("RequiredKeyMap")),
		field.Strings("workspaces").Optional().Default([]string{}),
		field.String("readme").Optional(),
		field.String("content_type"),
		field.String("data"),
		field.String("integrity"),
		field.String("shasum"),
		field.Int("length"),
		field.Int("publisher_id"),
		field.Int("package_id"),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Optional().Nillable(),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

func (v Version) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("publisher", User.Type).Ref("publishes").Unique().Required().Field("publisher_id"),
		edge.From("package", RepoPackage.Type).Ref("versions").Unique().Required().Field("package_id"),
	}
}
