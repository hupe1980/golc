package sqldb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	t.Run("TestIsSelect", func(t *testing.T) {
		testcase := []struct {
			query    string
			isSelect bool
		}{
			{"Select * from table", true},
			{"SELECT * from table", true},
			{"SeLEcT * from table", true},
			{`UPDATE Customers
			SET ContactName = 'Alfred Schmidt', City= 'Frankfurt'
			WHERE CustomerID = 1;`, false},
		}

		for _, tt := range testcase {
			p := NewParser(tt.query)
			require.Equal(t, tt.isSelect, p.IsSelect())
		}
	})
	t.Run("TestTableNames", func(t *testing.T) {
		testcase := []struct {
			query      string
			tableNames []string
		}{
			{
				"SELECT * FROM `table1`, `table2`, `table3` WHERE true;",
				[]string{"table1", "table2", "table3"},
			},
			{
				"SELECT * FROM table1, table2, table3 WHERE true;",
				[]string{"table1", "table2", "table3"},
			},
			{
				`SELECT * FROM table1, table2, table3 WHERE true
				UNION
				SELECT * FROM table4, table5, table6 WHERE true`,
				[]string{"table1", "table2", "table3", "table4", "table5", "table6"},
			},
			{
				`SELECT * FROM table1, table2, table3;`,
				[]string{"table1", "table2", "table3"},
			},
			{
				`SELECT * FROM table1, table2, table3`,
				[]string{"table1", "table2", "table3"},
			},
			{
				`SELECT column_name(s)
				FROM table1
				LEFT JOIN table2
				ON table1.column_name = table2.column_name;`,
				[]string{"table1", "table2"},
			},
			{
				`SELECT column_name(s)
				FROM table1
				LEFT JOIN table2
				ON table1.column_name = table2.column_name;`,
				[]string{"table1", "table2"},
			},
			{
				`UPDATE Customers
				SET ContactName = 'Alfred Schmidt', City= 'Frankfurt'
				WHERE CustomerID = 1;`,
				[]string{"Customers"},
			},
			{
				`DELETE FROM Customers;`,
				[]string{"Customers"},
			},
			{
				`INSERT INTO Customers(CustomerName, ContactName, Address, City, PostalCode, Country)
				VALUES ('Cardinal', 'Tom B. Erichsen', 'Skagen 21', 'Stavanger', '4006', 'Norway');`,
				[]string{"Customers"},
			},
			{
				`SELECT * FROM table1 
				RIGHT JOIN table2 ON table1.id = table2.id 
				LEFT JOIN table3 ON table1.id = table3.id 
				INNER JOIN table4 ON table1.id = table4.id 
				OUTER JOIN table5 ON table1.id = table5.id
				WHERE true`,
				[]string{"table1", "table2", "table3", "table4", "table5"},
			},
			{
				"SELECT * FROM table1 t1, table2 AS t2, table3 as t3 JOIN table4 as t4 on t1.id = t4.id;",
				[]string{"table1", "table2", "table3", "table4"},
			},
		}

		for _, tt := range testcase {
			p := NewParser(tt.query)
			require.Equal(t, tt.tableNames, p.TableNames())
		}
	})
}

func TestCleanQuery(t *testing.T) {
	testcase := []struct {
		query   string
		cleaned string
	}{
		{
			"SELECT * \r\nFROM table\r\n",
			"SELECT * FROM table",
		},
		{
			"SELECT *                   FROM               table",
			"SELECT * FROM table",
		},
		{
			`SELECT * /* some comment*/ 
      FROM table; /* another comment */`,
			"SELECT * FROM table;",
		},
		{
			`SELECT * -- some comment 
        FROM table; -- another comment`,
			"SELECT * FROM table;",
		},
	}

	for _, tt := range testcase {
		require.Equal(t, tt.cleaned, CleanQuery(tt.query))
	}
}
