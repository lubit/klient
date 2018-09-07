package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/urfave/cli"
)

type Kflags struct {
	Host    string
	Port    string
	User    string
	Pswd    string
	Db      string
	Shell   string
	Brokers string
	Topic   string
	Group   string
}

var kflags Kflags

func main() {

	LoadConfig()
	defer SaveConfig()

	app := cli.NewApp()
	app.Author = "罗发宣"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "host, H",
			Destination: &kflags.Host,
		},
		cli.StringFlag{
			Name:        "port, P",
			Destination: &kflags.Port,
		},
		cli.StringFlag{
			Name:        "user, u",
			Destination: &kflags.User,
		},
		cli.StringFlag{
			Name:        "password, p",
			Destination: &kflags.Pswd,
		},
		cli.StringFlag{
			Name:        "database, D",
			Destination: &kflags.Db,
		},
		cli.StringFlag{
			Name:        "brokers, B",
			Destination: &kflags.Brokers,
		},
		cli.StringFlag{
			Name:        "topic, T",
			Destination: &kflags.Topic,
		},
		cli.StringFlag{
			Name:        "group, G",
			Destination: &kflags.Group,
		},
		cli.StringFlag{
			Name:        "shell, e",
			Destination: &kflags.Shell,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "kafka",
			Flags:  app.Flags,
			Action: KafkaShell,
		},
		{
			Name:   "redis",
			Flags:  app.Flags,
			Action: RedisShell,
		},
		{
			Name:   "mysql",
			Flags:  app.Flags,
			Action: MysqlShell,
		},
		{
			Name:   "clickhouse",
			Flags:  app.Flags,
			Action: ClickhouseShell,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	fmt.Println(kflags)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
