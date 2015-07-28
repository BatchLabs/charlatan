package record

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	ch "github.com/batchlabs/charlatan"
)

// JSONRecord is a record for JSON objects
type JSONRecord struct {
	attrs map[string]json.RawMessage
}

var _ ch.Record = &JSONRecord{}

// NewJSONRecordFromReader creates a new JSONRecord from a JSON decoder
func NewJSONRecordFromReader(r io.Reader) (*JSONRecord, error) {
	attrs := make(map[string]json.RawMessage)

	if err := json.NewDecoder(r).Decode(&attrs); err != nil {
		return nil, err
	}

	return &JSONRecord{attrs: attrs}, nil
}

// Find implements the charlatan.Record interface
func (r *JSONRecord) Find(field *ch.Field) (*ch.Const, error) {
	var partial json.RawMessage
	var ok bool

	attrs := r.attrs
	parts := strings.Split(field.Name(), ".")

	for i, k := range parts {
		partial, ok = attrs[k]
		if !ok {
			return nil, fmt.Errorf("Unknown '%s' field (in '%s')", k, field.Name())
		}

		// update the attrs if we need to go deeper
		if i < len(parts)-1 {
			attrs = make(map[string]json.RawMessage)
			if err := json.Unmarshal(partial, &attrs); err != nil {
				return nil, err
			}
		}

	}

	return jsonToConst(partial)
}

func jsonToConst(partial json.RawMessage) (*ch.Const, error) {
	var value string

	asString := string(partial)

	if asString == "null" {
		return ch.NullConst(), nil
	}

	if err := json.Unmarshal(partial, &value); err != nil {
		if err, ok := err.(*json.UnmarshalTypeError); ok {
			// we failed to unmarshal into a string, let's try the other types
			switch err.Value {
			case "number":
				var n json.Number
				if err := json.Unmarshal(partial, &n); err != nil {
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
