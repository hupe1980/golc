package sqldb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCleanQuery(t *testing.T) {
	testcase := []struct {
		query   string
		cleaned string
	}{
		{
			"SELECT * \r\nFROM table\r\n",
			"SELECT * FROM table",
		},
		{
			"SELECT *                   FROM               table",
			"SELECT * FROM table",
		},
		{
			`SELECT * /* some comment*/ 
      FROM table; /* another comment */`,
			"SELECT * FROM table;",
		},
		{
			`SELECT * -- some comment 
        FROM table; -- another comment`,
			"SELECT * FROM table;",
		},
	}

	for _, tt := range testcase {
		require.Equal(t, tt.cleaned, CleanQuery(tt.query))
	}
}
