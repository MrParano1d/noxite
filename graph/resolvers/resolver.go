package resolvers

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/mrparano1d/noxite/ent"
	"github.com/mrparano1d/noxite/graph"
	"github.com/mrparano1d/noxite/pkg/core"

	noxqgql "github.com/mrparano1d/noxite/pkg/graphql"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	client *ent.Client
	core   *core.ApplicationCore
}

// NewSchema creates a graphql executable schema.
func NewSchema(client *ent.Client, core *core.ApplicationCore) graphql.ExecutableSchema {
	return graph.NewExecutableSchema(graph.Config{
		Resolvers: &Resolver{
			client: client,
			core:   core,
		},
		Directives: graph.DirectiveRoot{
			Auth: func(ctx context.Context, obj any, next graphql.Resolver, requires graph.AuthRole) (res any, err error) {

				if requires == graph.AuthRoleRestricted {
					token, exists := noxqgql.TokenFromContext(ctx)
					if !exists {
						return nil, fmt.Errorf("no token found in context")
					}
					if err := core.SessionService().ValidateToken(ctx, token); err != nil {
						return nil, err
					}
				}

				return next(ctx)
			},
		},
	})
}
