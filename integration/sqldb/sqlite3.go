package sqldb

import (
	"context"
	"database/sql"
	"fmt"

	"ariga.io/atlas/sql/sqlite"
)

// Compile time check to ensure SQLite3 satisfies the Engine interface.
var _ Engine = (*SQLite3)(nil)

// SQLite3Options holds options for the SQLite3 database engine.
type SQLite3Options struct {
	DriverName string
}

// SQLite3 represents the SQLite3 database engine.
type SQLite3 struct {
	db *sql.DB
	*atlas
	opts SQLite3Options
}

// NewSQLite3 creates a new instance of the SQLite3 database engine.
func NewSQLite3(dataSourceName string, optFns ...func(o *SQLite3Options)) (*SQLite3, error) {
	opts := SQLite3Options{
		DriverName: "sqlite3",
	}

	for _, fn := range optFns {
		fn(&opts)
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

// Dialect returns the dialect of the SQLite3 database engine.
func (e *SQLite3) Dialect() string {
	return "sqlite3"
}

// Exec executes an SQL query with the provided query string and arguments (args), returning the result and any errors encountered.
func (e *SQLite3) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return e.db.ExecContext(ctx, query, args...)
}

// Query executes an SQL query with the provided query string and arguments (args), returning the rows and any errors encountered.
func (e *SQLite3) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return e.db.QueryContext(ctx, query, args...)
}

// QueryRow executes an SQL query with the provided query string and arguments (args), returning a single row and any errors encountered.
func (e *SQLite3) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return e.db.QueryRowContext(ctx, query, args...)
}

// SampleRowsQuery returns the query to retrieve a sample of rows from the specified table (table) with a limit of (k) rows.
func (e *SQLite3) SampleRowsQuery(table string, k uint) string {
	return fmt.Sprintf("SELECT * FROM %s LIMIT %d;", table, k)
}

// Close closes the database connection.
func (e *SQLite3) Close() error {
	return e.db.Close()
}
