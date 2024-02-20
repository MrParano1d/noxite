package graphql

import (
	"entgo.io/contrib/entgql"
	"github.com/vektah/gqlparser/v2/ast"
)

func AuthDirective(requires Role) entgql.Directive {
	return entgql.NewDirective("auth", &ast.Argument{Name: "requires", Value: &ast.Value{Raw: requires.String(), Kind: ast.EnumValue}})
}
