package charlatan

import (
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
	fromName

	operationInitial
	leftOperand
	operator
	rightOperand

	startingInitial
	startingAt

	beforeEnd
	end
)

// Parser is the parser itself
type Parser struct {
	// the lexer to read tokens
	lexer *Lexer
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
	left     Operand
	operator TokenType
	right    Operand

	// the logical operation nodes
	first *LogicalOperation
	last  *LogicalOperation

	// the current logical operator
	logicalOperator TokenType

	// if stacked, in which state
	// this context should be restored
	// only leftOperand or rightOperand
	stackState state
}

// ParserFromString creates a new parser from the given string
func ParserFromString(s string) *Parser {
	return &Parser{
		LexerFromString(s),
		initial,
		make([]*Field, 0),
		nil,
		make([]*context, 0),
		newContext(),
	}
}

// Parse parses and returns the query
func (p *Parser) Parse() (*Query, error) {

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

		// select part
		case selectInitial:
			p.state, err = p.selectState(tok)
		case selectField:
			p.state, err = p.selectFieldState(tok)

		// from part
		case fromInitial:
			p.state, err = p.fromState(tok)
		case fromName:
			p.state, err = p.fromNameState(tok)

		// where part
		case operationInitial:
			p.state, err = p.operationState(tok)
		case leftOperand:
			p.state, err = p.operandLeftState(tok)
		case operator:
			p.state, err = p.operatorState(tok)
		case rightOperand:
			p.state, err = p.operandRightState(tok)

		// 'starting at' part
		case startingInitial:
			p.state, err = p.startingInitial(tok)
		case startingAt:
			p.state, err = p.startingAt(tok)

		case beforeEnd:
			p.state, err = p.beforeEnd(tok)

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
		p.query.SetWhere(p.current.first.simplify())
	}

	return p.query, nil
}

// We’re only waiting for the SELECT keyword
func (p *Parser) initialState(tok *Token) (state, error) {
	if tok.Type != TokSelect {
		return unexpected(tok, TokSelect)
	}

	return selectInitial, nil
}

// We’re waiting for a field
func (p *Parser) selectState(tok *Token) (state, error) {
	if tok.Type != TokField {
		return unexpected(tok, TokField)
	}

	p.fields = append(p.fields, NewField(tok.Value))

	return selectField, nil
}

// We’re waiting for a comma, or the FROM keyword
func (p *Parser) selectFieldState(tok *Token) (state, error) {
	if tok.Type == TokComma {
		// we just jump to the next field
		return selectInitial, nil
	}

	if tok.Type != TokFrom {
		return unexpected(tok, TokFrom)
	}

	return fromInitial, nil
}

// We’re waiting for a name of the from
func (p *Parser) fromState(tok *Token) (state, error) {
	if tok.Type != TokField {
		return unexpected(tok, TokField)
	}

	p.query = NewQuery(tok.Value)
	p.query.AddFields(p.fields)
	p.fields = nil

	return fromName, nil
}

// We’re waiting for the WHERE keyword, or the end, if there is no condition
func (p *Parser) fromNameState(tok *Token) (state, error) {
	if tok.Type == TokEnd {
		return end, nil
	}

	if tok.Type == TokWhere {
		return operationInitial, nil
	}

	if tok.Type == TokStarting {
		return startingInitial, nil
	}

	return unexpected(tok, TokWhere)
}

// We’re waiting for a left operand, or a (
func (p *Parser) operationState(tok *Token) (state, error) {
	if tok.Type == TokLeftParenthesis {
		// push the context and start a new operation
		p.pushContext(leftOperand)
		return operationInitial, nil
	}

	if tok.IsField() {
		p.current.left = NewField(tok.Value)
	} else if tok.IsConst() {
		c, err := tok.Const()
		if err != nil {
			return invalidState, err
		}
		p.current.left = NewConstOperand(c)
	} else {
		return unexpected(tok, TokInvalid)
	}

	// the left operand has been setted jump to the left operand state
	return leftOperand, nil
}

// We’re waiting for an operator:
//     - logical operator, we step to the next operation
//     - comparison operator, we continue to the operator state
//
// We can encounter a ), or the end
func (p *Parser) operandLeftState(tok *Token) (state, error) {

	if tok.IsEnd() {
		// end the previous operation
		if err := p.current.endOperation(); err != nil {
			return invalidState, err
		}

		return end, nil
	}

	if tok.Type == TokStarting {
		return startingInitial, nil
	}

	if tok.IsLogicalOperator() {

		// this is the end of the previous operation
		if err := p.current.endOperation(); err != nil {
			return invalidState, err
		}

		// set this operator and jump to a new operation
		p.current.logicalOperator = tok.Type
		return operationInitial, nil

	} else if tok.IsComparisonOperator() {

		// set the operator and jump to the operator state
		p.current.operator = tok.Type
		return operator, nil
	}

	// we close a context, pop it
	if tok.Type == TokRightParenthesis {
		return p.popContext()
	}

	return unexpected(tok, TokInvalid)
}

// We're waiting for a right operand, nothing else, or a (
func (p *Parser) operatorState(tok *Token) (state, error) {

	// handle the (
	if tok.Type == TokLeftParenthesis {
		// we push the state
		p.pushContext(rightOperand)
		// jump the the start of an operation
		return operationInitial, nil
	}

	if tok.IsField() {
		p.current.right = NewField(tok.Value)
	} else if tok.IsConst() {
		c, err := tok.Const()
		if err != nil {
			return invalidState, err
		}
		p.current.right = NewConstOperand(c)
	} else {
		return unexpected(tok, TokInvalid)
	}

	// the right operand has been setted jump to the left operand state
	return rightOperand, nil
}

// We're waiting for a logical operator, ), or the end
func (p *Parser) operandRightState(tok *Token) (state, error) {

	// end of a context
	// we pop, and continue to the correct state
	// (the pop will handle the end of the operation)
	if tok.Type == TokRightParenthesis {
		return p.popContext()
	}

	// in all following case we have to end the current operation
	if err := p.current.endOperation(); err != nil {
		return invalidState, err
	}

	if tok.Type == TokStarting {
		return startingInitial, nil
	}

	// the end of the expression
	if tok.IsEnd() {
		return end, nil
	}

	// we continue to chain
	if tok.IsLogicalOperator() {

		// set this operator and jump to a new operation
		p.current.logicalOperator = tok.Type
		return operationInitial, nil
	}

	return unexpected(tok, TokInvalid)
}

func (p *Parser) startingInitial(tok *Token) (state, error) {
	if tok.Type == TokAt {
		return startingAt, nil
	}
	return unexpected(tok, TokAt)
}

func (p *Parser) startingAt(tok *Token) (state, error) {

	if tok.IsNumeric() {
		c, err := tok.Const()
		if err != nil {
			return invalidState, err
		}

		p.query.startingAt = c.AsInt()

		return beforeEnd, nil
	}

	return unexpected(tok, TokInt)
}

func (p *Parser) beforeEnd(tok *Token) (state, error) {
	if tok.IsEnd() {
		return end, nil
	}

	return unexpected(tok, TokEnd)
}

// Push the context and create a new one
func (p *Parser) pushContext(state state) {
	// stack the context
	p.current.stackState = state
	p.stack = append(p.stack, p.current)

	// creates a new one and reset it
	p.current = newContext()
}

// End the current operation, and continue to the correct state (given on push)
func (p *Parser) popContext() (state, error) {

	// first we end the operation
	if err := p.current.endOperation(); err != nil {
		return invalidState, err
	}

	// creates a group operand
	g, err := NewGroupOperand(p.current.first.simplify())
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
		operator:        TokInvalid,
		logicalOperator: TokInvalid,
		stackState:      invalidState,
	}
}

// End the current operation
// wrap the current operation into a simple LogicalOperation, and chain it
func (c *context) endOperation() error {

	if c.left == nil {
		return fmt.Errorf("Empty operation")
	}

	// Creates the operand

	var op Operand
	var err error

	if c.operator != TokInvalid {
		op, err = NewComparison(c.left, OperatorTypeFromTokenType(c.operator), c.right)
		if err != nil {
			return err
		}
	} else {
		op = c.left
	}

	// reset the operation
	c.left = nil
	c.right = nil
	c.operator = TokInvalid

	// Creates a simple logical operation and chain it

	lo, err := newLogicalOperation(op)
	if err != nil {
		return err
	}

	if c.logicalOperator == TokInvalid {

		if c.first != nil || c.last != nil {
			return fmt.Errorf("Unexpected state, the logical operations should be not initialized")
		}

		c.first = lo
		c.last = lo

		return nil
	}

	if c.first == nil || c.last == nil {
		return fmt.Errorf("Unexpected state, no logical operation initialized")
	}

	// define how to chain according the logical operator
	switch c.logicalOperator {
	case TokAnd:
		c.last.chain(OperatorTypeFromTokenType(c.logicalOperator), lo)
		c.last = lo
	case TokOr:
		c.first, err = newLogicalOperation(c.first)
		if err != nil {
			return err
		}

		c.first.chain(OperatorTypeFromTokenType(c.logicalOperator), lo)
		c.last = lo
	}

	// reset the logical operator
	c.logicalOperator = TokInvalid

	return nil
}

// Helper to creates the unexpected error
func unexpected(tok *Token, expected TokenType) (state, error) {

	if expected != TokInvalid {

		if tok.IsEnd() {
			return invalidState, fmt.Errorf(
				"Expected '%s' at position %d, got end of stream",
				expected, tok.Pos)
		}

		return invalidState, fmt.Errorf(
			"Expected '%s' at position %d, got '%s'",
			expected, tok.Pos, tok.Value)
	}

	if tok.IsEnd() {
		return invalidState, fmt.Errorf(
			"Unexpected end of stream at pos %d", tok.Pos)
	}

	return invalidState, fmt.Errorf(
		"Unexpected '%s' at pos %d", tok.Value, tok.Pos)
}
