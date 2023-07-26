package sqldb

import (
	"context"
	"testing"

	"ariga.io/atlas/sql/schema"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLite3(t *testing.T) {
	sql, err := NewSQLite3(":memory:")
	require.NoError(t, err)

	defer sql.Close()

	_, err = sql.Exec(context.Background(), "create table example ( id int not null );")
	require.NoError(t, err)

	for i := 0; i < 4; i++ {
		_, iErr := sql.Exec(context.Background(), "INSERT INTO example (id) VALUES (?) ;", i)
		if iErr != nil {
			require.NoError(t, iErr)
		}
	}

	t.Run("TestInspect", func(t *testing.T) {
		infos, err := sql.Inspect(context.Background(), "", &schema.InspectOptions{
			Mode: schema.InspectTables,
		})
		require.NoError(t, err)
		require.Equal(t, "CREATE TABLE `example` (`id` int NOT NULL)", infos["example"])
	})

	t.Run("TestQuery", func(t *testing.T) {
		rows, err := sql.Query(context.Background(), "SELECT * FROM example")
		require.NoError(t, err)

		defer rows.Close()

		cols, err := rows.Columns()
		require.NoError(t, err)
		require.Equal(t, []string{"id"}, cols)
	})

	t.Run("TestQueryRow", func(t *testing.T) {
		row := sql.QueryRow(context.Background(), "SELECT COUNT(*) FROM example")

		var count int
		err := row.Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 4, count)
	})

	t.Run("TestDialect", func(t *testing.T) {
		assert.Equal(t, "sqlite3", sql.Dialect())
	})

	t.Run("TestSampleRowsQuery", func(t *testing.T) {
		assert.Equal(t, "SELECT * FROM eample LIMIT 5;", sql.SampleRowsQuery("eample", 5))
	})
}
