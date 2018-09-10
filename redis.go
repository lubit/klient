package main

import (
	"fmt"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/urfave/cli"
)

func RedisShell(c *cli.Context) error {
	fmt.Println(kflags)
	config := &KlientConfig.Redis
	if len(kflags.Host) > 0 {
		config.Host = kflags.Host
		config.Port = kflags.Port
		config.User = kflags.User
		config.Pswd = kflags.Pswd
		config.Db = kflags.Db
	}

	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	con, err := redis.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer con.Close()

	args := strings.Split(kflags.Shell, " ")
	if len(args) > 1 {
		params := make([]interface{}, len(args)-1)
		for i, v := range args[1:] {
			fmt.Println(i, v)
		}
		res, err := redis.String(con.Do(args[0], params))
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	} else {
		res, err := redis.String(con.Do(args[0]))
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}
	return nil
}
