package record

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	ch "github.com/BatchLabs/charlatan"
)

// JSONRecord is a record for JSON objects.
//
// It supports the special field "*", as in "SELECT * FROM x WHERE y", which
// returns the JSON as-is, except that the keys order is not garanteed.
//
// If the SoftMatching attribute is set to true, non-existing fields are
// returned as null contants instead of failing with an error.
type JSONRecord struct {
	attrs        map[string]*json.RawMessage
	SoftMatching bool
}

var _ ch.Record = &JSONRecord{}

var errEmptyField = errors.New("Empty field name")

// NewJSONRecordFromDecoder creates a new JSONRecord from a JSON decoder
func NewJSONRecordFromDecoder(dec *json.Decoder) (*JSONRecord, error) {
	attrs := make(map[string]*json.RawMessage)

	if err := dec.Decode(&attrs); err != nil {
		return nil, err
	}

	return &JSONRecord{attrs: attrs}, nil
}

// Find implements the charlatan.Record interface
func (r *JSONRecord) Find(field *ch.Field) (*ch.Const, error) {
	var ok bool
	var partial *json.RawMessage
	var name string

	if name = field.Name(); len(name) == 0 {
		return nil, errEmptyField
	}

	// support for "SELECT *"
	if name == "*" {
		b, err := json.Marshal(r.attrs)
		if err != nil {
			return nil, err
		}
		return ch.StringConst(string(b)), nil
	}

	attrs := r.attrs
	parts := strings.Split(name, ".")

	for i, k := range parts {
		partial, ok = attrs[k]

		if !ok {
			if r.SoftMatching {
				return ch.NullConst(), nil
			}

			return nil, fmt.Errorf("Unknown '%s' field (in '%s')", k, field.Name())
		}

		// update the attrs if we need to go deeper
		if i < len(parts)-1 {
			attrs = make(map[string]*json.RawMessage)
			if err := json.Unmarshal(*partial, &attrs); err != nil {
				if r.SoftMatching {
					return ch.NullConst(), nil
				}
				return nil, err
			}
		}
	}

	return jsonToConst(partial)
}

func jsonToConst(partial *json.RawMessage) (*ch.Const, error) {
	var value string

	if partial == nil {
		return ch.NullConst(), nil
	}

	asString := string(*partial)

	if asString == "null" {
		return ch.NullConst(), nil
	}

	if err := json.Unmarshal(*partial, &value); err != nil {
		if err, ok := err.(*json.UnmarshalTypeError); ok {
			// we failed to unmarshal into a string, let's try the other types
			switch err.Value {
			case "number":
				var n json.Number
				if err := json.Unmarshal(*partial, &n); err != nil {
					return nil, err
				}

				value = n.String()

			case "bool":
				value = asString

			default:
				return nil, err
			}
		}
	}

	return ch.ConstFromString(value), nil
}
