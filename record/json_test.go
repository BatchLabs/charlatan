package record

import (
	"io"
	"strings"
	"testing"

	ch "github.com/BatchLabs/charlatan"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// func (r *JSONRecord) Find(field *ch.Field) (*ch.Const, error) {

func testJSONReader() io.Reader {
	return strings.NewReader(`
	{
		"name": "Michel",
		"b": true,
		"age": 92,
		"n": null,
		"a": [],
		"we":{"need": {"to": {"go": {"deeper": 1, "a": "d"}}}}
	}
	`)
}

func TestFindUnexistingField(t *testing.T) {
	r, err := NewJSONRecordFromReader(testJSONReader())
	require.Nil(t, err)
	require.NotNil(t, r)

	_, err = r.Find(ch.NewField("yolo"))
	assert.NotNil(t, err)
}

func TestFindNotAConstField(t *testing.T) {
	r, err := NewJSONRecordFromReader(testJSONReader())
	require.Nil(t, err)
	require.NotNil(t, r)

	_, err = r.Find(ch.NewField("we"))
	assert.NotNil(t, err)

	_, err = r.Find(ch.NewField("a"))
	assert.NotNil(t, err)
}

func TestFindTopLevelStringField(t *testing.T) {
	r, err := NewJSONRecordFromReader(testJSONReader())
	require.Nil(t, err)
	require.NotNil(t, r)

	c, err := r.Find(ch.NewField("name"))
	require.Nil(t, err)
	require.NotNil(t, c)

	assert.True(t, c.IsString())
	assert.Equal(t, "Michel", c.AsString())
}

func TestFindTopLevelIntField(t *testing.T) {
	r, err := NewJSONRecordFromReader(testJSONReader())
	require.Nil(t, err)
	require.NotNil(t, r)

	c, err := r.Find(ch.NewField("age"))
	require.Nil(t, err)
	require.NotNil(t, c)

	assert.True(t, c.IsNumeric())
	assert.Equal(t, int64(92), c.AsInt())
}

func TestFindTopLevelBoolField(t *testing.T) {
	r, err := NewJSONRecordFromReader(testJSONReader())
	require.Nil(t, err)
	require.NotNil(t, r)

	c, err := r.Find(ch.NewField("b"))
	require.Nil(t, err)
	require.NotNil(t, c)

	assert.True(t, c.IsBool())
	assert.Equal(t, true, c.AsBool())
}

func TestFindTopLevelNullField(t *testing.T) {
	r, err := NewJSONRecordFromReader(testJSONReader())
	require.Nil(t, err)
	require.NotNil(t, r)

	c, err := r.Find(ch.NewField("n"))
	require.Nil(t, err)
	require.NotNil(t, c)

	assert.True(t, c.IsNull())
}

func TestFindTopLevelEmptyStringField(t *testing.T) {
	r, err := NewJSONRecordFromReader(strings.NewReader(`{"foo": ""}`))
	require.Nil(t, err)
	require.NotNil(t, r)

	c, err := r.Find(ch.NewField("foo"))
	require.Nil(t, err)
	require.NotNil(t, c)

	assert.False(t, c.IsNull())
	assert.True(t, c.IsString())
}

func TestFindDeepStringField(t *testing.T) {
	r, err := NewJSONRecordFromReader(testJSONReader())
	require.Nil(t, err)
	require.NotNil(t, r)

	c, err := r.Find(ch.NewField("we.need.to.go.deeper"))
	require.Nil(t, err)
	require.NotNil(t, c)

	assert.True(t, c.IsNumeric())
	assert.Equal(t, int64(1), c.AsInt())
}
