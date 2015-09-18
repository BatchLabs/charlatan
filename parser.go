package charlatan

import (
	"errors"
	"fmt"
)

// the automate state
type state int

const (
	invalidState state = -1

	initial state = iota

	selectInitial
	selectField

	fromInitial

	operationInitial
	leftOperand
	operator
	rightOperand

	startingInitial
	startingAt

	limitInitial
	limitSep
	limitMax

	rangeMin
	rangeAnd
	rangeMax

	clauseEnd
	end
)

// parser is the parser itself
type parser struct {
	// the lexer to read tokens
	lexer *lexer
	// the current state
	state state

	// the array of fields,
	// before the query is initialized
	fields []*Field

	// the query
	query *Query

	// the operation context stack
	stack []*context
	// and the current context
	current *context
}

// the context, that contains informations
// about the current operations
type context struct {
	// the current operation
	left     operand
	operator tokenType
	right    operand

	// the logical operation nodes
	first *logicalOperation
	last  *logicalOperation

	// the current range test, if there's one
	rangeTest *rangeTestOperation

	// the current logical operator
	logicalOperator tokenType

	// if stacked, in which state
	// this context should be restored
	// only leftOperand or rightOperand
	stackState state
}

// parserFromString creates a new parser from the given string
func parserFromString(s string) *parser {
	return &parser{
		lexer:   lexerFromString(s),
		state:   initial,
		fields:  make([]*Field, 0),
		stack:   make([]*context, 0),
		current: newContext(),
	}
}

// Parse parses and returns the query
func (p *parser) Parse() (*Query, error) {

	// Read tokens, move step by step
	// until the next state is the end

	for p.state != end {
		tok, err := p.lexer.NextToken()
		if err != nil {
			return nil, err
		}

		switch p.state {

		// the very begining
		case initial:
			p.state, err = p.initialState(tok)

		// SELECT
		case selectInitial:
			p.state, err = p.selectState(tok)
		case selectField:
			p.state, err = p.selectFieldState(tok)

		// FROM
		case fromInitial:
			p.state, err = p.fromState(tok)

		// WHERE
		case operationInitial:
			p.state, err = p.operationState(tok)
		case leftOperand:
			p.state, err = p.operandLeftState(tok)
		case operator:
			p.state, err = p.operatorState(tok)
		case rightOperand:
			p.state, err = p.operandRightState(tok)

		// range
		case rangeMin:
			p.state, err = p.rangeMin(tok)
		case rangeAnd:
			p.state, err = p.expect(tokAnd, tok, rangeMax)
		case rangeMax:
			p.state, err = p.rangeMax(tok)

		// STARTING
		case startingInitial:
			p.state, err = p.expect(tokAt, tok, startingAt)
		// AT
		case startingAt:
			p.state, err = p.startingAt(tok)

		// LIMIT N
		//       ^
		case limitInitial:
			p.state, err = p.limit(tok)

		// LIMIT N, M
		//        ^
		case limitSep:
			p.state, err = p.limitSep(tok)

		// LIMIT N, M
		//          ^
		case limitMax:
			p.state, err = p.limitMax(tok)

		case clauseEnd:
			p.state, err = p.clauseEnd(tok)

		// unknown
		default:
			err = fmt.Errorf("Unknown state %d", p.state)
		}

		// an error occured during handling the state
		if err != nil {
			return nil, err
		}
	}

	// Handle the operation context

	if len(p.stack) > 0 {
		return nil, fmt.Errorf("Unbalanced parenthesis")
	}

	if p.current.first != nil {
		// affect the logical node as the expression
		// FIXME should we finish the operation ???
		p.query.setWhere(p.current.first.simplify())
	}

	return p.query, nil
}

// We’re only waiting for the SELECT keyword
func (p *parser) initialState(tok *token) (state, error) {
	if tok.Type != tokSelect {
		return unexpected(tok, tokSelect)
	}

	return selectInitial, nil
}

// We’re waiting for a field
func (p *parser) selectState(tok *token) (state, error) {
	if tok.Type != tokField {
		return unexpected(tok, tokField)
	}

	p.fields = append(p.fields, NewField(tok.Value))

	return selectField, nil
}

// We’re waiting for a comma, or the FROM keyword
func (p *parser) selectFieldState(tok *token) (state, error) {
	if tok.Type == tokComma {
		// we just jump to the next field
		return selectInitial, nil
	}

	if tok.Type != tokFrom {
		return unexpected(tok, tokFrom)
	}

	return fromInitial, nil
}

// We’re waiting for a name of the from
func (p *parser) fromState(tok *token) (state, error) {
	if tok.Type != tokField {
		return unexpected(tok, tokField)
	}

	p.query = NewQuery(tok.Value)
	p.query.AddFields(p.fields)
	p.fields = nil

	return clauseEnd, nil
}

