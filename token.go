package charlatan

import (
	"fmt"
)

// TokenType is a token type
type TokenType int

// token types are the types of the tokens used by the lexer
const (
	TokInvalid TokenType = iota

	// TokField is an alpha-numeric field
	TokField

	// keywords

	tokKeywordStart
	TokSelect   // SELECT
	TokFrom     // FROM
	TokWhere    // WHERE
	TokStarting // STARTING
	TokAt       // AT
	tokKeywordEnd

	// operators

	tokLogicalOperatorStart
	TokAnd // && or AND
	TokOr  // || or OR
	tokLogicalOperatorEnd

	tokComparisonOperatorStart
	TokEq  // =
	TokNeq // !=
	TokLt  // <
	TokLte // <=
	TokGt  // >
	TokGte // >=
	tokComparisonOperatorEnd

	TokInt
	TokFloat
	TokString

	// special values

	TokTrue  // true
	TokFalse // false
	TokNull  // null

	// misc

	TokLeftParenthesis  // (
	TokRightParenthesis // )
	TokComma            // ,

	// TokEnd is the end token
	TokEnd = -1
)

// Token is the token found by the lexer
type Token struct {
	// the token type
	Type TokenType
	// the string value of this token
	Value string
	// the position into the parsed string
	Pos int
}

// Const returns the token's value as a Const
func (tok Token) Const() (*Const, error) {
	switch tok.Type {
	case TokTrue:
		return BoolConst(true), nil
	case TokFalse:
		return BoolConst(false), nil
	case TokNull:
		return NullConst(), nil
	case TokString:
		return StringConst(tok.Value), nil
	case TokInt, TokFloat:
		return ConstFromString(tok.Value), nil
	default:
		return nil, fmt.Errorf("Token type %v isn't a const", tok.Type)
	}
}

// IsEnd checks if the token is a TokEnd
func (tok Token) IsEnd() bool { return tok.Type == TokEnd }

// IsNumeric checks if the token is numeric
func (tok Token) IsNumeric() bool { return tok.Type == TokInt || tok.Type == TokFloat }

// IsKeyword checks if the token is a keyword
func (tok Token) IsKeyword() bool { return tok.Type > tokKeywordStart && tok.Type < tokKeywordEnd }

// IsOperator checks if the token is an operator
func (tok Token) IsOperator() bool { return tok.IsLogicalOperator() || tok.IsComparisonOperator() }

// IsLogicalOperator checks if the token is a logical operator
func (tok Token) IsLogicalOperator() bool {
	return tok.Type > tokLogicalOperatorStart && tok.Type < tokLogicalOperatorEnd
}

// IsComparisonOperator checks if the token is a comparison operator
func (tok Token) IsComparisonOperator() bool {
	return tok.Type > tokComparisonOperatorStart && tok.Type < tokComparisonOperatorEnd
}

// IsConst tests if the token represents a constant value. If so, one can use
// the Const() method to get the const value.
func (tok Token) IsConst() bool {
	return tok.IsNumeric() || tok.Type == TokString || tok.Type == TokTrue ||
		tok.Type == TokFalse || tok.Type == TokNull
}

// IsField tests if the token represents a field
func (tok Token) IsField() bool {
	return tok.Type == TokField
}

func (tok Token) String() string {
	return fmt.Sprintf("%d:%s(%s)", tok.Pos, tok.Type, tok.Value)
}

func (t TokenType) String() string {
	switch t {
	case TokEnd:
		return "End"
	case TokField:
		return "Field"
	case TokTrue:
		return "True"
	case TokFalse:
		return "False"
	case TokNull:
		return "Null"
	case TokInt:
		return "Int"
	case TokFloat:
		return "Float"
	case TokSelect:
		return "Select"
	case TokFrom:
		return "From"
	case TokWhere:
		return "Where"
	case TokStarting:
		return "Starting"
	case TokAt:
		return "At"
	case TokAnd:
		return "And"
	case TokOr:
		return "Or"
	case TokEq:
		return "Eq"
	case TokNeq:
		return "Neq"
	case TokLt:
		return "Lt"
	case TokLte:
		return "Lte"
	case TokGt:
		return "Gt"
	case TokGte:
		return "Gte"
	case TokLeftParenthesis:
		return "TokLeftParenthesis"
	case TokRightParenthesis:
		return "TokRightParenthesis"
	case TokComma:
		return "Comma"
	}

	return "UNKNOWN"
}
