package charlatan

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertNextToken(t *testing.T, l *lexer, expected tokenType) {
	tok, err := l.NextToken()
	assert.Nil(t, err)
	assert.Equal(t, expected, tok.Type,
		fmt.Sprintf("Expected %s, got %s", expected, tok.Type))
}

func assertNextTokens(t *testing.T, l *lexer, expected ...tokenType) {
	for _, e := range expected {
		assertNextToken(t, l, e)
	}
}

func TestLexerEOF(t *testing.T) {
	l := lexerFromString("")
	assertNextToken(t, l, tokEnd)
}

func TestLexerSimpleSelectFrom(t *testing.T) {
	l := lexerFromString("SELECT foo FROM bar")
	assertNextTokens(t, l, tokSelect, tokField, tokFrom, tokField, tokEnd)
}

func TestLexerSimpleSelectFromWhere(t *testing.T) {
	l := lexerFromString("SELECT foo FROM bar WHERE 1 > 2")
	assertNextTokens(t, l, tokSelect, tokField, tokFrom, tokField, tokWhere,
		tokInt, tokGt, tokInt, tokEnd)
}

func TestLexerOperators(t *testing.T) {
	for s, tokType := range map[string]tokenType{
		"1 = 2":   tokEq,
		"1 != 2":  tokNeq,
		"1 < 2":   tokLt,
		"1 > 2":   tokGt,
		"1 <= 2":  tokLte,
		"1 >= 2":  tokGte,
		"1 OR 2":  tokOr,
		"1 || 2":  tokOr,
		"1 AND 2": tokAnd,
		"1 && 2":  tokAnd,
		"1 IN 2":  tokIn,
	} {
		l := lexerFromString(s)
		assertNextTokens(t, l, tokInt, tokType, tokInt, tokEnd)
	}
}

func TestLexerStringDoubleQuotes(t *testing.T) {
	l := lexerFromString(`"some string"`)
	assertNextTokens(t, l, tokString, tokEnd)
}

func TestLexerStringSingleQuotes(t *testing.T) {
	l := lexerFromString(`'some string'`)
	assertNextTokens(t, l, tokString, tokEnd)
}
