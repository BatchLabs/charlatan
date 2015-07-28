package charlatan

import "fmt"

// Operand is an operand, can be evaluated and have to return a constant.
// Returns a error, if the evaluation is not possible
type Operand interface {
	Evaluate(Record) (*Const, error)
	String() string
}

// ConstOperand is the constant operand
type ConstOperand struct {
	constant *Const
}

// Comparison is the comparison operation
type Comparison struct {
	left     Operand
	operator OperatorType
	right    Operand
}

// LogicalOperation is the logical operation
type LogicalOperation struct {
	left     Operand
	operator OperatorType
	right    Operand
}

// GroupOperand is the group operand
// Just keep in mind that there was () surrounding this operation
type GroupOperand struct {
	operand Operand
}

// NewConstOperand returns a new ConstOperand from the given Const
func NewConstOperand(constant *Const) *ConstOperand {
	return &ConstOperand{constant}
}

// Evaluate evaluates the constant against a record.
func (c *ConstOperand) Evaluate(record Record) (*Const, error) {
	return c.constant, nil
}

func (c *ConstOperand) String() string {
	s := c.constant.String()
	if c.constant.IsString() {
		return fmt.Sprintf("'%s'", s)
	}
	return s
}

// NewLogicalOperation creates a new logicial operation from the given operator
// and operands
func NewLogicalOperation(left Operand, operator OperatorType, right Operand) (*LogicalOperation, error) {

	o, err := newLogicalOperation(left)
	if err != nil {
		return nil, err
	}

	if err := o.chain(operator, right); err != nil {
		return nil, err
	}

	return o, nil
}

// An internal constructor, just with the left operand
func newLogicalOperation(left Operand) (*LogicalOperation, error) {

	if left == nil {
		return nil, fmt.Errorf("Can't creates a new comparison with the left operand nil")
	}

	return &LogicalOperation{left, -1, nil}, nil
}

// Chain a right operand with the given operator
func (o *LogicalOperation) chain(operator OperatorType, right Operand) error {

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
func (o *LogicalOperation) simplify() Operand {

	if left, ok := o.left.(*LogicalOperation); ok {
		o.left = left.simplify()
	}

	if o.right == nil {
		return o.left
	}

	if right, ok := o.right.(*LogicalOperation); ok {
		o.right = right.simplify()
	}

	return o
}

// Evaluate evaluates the logical operation against the given record
func (o *LogicalOperation) Evaluate(record Record) (*Const, error) {

	var err error
	var leftValue, rightValue *Const

	leftValue, err = o.left.Evaluate(record)
	if err != nil {
		return nil, err
	}

	leftBool := leftValue.AsBool()

	// AND
	if !leftBool && o.operator == OperatorAnd {
		return BoolConst(false), nil
	}

	// OR
	if leftBool && o.operator == OperatorOr {
		return BoolConst(true), nil
	}

	rightValue, err = o.right.Evaluate(record)
	if err != nil {
		return nil, err
	}

	return BoolConst(rightValue.AsBool()), nil
}

func (o *LogicalOperation) String() string {
	switch o.operator {
	case OperatorAnd:
		return fmt.Sprintf("%s AND %s", o.left, o.right)
	case OperatorOr:
		return fmt.Sprintf("%s OR %s", o.left, o.right)
	default:
		return "Unknown operator"
	}
}

// NewGroupOperand returns a new group operand from the given operand
func NewGroupOperand(operand Operand) (*GroupOperand, error) {
	if operand == nil {
		return nil, fmt.Errorf("Can't creates a new group with the an operand nil")
	}

	return &GroupOperand{operand}, nil
}

// Evaluate evaluates the group operand against the given record
func (o *GroupOperand) Evaluate(record Record) (*Const, error) {
	return o.operand.Evaluate(record)
}

func (o *GroupOperand) String() string {
	return fmt.Sprintf("(%s)", o.operand)
}
