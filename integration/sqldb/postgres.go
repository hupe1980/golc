package sqldb

import (
	"context"
	"database/sql"
	"fmt"

	"ariga.io/atlas/sql/postgres"
)

// Compile time check to ensure Postgres satisfies the Engine interface.
var _ Engine = (*Postgres)(nil)

type PostgresOptions struct {
	DriverName string
}

type Postgres struct {
	db *sql.DB
	*atlas
	opts PostgresOptions
}

func NewPostgres(dataSourceName string) (*Postgres, error) {
	opts := PostgresOptions{
		DriverName: "pgx",
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

func (e *Postgres) Dialect() string {
	return "Postgres"
}

func (e *Postgres) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return e.db.ExecContext(ctx, query, args...)
}

func (e *Postgres) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return e.db.QueryContext(ctx, query, args...)
}

func (e *Postgres) SampleRowsQuery(table string, k uint) string {
	return fmt.Sprintf("SELECT * FROM %s LIMIT %d", table, k)
}

func (e *Postgres) Close() error {
	return e.db.Close()
}
