package charlatan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenIsEnd(t *testing.T) {
	assert.True(t, Token{Type: TokEnd}.IsEnd())
}

func TestTokenIsNumeric(t *testing.T) {
	assert.True(t, Token{Type: TokInt}.IsNumeric())
	assert.True(t, Token{Type: TokFloat}.IsNumeric())
}

func TestTokenIsKeyword(t *testing.T) {
	assert.True(t, Token{Type: TokSelect}.IsKeyword())
	assert.True(t, Token{Type: TokFrom}.IsKeyword())
	assert.True(t, Token{Type: TokWhere}.IsKeyword())
	assert.True(t, Token{Type: TokStarting}.IsKeyword())

	for _, ty := range []TokenType{
		TokSelect, TokFrom, TokWhere, TokStarting,
	} {
		assert.True(t, Token{Type: ty}.IsKeyword())
	}
}

func TestTokenIsOperator(t *testing.T) {
	for _, ty := range []TokenType{
		TokAnd, TokOr, TokEq, TokNeq, TokLt, TokLte, TokGt, TokGte,
	} {
		assert.True(t, Token{Type: ty}.IsOperator())
	}
}

func TestTokenIsLogicalOperator(t *testing.T) {
	for _, ty := range []TokenType{TokAnd, TokOr} {
		assert.True(t, Token{Type: ty}.IsLogicalOperator())
	}
}

func TestTokenIsComparisonOperator(t *testing.T) {
	for _, ty := range []TokenType{
		TokEq, TokNeq, TokLt, TokLte, TokGt, TokGte,
	} {

		assert.True(t, Token{Type: ty}.IsComparisonOperator())
	}
}

func TestTokenTypeString(t *testing.T) {
	for _, ty := range []TokenType{
		TokField, TokInt, TokFloat, TokTrue, TokFalse, TokNull, TokSelect,
		TokFrom, TokWhere, TokStarting, TokAt, TokAnd, TokOr, TokEq, TokNeq,
		TokLt, TokLte, TokGt, TokGte, TokLeftParenthesis, TokRightParenthesis,
		TokComma, TokEnd,
	} {
		assert.NotEqual(t, "", ty.String())
		assert.NotEqual(t, "UNKNOWN", ty.String())
	}
}
