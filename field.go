package charlatan

// Field is a field, contained into the SELECT part and the condition.
// A field is an operand, it can return the value extracted into the Record.
type Field struct {
	name string
}

// NewField returns a new field from the given string
func NewField(name string) *Field {
	return &Field{name}
}

// Evaluate evaluates the field on a record
func (f *Field) Evaluate(record Record) (*Const, error) {
	return record.Find(f)
}

// Name returns the field's name
func (f *Field) Name() string {
	return f.name
}

func (f *Field) String() string {
	return f.name
}
