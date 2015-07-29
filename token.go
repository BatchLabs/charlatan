package charlatan

import (
	"fmt"
)

// tokenType is a token type
type tokenType int

// token types are the types of the tokens used by the lexer
const (
	tokInvalid tokenType = iota

	// tokField is an alpha-numeric field
	tokField

	// keywords

	tokKeywordStart
	tokSelect   // SELECT
	tokFrom     // FROM
	tokWhere    // WHERE
	tokStarting // STARTING
	tokAt       // AT
	tokKeywordEnd

	// operators

	tokLogicalOperatorStart
	tokAnd // && or AND
	tokOr  // || or OR
	tokLogicalOperatorEnd

	tokComparisonOperatorStart
	tokEq  // =
	tokNeq // !=
	tokLt  // <
	tokLte // <=
	tokGt  // >
	tokGte // >=
	tokComparisonOperatorEnd

	tokInt
	tokFloat
	tokString

	// special values

	tokTrue  // true
	tokFalse // false
	tokNull  // null

	// misc

	tokLeftParenthesis  // (
	tokRightParenthesis // )
	tokComma            // ,

	// tokEnd is the end token
	tokEnd = -1
)

// token is the token found by the lexer
type token struct {
	// the token type
	Type tokenType
	// the string value of this token
	Value string
	// the position into the parsed string
	Pos int
}

// Const returns the token's value as a Const
func (tok token) Const() (*Const, error) {
	switch tok.Type {
	case tokTrue:
		return BoolConst(true), nil
	case tokFalse:
		return BoolConst(false), nil
	case tokNull:
		return NullConst(), nil
	case tokString:
		return StringConst(tok.Value), nil
	case tokInt, tokFloat:
		return ConstFromString(tok.Value), nil
	default:
		return nil, fmt.Errorf("Token type %v isn't a const", tok.Type)
	}
}

// isEnd checks if the token is a tokEnd
func (tok token) isEnd() bool { return tok.Type == tokEnd }

// isNumeric checks if the token is numeric
func (tok token) isNumeric() bool { return tok.Type == tokInt || tok.Type == tokFloat }

// isKeyword checks if the token is a keyword
func (tok token) isKeyword() bool { return tok.Type > tokKeywordStart && tok.Type < tokKeywordEnd }

// isOperator checks if the token is an operator
func (tok token) isOperator() bool { return tok.isLogicalOperator() || tok.isComparisonOperator() }

// isLogicalOperator checks if the token is a logical operator
func (tok token) isLogicalOperator() bool {
	return tok.Type > tokLogicalOperatorStart && tok.Type < tokLogicalOperatorEnd
}

// isComparisonOperator checks if the token is a comparison operator
func (tok token) isComparisonOperator() bool {
	return tok.Type > tokComparisonOperatorStart && tok.Type < tokComparisonOperatorEnd
}

// isConst tests if the token represents a constant value. If so, one can use
// the Const() method to get the const value.
func (tok token) isConst() bool {
	return tok.isNumeric() || tok.Type == tokString || tok.Type == tokTrue ||
		tok.Type == tokFalse || tok.Type == tokNull
}

// isField tests if the token represents a field
func (tok token) isField() bool {
	return tok.Type == tokField
}

func (tok token) String() string {
	return fmt.Sprintf("%d:%s(%s)", tok.Pos, tok.Type, tok.Value)
}

func (t tokenType) String() string {
	switch t {
	case tokEnd:
		return "End"
	case tokField:
		return "Field"
	case tokTrue:
		return "True"
	case tokFalse:
		return "False"
	case tokNull:
		return "Null"
	case tokInt:
		return "Int"
	case tokFloat:
		return "Float"
	case tokSelect:
		return "Select"
	case tokFrom:
		return "From"
	case tokWhere:
		return "Where"
	case tokStarting:
		return "Starting"
	case tokAt:
		return "At"
	case tokAnd:
		return "And"
	case tokOr:
		return "Or"
	case tokEq:
		return "Eq"
	case tokNeq:
		return "Neq"
	case tokLt:
		return "Lt"
	case tokLte:
		return "Lte"
	case tokGt:
		return "Gt"
	case tokGte:
		return "Gte"
	case tokLeftParenthesis:
		return "tokLeftParenthesis"
	case tokRightParenthesis:
		return "tokRightParenthesis"
	case tokComma:
		return "Comma"
	}

	return "UNKNOWN"
}
