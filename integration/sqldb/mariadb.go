package sqldb

// Compile time check to ensure MariaDB satisfies the Engine interface.
var _ Engine = (*MariaDB)(nil)

// MariaDBOptions holds options for the MariaDB database engine.
type MariaDBOptions struct {
	DriverName string
}

// MariaDB represents the MariaDB database engine.
type MariaDB struct {
	*MySQL
}

// NewMariaDB creates a new instance of the MariaDB database engine.
func NewMariaDB(dataSourceName string, optFns ...func(o *MariaDBOptions)) (*MariaDB, error) {
	opts := MariaDBOptions{
		DriverName: "mysql",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	mysql, err := NewMySQL(dataSourceName, func(o *MySQLOptions) {
		o.DriverName = opts.DriverName
	})
	if err != nil {
		return nil, err
	}

	return &MariaDB{
		MySQL: mysql,
	}, nil
}

// Dialect returns the dialect of the MariaDB database engine.
func (e *MariaDB) Dialect() string {
	return "MariaDB"
}
