package graph

import (
	"encoding/json"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/mrparano1d/noxite/pkg/core/fields"
)

func MarshalRequiredMap(rm map[fields.RequiredString]fields.RequiredString) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		json.NewEncoder(w).Encode(rm)
	})
}

func UnmarshalRequiredMap(v any) (map[fields.RequiredString]fields.RequiredString, error) {
	switch v := v.(type) {
	case map[fields.RequiredString]fields.RequiredString:
		return v, nil
	case *map[fields.RequiredString]fields.RequiredString:
		if v == nil {
			return map[fields.RequiredString]fields.RequiredString{}, nil
		}
		return *v, nil
	case []byte:
		var rm map[fields.RequiredString]fields.RequiredString
		err := json.Unmarshal(v, &rm)
		return rm, err
	default:
		return map[fields.RequiredString]fields.RequiredString{}, nil
	}
}

func MarshalRequiredMapRequiredKeyMap(rm map[fields.RequiredString]map[fields.RequiredString]interface{}) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		json.NewEncoder(w).Encode(rm)
	})
}

func UnmarshalRequiredMapRequiredKeyMap(v any) (map[fields.RequiredString]map[fields.RequiredString]interface{}, error) {
	switch v := v.(type) {
	case map[fields.RequiredString]map[fields.RequiredString]interface{}:
		return v, nil
	case *map[fields.RequiredString]map[fields.RequiredString]interface{}:
		if v == nil {
			return map[fields.RequiredString]map[fields.RequiredString]interface{}{}, nil
		}
		return *v, nil
	case []byte:
		var rm map[fields.RequiredString]map[fields.RequiredString]interface{}
		err := json.Unmarshal(v, &rm)
		return rm, err
	default:
		return map[fields.RequiredString]map[fields.RequiredString]interface{}{}, nil
	}
}

func MarshalRequiredKeyMap(rm map[fields.RequiredString]interface{}) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		json.NewEncoder(w).Encode(rm)
	})
}

func UnmarshalRequiredKeyMap(v any) (map[fields.RequiredString]interface{}, error) {
	switch v := v.(type) {
	case map[fields.RequiredString]interface{}:
		return v, nil
	case *map[fields.RequiredString]interface{}:
		if v == nil {
			return map[fields.RequiredString]interface{}{}, nil
		}
		return *v, nil
	case []byte:
		var rm map[fields.RequiredString]interface{}
		err := json.Unmarshal(v, &rm)
		return rm, err
	default:
		return map[fields.RequiredString]interface{}{}, nil
	}
}

func MarshalRequiredList(rl []fields.RequiredString) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		json.NewEncoder(w).Encode(rl)
	})
}

func UnmarshalRequiredList(v any) ([]fields.RequiredString, error) {
	switch v := v.(type) {
	case []fields.RequiredString:
		return v, nil
	case *[]fields.RequiredString:
		if v == nil {
			return []fields.RequiredString{}, nil
		}
		return *v, nil
	case []byte:
		var rl []fields.RequiredString
		err := json.Unmarshal(v, &rl)
		return rl, err
	default:
		return []fields.RequiredString{}, nil
	}
}
