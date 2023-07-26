package sqldb

import (
	"context"
	"database/sql"
	"fmt"

	"ariga.io/atlas/sql/sqlite"
)

// Compile time check to ensure SQLite3 satisfies the Engine interface.
var _ Engine = (*SQLite3)(nil)

type SQLLite3Options struct {
	DriverName string
}

type SQLite3 struct {
	db *sql.DB
	*atlas
	opts SQLLite3Options
}

func NewSQLite3(dataSourceName string) (*SQLite3, error) {
	opts := SQLLite3Options{
		DriverName: "sqlite3",
	}

	db, err := sql.Open(opts.DriverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	driver, err := sqlite.Open(db)
	if err != nil {
		return nil, fmt.Errorf("failed opening atlas driver: %s", err)
	}

	return &SQLite3{
		db:    db,
		atlas: &atlas{driver: driver},
		opts:  opts,
	}, nil
}

func (e *SQLite3) Dialect() string {
	return "sqlite3"
}

func (e *SQLite3) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return e.db.ExecContext(ctx, query, args...)
}

func (e *SQLite3) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return e.db.QueryContext(ctx, query, args...)
}

func (e *SQLite3) SampleRowsQuery(table string, k uint) string {
	return fmt.Sprintf("SELECT * FROM %s LIMIT %d", table, k)
}

func (e *SQLite3) Close() error {
	return e.db.Close()
}
