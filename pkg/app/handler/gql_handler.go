package handler

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/mrparano1d/noxite/ent"
	"github.com/mrparano1d/noxite/graph/resolvers"
	"github.com/mrparano1d/noxite/pkg/core"
)

func GQLHandler(r chi.Router, app *core.ApplicationCore, client *ent.Client) {
	srv := handler.NewDefaultServer(resolvers.NewSchema(client, app))
	r.Handle("/graphql/query", srv)
	r.Handle("/graphql", playground.Handler("Noxite Playground", "/graphql/query"))
}
