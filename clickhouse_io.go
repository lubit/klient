package main

import (
	"database/sql"
	"fmt"

	_ "github.com/kshvakov/clickhouse"
	"github.com/urfave/cli"
)

func ClickhouseShell(c *cli.Context) error {
	fmt.Println(kflags)
	dsn := fmt.Sprintf("tcp://%s:%s?username=%s&password=%s&database=%s",
		kflags.Host,
		kflags.Port,
		kflags.User,
		kflags.Pswd,
		kflags.Db)
	con, err := sql.Open("clickhouse", dsn)
	if err != nil {
		panic(err)
	} else if err = con.Ping(); err != nil {
		panic(err)
	}

	rows, err := con.Query(kflags.Shell)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	rows.ColumnTypes()
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	count := 0
	for rows.Next() {
		if count > 100 {
			return nil
		}
		err = rows.Scan(scanArgs...)

		for _, col := range values {
			fmt.Printf(" %v \t", col)
			/*
				switch col.(type) {
				case string:
					record[columns[i]] = string(col.([]byte))
				case int64:
					fmt.Println(col)
				default:
					fmt.Println(col)
				}

					if col != nil {
						record[columns[i]] = string(col.([]byte))
					}
			*/
		}
		fmt.Println("")
		count += 1
	}
	return nil
}
