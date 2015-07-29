package charlatan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenIsEnd(t *testing.T) {
	assert.True(t, token{Type: tokEnd}.isEnd())
}

func TestTokenIsNumeric(t *testing.T) {
	assert.True(t, token{Type: tokInt}.isNumeric())
	assert.True(t, token{Type: tokFloat}.isNumeric())
}

func TestTokenIsKeyword(t *testing.T) {
	assert.True(t, token{Type: tokSelect}.isKeyword())
	assert.True(t, token{Type: tokFrom}.isKeyword())
	assert.True(t, token{Type: tokWhere}.isKeyword())
	assert.True(t, token{Type: tokStarting}.isKeyword())

	for _, ty := range []tokenType{
		tokSelect, tokFrom, tokWhere, tokStarting,
	} {
		assert.True(t, token{Type: ty}.isKeyword())
	}
}

func TestTokenIsOperator(t *testing.T) {
	for _, ty := range []tokenType{
		tokAnd, tokOr, tokEq, tokNeq, tokLt, tokLte, tokGt, tokGte,
	} {
		assert.True(t, token{Type: ty}.isOperator())
	}
}

func TestTokenIsLogicalOperator(t *testing.T) {
	for _, ty := range []tokenType{tokAnd, tokOr} {
		assert.True(t, token{Type: ty}.isLogicalOperator())
	}
}

func TestTokenIsComparisonOperator(t *testing.T) {
	for _, ty := range []tokenType{
		tokEq, tokNeq, tokLt, tokLte, tokGt, tokGte,
	} {

		assert.True(t, token{Type: ty}.isComparisonOperator())
	}
}

func TestTokenTypeString(t *testing.T) {
	for _, ty := range []tokenType{
		tokField, tokInt, tokFloat, tokTrue, tokFalse, tokNull, tokSelect,
		tokFrom, tokWhere, tokStarting, tokAt, tokAnd, tokOr, tokEq, tokNeq,
		tokLt, tokLte, tokGt, tokGte, tokLeftParenthesis, tokRightParenthesis,
		tokComma, tokLeftSquareBracket, tokRightSquareBracket, tokIn, tokEnd,
	} {
		assert.NotEqual(t, "", ty.String())
		assert.NotEqual(t, "UNKNOWN", ty.String())
	}
}
