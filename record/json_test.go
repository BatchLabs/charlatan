package record

import (
	"encoding/json"
	"io"
	"strings"
	"testing"

	ch "github.com/BatchLabs/charlatan"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testJSONDecoder() *json.Decoder {
	return json.NewDecoder(strings.NewReader(`
	{
		"name": "Michel",
		"b": true,
		"age": 92,
		"n": null,
		"a": [],
		"we":{"need": {"to": {"go": {"deeper": 1, "a": "d"}}}}
	}
	`))
}

func TestFindUnexistingField(t *testing.T) {
	r, err := NewJSONRecordFromDecoder(testJSONDecoder())
	require.Nil(t, err)
	require.NotNil(t, r)

	_, err = r.Find(ch.NewField("yolo"))
	assert.NotNil(t, err)
}

func TestFindArrayField(t *testing.T) {
	r, err := NewJSONRecordFromDecoder(testJSONDecoder())
	require.Nil(t, err)
	require.NotNil(t, r)

	c, err := r.Find(ch.NewField("a"))
	require.Nil(t, err)
	require.NotNil(t, c)

	require.True(t, c.IsString())
	require.Equal(t, `[]`, c.AsString())
}

func TestFindObjectField(t *testing.T) {
	r, err := NewJSONRecordFromDecoder(testJSONDecoder())
	require.Nil(t, err)
	require.NotNil(t, r)

	c, err := r.Find(ch.NewField("we"))
	require.Nil(t, err)
	require.NotNil(t, c)

	require.True(t, c.IsString())
	require.Equal(t, `{"need": {"to": {"go": {"deeper": 1, "a": "d"}}}}`, c.AsString())
}

func TestFindTopLevelStringField(t *testing.T) {
	r, err := NewJSONRecordFromDecoder(testJSONDecoder())
	require.Nil(t, err)
	require.NotNil(t, r)

	c, err := r.Find(ch.NewField("name"))
	require.Nil(t, err)
	require.NotNil(t, c)

	assert.True(t, c.IsString())
	assert.Equal(t, "Michel", c.AsString())
}

func TestFindTopLevelIntField(t *testing.T) {
	r, err := NewJSONRecordFromDecoder(testJSONDecoder())
	require.Nil(t, err)
	require.NotNil(t, r)

	c, err := r.Find(ch.NewField("age"))
	require.Nil(t, err)
	require.NotNil(t, c)

	assert.True(t, c.IsNumeric())
	assert.Equal(t, int64(92), c.AsInt())
}

func TestFindTopLevelBoolField(t *testing.T) {
	r, err := NewJSONRecordFromDecoder(testJSONDecoder())
	require.Nil(t, err)
	require.NotNil(t, r)

	c, err := r.Find(ch.NewField("b"))
	require.Nil(t, err)
	require.NotNil(t, c)

	assert.True(t, c.IsBool())
	assert.Equal(t, true, c.AsBool())
}

func TestFindTopLevelNullField(t *testing.T) {
	r, err := NewJSONRecordFromDecoder(testJSONDecoder())
	require.Nil(t, err)
	require.NotNil(t, r)

	c, err := r.Find(ch.NewField("n"))
	require.Nil(t, err)
	require.NotNil(t, c)

	assert.True(t, c.IsNull())
}

func TestFindTopLevelEmptyStringField(t *testing.T) {
	r, err := NewJSONRecordFromDecoder(json.NewDecoder(strings.NewReader(`{"foo": ""}`)))
	require.Nil(t, err)
	require.NotNil(t, r)

	c, err := r.Find(ch.NewField("foo"))
	require.Nil(t, err)
	require.NotNil(t, c)

	assert.False(t, c.IsNull())
	assert.True(t, c.IsString())
}

func TestFindDeepStringField(t *testing.T) {
	r, err := NewJSONRecordFromDecoder(testJSONDecoder())
	require.Nil(t, err)
	require.NotNil(t, r)

	c, err := r.Find(ch.NewField("we.need.to.go.deeper"))
	require.Nil(t, err)
	require.NotNil(t, c)

	assert.True(t, c.IsNumeric())
	assert.Equal(t, int64(1), c.AsInt())
}

func TestJSONDecoderMultipleRecords(t *testing.T) {
	r := json.NewDecoder(strings.NewReader(`
	{"age": 42}
	{"age": 19}
	`))
	require.NotNil(t, r)

	_, err := NewJSONRecordFromDecoder(r)
	require.Nil(t, err)

	_, err = NewJSONRecordFromDecoder(r)
	require.Nil(t, err)

	_, err = NewJSONRecordFromDecoder(r)
	assert.Equal(t, io.EOF, err)
}

func TestJSONRecordSelectStar(t *testing.T) {
	r := json.NewDecoder(strings.NewReader(`{"foo": 42}`))
	require.NotNil(t, r)

	rec, err := NewJSONRecordFromDecoder(r)
	require.Nil(t, err)
	require.NotNil(t, rec)

	all, err := rec.Find(ch.NewField("*"))
	require.Nil(t, err)
	require.NotNil(t, all)

	require.True(t, all.IsString())
	require.Equal(t, `{"foo":42}`, all.AsString())
}

func TestJSONRecordTopLevelSoftFind(t *testing.T) {
	rec, err := NewJSONRecordFromDecoder(json.NewDecoder(strings.NewReader(`{"foo":42}`)))
	require.Nil(t, err)
	require.NotNil(t, rec)

	_, err = rec.Find(ch.NewField("yo"))
	assert.NotNil(t, err)

	rec.SoftMatching = true
	v, err := rec.Find(ch.NewField("yo"))
	assert.Nil(t, err)
	assert.True(t, v.IsNull())
}

func TestJSONRecordSoftFindNotAnObject(t *testing.T) {
	rec, err := NewJSONRecordFromDecoder(json.NewDecoder(strings.NewReader(`{"foo":42}`)))
	require.Nil(t, err)
	require.NotNil(t, rec)

	_, err = rec.Find(ch.NewField("foo.bar.qux"))
	assert.NotNil(t, err)

	rec.SoftMatching = true
	v, err := rec.Find(ch.NewField("foo.bar.qux"))
	assert.Nil(t, err)
	require.NotNil(t, v)
	assert.True(t, v.IsNull())
}

func TestJSONRecordSoftFind(t *testing.T) {
	rec, err := NewJSONRecordFromDecoder(json.NewDecoder(strings.NewReader(`
		{"foo": {"a": 2}}
	`)))
	require.Nil(t, err)
	require.NotNil(t, rec)

	_, err = rec.Find(ch.NewField("foo.bar.qux"))
	assert.NotNil(t, err)

	rec.SoftMatching = true
	v, err := rec.Find(ch.NewField("foo.bar.qux"))
	assert.Nil(t, err)
	require.NotNil(t, v)
	assert.True(t, v.IsNull())
}
