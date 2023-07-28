package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/integration/sqldb"
	"github.com/hupe1980/golc/model/llm"

	// Add your sql db driver, see https://github.com/golang/go/wiki/SQLDrivers
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx := context.Background()

	openai, err := llm.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	engine, err := sqldb.NewSQLite3(":memory:")
	if err != nil {
		log.Fatal(err)
	}

	defer engine.Close()

	_, exErr := engine.Exec(ctx, "CREATE TABLE IF NOT EXISTS employee ( id int not null );")
	if exErr != nil {
		log.Fatal(exErr)
	}

	for i := 0; i < 4; i++ {
		_, iErr := engine.Exec(ctx, "INSERT INTO employee (id) VALUES (?);", i)
		if iErr != nil {
			log.Fatal(iErr)
		}
	}

	sql, err := chain.NewSQL(openai, engine, func(o *chain.SQLOptions) {
		o.Tables = []string{"employee"}
	})
	if err != nil {
		log.Fatal(err)
	}

	result, err := golc.SimpleCall(ctx, sql, "How many employees are there?")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
