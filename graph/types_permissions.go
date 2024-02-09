package graph

import (
	"encoding/json"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/mrparano1d/noxite/pkg/core/entities"
)

func MarshalPermissions(p entities.Permissions) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		json.NewEncoder(w).Encode(p)
	})
}

func UnmarshalPermissions(v any) (entities.Permissions, error) {
	switch v := v.(type) {
	case entities.Permissions:
		return v, nil
	case *entities.Permissions:
		if v == nil {
			return entities.Permissions{}, nil
		}
		return *v, nil
	case []byte:
		var p entities.Permissions
		err := json.Unmarshal(v, &p)
		return p, err
	default:
		return entities.Permissions{}, nil
	}
}
