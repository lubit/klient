package main

import (
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

	app := cli.NewApp()
	app.Author = "罗发宣"
	app.Version = "0.0.1"
	flags := []cli.Flag{
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
	app.Flags = flags

	// Default Command: Kafka
	app.Action = KafkaShell
	app.Commands = []cli.Command{
		{
			Name:   "kafka",
			Flags:  app.Flags,
			Action: KafkaShell,
		},
		{
			Name:   "redis",
			Flags:  flags,
			Action: RedisShell,
		},
		{
			Name:   "mysql",
			Flags:  flags,
			Action: MysqlShell,
		},
		{
			Name:   "clickhouse",
			Flags:  flags,
			Action: ClickhouseShell,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
