package graph

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/mrparano1d/noxite/pkg/core/fields"
)

func MarshalBugs(b fields.Bugs) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		json.NewEncoder(w).Encode(b)
	})
}

func UnmarshalBugs(v interface{}) (fields.Bugs, error) {
	switch v := v.(type) {
	case fields.Bugs:
		return v, nil
	case *fields.Bugs:
		if v == nil {
			return fields.Bugs{}, nil
		}
		return *v, nil
	case []byte:
		var b fields.Bugs
		err := json.Unmarshal(v, &b)
		return b, err
	default:
		return fields.Bugs{}, fmt.Errorf("%T is not a valid type for Bugs", v)
	}
}

// scalar MixedAuthor
// scalar MixedAuthors
// scalar UrlType
// scalar Directories
// scalar Repository

func MarshalMixedAuthor(ma fields.MixedAuthor) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		jsonBytes, err := json.Marshal(ma)
		if err != nil {
			w.Write([]byte("null"))
			return
		}

		w.Write(jsonBytes)
	})
}

func UnmarshalMixedAuthor(v interface{}) (fields.MixedAuthor, error) {
	switch v := v.(type) {
	case fields.MixedAuthor:
		return v, nil
	case *fields.MixedAuthor:
		if v == nil {
			return fields.MixedAuthor{}, nil
		}
		return *v, nil
	case []byte:
		var ma fields.MixedAuthor
		err := json.Unmarshal(v, &ma)
		return ma, err
	default:
		return fields.MixedAuthor{}, fmt.Errorf("%T is not a valid type for MixedAuthor", v)
	}
}

func MarshalMixedAuthors(ma fields.MixedAuthors) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		json.NewEncoder(w).Encode(ma)
	})
}

func UnmarshalMixedAuthors(v interface{}) (fields.MixedAuthors, error) {
	switch v := v.(type) {
	case fields.MixedAuthors:
		return v, nil
	case *fields.MixedAuthors:
		if v == nil {
			return fields.MixedAuthors{}, nil
		}
		return *v, nil
	case []byte:
		var ma fields.MixedAuthors
		err := json.Unmarshal(v, &ma)
		return ma, err
	default:
		return fields.MixedAuthors{}, fmt.Errorf("%T is not a valid type for MixedAuthors", v)
	}
}

func MarshalUrlType(u fields.UrlType) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		json.NewEncoder(w).Encode(u)
	})
}

func UnmarshalUrlType(v interface{}) (fields.UrlType, error) {
	switch v := v.(type) {
	case fields.UrlType:
		return v, nil
	case *fields.UrlType:
		if v == nil {
			return fields.UrlType{}, nil
		}
		return *v, nil
	case []byte:
		var u fields.UrlType
		err := json.Unmarshal(v, &u)
		return u, err
	default:
		return fields.UrlType{}, fmt.Errorf("%T is not a valid type for UrlType", v)
	}
}

func MarshalDirectories(d fields.Directories) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		json.NewEncoder(w).Encode(d)
	})
}

func UnmarshalDirectories(v interface{}) (fields.Directories, error) {
	switch v := v.(type) {
	case fields.Directories:
		return v, nil
	case *fields.Directories:
		if v == nil {
			return fields.Directories{}, nil
		}
		return *v, nil
	case []byte:
		var d fields.Directories
		err := json.Unmarshal(v, &d)
		return d, err
	default:
		return fields.Directories{}, fmt.Errorf("%T is not a valid type for Directories", v)
	}
}

func MarshalRepository(r fields.Repository) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		json.NewEncoder(w).Encode(r)
	})
}

func UnmarshalRepository(v interface{}) (fields.Repository, error) {
	switch v := v.(type) {
	case fields.Repository:
		return v, nil
	case *fields.Repository:
		if v == nil {
			return fields.Repository{}, nil
		}
		return *v, nil
	case []byte:
		var r fields.Repository
		err := json.Unmarshal(v, &r)
		return r, err
	default:
		return fields.Repository{}, fmt.Errorf("%T is not a valid type for Repository", v)
	}
}
