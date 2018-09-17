package main

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

type Kflags struct {
	Host          string
	Port          string
	User          string
	Pswd          string
	Db            string
	Shell         string
	Brokers       string
	Topic         string
	Group         string
	FromBeginning bool
}

var (
	kflags Kflags
	// error
	k_red func(a ...interface{}) string
	// warn
	k_yellow func(a ...interface{}) string
	// success
	k_green func(a ...interface{}) string
	// key
	k_cyan func(a ...interface{}) string
	// value
	k_blue func(a ...interface{}) string
)

func main() {

	k_yellow = color.New(color.FgYellow).SprintFunc()
	k_red = color.New(color.FgRed).SprintFunc()
	k_cyan = color.New(color.FgCyan).SprintFunc()
	k_blue = color.New(color.FgBlue).SprintFunc()
	k_green = color.New(color.FgGreen).SprintFunc()

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
		cli.BoolFlag{
			Name:        "from-beginning",
			Destination: &kflags.FromBeginning,
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
			Name:   "pgsql",
			Flags:  app.Flags,
			Action: PgsqlShell,
		},
		{
			Name:   "clickhouse",
			Flags:  app.Flags,
			Action: ClickhouseShell,
		},
	}

	//sort.Sort(cli.FlagsByName(app.Flags))
	//sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
