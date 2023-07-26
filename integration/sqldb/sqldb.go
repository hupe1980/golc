package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/schema"
)

type Engine interface {
	Dialect() string
	SampleRowsQuery(table string, k uint) string
	Inspect(ctx context.Context, name string, opts *schema.InspectOptions) (map[string]string, error)
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	Close() error
}

type SQLDBOptions struct {
	Schema                string
	Tables                []string
	Exclude               []string
	SampleRowsinTableInfo uint
}

type SQLDB struct {
	engine Engine
	opts   SQLDBOptions
}

func New(engine Engine) (*SQLDB, error) {
	opts := SQLDBOptions{}

	return &SQLDB{
		engine: engine,
		opts:   opts,
	}, nil
}

func (db *SQLDB) Dialect() string {
	return db.engine.Dialect()
}

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

type QueryResult struct {
	Columns []string
	Rows    [][]string
}

func (qr *QueryResult) String() string {
	str := strings.Join(qr.Columns, "\t") + "\n"
	for _, row := range qr.Rows {
		str += strings.Join(row, "\t") + "\n"
	}

	return str
}

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
		rowPtrs := make([]any, len(cols))

		for i := range row {
			rowPtrs[i] = &row[i]
		}

		err = rows.Scan(rowPtrs...)
		if err != nil {
			return nil, err
		}

		results = append(results, row)
	}

	return &QueryResult{
		Columns: cols,
		Rows:    results,
	}, nil
}

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

func (db *SQLDB) Close() error {
	return db.engine.Close()
}

type atlas struct {
	driver migrate.Driver
}

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
