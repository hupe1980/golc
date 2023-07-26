package sqldb

import (
	"context"
	"database/sql"
	"fmt"

	"ariga.io/atlas/sql/mysql"
)

// Compile time check to ensure MySQL satisfies the Engine interface.
var _ Engine = (*MySQL)(nil)

// MySQLOptions holds options for the MySQL database engine.
type MySQLOptions struct {
	DriverName string
}

// MySQL represents the MySQL database engine.
type MySQL struct {
	db *sql.DB
	*atlas
	opts MySQLOptions
}

// NewMySQL creates a new instance of the MySQL database engine.
func NewMySQL(dataSourceName string, optFns ...func(o *MySQLOptions)) (*MySQL, error) {
	opts := MySQLOptions{
		DriverName: "mysql",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	db, err := sql.Open(opts.DriverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	driver, err := mysql.Open(db)
	if err != nil {
		return nil, fmt.Errorf("failed opening atlas driver: %s", err)
	}

	return &MySQL{
		db:    db,
		atlas: &atlas{driver: driver},
		opts:  opts,
	}, nil
}

// Dialect returns the dialect of the MySQL database engine.
func (e *MySQL) Dialect() string {
	return "MySQL"
}

// Exec executes an SQL query with the provided query string and arguments (args), returning the result and any errors encountered.
func (e *MySQL) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return e.db.ExecContext(ctx, query, args...)
}

// Query executes an SQL query with the provided query string and arguments (args), returning the rows and any errors encountered.
func (e *MySQL) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return e.db.QueryContext(ctx, query, args...)
}

// QueryRow executes an SQL query with the provided query string and arguments (args), returning a single row and any errors encountered.
func (e *MySQL) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return e.db.QueryRowContext(ctx, query, args...)
}

// SampleRowsQuery returns the query to retrieve a sample of rows from the specified table (table) with a limit of (k) rows.
func (e *MySQL) SampleRowsQuery(table string, k uint) string {
	return fmt.Sprintf("SELECT * FROM %s LIMIT %d", table, k)
}

// Close closes the database connection.
func (e *MySQL) Close() error {
	return e.db.Close()
}
