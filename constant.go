package charlatan

import (
	"fmt"
	"strconv"
	"strings"
)

// constType represents the constant types
type constType int

// types are either null, int, float, bool or string
const (
	constNull constType = iota
	constInt
	constFloat
	constBool
	constString
)

// Const represents a Constant
type Const struct {
	// the type of this constant
	constType constType
	// the values, stored into the right var
	intValue    int64
	floatValue  float64
	boolValue   bool
	stringValue string
}

// NewConst creates a new constant whatever the type is
func NewConst(value interface{}) (*Const, error) {
	if value == nil {
		return NullConst(), nil
	}
	switch value := value.(type) {
	case int:
		return IntConst(int64(value)), nil
	case int8:
		return IntConst(int64(value)), nil
	case int16:
		return IntConst(int64(value)), nil
	case int32:
		return IntConst(int64(value)), nil
	case int64:
		return IntConst(value), nil
	case float64:
		return FloatConst(float64(value)), nil
	case float32:
		return FloatConst(float64(value)), nil
	case bool:
		return BoolConst(value), nil
	case string:
		return StringConst(value), nil

	case *int:
		return IntConst(int64(*value)), nil
	case *int8:
		return IntConst(int64(*value)), nil
	case *int16:
		return IntConst(int64(*value)), nil
	case *int32:
		return IntConst(int64(*value)), nil
	case *int64:
		return IntConst(*value), nil
	case *float32:
		return FloatConst(float64(*value)), nil
	case *float64:
		return FloatConst(*value), nil
	case *bool:
		return BoolConst(*value), nil
	case *string:
		return StringConst(*value), nil
	default:
		return nil, fmt.Errorf("unexpected constant type %T", value)
	}
}

// NullConst returns a new const of type null
func NullConst() *Const {
	return &Const{constType: constNull}
}

// IntConst returns a new Const of type int
func IntConst(value int64) *Const {
	return &Const{intValue: value, constType: constInt}
}

// FloatConst returns a new Const of type float
func FloatConst(value float64) *Const {
	return &Const{floatValue: value, constType: constFloat}
}

// BoolConst returns a new Const of type Bool
func BoolConst(value bool) *Const {
	return &Const{boolValue: value, constType: constBool}
}

// StringConst returns a new Const of type String
func StringConst(value string) *Const {
	return &Const{stringValue: value, constType: constString}
}

// ConstFromString parses a Const from a string
func ConstFromString(s string) *Const {

	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return IntConst(i)
	}

	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return FloatConst(f)
	}

	if b, err := parseBool(s); err == nil {
		return BoolConst(b)
	}

	return StringConst(s)
}

// Evaluate evaluates a const against a record. In practice it always returns a
// pointer on itself
func (c Const) Evaluate(r Record) (*Const, error) {
	return &c, nil
}

func parseBool(s string) (bool, error) {
	switch strings.ToUpper(s) {
	case "TRUE":
		return true, nil
	case "FALSE":
		return false, nil
	default:
		return false, fmt.Errorf("Unrecognized boolean: %s", s)
	}
}

// IsNumeric tests if a const has a numeric type (int or float)
func (c Const) IsNumeric() bool {
	return c.constType == constInt || c.constType == constFloat
}

// IsBool tests if a const is a bool
func (c Const) IsBool() bool {
	return c.constType == constBool
}

// IsString tests if a const is a string
func (c Const) IsString() bool {
	return c.constType == constString
}

// IsNull tests if a const is null
func (c Const) IsNull() bool {
	return c.constType == constNull
}

// Value returns the value of a const
func (c Const) Value() interface{} {
	switch c.constType {
	case constInt:
		return c.intValue
	case constFloat:
		return c.floatValue
	case constBool:
		return c.boolValue
	case constString:
		return c.stringValue
	}
	return nil
}

func (c Const) String() string {
	switch c.constType {
	case constString:
		return "\"" + c.stringValue + "\""
	default:
		return c.AsString()
	}
}

// AsFloat converts into a float64
// Returns 0 if the const is a string or null
func (c Const) AsFloat() float64 {
	switch c.constType {
	case constInt:
		return float64(c.intValue)
	case constFloat:
		return c.floatValue
	case constBool:
		if c.boolValue {
			return 1.0
		}
		return 0.0
	}
	return 0
}

// AsInt converts into an int64
// Returns 0 if the const is a string or null
func (c Const) AsInt() int64 {
	switch c.constType {
	case constInt:
		return c.intValue
	case constFloat:
		return int64(c.floatValue)
	case constBool:
		if c.boolValue {
			return 1
		}
		return 0
	}
	return 0
}

// AsBool converts into a bool
//     - for bool, returns the value
//     - for null, returns false
//     - for numeric, returns true if not 0
//     - for strings, return true (test existence)
func (c Const) AsBool() bool {
	switch c.constType {
	case constNull:
		return false
	case constInt:
		return c.intValue != 0
	case constFloat:
		return c.floatValue != 0
	case constBool:
		return c.boolValue
	case constString:
		return true
	}
	return false
}

// AsString converts into a string
func (c Const) AsString() string {
	switch c.constType {
	case constNull:
		return "null"
	case constInt:
		return strconv.FormatInt(c.intValue, 10)
	case constFloat:
		return strconv.FormatFloat(c.floatValue, 'f', 2, 64)
	case constBool:
		return strconv.FormatBool(c.boolValue)
	case constString:
		return c.stringValue
	}

	// fallback to sprintf .... should never append
	return fmt.Sprintf("%v", c.Value())
}

// CompareTo returns:
//     - a positive integer if this Constant is greater than the given one
//     - a negative integer if this Constant is lower than the given one
//     - zero, is this constant is equals to the given one
//
// If the comparison is not possible (incompatible types), an error will be
// returned
func (c Const) CompareTo(c2 *Const) (int, error) {

	if c.constType == c2.constType {
		switch c.constType {
		case constNull:
			return 0, nil
		case constInt:
			return int(c.intValue - c2.intValue), nil
		case constFloat:
			return int(c.floatValue - c2.floatValue), nil
		case constBool:
			return cmpBools(c.boolValue, c2.boolValue), nil
		case constString:
			return cmpStrings(c.stringValue, c2.stringValue), nil
		default:
			return 0, fmt.Errorf("Unknown const type: %v", c.constType)
		}

	}
	if c.IsNull() {
		return -1, nil
	}
	if c2.IsNull() {
		return 1, nil
	}
	if c.IsNumeric() && c2.IsNumeric() {
		return int(c.AsFloat() - c2.AsFloat()), nil

	}
	if c.IsBool() || c2.IsBool() {
		return cmpBools(c.AsBool(), c2.AsBool()), nil

	}
	if c.IsString() || c2.IsString() {
		return cmpStrings(c.AsString(), c2.AsString()), nil
	}

	return 0, fmt.Errorf("Can't compare the two constants %s(%v), %s(%v)",
		c.constType, c.Value(), c2.constType, c2.Value())
}

func cmpBools(b1, b2 bool) int {
	if b1 == b2 {
		return 0
	}
	if b1 {
		return 1
	}
	return -1
}

func cmpStrings(s1, s2 string) int {
	if s1 == s2 {
		return 0
	}
	if s1 > s2 {
		return 1
	}
	return -1
}

func (t constType) String() string {
	switch t {
	case constNull:
		return "NULL"
	case constInt:
		return "INT"
	case constFloat:
		return "FLOAT"
	case constBool:
		return "BOOL"
	case constString:
		return "STRING"
	default:
		return "UNDEFINED"
	}
}
