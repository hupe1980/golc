package chain

import (
	"context"
	"fmt"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/integration/sqldb"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

// defaultSQLTemplate defines the default template for generating SQL queries.
const defaultSQLTemplate = `Given an input question, first create a syntactically correct {{.dialect}} query to run, then look at the results of the query and return the answer. Unless the user specifies in his question a specific number of examples he wishes to obtain, always limit your query to at most {{.topK}} results. You can order the results by a relevant column to return the most interesting examples in the database.

Never query for all the columns from a specific table, only ask for a the few relevant columns given the question.

Pay attention to use only the column names that you can see in the schema description. Be careful to not query for columns that do not exist. Also, pay attention to which column is in which table.

Use the following format:

Question: Question here
SQLQuery: SQL Query to run
SQLResult: Result of the SQLQuery
Answer: Final answer here

Only use the following tables:
{{.tableInfo}}

Question: {{.input}}`

// Compile time check to ensure SQL satisfies the Chain interface.
var _ schema.Chain = (*SQL)(nil)

// VerifySQL is a function signature used to verify the validity of the generated SQL query before execution.
type VerifySQL func(sqlQuery string) bool

// SQLOptions contains options for the SQL chain.
type SQLOptions struct {
	// CallbackOptions contains options for the chain callbacks.
	*schema.CallbackOptions

	// InputKey is the key to access the input value containing the user SQL query.
	InputKey string

	// OutputKey is the key to access the output value containing the SQL query result.
	OutputKey string

	// TopK specifies the maximum number of results to return from the SQL query.
	TopK uint

	// Schema represents the database schema information.
	Schema string

	// Tables is the list of tables to consider when executing the SQL query.
	Tables []string

	// Exclude is the list of tables to exclude when executing the SQL query.
	Exclude []string

	// SampleRowsinTableInfo specifies the number of sample rows to include in the table information.
	SampleRowsinTableInfo uint

	// VerifySQL is a function used to verify the validity of the generated SQL query before execution.
	// It should return true if the SQL query is valid, false otherwise.
	VerifySQL VerifySQL
}

// SQL is a chain implementation that prompts the user to provide an SQL query
// to run on a database. It then verifies and executes the provided SQL query
// and returns the result.
type SQL struct {
	sqldb    *sqldb.SQLDB
	llmChain *LLM
	opts     SQLOptions
}

// NewSQL creates a new instance of the SQL chain.
func NewSQL(llm schema.Model, engine sqldb.Engine, optFns ...func(o *SQLOptions)) (*SQL, error) {
	opts := SQLOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		InputKey:              "query",
		OutputKey:             "result",
		TopK:                  5,
		SampleRowsinTableInfo: 3,
		VerifySQL:             func(sqlQuery string) bool { return true },
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	sqldb, err := sqldb.New(engine, func(o *sqldb.SQLDBOptions) {
		o.Tables = opts.Tables
		o.Exclude = opts.Exclude
		o.SampleRowsinTableInfo = opts.SampleRowsinTableInfo
	})
	if err != nil {
		return nil, err
	}

	llmChain, err := NewLLM(llm, prompt.NewTemplate(defaultSQLTemplate))
	if err != nil {
		return nil, err
	}

	return &SQL{
		sqldb:    sqldb,
		llmChain: llmChain,
		opts:     opts,
	}, nil
}

// Call executes the sql chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *SQL) Call(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
	query, ok := inputs[c.opts.InputKey].(string)
	if !ok {
		return nil, fmt.Errorf("%w: no value for inputKey %s", ErrInvalidInputValues, c.opts.InputKey)
	}

	tableInfo, err := c.sqldb.TableInfo(ctx)
	if err != nil {
		return nil, err
	}

	input := fmt.Sprintf("%s\nSQLQuery:", query)

	sqlQuery, err := golc.SimpleCall(ctx, c.llmChain, schema.ChainValues{
		"dialect":   c.sqldb.Dialect(),
		"input":     input,
		"tableInfo": tableInfo,
		"topK":      c.opts.TopK,
	}, func(sco *golc.SimpleCallOptions) {
		sco.Stop = []string{"\nSQLResult:"}
	})
	if err != nil {
		return nil, err
	}

	sqlQuery = sqldb.CleanQuery(sqlQuery)

	if ok := c.opts.VerifySQL(sqlQuery); !ok {
		return nil, fmt.Errorf("invalid sql query: %s", sqlQuery)
	}

	queryResult, err := c.sqldb.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}

	input += fmt.Sprintf("%s\nSQLResult: %s\nAnswer:", sqlQuery, queryResult)

	result, err := golc.SimpleCall(ctx, c.llmChain, schema.ChainValues{
		"dialect":   c.sqldb.Dialect(),
		"input":     input,
		"tableInfo": tableInfo,
		"topK":      c.opts.TopK,
	}, func(sco *golc.SimpleCallOptions) {
		sco.Stop = []string{"\nSQLResult:"}
	})
	if err != nil {
		return nil, err
	}

	return schema.ChainValues{
		c.opts.OutputKey: strings.TrimSpace(result),
	}, nil
}

// Memory returns the memory associated with the chain.
func (c *SQL) Memory() schema.Memory {
	return nil
}

// Type returns the type of the chain.
func (c *SQL) Type() string {
	return "SQL"
}

// Verbose returns the verbosity setting of the chain.
func (c *SQL) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (c *SQL) Callbacks() []schema.Callback {
	return c.opts.CallbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *SQL) InputKeys() []string {
	return []string{c.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (c *SQL) OutputKeys() []string {
	return []string{c.opts.OutputKey}
}
