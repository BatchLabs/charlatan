package record

import (
	"fmt"
	"strconv"

	ch "github.com/BatchLabs/charlatan"
)

// CSVRecord is a CSV record, i.e. a line from a CSV file
//
// If created with a header its fields can be retrieved with their column name.
// In any case, one can use a "$N" field name, where N is the column index,
// starting at 0.
//
// All values are retrieved as strings, and the special field "*" can be used
// to get a string representation of the record values.
type CSVRecord struct {
	header, record []string
}

var _ ch.Record = &CSVRecord{}

// NewCSVRecord returns a new CSVRecord
func NewCSVRecord(record []string) *CSVRecord {
	return &CSVRecord{record: record}
}

// NewCSVRecordWithHeader returns a new CSVRecord with the given header
func NewCSVRecordWithHeader(record []string, header []string) *CSVRecord {
	return &CSVRecord{header: header, record: record}
}

// Find implements the charlatan.Record interface
func (r *CSVRecord) Find(field *ch.Field) (*ch.Const, error) {

	name := field.Name()

	if name == "*" {
		return ch.StringConst(fmt.Sprintf("%v", r.record)), nil
	}

	// Column index

	if name[0] == '$' {
		index, err := strconv.ParseInt(name[1:], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid column index %s: %s", name, err)
		}

		return r.AtIndex(int(index))
	}

	// Column name

	index := r.ColumnNameIndex(name)
	if index < 0 {
		return nil, fmt.Errorf("Can't find field name: %s", name)
	}

	return r.AtIndex(index)
}

// AtIndex gets the value at the given index
func (r *CSVRecord) AtIndex(index int) (*ch.Const, error) {

	if index < 0 || index > len(r.record) {
		return nil, fmt.Errorf("index out of bounds %d", index)
	}

	// FIXME should we accept NULL values?
	value := r.record[index]
	if value == "NULL" {
		return ch.NullConst(), nil
	}

	return ch.ConstFromString(value), nil
}

// ColumnNameIndex searches the index of the column name in this record’s
// header. If it doesn’t have a header or if the column wasn’t found, the
// method returns -1. The column name match is case-sensitive, the first
// matching one is used.
func (r *CSVRecord) ColumnNameIndex(name string) int {
	for index, element := range r.header {
		if element == name {
			return index
		}
	}

	return -1
}