// We’re waiting for any of the WHERE, STARTING, or LIMIT keywords, or the end
func (p *parser) clauseEnd(tok *token) (state, error) {
	switch tok.Type {
	case tokEnd:
		return end, nil
	case tokWhere:
		return operationInitial, nil
	case tokStarting:
		return startingInitial, nil
	case tokLimit:
		return limitInitial, nil
	default:
		if p.query != nil && p.query.expression != nil {
			if p.query.startingAt == 0 {
				return unexpected(tok, tokStarting)
			}

			if !p.query.HasLimit() {
				return unexpected(tok, tokLimit)
			}
		}

		return unexpected(tok, tokWhere)
	}
}

func tok2operand(tok *token) (operand, error) {
	if tok.isField() {
		return NewField(tok.Value), nil
	}
	if tok.isConst() {
		return tok.Const()
	}
	return nil, fmt.Errorf("Invalid token: not an operand: %v", tok)
}

// We’re waiting for a left operand, or a (
func (p *parser) operationState(tok *token) (state, error) {
	if tok.Type == tokLeftParenthesis {
		// push the context and start a new operation
		p.pushContext(leftOperand)
		return operationInitial, nil
	}

	c, err := tok2operand(tok)
	if err != nil {
		return invalidState, err
	}
	p.current.left = c

	// the left operand has been setted jump to the left operand state
	return leftOperand, nil
}

// We’re waiting for an operator:
//     - logical operator, we step to the next operation
//     - comparison operator, we continue to the operator state
//     - "BETWEEN", we continue to the range operator state
//
// We can encounter a ), or the end
func (p *parser) operandLeftState(tok *token) (state, error) {

	if tok.isEnd() {
		// end the previous operation
		if err := p.current.endOperation(); err != nil {
			return invalidState, err
		}

		return end, nil
	}

	switch tok.Type {
	case tokStarting:
		return startingInitial, nil
	case tokLimit:
		return limitInitial, nil
	}

	if tok.isLogicalOperator() {

		// this is the end of the previous operation
		if err := p.current.endOperation(); err != nil {
			return invalidState, err
		}

		// set this operator and jump to a new operation
		p.current.logicalOperator = tok.Type
		return operationInitial, nil

	}

	if tok.isComparisonOperator() {

		// set the operator and jump to the operator state
		p.current.operator = tok.Type
		return operator, nil
	}

	// we close a context, pop it
	if tok.Type == tokRightParenthesis {
		return p.popContext()
	}

	if tok.Type == tokBetween {
		p.current.rangeTest = &rangeTestOperation{
			test: p.current.left,
		}
		p.current.left = p.current.rangeTest

		return rangeMin, nil
	}

	return unexpected(tok, tokInvalid)
}

func (p *parser) rangeMin(tok *token) (state, error) {
	c, err := tok2operand(tok)
	if err != nil {
		return invalidState, err
	}

	if p.current.rangeTest == nil {
		return invalidState, errors.New("Nil range test")
	}

	p.current.rangeTest.min = c

	return rangeAnd, nil
}

func (p *parser) rangeMax(tok *token) (state, error) {
	c, err := tok2operand(tok)
	if err != nil {
		return invalidState, err
	}

	if p.current.rangeTest == nil {
		return invalidState, errors.New("Nil range test")
	}

	p.current.rangeTest.max = c

	return rightOperand, nil
}

// We're waiting for a right operand, nothing else, or a (
func (p *parser) operatorState(tok *token) (state, error) {

	// handle the (
	if tok.Type == tokLeftParenthesis {
		// we push the state
		p.pushContext(rightOperand)
		// jump the the start of an operation
		return operationInitial, nil
	}

	if tok.isField() {
		p.current.right = NewField(tok.Value)
	} else if tok.isConst() {
		c, err := tok.Const()
		if err != nil {
			return invalidState, err
		}
		p.current.right = c
	} else {
		return unexpected(tok, tokInvalid)
	}

	// the right operand has been setted jump to the left operand state
	return rightOperand, nil
}

// We're waiting for a logical operator, ), or the end
func (p *parser) operandRightState(tok *token) (state, error) {

	// end of a context
	// we pop, and continue to the correct state
	// (the pop will handle the end of the operation)
	if tok.Type == tokRightParenthesis {
		return p.popContext()
	}

	// in all following case we have to end the current operation
	if err := p.current.endOperation(); err != nil {
		return invalidState, err
	}

	switch tok.Type {
	case tokStarting:
		return startingInitial, nil
	case tokLimit:
		return limitInitial, nil
	}

	// the end of the expression
	if tok.isEnd() {
		return end, nil
	}

	// we continue to chain
	if tok.isLogicalOperator() {

		// set this operator and jump to a new operation
		p.current.logicalOperator = tok.Type
		return operationInitial, nil
	}

	return unexpected(tok, tokInvalid)
}

