package charlatan

// operatorType is the type of an operator
type operatorType int

// operators can be either logical or comparison-al
const (
	operatorInvalid operatorType = iota

	operatorAnd
	operatorOr

	operatorEq
	operatorNeq
	operatorLt
	operatorLte
	operatorGt
	operatorGte
)

// operatorTypeFromTokenType converts a TokenType to an operatorType
func operatorTypeFromTokenType(ty tokenType) operatorType {
	switch ty {
	case tokAnd:
		return operatorAnd
	case tokOr:
		return operatorOr
	case tokEq:
		return operatorEq
	case tokNeq:
		return operatorNeq
	case tokLt:
		return operatorLt
	case tokLte:
		return operatorLte
	case tokGt:
		return operatorGt
	case tokGte:
		return operatorGte
	default:
		return operatorInvalid
	}
}

// IsLogical tests if an operator is logical
func (o operatorType) IsLogical() bool {
	return o == operatorAnd || o == operatorOr
}

// isComparison tests if an operator is a comparison
func (o operatorType) isComparison() bool {
	return !o.IsLogical() && o != operatorInvalid
}

func (o operatorType) String() string {
	switch o {
	case operatorAnd:
		return "&&"
	case operatorOr:
		return "||"
	case operatorEq:
		return "="
	case operatorNeq:
		return "!="
	case operatorLt:
		return "<"
	case operatorLte:
		return "<="
	case operatorGt:
		return ">"
	case operatorGte:
		return ">="
	default:
		return "<unknown operator>"
	}
}
