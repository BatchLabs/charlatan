package charlatan

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dummyPerson struct {
	name string
	age  int
}

func (d *dummyPerson) Find(f *Field) (*Const, error) {
	switch f.Name() {
	case "name":
		return StringConst(d.name), nil
	case "age":
		return IntConst(int64(d.age)), nil
	default:
		return nil, errors.New("wrong field")
	}
}

var _ Record = &dummyPerson{}

func TestQueryInRange(t *testing.T) {
	q := &Query{
		fields: []*Field{
			&Field{"name"},
		},
		expression: &rangeTestOperation{
			test: Field{"age"},
			min:  IntConst(10),
			max:  IntConst(20),
		},
	}

	m, err := q.Evaluate(&dummyPerson{name: "A", age: 15})
	assert.Nil(t, err)
	assert.True(t, m)

	m, err = q.Evaluate(&dummyPerson{name: "A", age: 25})
	assert.Nil(t, err)
	assert.False(t, m)

	vals, err := q.FieldsValues(&dummyPerson{name: "A", age: 12})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(vals))
	assert.Equal(t, "A", vals[0].AsString())
}

func TestQueryFromStringInRange(t *testing.T) {
	q, err := QueryFromString("SELECT a FROM b WHERE a IN [1, 2]")
	assert.Nil(t, err)
	assert.NotNil(t, q)
	assert.NotNil(t, q.fields)
	assert.NotNil(t, q.expression)

	assert.Equal(t, "b", q.from)
	assert.Equal(t, int64(0), q.startingAt)

	assert.Equal(t, 1, len(q.fields))
	assert.NotNil(t, q.fields[0])
	assert.Equal(t, "a", q.fields[0].name)

	expr, ok := q.expression.(*rangeTestOperation)
	require.True(t, ok, "%v should be a range test", q.expression)

	assert.NotNil(t, expr.test)
	assert.NotNil(t, expr.min)
	assert.NotNil(t, expr.max)

	f, ok := expr.test.(*Field)
	require.True(t, ok)
	assert.Equal(t, "a", f.name)

	min, ok := expr.min.(*Const)
	require.True(t, ok, "%v should be a const (not %s)",
		expr.min, reflect.TypeOf(expr.min))
	assert.True(t, min.IsNumeric())
	assert.Equal(t, int64(1), min.AsInt())

	max, ok := expr.max.(*Const)
	require.True(t, ok, "%v should be a const (not %s)",
		expr.min, reflect.TypeOf(expr.max))
	assert.True(t, max.IsNumeric())
	assert.Equal(t, int64(2), max.AsInt())
}
