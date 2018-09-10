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
	config := &KlientConfig.Mysql
	if len(kflags.Host) > 0 {
		config.Host = kflags.Host
		config.Port = kflags.Port
		config.User = kflags.User
		config.Pswd = kflags.Pswd
		config.Db = kflags.Db
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.User, config.Pswd, config.Host, config.Port, config.Db)
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
