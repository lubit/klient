package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/Shopify/sarama"
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
			Name:        "shell, e",
			Destination: &kflags.Shell,
		},
	}
	app.Flags = flags

	// Default Command: Kafka
	app.Action = kafkaShell
	app.Commands = []cli.Command{
		{
			Name:   "kafka",
			Flags:  app.Flags,
			Action: kafkaShell,
		},
		{
			Name:   "redis",
			Flags:  flags,
			Action: redisShell,
		},
		{
			Name:   "mysql",
			Flags:  flags,
			Action: mysqlShell,
		},
		{
			Name:   "clickhouse",
			Flags:  flags,
			Action: clickhouseShell,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type KafkaConfig struct {
	Brokers  []string
	Topic    string
	User     string
	Password string
	consumer sarama.Consumer
	client   sarama.Client
}

func kafkaShell(c *cli.Context) (err error) {

	config := KafkaConfig{}
	config.Brokers = strings.Split(kflags.Brokers, ",")
	config.Topic = kflags.Topic
	config.User = kflags.User
	config.Password = kflags.Pswd

	consumerConfig := sarama.NewConfig()
	consumerConfig.Net.SASL.User = kflags.User
	consumerConfig.Net.SASL.Password = kflags.Pswd
	consumerConfig.Net.SASL.Handshake = true
	consumerConfig.Net.SASL.Enable = true

	consumerConfig.Net.TLS.Enable = true
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ClientAuth:         0,
	}
	consumerConfig.Net.TLS.Config = tlsConfig

	fmt.Println("subcommands:", c.Args())
	config.client, err = sarama.NewClient(config.Brokers, consumerConfig)
	if err != nil {
		fmt.Println("Unable to create kafka client " + err.Error())
		return
	}

	config.consumer, err = sarama.NewConsumerFromClient(config.client)
	if err != nil {
		fmt.Println("Unable to create new kafka consumer", err, config.client)
		return
	}

	partitions, err := config.client.Partitions(config.Topic)

	if err != nil {
		fmt.Println("Unable to fetch partition IDs for the topic", err, config.client, config.Topic)
		return
	}

	fmt.Println("Partitions:", partitions)

	topics, err := config.client.Topics()
	if err != nil {
		fmt.Println("Unable to fetch topics", err, config.client)
		return
	}
	for _, v := range topics {
		for _, p := range partitions {
			offset, err := config.client.GetOffset(v, p, sarama.OffsetNewest)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Topic[%s] Partition[%.2d] Offset : %d \n", v, p, offset)
		}
	}

	return

}

func redisShell(c *cli.Context) error {
	fmt.Println(kflags)
	return nil
}

func mysqlShell(c *cli.Context) error {
	fmt.Println(kflags)
	return nil
}

func clickhouseShell(c *cli.Context) error {
	fmt.Println(kflags)
	return nil
}
