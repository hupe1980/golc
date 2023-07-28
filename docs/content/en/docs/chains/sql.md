---
title: SQL
description: All about sql chains.
weight: 90
---
{{% alert title="Warning" color="warning" %}}
The SQL chain is a powerful tool for executing SQL queries dynamically. However, it should be used with caution to prevent potential SQL injection vulnerabilities. SQL injection is a serious security risk that can lead to unauthorized access, data manipulation, and potentially compromising the entire database.

To mitigate the risks of SQL injection, it is crucial to follow these best practices while using the SQL chain:

- Least Privilege Principle: Ensure that the database user used in the application has the least privilege necessary to perform its required tasks. Restrict the user's permissions to only the required tables and operations.

- Table Whitelisting or Blacklisting: Use the Tables or Exclude options to reduce the allowed tables that can be accessed via the SQL chain. This will limit the potential impact of any SQL injection attack by restricting the scope of accessible tables.

- VerifySQL Hook: Implement the VerifySQL hook diligently to validate and sanitize user input. This hook should be used to check and ensure that the generated SQL queries are safe and adhere to the allowed tables and queries.

It is the responsibility of the application developers and administrators to ensure the secure usage of the SQL chain. Failure to do so can lead to severe security breaches and compromise the integrity of the application and database. We strongly recommend thorough testing, security reviews, and adherence to secure coding practices to protect against SQL injection and other security threats.

See an example below.
{{% /alert %}}

```go
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

	// Only for demonstration
    _, exErr := engine.Exec(ctx, "CREATE TABLE IF NOT EXISTS employee ( id int not null );")
	if exErr != nil {
		log.Fatal(exErr)
	}

	// Only for demonstration
    for i := 0; i < 4; i++ {
		_, qErr := engine.Exec(ctx, "INSERT INTO employee (id) VALUES (?) ;", i)
		if qErr != nil {
			log.Fatal(qErr)
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
```
Output:
```text
There are 4 employees.
```

## Supported databases
MySQL, MariaDB, PostgresSQL, SQLite, CockroachDB

## Golang SQL Drivers
https://github.com/golang/go/wiki/SQLDrivers