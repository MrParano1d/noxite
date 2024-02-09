package resolvers

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/mrparano1d/noxite/ent"
	"github.com/mrparano1d/noxite/graph"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	client *ent.Client
}

// NewSchema creates a graphql executable schema.
func NewSchema(client *ent.Client) graphql.ExecutableSchema {
	return graph.NewExecutableSchema(graph.Config{
		Resolvers: &Resolver{client},
	})
}
