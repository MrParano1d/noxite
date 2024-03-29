package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.43

import (
	"context"

	"entgo.io/contrib/entgql"
	"github.com/mrparano1d/noxite/ent"
	"github.com/mrparano1d/noxite/graph"
)

// Node is the resolver for the node field.
func (r *queryResolver) Node(ctx context.Context, id int) (ent.Noder, error) {
	return r.client.Noder(ctx, id)
}

// Nodes is the resolver for the nodes field.
func (r *queryResolver) Nodes(ctx context.Context, ids []int) ([]ent.Noder, error) {
	return r.client.Noders(ctx, ids)
}

// RepoPackages is the resolver for the repoPackages field.
func (r *queryResolver) RepoPackages(ctx context.Context, after *entgql.Cursor[int], first *int, before *entgql.Cursor[int], last *int) (*ent.RepoPackageConnection, error) {
	return r.client.RepoPackage.Query().Paginate(ctx, after, first, before, last)
}

// Roles is the resolver for the roles field.
func (r *queryResolver) Roles(ctx context.Context, after *entgql.Cursor[int], first *int, before *entgql.Cursor[int], last *int) (*ent.RoleConnection, error) {
	return r.client.Role.Query().Paginate(ctx, after, first, before, last)
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, after *entgql.Cursor[int], first *int, before *entgql.Cursor[int], last *int) (*ent.UserConnection, error) {
	return r.client.User.Query().Paginate(ctx, after, first, before, last)
}

// Versions is the resolver for the versions field.
func (r *queryResolver) Versions(ctx context.Context, after *entgql.Cursor[int], first *int, before *entgql.Cursor[int], last *int) (*ent.VersionConnection, error) {
	return r.client.Version.Query().Paginate(ctx, after, first, before, last)
}

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
