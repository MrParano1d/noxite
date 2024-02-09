package graph

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

// MarshalBytes transforms []byte to string that can be used in a graphql response
func MarshalBytes(b []byte) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = fmt.Fprintf(w, "%q", string(b))
	})
}

// UnmarshalBytes transforms the graphql request to []byte
func UnmarshalBytes(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case string:
		return []byte(v), nil
	case *string:
		return []byte(*v), nil
	case []byte:
		return v, nil
	case map[string]interface{}:
		bytes := make([]byte, len(v))
		for k, val := range v {
			idx, err := strconv.Atoi(k)
			if err != nil {
				return nil, fmt.Errorf("invalid byte index: %v", err)
			}
			n, ok := val.(json.Number)
			if !ok {
				return nil, fmt.Errorf("invalid byte json value: %T", val)
			}
			i, err := n.Int64()
			if err != nil {
				return nil, fmt.Errorf("invalid byte json int: %v", err)
			}
			if (i < 0) || (i > 255) {
				return nil, fmt.Errorf("invalid byte value: %v", i)
			}
			bytes[idx] = uint8(i)
		}
		return bytes, nil
	case []any:
		bytes := make([]byte, len(v))
		for idx, val := range v {
			n, ok := val.(json.Number)
			if !ok {
				return nil, fmt.Errorf("invalid byte json value: %T", val)
			}
			i, err := n.Int64()
			if err != nil {
				return nil, fmt.Errorf("invalid byte json int: %v", err)
			}
			if (i < 0) || (i > 255) {
				return nil, fmt.Errorf("invalid byte value: %v", i)
			}
			bytes[idx] = uint8(i)
		}
	default:
		return nil, fmt.Errorf("%T is not []byte", v)
	}

	return nil, fmt.Errorf("invalid type for []byte: %T", v)
}
