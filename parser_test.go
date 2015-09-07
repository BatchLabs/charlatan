package charlatan

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func okQuery(t *testing.T, q string) {
	qry, err := parserFromString(q).Parse()
	require.Nil(t, err)
	require.NotNil(t, qry)
}

func TestParserParseLimit(t *testing.T) {
	for _, s := range []string{
		"SELECT x FROM y limit 42",
		"SELECT x FROM y starting at 3 limit 42",
		"SELECT x FROM y limit 42 starting at 3",
		"SELECT x FROM y WHERE z limit 42",
		"SELECT x FROM y WHERE z = 2 limit 42",
		"SELECT x FROM y WHERE (z = 2) limit 42",
		"SELECT x FROM y WHERE (z = 2 && 43) limit 42",
		"SELECT x FROM y WHERE z starting at 2 limit 42",
		"SELECT x FROM y WHERE z limit 42 starting at 2",
	} {
		okQuery(t, s)
	}
}
