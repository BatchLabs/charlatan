package record

import (
	"testing"

	ch "github.com/BatchLabs/charlatan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCSVRecordWithoutHeader(t *testing.T) {
	c := NewCSVRecord([]string{"a", "b", "c"})
	require.NotNil(t, c)
	assert.Nil(t, c.header)
}

func TestFindNameWithoutHeader(t *testing.T) {
	c := NewCSVRecord([]string{"a", "b", "c"})
	require.NotNil(t, c)

	_, err := c.Find(ch.NewField("foo"))
	assert.NotNil(t, err)
}

func TestFindColumnIndex(t *testing.T) {
	c := NewCSVRecord([]string{"a", "b", "c"})
	require.NotNil(t, c)

	_, err := c.Find(ch.NewField("$-1"))
	assert.NotNil(t, err)

	_, err = c.Find(ch.NewField("$42"))
	assert.NotNil(t, err)

	v, err := c.Find(ch.NewField("$1"))
	assert.Nil(t, err)
	assert.True(t, v.IsString())
	assert.Equal(t, "b", v.AsString())
}

func TestFindColumnName(t *testing.T) {
	c := NewCSVRecordWithHeader([]string{"a", "b", "c"}, []string{"id", "x", "y"})
	require.NotNil(t, c)

	_, err := c.Find(ch.NewField("yo"))
	assert.NotNil(t, err)

	_, err = c.Find(ch.NewField("xy"))
	assert.NotNil(t, err)

	v, err := c.Find(ch.NewField("y"))
	assert.Nil(t, err)
	assert.Equal(t, "c", v.AsString())
}

func TestFindStar(t *testing.T) {
	c := NewCSVRecord([]string{"x", "y", "z"})
	require.NotNil(t, c)

	v, err := c.Find(ch.NewField("*"))
	assert.Nil(t, err)
	assert.True(t, v.IsString())
	assert.Equal(t, "[x y z]", v.AsString())
}
