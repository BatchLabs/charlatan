package charlatan

import "fmt"

// NewComparison creates a new Comparison from the given operands
func NewComparison(left Operand, operator OperatorType, right Operand) (*Comparison, error) {

	if left == nil {
		return nil, fmt.Errorf("Can't creates a new comparison with a nil left operand")
	}

	if right == nil {
		return nil, fmt.Errorf("Can't creates a new comparison with a nil right operand")
	}

	if !operator.IsComparison() {
		return nil, fmt.Errorf("The operator should be a comparison operator")
	}

	return &Comparison{left, operator, right}, nil
}

// Evaluate evaluates the comparison against a given record and return the
// resulting value
func (c *Comparison) Evaluate(record Record) (*Const, error) {
	var err error
	var leftValue, rightValue *Const

	leftValue, err = c.left.Evaluate(record)
	if err != nil {
		return nil, err
	}

	rightValue, err = c.right.Evaluate(record)
	if err != nil {
		return nil, err
	}

	r, err := leftValue.CompareTo(rightValue)
	if err != nil {
		return nil, err
	}

	switch c.operator {
	case OperatorEq:
		return BoolConst(r == 0), nil
	case OperatorNeq:
		return BoolConst(r != 0), nil
	case OperatorLt:
		return BoolConst(r < 0), nil
	case OperatorLte:
		return BoolConst(r <= 0), nil
	case OperatorGt:
		return BoolConst(r > 0), nil
	case OperatorGte:
		return BoolConst(r >= 0), nil
	}

	return nil, fmt.Errorf("Unknown operator %s", c.operator)
}

func (c *Comparison) String() string {
	return fmt.Sprintf("%s %s %s", c.left, c.operator, c.right)
}
