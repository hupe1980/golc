package sqldb

import (
	"context"
	"testing"

	"ariga.io/atlas/sql/schema"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestInspect(t *testing.T) {
	sql, err := NewSQLite3(":memory:")
	require.NoError(t, err)

	defer sql.Close()

	ctx := context.Background()

	_, err = sql.Exec(ctx, "create table example ( id int not null );")
	require.NoError(t, err)

	infos, err := sql.Inspect(ctx, "", &schema.InspectOptions{
		Mode: schema.InspectTables,
	})
	require.NoError(t, err)
	require.Equal(t, "CREATE TABLE `example` (`id` int NOT NULL)", infos["example"])
}
