package charlatan

import "bytes"

// Query is a query
type Query struct {
	// the fields to select if condition match the object
	fields []*Field
	// the resource from wich we want to evaluate and select fields
	from string
	// the expression to evaluate on each record. The resulting constant will
	// always be converted as a bool
	expression operand
	// the record index to start from
	startingAt int64
	// the record index to stop at
	limit *int64
}

// A Record is a record
type Record interface {
	Find(*Field) (*Const, error)
}

// QueryFromString creates a query from the given string
func QueryFromString(s string) (*Query, error) {
	return parserFromString(s).Parse()
}

// NewQuery creates a new query with the given from part
func NewQuery(from string) *Query {
	return &Query{from: from}
}

// From returns the FROM part of this query
func (q *Query) From() string {
	return q.from
}

// StartingAt returns the 'STARTING AT' part of the query, or 0 if it's not
// present
func (q *Query) StartingAt() int64 {
	return q.startingAt
}

// HasLimit tests if the query has a limit
func (q *Query) HasLimit() bool {
	return q.limit != nil
}

// Limit returns the 'LIMIT' part of the query, or 0 if it's not present
func (q *Query) Limit() int64 {
	if q.limit == nil {
		return 0
	}
	return *q.limit
}

// AddField adds one field
func (q *Query) AddField(field *Field) {
	if field != nil {
		q.fields = append(q.fields, field)
	}
}

// AddFields adds multiple fields
func (q *Query) AddFields(fields []*Field) {
	for _, field := range fields {
		q.AddField(field)
	}
}

// Fields returns the fields
func (q *Query) Fields() []*Field {
	return q.fields
}

// setWhere sets the where condition
func (q *Query) setWhere(op operand) {

	if op == nil {
		return
	}

	q.expression = op
}

func (q *Query) setLimit(limit int64) {
	q.limit = &limit
}

// FieldsValues extracts the values of each fields into the given record
// Note that you should evaluate the query first
func (q *Query) FieldsValues(record Record) ([]*Const, error) {
	values := make([]*Const, len(q.fields))

	for i, field := range q.fields {

		value, err := field.Evaluate(record)
		if err != nil {
			return nil, err
		}

		values[i] = value
	}

	return values, nil
}

// Evaluate evaluates the query against the given record
func (q *Query) Evaluate(record Record) (bool, error) {

	// no expression, always valid
	if q.expression == nil {
		return true, nil
	}

	constant, err := q.expression.Evaluate(record)
	if err != nil {
		return false, err
	}

	return constant.AsBool(), nil
}

// String returns a string representation of this query
func (q *Query) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("SELECT ")

	for i, field := range q.fields {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(field.Name())
	}

	buffer.WriteString(" FROM ")
	buffer.WriteString(q.from)

	if q.expression != nil {
		buffer.WriteString(" WHERE ")
		buffer.WriteString(q.expression.String())
	}

	return buffer.String()
}
