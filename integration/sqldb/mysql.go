package sqldb

import (
	"context"
	"database/sql"
	"fmt"

	"ariga.io/atlas/sql/mysql"
)

// Compile time check to ensure MySQL satisfies the Engine interface.
var _ Engine = (*MySQL)(nil)

type MySQLOptions struct {
	DriverName string
}

type MySQL struct {
	db *sql.DB
	*atlas
	opts MySQLOptions
}

func NewMySQL(dataSourceName string) (*MySQL, error) {
	opts := MySQLOptions{
		DriverName: "mysql",
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

func (e *MySQL) Dialect() string {
	return "MySQL"
}

func (e *MySQL) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return e.db.ExecContext(ctx, query, args...)
}

func (e *MySQL) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return e.db.QueryContext(ctx, query, args...)
}

func (e *MySQL) SampleRowsQuery(table string, k uint) string {
	return fmt.Sprintf("SELECT * FROM %s LIMIT %d", table, k)
}

func (e *MySQL) Close() error {
	return e.db.Close()
}
