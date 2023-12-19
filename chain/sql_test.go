package chain

import (
	"context"
	"strings"
	"testing"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/integration/sqldb"
	"github.com/hupe1980/golc/model/llm"
	"github.com/hupe1980/golc/schema"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestSQL(t *testing.T) {
	ctx := context.Background()

	engine, err := sqldb.NewSQLite3(":memory:")
	assert.NoError(t, err)

	_, err = engine.Exec(ctx, "CREATE TABLE IF NOT EXISTS employee ( id int not null );")
	assert.NoError(t, err)

	for i := 0; i < 4; i++ {
		_, err := engine.Exec(ctx, "INSERT INTO employee (id) VALUES (?) ;", i)
		assert.NoError(t, err)
	}

	t.Run("Valid Question", func(t *testing.T) {
		fake := llm.NewFake(func(ctx context.Context, prompt string) (*schema.ModelResult, error) {
			text := "There are 4 employees."
			if strings.HasSuffix(prompt, "SQLQuery:") {
				text = "SELECT count(*) FROM employee;"
			}

			return &schema.ModelResult{
				Generations: []schema.Generation{{Text: text}},
				LLMOutput:   map[string]any{},
			}, nil
		})

		sqlChain, err := NewSQL(fake, engine)
		assert.NoError(t, err)

		output, err := golc.SimpleCall(ctx, sqlChain, "How many employees are there?")
		assert.NoError(t, err)
		assert.Equal(t, "There are 4 employees.", output)
	})

	t.Run("Invalid Input Key", func(t *testing.T) {
		fake := llm.NewFake(func(ctx context.Context, prompt string) (*schema.ModelResult, error) {
			text := "There are 4 employees."
			if strings.HasSuffix(prompt, "SQLQuery:") {
				text = "SELECT count(*) FROM employee;"
			}

			return &schema.ModelResult{
				Generations: []schema.Generation{{Text: text}},
				LLMOutput:   map[string]any{},
			}, nil
		})

		sqlChain, err := NewSQL(fake, engine)
		assert.NoError(t, err)

		_, err = golc.Call(context.Background(), sqlChain, schema.ChainValues{"invalid_key": "foo"})
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid chain values: no value for key query")
	})

	t.Run("Invalid sql query", func(t *testing.T) {
		fake := llm.NewFake(func(ctx context.Context, prompt string) (*schema.ModelResult, error) {
			text := "There are 4 employees."
			if strings.HasSuffix(prompt, "SQLQuery:") {
				text = "SELECT count(*) FROM employee;"
			}

			return &schema.ModelResult{
				Generations: []schema.Generation{{Text: text}},
				LLMOutput:   map[string]any{},
			}, nil
		})

		sqlChain, err := NewSQL(fake, engine, func(o *SQLOptions) {
			o.VerifySQL = func(sqlQuery string) bool { return false }
		})
		assert.NoError(t, err)

		_, err = golc.SimpleCall(ctx, sqlChain, "How many employees are there?")
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid sql query: SELECT count(*) FROM employee;")
	})

	t.Run("No select sql query", func(t *testing.T) {
		fake := llm.NewFake(func(ctx context.Context, prompt string) (*schema.ModelResult, error) {
			text := "There are 4 employees."
			if strings.HasSuffix(prompt, "SQLQuery:") {
				text = "DROP TABLE employee;"
			}

			return &schema.ModelResult{
				Generations: []schema.Generation{{Text: text}},
				LLMOutput:   map[string]any{},
			}, nil
		})

		sqlChain, err := NewSQL(fake, engine)
		assert.NoError(t, err)

		_, err = golc.SimpleCall(ctx, sqlChain, "How many employees are there?")
		assert.Error(t, err)
		assert.EqualError(t, err, "unsupported sql query: DROP TABLE employee;")
	})

	t.Run("table exlude", func(t *testing.T) {
		fake := llm.NewFake(func(ctx context.Context, prompt string) (*schema.ModelResult, error) {
			text := "There are 4 employees."
			if strings.HasSuffix(prompt, "SQLQuery:") {
				text = "SELECT count(*) FROM employee;"
			}

			return &schema.ModelResult{
				Generations: []schema.Generation{{Text: text}},
				LLMOutput:   map[string]any{},
			}, nil
		})

		sqlChain, err := NewSQL(fake, engine, func(o *SQLOptions) {
			o.Exclude = []string{"employee"}
		})
		assert.NoError(t, err)

		_, err = golc.SimpleCall(ctx, sqlChain, "How many employees are there?")
		assert.Error(t, err)
		assert.EqualError(t, err, "not allowed table: employee")
	})

	t.Run("not in whitelist", func(t *testing.T) {
		fake := llm.NewFake(func(ctx context.Context, prompt string) (*schema.ModelResult, error) {
			text := "There are 4 employees."
			if strings.HasSuffix(prompt, "SQLQuery:") {
				text = "SELECT count(*) FROM employee;"
			}

			return &schema.ModelResult{
				Generations: []schema.Generation{{Text: text}},
				LLMOutput:   map[string]any{},
			}, nil
		})

		sqlChain, err := NewSQL(fake, engine, func(o *SQLOptions) {
			o.Tables = []string{"table"}
		})
		assert.NoError(t, err)

		_, err = golc.SimpleCall(ctx, sqlChain, "How many employees are there?")
		assert.Error(t, err)
		assert.EqualError(t, err, "not allowed table: employee")
	})
}
