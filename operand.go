package charlatan

import "fmt"

// operand is an operand, can be evaluated and have to return a constant.
// Returns a error, if the evaluation is not possible
type operand interface {
	Evaluate(Record) (*Const, error)
	String() string
}

// constOperand is the constant operand
type constOperand struct {
	constant *Const
}

// comparison is the comparison operation
type comparison struct {
	left     operand
	operator operatorType
	right    operand
}

// logicalOperation is the logical operation
type logicalOperation struct {
	left     operand
	operator operatorType
	right    operand
}

// groupOperand is the group operand
// Just keep in mind that there was () surrounding this operation
type groupOperand struct {
	operand operand
}

// newConstOperand returns a new constOperand from the given Const
func newConstOperand(constant *Const) *constOperand {
	return &constOperand{constant}
}

// Evaluate evaluates the constant against a record.
func (c *constOperand) Evaluate(record Record) (*Const, error) {
	return c.constant, nil
}

func (c *constOperand) String() string {
	s := c.constant.String()
	if c.constant.IsString() {
		return fmt.Sprintf("'%s'", s)
	}
	return s
}

// newLogicalOperation creates a new logicial operation from the given operator
// and operands
func newLogicalOperation(left operand, operator operatorType, right operand) (*logicalOperation, error) {

	o, err := newLeftLogicalOperation(left)
	if err != nil {
		return nil, err
	}

	if err := o.chain(operator, right); err != nil {
		return nil, err
	}

	return o, nil
}

// A constructor with the left operand only
func newLeftLogicalOperation(left operand) (*logicalOperation, error) {

	if left == nil {
		return nil, fmt.Errorf("Can't creates a new comparison with the left operand nil")
	}

	return &logicalOperation{left, -1, nil}, nil
}

// Chain a right operand with the given operator
func (o *logicalOperation) chain(operator operatorType, right operand) error {

	if right == nil {
		return fmt.Errorf("Can't creates a new comparison with the right operand nil")
	}

	if !operator.IsLogical() {
		return fmt.Errorf("The operator should be a logical operator")
	}

	o.operator = operator
	o.right = right

	return nil
}

// Simplify this operation
// In case of the right operand is missing, just return the left one
func (o *logicalOperation) simplify() operand {

	if left, ok := o.left.(*logicalOperation); ok {
		o.left = left.simplify()
	}

	if o.right == nil {
		return o.left
	}

	if right, ok := o.right.(*logicalOperation); ok {
		o.right = right.simplify()
	}

	return o
}

// Evaluate evaluates the logical operation against the given record
func (o *logicalOperation) Evaluate(record Record) (*Const, error) {

	var err error
	var leftValue, rightValue *Const

	leftValue, err = o.left.Evaluate(record)
	if err != nil {
		return nil, err
	}

	leftBool := leftValue.AsBool()

	// AND
	if !leftBool && o.operator == operatorAnd {
		return BoolConst(false), nil
	}

	// OR
	if leftBool && o.operator == operatorOr {
		return BoolConst(true), nil
	}

	rightValue, err = o.right.Evaluate(record)
	if err != nil {
		return nil, err
	}

	return BoolConst(rightValue.AsBool()), nil
}

func (o *logicalOperation) String() string {
	switch o.operator {
	case operatorAnd:
		return fmt.Sprintf("%s AND %s", o.left, o.right)
	case operatorOr:
		return fmt.Sprintf("%s OR %s", o.left, o.right)
	default:
		return "Unknown operator"
	}
}

// newGroupOperand returns a new group operand from the given operand
func newGroupOperand(operand operand) (*groupOperand, error) {
	if operand == nil {
		return nil, fmt.Errorf("Can't creates a new group with the an operand nil")
	}

	return &groupOperand{operand}, nil
}

// Evaluate evaluates the group operand against the given record
func (o *groupOperand) Evaluate(record Record) (*Const, error) {
	return o.operand.Evaluate(record)
}

func (o *groupOperand) String() string {
	return fmt.Sprintf("(%s)", o.operand)
}
