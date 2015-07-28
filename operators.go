package charlatan

// OperatorType is the type of an operator
type OperatorType int

// Operators can be either logical or comparison-al
const (
	OperatorInvalid OperatorType = iota

	OperatorAnd
	OperatorOr

	OperatorEq
	OperatorNeq
	OperatorLt
	OperatorLte
	OperatorGt
	OperatorGte
)

// OperatorTypeFromTokenType converts a TokenType to an OperatorType
func OperatorTypeFromTokenType(ty TokenType) OperatorType {
	switch ty {
	case TokAnd:
		return OperatorAnd
	case TokOr:
		return OperatorOr
	case TokEq:
		return OperatorEq
	case TokNeq:
		return OperatorNeq
	case TokLt:
		return OperatorLt
	case TokLte:
		return OperatorLte
	case TokGt:
		return OperatorGt
	case TokGte:
		return OperatorGte
	default:
		return OperatorInvalid
	}
}

// IsLogical tests if an operator is logical
func (o OperatorType) IsLogical() bool {
	return o == OperatorAnd || o == OperatorOr
}

// IsComparison tests if an operator is a comparison
func (o OperatorType) IsComparison() bool {
	return !o.IsLogical() && o != OperatorInvalid
}

func (o OperatorType) String() string {
	switch o {
	case OperatorAnd:
		return "&&"
	case OperatorOr:
		return "||"
	case OperatorEq:
		return "="
	case OperatorNeq:
		return "!="
	case OperatorLt:
		return "<"
	case OperatorLte:
		return "<="
	case OperatorGt:
		return ">"
	case OperatorGte:
		return ">="
	default:
		return "<unknown operator>"
	}
}
