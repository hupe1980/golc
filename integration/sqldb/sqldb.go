// Package sqldb provides an SQL database abstraction for performing queries and interacting with the database.
package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/schema"
)

// Engine defines the interface for an SQL database engine.
type Engine interface {
	// Dialect returns the dialect of the SQL database engine.
	Dialect() string

	// SampleRowsQuery returns the query to retrieve a sample of rows from the specified table (table) with a limit of (k) rows.
	SampleRowsQuery(table string, k uint) string

	// Inspect retrieves information about the database schema with the specified name (name) and options (opts).
	// The returned map contains table names as keys and their corresponding CREATE TABLE statements as values.
	Inspect(ctx context.Context, name string, opts *schema.InspectOptions) (map[string]string, error)

	// Exec executes an SQL query with the provided query string and arguments (args), returning the result and any errors encountered.
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)

	// Query executes an SQL query with the provided query string and arguments (args), returning the rows and any errors encountered.
	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)

	// QueryRow executes an SQL query with the provided query string and arguments (args), returning a single row and any errors encountered.
	QueryRow(ctx context.Context, query string, args ...any) *sql.Row

	// Close closes the database connection.
	Close() error
}

// SQLDBOptions holds options for the SQLDB.
type SQLDBOptions struct {
	Schema                string
	Tables                []string
	Exclude               []string
	SampleRowsinTableInfo uint
}

// SQLDB represents an SQL database.
type SQLDB struct {
	engine Engine
	opts   SQLDBOptions
}

// New creates a new SQLDB instance.
func New(engine Engine, optFns ...func(o *SQLDBOptions)) (*SQLDB, error) {
	opts := SQLDBOptions{
		SampleRowsinTableInfo: 3,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &SQLDB{
		engine: engine,
		opts:   opts,
	}, nil
}

// Dialect returns the dialect of the SQL database engine.
func (db *SQLDB) Dialect() string {
	return db.engine.Dialect()
}

// TableInfo retrieves information about the tables in the database.
func (db *SQLDB) TableInfo(ctx context.Context) (string, error) {
	createStmts, err := db.engine.Inspect(ctx, db.opts.Schema, &schema.InspectOptions{
		Tables:  db.opts.Tables,
		Exclude: db.opts.Exclude,
	})
	if err != nil {
		return "", err
	}

	info := ""
	for k, v := range createStmts {
		info += fmt.Sprintf("%s\n\n", v)

		if db.opts.SampleRowsinTableInfo > 0 {
			sampleRows, err := db.sampleRows(ctx, k, db.opts.SampleRowsinTableInfo)
			if err != nil {
				return "", err
			}

			info += "/*\n" + sampleRows + "*/ \n\n"
		}
	}

	return info, nil
}

// QueryResult holds the result of an SQL query.
type QueryResult struct {
	Columns []string
	Rows    [][]string
}

// String returns the string representation of the QueryResult.
func (qr *QueryResult) String() string {
	str := strings.Join(qr.Columns, "\t") + "\n"
	for _, row := range qr.Rows {
		str += strings.Join(row, "\t") + "\n"
	}

	return str
}

// Query executes an SQL query and returns the result.
func (db *SQLDB) Query(ctx context.Context, query string, args ...any) (*QueryResult, error) {
	rows, err := db.engine.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	if rErr := rows.Err(); rErr != nil {
		return nil, rErr
	}

	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	results := make([][]string, 0)

	for rows.Next() {
		row := make([]string, len(cols))
		rowNullable := make([]sql.NullString, len(cols))
		rowPtrs := make([]any, len(cols))

		for i := range row {
			rowPtrs[i] = &rowNullable[i]
		}

		err = rows.Scan(rowPtrs...)
		if err != nil {
			return nil, err
		}

		for i := range rowNullable {
			if rowNullable[i].Valid {
				row[i] = rowNullable[i].String
			}
		}

		results = append(results, row)
	}

	return &QueryResult{
		Columns: cols,
		Rows:    results,
	}, nil
}

// sampleRows retrieves a sample of rows from the given table.
func (db *SQLDB) sampleRows(ctx context.Context, table string, k uint) (string, error) {
	query := db.engine.SampleRowsQuery(table, k)

	result, err := db.Query(ctx, query)
	if err != nil {
		return "", err
	}

	ret := fmt.Sprintf("%d rows from %s table:\n", k, table)

	ret += result.String()

	return ret, nil
}

// Close closes the database connection.
func (db *SQLDB) Close() error {
	return db.engine.Close()
}

// atlas represents the atlas migration driver.
type atlas struct {
	driver migrate.Driver
}

// Inspect retrieves information about the schema.
func (e *atlas) Inspect(ctx context.Context, name string, opts *schema.InspectOptions) (map[string]string, error) {
	s, err := e.driver.InspectSchema(ctx, name, opts)
	if err != nil {
		return nil, err
	}

	var changes schema.Changes
	for _, t := range s.Tables {
		changes = append(changes, &schema.AddTable{T: t})
	}

	plan, err := e.driver.PlanChanges(ctx, "table", changes)
	if err != nil {
		return nil, err
	}

	infos := make(map[string]string, len(s.Tables))
	for i, t := range s.Tables {
		infos[t.Name] = plan.Changes[i].Cmd
	}

	return infos, nil
}
