package sqldb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLDB(t *testing.T) {
	sql, err := NewSQLite3(":memory:")
	require.NoError(t, err)

	defer sql.Close()

	_, err = sql.Exec(context.Background(), "create table example ( id int not null );")
	require.NoError(t, err)

	sqldb, err := New(sql, func(o *SQLDBOptions) {
		o.SampleRowsinTableInfo = 2
	})
	require.NoError(t, err)

	for i := 0; i < 4; i++ {
		_, iErr := sql.Exec(context.Background(), "INSERT INTO example (id) VALUES (?) ;", i)
		if iErr != nil {
			require.NoError(t, iErr)
		}
	}

	t.Run("TestTableInfo", func(t *testing.T) {
		info, err := sqldb.TableInfo(context.Background())
		require.NoError(t, err)

		require.Equal(t, "CREATE TABLE `example` (`id` int NOT NULL)\n\n/*\n2 rows from example table:\nid\n0\n1\n*/ \n\n", info)
	})

	t.Run("TestDialect", func(t *testing.T) {
		require.Equal(t, "sqlite3", sqldb.Dialect())
	})

	t.Run("TestQuery", func(t *testing.T) {
		result, err := sqldb.Query(context.Background(), "SELECT COUNT(*) FROM example")
		require.NoError(t, err)
		require.Equal(t, "COUNT(*)\n4\n", result.String())
	})
}
