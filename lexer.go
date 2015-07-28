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

// Lexer is a lexer
type Lexer struct {
	r     *bufio.Reader
	index int
}

// LexerFromString creates a new lexer from the given string
func LexerFromString(s string) *Lexer {
	return &Lexer{r: bufio.NewReader(strings.NewReader(s))}
}

func (l *Lexer) readRune() (rune, error) {
	l.index++
	r, _, err := l.r.ReadRune()
	return r, err
}

func (l *Lexer) unread() error {
	l.index--
	return l.r.UnreadRune()
}

func (l *Lexer) skipWhiteSpaces() error {
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
func (l *Lexer) NextToken() (*Token, error) {
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
		return l.simpleToken(TokLeftParenthesis, index)
	case ')':
		return l.simpleToken(TokRightParenthesis, index)
	case ',':
		return l.simpleToken(TokComma, index)
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
		return l.token(TokSelect, k, index)
	case "FROM":
		return l.token(TokFrom, k, index)
	case "WHERE":
		return l.token(TokWhere, k, index)
	case "STARTING":
		return l.token(TokStarting, k, index)
	case "AT":
		return l.token(TokAt, k, index)
	case "AND":
		return l.token(TokAnd, k, index)
	case "OR":
		return l.token(TokOr, k, index)
	}

	// special values
	switch w {
	case "true":
		return l.token(TokTrue, "true", index)
	case "false":
		return l.token(TokFalse, "false", index)
	case "null", "NULL":
		return l.token(TokNull, "null", index)
	}

	if _, err := strconv.ParseInt(w, 10, 64); err == nil {
		return l.token(TokInt, w, index)
	}

	if _, err := strconv.ParseFloat(w, 10); err == nil {
		return l.token(TokFloat, w, index)
	}

	if w != "" {
		return l.token(TokField, w, index)
	}

	// operators

	op, err := l.readOperator()
	if err != nil {
		return nil, err
	}

	switch op {
	case "=":
		return l.token(TokEq, op, index)
	case "!=":
		return l.token(TokNeq, op, index)
	case "<":
		return l.token(TokLt, op, index)
	case ">":
		return l.token(TokGt, op, index)
	case "<=":
		return l.token(TokLte, op, index)
	case ">=":
		return l.token(TokGte, op, index)
	case "&&":
		return l.token(TokAnd, op, index)
	case "||":
		return l.token(TokOr, op, index)
	}

	if op != "" {
		return nil, fmt.Errorf("Invalid operator '%s'", op)
	}

	return nil, fmt.Errorf("No known alternative at input %d", index)
}

func (l *Lexer) token(typ TokenType, v string, index int) (*Token, error) {
	return &Token{Type: typ, Value: v, Pos: index}, nil
}

func (l *Lexer) eof() (*Token, error) {
	return l.token(TokEnd, "", l.index)
}

func (l *Lexer) field(v string, index int) (*Token, error) {
	return l.token(TokField, v, index)
}

func (l *Lexer) str(v string, index int) (*Token, error) {
	return l.token(TokString, v, index)
}

func (l *Lexer) simpleToken(typ TokenType, index int) (*Token, error) {
	return l.token(typ, "", index)
}

func (l *Lexer) readUntil(delim rune) (string, error) {
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

func (l *Lexer) readWord() (string, error)     { return l.readWhile(isWordRune) }
func (l *Lexer) readOperator() (string, error) { return l.readWhile(isOperatorRune) }

func (l *Lexer) readWhile(cond func(rune) bool) (string, error) {
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
