package sqldb

// Compile time check to ensure CockroachDB satisfies the Engine interface.
var _ Engine = (*CockroachDB)(nil)

// CockroachDBOptions holds options for the CockroachDB database engine.
type CockroachDBOptions struct {
	DriverName string
}

// CockroachDB represents the CockroachDB database engine.
type CockroachDB struct {
	*Postgres
}

// NewCockroachDB creates a new instance of the CockroachDB database engine.
func NewCockroachDB(dataSourceName string, optFns ...func(o *CockroachDBOptions)) (*CockroachDB, error) {
	opts := CockroachDBOptions{
		DriverName: "pgx",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	postgres, err := NewPostgres(dataSourceName, func(o *PostgresOptions) {
		o.DriverName = opts.DriverName
	})
	if err != nil {
		return nil, err
	}

	return &CockroachDB{
		Postgres: postgres,
	}, nil
}

// Dialect returns the dialect of the CockroachDB database engine.
func (e *CockroachDB) Dialect() string {
	return "CockroachDB"
}
