package sqldb

import (
	"context"
	"database/sql"
	"fmt"

	"ariga.io/atlas/sql/postgres"
)

// Compile time check to ensure Postgres satisfies the Engine interface.
var _ Engine = (*Postgres)(nil)

// PostgresOptions holds options for the Postgres database engine.
type PostgresOptions struct {
	DriverName string
}

// Postgres represents the Postgres database engine.
type Postgres struct {
	db *sql.DB
	*atlas
	opts PostgresOptions
}

// NewPostgres creates a new instance of the Postgres database engine.
func NewPostgres(dataSourceName string, optFns ...func(o *PostgresOptions)) (*Postgres, error) {
	opts := PostgresOptions{
		DriverName: "pgx",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	db, err := sql.Open(opts.DriverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	driver, err := postgres.Open(db)
	if err != nil {
		return nil, fmt.Errorf("failed opening atlas driver: %s", err)
	}

	return &Postgres{
		db:    db,
		atlas: &atlas{driver: driver},
		opts:  opts,
	}, nil
}

// Dialect returns the dialect of the Postgres database engine.
func (e *Postgres) Dialect() string {
	return "Postgres"
}

// Exec executes an SQL query with the provided query string and arguments (args), returning the result and any errors encountered.
func (e *Postgres) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return e.db.ExecContext(ctx, query, args...)
}

// Query executes an SQL query with the provided query string and arguments (args), returning the rows and any errors encountered.
func (e *Postgres) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return e.db.QueryContext(ctx, query, args...)
}

// QueryRow executes an SQL query with the provided query string and arguments (args), returning a single row and any errors encountered.
func (e *Postgres) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return e.db.QueryRowContext(ctx, query, args...)
}

// SampleRowsQuery returns the query to retrieve a sample of rows from the specified table (table) with a limit of (k) rows.
func (e *Postgres) SampleRowsQuery(table string, k uint) string {
	return fmt.Sprintf("SELECT * FROM %s LIMIT %d", table, k)
}

// Close closes the database connection.
func (e *Postgres) Close() error {
	return e.db.Close()
}
