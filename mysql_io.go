package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli"
)

func MysqlShell(c *cli.Context) error {
	fmt.Println(kflags)
	if len(kflags.Shell) < 0 {
		fmt.Println("mysql sql is nil")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", kflags.User, kflags.Pswd, kflags.Host, kflags.Port, kflags.Db)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	if err = conn.Ping(); err != nil {
		panic(err)
	}
	rows, err := conn.Query(kflags.Shell)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

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
		}
		fmt.Println("")
		count += 1
	}

	return nil
}
