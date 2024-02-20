package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.43

import (
	"context"

	"github.com/mrparano1d/noxite/graph"
)

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, usernameOrEmail string, password string) (*graph.AuthPayload, error) {
	sess, err := r.core.AuthService().Login(ctx, usernameOrEmail, password)
	if err != nil {
		return nil, err
	}
	return &graph.AuthPayload{
		Token:     sess.Token.String(),
		ExpiresAt: sess.ExpiresAt,
	}, nil
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }