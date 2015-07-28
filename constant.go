package charlatan

import (
	"fmt"
	"strconv"
	"strings"
)

// ConstType represents the constant types
type ConstType int

// types are either null, int, float, bool or string
const (
	ConstNull ConstType = iota
	ConstInt
	ConstFloat
	ConstBool
	ConstString
)

// Const represents a Constant
type Const struct {
	// the type of this constant
	constType ConstType
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

// NullConst returns a new const of type ConstNull
func NullConst() *Const {
	return &Const{constType: ConstNull}
}

// IntConst returns a new Const of type ConstInt
func IntConst(value int64) *Const {
	return &Const{intValue: value, constType: ConstInt}
}

// FloatConst returns a new Const of type ConstFloat
func FloatConst(value float64) *Const {
	return &Const{floatValue: value, constType: ConstFloat}
}

// BoolConst returns a new Const of type Bool
func BoolConst(value bool) *Const {
	return &Const{boolValue: value, constType: ConstBool}
}

// StringConst returns a new Const of type String
func StringConst(value string) *Const {
	return &Const{stringValue: value, constType: ConstString}
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

// GetType returns the type of a const
func (c Const) GetType() ConstType {
	return c.constType
}

// IsNumeric tests if a const has a numeric type (int or float)
func (c Const) IsNumeric() bool {
	return c.constType == ConstInt || c.constType == ConstFloat
}

// IsBool tests if a const is of ConstBool type
func (c Const) IsBool() bool {
	return c.constType == ConstBool
}

// IsString tests if a const is of ConstString type
func (c Const) IsString() bool {
	return c.constType == ConstString
}

// IsNull tests if a consts is of ConstNull type
func (c Const) IsNull() bool {
	return c.constType == ConstNull
}

// Value returns the value of a const
func (c Const) Value() interface{} {
	switch c.constType {
	case ConstInt:
		return c.intValue
	case ConstFloat:
		return c.floatValue
	case ConstBool:
		return c.boolValue
	case ConstString:
		return c.stringValue
	}
	return nil
}

func (c Const) String() string {
	return c.AsString()
}

// AsFloat converts into a float64
// Returns 0 for the const types ConstNull and ConstString
func (c Const) AsFloat() float64 {
	switch c.constType {
	case ConstInt:
		return float64(c.intValue)
	case ConstFloat:
		return c.floatValue
	case ConstBool:
		if c.boolValue {
			return 1.0
		}
		return 0.0
	}
	return 0
}

// AsInt converts into an int64
// Returns 0 for the cont types ConstNull and ConstString
func (c Const) AsInt() int64 {
	switch c.constType {
	case ConstInt:
		return c.intValue
	case ConstFloat:
		return int64(c.floatValue)
	case ConstBool:
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
//     - for strings, return true (test existance)
func (c Const) AsBool() bool {
	switch c.constType {
	case ConstNull:
		return false
	case ConstInt:
		return c.intValue != 0
	case ConstFloat:
		return c.floatValue != 0
	case ConstBool:
		return c.boolValue
	case ConstString:
		return true
	}
	return false
}

// AsString converts into a string
func (c Const) AsString() string {
	switch c.constType {
	case ConstNull:
		return "null"
	case ConstInt:
		return strconv.FormatInt(c.intValue, 10)
	case ConstFloat:
		return strconv.FormatFloat(c.floatValue, 'f', 2, 64)
	case ConstBool:
		return strconv.FormatBool(c.boolValue)
	case ConstString:
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
		case ConstNull:
			return 0, nil
		case ConstInt:
			return int(c.intValue - c2.intValue), nil
		case ConstFloat:
			return int(c.floatValue - c2.floatValue), nil
		case ConstBool:
			return cmpBools(c.boolValue, c2.boolValue), nil
		case ConstString:
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

func (t ConstType) String() string {
	switch t {
	case ConstNull:
		return "NULL"
	case ConstInt:
		return "INT"
	case ConstFloat:
		return "FLOAT"
	case ConstBool:
		return "BOOL"
	case ConstString:
		return "STRING"
	default:
		return "UNDEFINED"
	}
}
