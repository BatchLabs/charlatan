package charlatan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldName(t *testing.T) {
	assert.Equal(t, "yo", NewField("yo").Name())
}
