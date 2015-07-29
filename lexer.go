package charlatan

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

// lexer is a lexer
type lexer struct {
	r     *bufio.Reader
	index int
}

// lexerFromString creates a new lexer from the given string
func lexerFromString(s string) *lexer {
	return &lexer{r: bufio.NewReader(strings.NewReader(s))}
}

func (l *lexer) readRune() (rune, error) {
	l.index++
	r, _, err := l.r.ReadRune()
	return r, err
}

func (l *lexer) unread() error {
	l.index--
	return l.r.UnreadRune()
}

func (l *lexer) skipWhiteSpaces() error {
	for {
		r, err := l.readRune()
		if err != nil {
			return err
		}
		if !unicode.IsSpace(r) {
			l.unread()
			break
		}
	}
	return nil
}

// NextToken readss the next token and returns it
func (l *lexer) NextToken() (*token, error) {
	if err := l.skipWhiteSpaces(); err != nil {
		if err == io.EOF {
			return l.eof()
		}

		return nil, err
	}

	r, err := l.readRune()
	if err != nil {
		return nil, err
	}

	index := l.index

	// delimiters: `, ", '
	switch r {
	case '`', '"', '\'':
		v, err := l.readUntil(r)
		if err != nil {
			return nil, err
		}
		// consume the trailing delimiter
		if _, err := l.readRune(); err != nil {
			return nil, err
		}
		// `foo`
		if r == '`' {
			return l.field(v, index)
		}
		// "foo" or 'foo'
		return l.str(v, index)
	case '(':
		return l.simpleToken(tokLeftParenthesis, index)
	case ')':
		return l.simpleToken(tokRightParenthesis, index)
	case ',':
		return l.simpleToken(tokComma, index)
	}

	if err := l.unread(); err != nil {
		return nil, err
	}

	// one char is no longer enough, read the next word instead
	w, err := l.readWord()
	if err != nil {
		return nil, err
	}

	// keywords
	switch k := strings.ToUpper(w); k {
	case "SELECT":
		return l.token(tokSelect, k, index)
	case "FROM":
		return l.token(tokFrom, k, index)
	case "WHERE":
		return l.token(tokWhere, k, index)
	case "STARTING":
		return l.token(tokStarting, k, index)
	case "AT":
		return l.token(tokAt, k, index)
	case "AND":
		return l.token(tokAnd, k, index)
	case "OR":
		return l.token(tokOr, k, index)
	}

	// special values
	switch w {
	case "true":
		return l.token(tokTrue, "true", index)
	case "false":
		return l.token(tokFalse, "false", index)
	case "null", "NULL":
		return l.token(tokNull, "null", index)
	}

	if _, err := strconv.ParseInt(w, 10, 64); err == nil {
		return l.token(tokInt, w, index)
	}

	if _, err := strconv.ParseFloat(w, 10); err == nil {
		return l.token(tokFloat, w, index)
	}

	if w != "" {
		return l.token(tokField, w, index)
	}

	// operators

	op, err := l.readOperator()
	if err != nil {
		return nil, err
	}

	switch op {
	case "=":
		return l.token(tokEq, op, index)
	case "!=":
		return l.token(tokNeq, op, index)
	case "<":
		return l.token(tokLt, op, index)
	case ">":
		return l.token(tokGt, op, index)
	case "<=":
		return l.token(tokLte, op, index)
	case ">=":
		return l.token(tokGte, op, index)
	case "&&":
		return l.token(tokAnd, op, index)
	case "||":
		return l.token(tokOr, op, index)
	}

	if op != "" {
		return nil, fmt.Errorf("Invalid operator '%s'", op)
	}

	return nil, fmt.Errorf("No known alternative at input %d", index)
}

func (l *lexer) token(typ tokenType, v string, index int) (*token, error) {
	return &token{Type: typ, Value: v, Pos: index}, nil
}

func (l *lexer) eof() (*token, error) {
	return l.token(tokEnd, "", l.index)
}

func (l *lexer) field(v string, index int) (*token, error) {
	return l.token(tokField, v, index)
}

func (l *lexer) str(v string, index int) (*token, error) {
	return l.token(tokString, v, index)
}

func (l *lexer) simpleToken(typ tokenType, index int) (*token, error) {
	return l.token(typ, "", index)
}

func (l *lexer) readUntil(delim rune) (string, error) {
	var buf bytes.Buffer

	for {
		r, err := l.readRune()
		if err != nil {
			return "", err
		}

		if r == delim {
			l.unread()
			break
		}

		buf.WriteRune(r)
	}

	return buf.String(), nil
}

func (l *lexer) readWord() (string, error)     { return l.readWhile(isWordRune) }
func (l *lexer) readOperator() (string, error) { return l.readWhile(isOperatorRune) }

func (l *lexer) readWhile(cond func(rune) bool) (string, error) {
	var buf bytes.Buffer

	for {
		r, err := l.readRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if !cond(r) {
			l.unread()
			break
		}
		buf.WriteRune(r)
	}

	return buf.String(), nil
}

func isWordRune(r rune) bool {
	return !unicode.IsSpace(r) && strings.IndexRune("(),`'\"|&=!<>", r) == -1
}

func isOperatorRune(r rune) bool {
	return strings.IndexRune("<>=!&|", r) > -1
}
