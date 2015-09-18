package charlatan

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func okQuery(t *testing.T, q string) {
	qry, err := parserFromString(q).Parse()
	require.Nil(t, err, "There should be no error parsing '%s'", q)
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

func TestParserParseLimitWithOffset(t *testing.T) {
	for _, s := range []string{
		"SELECT x FROM y limit 0, 42",
		"SELECT x FROM y starting at 3 limit 17, 42",
		"SELECT x FROM y limit 42, 78 starting at 3",
		"SELECT x FROM y WHERE z limit 1, 42",
		"SELECT x FROM y WHERE z = 2 limit 4, 42",
		"SELECT x FROM y WHERE (z = 2) limit 45, 42",
		"SELECT x FROM y WHERE (z = 2 && 43) limit 2, 42",
		"SELECT x FROM y WHERE z starting at 2 limit 3, 42",
		"SELECT x FROM y WHERE z limit 44,2 starting at 45",
	} {
		okQuery(t, s)
	}
}
