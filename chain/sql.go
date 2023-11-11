package chain

import (
	"context"
	"fmt"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/integration/sqldb"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/util"
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
//
// WARNING: The SQL chain is a powerful tool for executing SQL queries dynamically. However, it should be used with caution
// to prevent potential SQL injection vulnerabilities. SQL injection is a serious security risk that can lead to unauthorized
// access, data manipulation, and potentially compromising the entire database.
//
// To mitigate the risks of SQL injection, it is crucial to follow these best practices while using the SQL chain:
//
//   - Least Privilege Principle: Ensure that the database user used in the application has the least privilege necessary
//     to perform its required tasks. Restrict the user's permissions to only the required tables and operations.
//
//   - Table Whitelisting or Blacklisting: Use the Tables or Exclude options to reduce the allowed tables that can be accessed
//     via the SQL chain. This will limit the potential impact of any SQL injection attack by restricting the scope of accessible tables.
//
//   - VerifySQL Hook: Implement the VerifySQL hook diligently to validate and sanitize user input. This hook should be used to check
//     and ensure that the generated SQL queries are safe and adhere to the allowed tables and queries.
//
// It is the responsibility of the application developers and administrators to ensure the secure usage of the SQL chain. Failure
// to do so can lead to severe security breaches and compromise the integrity of the application and database. We strongly recommend
// thorough testing, security reviews, and adherence to secure coding practices to protect against SQL injection and other security threats.
type SQL struct {
	sqldb    *sqldb.SQLDB
	llmChain *LLM
	opts     SQLOptions
}

// NewSQL creates a new instance of the SQL chain.
func NewSQL(model schema.Model, engine sqldb.Engine, optFns ...func(o *SQLOptions)) (*SQL, error) {
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
		o.Tables = append([]string(nil), opts.Tables...)
		o.Exclude = append([]string(nil), opts.Exclude...)
		o.SampleRowsinTableInfo = opts.SampleRowsinTableInfo
	})
	if err != nil {
		return nil, err
	}

	llmChain, err := NewLLM(model, prompt.NewTemplate(defaultSQLTemplate))
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
	opts := schema.CallOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	query, err := inputs.GetString(c.opts.InputKey)
	if err != nil {
		return nil, err
	}

	if cbErr := opts.CallbackManger.OnText(ctx, &schema.TextManagerInput{
		Text: query,
	}); cbErr != nil {
		return nil, cbErr
	}

	tableInfo, err := c.sqldb.TableInfo(ctx)
	if err != nil {
		return nil, err
	}

	if cbErr := opts.CallbackManger.OnText(ctx, &schema.TextManagerInput{
		Text: tableInfo,
	}); cbErr != nil {
		return nil, cbErr
	}

	input := fmt.Sprintf("%s\nSQLQuery:", query)

	sqlQuery, err := golc.SimpleCall(ctx, c.llmChain, schema.ChainValues{
		"dialect":   c.sqldb.Dialect(),
		"input":     input,
		"tableInfo": tableInfo,
		"topK":      c.opts.TopK,
	}, func(sco *golc.SimpleCallOptions) {
		sco.Callbacks = opts.CallbackManger.GetInheritableCallbacks()
		sco.ParentRunID = opts.CallbackManger.RunID()
		sco.Stop = []string{"\nSQLResult:"}
	})
	if err != nil {
		return nil, err
	}

	sqlQuery = sqldb.CleanQuery(sqlQuery)

	p := sqldb.NewParser(sqlQuery)

	if !p.IsSelect() {
		return nil, fmt.Errorf("unsupported sql query: %s", sqlQuery)
	}

	if ctErr := c.checkTables(p.TableNames()); ctErr != nil {
		return nil, ctErr
	}

	if ok := c.opts.VerifySQL(sqlQuery); !ok {
		return nil, fmt.Errorf("invalid sql query: %s", sqlQuery)
	}

	if cbErr := opts.CallbackManger.OnText(ctx, &schema.TextManagerInput{
		Text: sqlQuery,
	}); cbErr != nil {
		return nil, cbErr
	}

	queryResult, err := c.sqldb.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}

	if cbErr := opts.CallbackManger.OnText(ctx, &schema.TextManagerInput{
		Text: queryResult.String(),
	}); cbErr != nil {
		return nil, cbErr
	}

	input += fmt.Sprintf("%s\nSQLResult: %s\nAnswer:", sqlQuery, queryResult)

	result, err := golc.SimpleCall(ctx, c.llmChain, schema.ChainValues{
		"dialect":   c.sqldb.Dialect(),
		"input":     input,
		"tableInfo": tableInfo,
		"topK":      c.opts.TopK,
	}, func(sco *golc.SimpleCallOptions) {
		sco.Callbacks = opts.CallbackManger.GetInheritableCallbacks()
		sco.ParentRunID = opts.CallbackManger.RunID()
	})
	if err != nil {
		return nil, err
	}

	result = strings.TrimSpace(result)

	if cbErr := opts.CallbackManger.OnText(ctx, &schema.TextManagerInput{
		Text: result,
	}); cbErr != nil {
		return nil, cbErr
	}

	return schema.ChainValues{
		c.opts.OutputKey: result,
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

// checkTables checks if the provided tables are allowed based on the options specified in SQLOptions.
// If the Tables option is set, it verifies that the tables are present in the allowed list.
// If the Exclude option is set, it verifies that the tables are not present in the excluded list.
// If a table is not allowed, it returns an error indicating that the table is not allowed for use in the SQL query.
func (c *SQL) checkTables(tables []string) error {
	if len(c.opts.Tables) > 0 {
		for _, t := range tables {
			if !util.Contains(c.opts.Tables, strings.ToLower(t)) {
				return fmt.Errorf("not allowed table: %s", t)
			}
		}
	}

	if len(c.opts.Exclude) > 0 {
		for _, e := range c.opts.Exclude {
			if util.Contains(tables, strings.ToLower(e)) {
				return fmt.Errorf("not allowed table: %s", e)
			}
		}
	}

	return nil
}
