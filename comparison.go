package charlatan

import "fmt"

// newComparison creates a new comparison from the given operands
func newComparison(left operand, operator operatorType, right operand) (*comparison, error) {

	if left == nil {
		return nil, fmt.Errorf("Can't creates a new comparison with a nil left operand")
	}

	if right == nil {
		return nil, fmt.Errorf("Can't creates a new comparison with a nil right operand")
	}

	if !operator.isComparison() {
		return nil, fmt.Errorf("The operator should be a comparison operator")
	}

	return &comparison{left, operator, right}, nil
}

// Evaluate evaluates the comparison against a given record and return the
// resulting value
func (c *comparison) Evaluate(record Record) (*Const, error) {
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
	case operatorEq:
		return BoolConst(r == 0), nil
	case operatorNeq:
		return BoolConst(r != 0), nil
	case operatorLt:
		return BoolConst(r < 0), nil
	case operatorLte:
		return BoolConst(r <= 0), nil
	case operatorGt:
		return BoolConst(r > 0), nil
	case operatorGte:
		return BoolConst(r >= 0), nil
	}

	return nil, fmt.Errorf("Unknown operator %s", c.operator)
}

func (c *comparison) String() string {
	return fmt.Sprintf("%s %s %s", c.left, c.operator, c.right)
}
