package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/kshvakov/clickhouse"
	_ "github.com/lib/pq"
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

	SQLPrettyPrint(rows)

	return nil
}

func PgsqlShell(c *cli.Context) error {

	fmt.Println(kflags)
	if len(kflags.Shell) < 0 {
		fmt.Println("postgres sql is nil")
	}
	config := &KlientConfig.Mysql
	if len(kflags.Host) > 0 {
		config.Host = kflags.Host
		config.Port = kflags.Port
		config.User = kflags.User
		config.Pswd = kflags.Pswd
		config.Db = kflags.Db
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.User, config.Pswd, config.Host, config.Port, config.Db)
	conn, err := sql.Open("postgres", dsn)
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

	SQLPrettyPrint(rows)

	return nil
}

func ClickhouseShell(c *cli.Context) error {
	fmt.Println(kflags)
	config := &KlientConfig.Clickhouse
	if len(kflags.Host) > 0 {
		config.Host = kflags.Host
		config.Port = kflags.Port
		config.User = kflags.User
		config.Pswd = kflags.Pswd
		config.Db = kflags.Db
	}
	fmt.Println(config)
	dsn := fmt.Sprintf("tcp://%s:%s?username=%s&password=%s&database=%s",
		config.Host,
		config.Port,
		config.User,
		config.Pswd,
		config.Db)
	con, err := sql.Open("clickhouse", dsn)
	if err != nil {
		panic(k_red(err))
	} else if err = con.Ping(); err != nil {
		panic(k_red(err))
	}

	rows, err := con.Query(kflags.Shell)
	if err != nil {
		panic(k_red(err))
	}
	defer rows.Close()

	SQLPrettyPrint(rows)

	return nil
}

func SQLPrettyPrint(rows *sql.Rows) {
	cols, err := rows.Columns()
	if err != nil {
		fmt.Println("Error: ", k_red(err))
		return
	}
	vals := make([]sql.RawBytes, len(cols))
	args := make([]interface{}, len(cols))
	for i, _ := range cols {
		args[i] = &vals[i]
	}
	count := 0
	for rows.Next() {
		if count > 100 {
			return
		}
		fmt.Println(k_yellow("******************************* ", count, ".row", " ******************************"))
		err = rows.Scan(args...)
		var val string
		for i, col := range vals {
			if col == nil {
				val = "nil"
			} else {
				val = string(col)
			}
			fmt.Printf("%40s : %s \n", k_cyan(cols[i]), k_green(val))
		}
		count += 1
	}

	return
}