func (p *parser) startingInitial(tok *token) (state, error) {
	return p.expect(tokAt, tok, startingAt)
}

func (p *parser) startingAt(tok *token) (state, error) {

	if tok.isNumeric() {
		c, err := tok.Const()
		if err != nil {
			return invalidState, err
		}

		p.query.startingAt = c.AsInt()

		return clauseEnd, nil
	}

	return unexpected(tok, tokInt)
}

func (p *parser) limit(tok *token) (state, error) {
	if tok.isNumeric() {
		c, err := tok.Const()
		if err != nil {
			return invalidState, err
		}

		p.query.setLimit(c.AsInt())

		return limitSep, nil
	}

	return unexpected(tok, tokInt)
}

func (p *parser) limitSep(tok *token) (state, error) {
	if tok.Type == tokComma {
		return limitMax, nil
	}

	return p.clauseEnd(tok)
}

func (p *parser) limitMax(tok *token) (state, error) {
	if tok.isNumeric() {
		c, err := tok.Const()
		if err != nil {
			return invalidState, err
		}

		p.query.startingAt = p.query.Limit()
		p.query.setLimit(c.AsInt())

		return limitSep, nil
	}

	return unexpected(tok, tokInt)
}

func (p *parser) expect(expected tokenType, tok *token, state state) (state, error) {
	if tok.Type != expected {
		return unexpected(tok, expected)
	}
	return state, nil
}

// Push the context and create a new one
func (p *parser) pushContext(state state) {
	// stack the context
	p.current.stackState = state
	p.stack = append(p.stack, p.current)

	// creates a new one and reset it
	p.current = newContext()
}

// End the current operation, and continue to the correct state (given on push)
func (p *parser) popContext() (state, error) {

	// first we end the operation
	if err := p.current.endOperation(); err != nil {
		return invalidState, err
	}

	// creates a group operand
	g, err := newGroupOperand(p.current.first.simplify())
	if err != nil {
		return invalidState, err
	}

	// pop
	l := len(p.stack)
	p.current = p.stack[l-1]
	p.stack = p.stack[:l-1]

	// according the stack state, put the group it on the correct side
	switch p.current.stackState {
	case leftOperand:
		p.current.left = g
	case rightOperand:
		p.current.right = g
	}

	// clear the current stackState
	s := p.current.stackState
	p.current.stackState = -1

	// return the next state
	return s, nil
}

// Creates a new and fresh context
func newContext() *context {
	return &context{
		operator:        tokInvalid,
		logicalOperator: tokInvalid,
		stackState:      invalidState,
	}
}

// End the current operation
// wrap the current operation into a simple LogicalOperation, and chain it
func (c *context) endOperation() error {

	if c.left == nil {
		return errors.New("Empty operation")
	}

	// Creates the operand

	var op operand
	var err error

	if c.rangeTest != nil {
		op = c.rangeTest
	}

	if c.operator != tokInvalid {
		op, err = newComparison(c.left, operatorTypeFromTokenType(c.operator), c.right)
		if err != nil {
			return err
		}
	} else {
		op = c.left
	}

	// reset the operation
	c.left = nil
	c.right = nil
	c.operator = tokInvalid

	// Creates a simple logical operation and chain it

	lo, err := newLeftLogicalOperation(op)
	if err != nil {
		return err
	}

	if c.logicalOperator == tokInvalid {

		if c.first != nil || c.last != nil {
			return errors.New("Unexpected state, the logical operations should be not initialized")
		}

		c.first = lo
		c.last = lo

		return nil
	}

	if c.first == nil || c.last == nil {
		return errors.New("Unexpected state, no logical operation initialized")
	}

	// define how to chain according the logical operator
	switch c.logicalOperator {
	case tokAnd:
		c.last.chain(operatorTypeFromTokenType(c.logicalOperator), lo)
		c.last = lo
	case tokOr:
		c.first, err = newLeftLogicalOperation(c.first)
		if err != nil {
			return err
		}

		c.first.chain(operatorTypeFromTokenType(c.logicalOperator), lo)
		c.last = lo
	}

	// reset the logical operator
	c.logicalOperator = tokInvalid

	return nil
}

// Helper to creates the unexpected error
func unexpected(tok *token, expected tokenType) (state, error) {

	if expected != tokInvalid {

		if tok.isEnd() {
			return invalidState, fmt.Errorf(
				"Expected '%s' at position %d, got end of stream",
				expected, tok.Pos)
		}

		return invalidState, fmt.Errorf(
			"Expected '%s' at position %d, got '%s'",
			expected, tok.Pos, tok.Value)
	}

	if tok.isEnd() {
		return invalidState, fmt.Errorf(
			"Unexpected end of stream at pos %d", tok.Pos)
	}

	return invalidState, fmt.Errorf(
		"Unexpected '%s' at pos %d", tok.Value, tok.Pos)
}
