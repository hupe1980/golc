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

type SQLOptions struct {
	*schema.CallbackOptions
	InputKey              string
	TablesInputKey        string
	OutputKey             string
	TopK                  uint
	Schema                string
	Tables                []string
	Exclude               []string
	SampleRowsinTableInfo uint
}

type SQL struct {
	sqldb    *sqldb.SQLDB
	llmChain *LLM
	opts     SQLOptions
}

func NewSQL(llm schema.Model, engine sqldb.Engine, optFns ...func(o *SQLOptions)) (*SQL, error) {
	opts := SQLOptions{
		InputKey:              "query",
		OutputKey:             "result",
		TopK:                  5,
		SampleRowsinTableInfo: 3,
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
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
