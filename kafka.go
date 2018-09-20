package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/urfave/cli"
)

func KafkaShell(c *cli.Context) (err error) {

	config := &KlientConfig.Kafka
	if len(kflags.Brokers) > 0 {
		config.Brokers = strings.Split(kflags.Brokers, ",")
		config.Topics = strings.Split(kflags.Topic, ",")
		config.User = kflags.User
		config.Group = kflags.Group
		config.Pswd = kflags.Pswd
	}

	if len(config.Group) <= 0 {
		config.Group = KafkaDefaultGroup
	}

	sub := c.Args().First()
	switch sub {
	case "status":
		kafkaShellStatus(c, config)
	case "consume":
		kafkaShellConsume(c, config)
	case "produce":
		kafkaShellProduce(c, config)
	}

	return

}

func kafkaShellStatus(c *cli.Context, kf *KafkaConfigSection) error {

	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumerConfig := sarama.NewConfig()
	if len(kf.User) > 1 {
		consumerConfig.Net.SASL.Enable = true
		consumerConfig.Net.SASL.User = kf.User
		consumerConfig.Net.SASL.Password = kf.Pswd
	}
	client, err := sarama.NewClient(kf.Brokers, consumerConfig)
	if err != nil {
		fmt.Println("Unable to create kafka client err:  " + k_red(err.Error()))
		return err
	}
	// Status Offset
	if len(kf.Topics) < 0 {
		fmt.Println("topic is needed")
	}

	partitions, err := client.Partitions(kf.Topics[0])
	if err != nil {
		fmt.Println("Unable to fetch partition IDs for the topic", k_red(err), k_yellow(client), k_green(kf.Topics))
		return err
	}

	ofmg, err := sarama.NewOffsetManagerFromClient(kf.Group, client)
	if err != nil {
		panic(err)
	}

	for _, p := range partitions {
		pom, err := ofmg.ManagePartition(kf.Topics[0], p)
		if err != nil {
			panic(err)
		}
		offset, _ := pom.NextOffset()

		partOff, err := client.GetOffset(kf.Topics[0], p, sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}
		lag := partOff - offset
		fmt.Printf("Topic[%s] Partition[%2s] Offset[%s] : lag %s \n", k_yellow(kf.Topics), k_cyan(p), k_green(offset), k_cyan(lag))
	}
	return nil

}

func kafkaShellConsume(c *cli.Context, kf *KafkaConfigSection) error {

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	if len(kf.User) > 0 {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = kf.User
		config.Net.SASL.Password = kf.Pswd
	}
	if kflags.FromBeginning {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	consumer, err := cluster.NewConsumer(kf.Brokers, kf.Group, kf.Topics, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume msg
	// consume messages, watch signals
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if !ok {
				fmt.Println("kafka no more msg")
				return nil
			}
			fmt.Printf("topic[%s] parition[%2s] offset[%s]: \tkey:%s\t value: %s \n", k_green(msg.Topic), k_green(msg.Partition), k_green(msg.Offset), k_cyan(string(msg.Key)), k_cyan(string(msg.Value)))
			consumer.MarkOffset(msg, "") // mark message as processed
		case err := <-consumer.Errors():
			fmt.Printf("Error: %s\n", k_red(err.Error()))
		case ntf := <-consumer.Notifications():
			fmt.Printf("Rebalanced: %+v\n", k_yellow(ntf))
		case <-signals:
			return nil
		}
	}

	return nil
}

func kafkaShellProduce(c *cli.Context, kf *KafkaConfigSection) {
	config := sarama.NewConfig()
	if len(kf.User) > 0 {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = kf.User
		config.Net.SASL.Password = kf.Pswd
	}

	producer, err := sarama.NewAsyncProducer(kf.Brokers, config)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	// input from terminal
	in_ch := make(chan string, 10)
	in_ch <- c.Args().Get(1)
	go func() {
		in_reader := bufio.NewReader(os.Stdin)
		for {
			msg, err := in_reader.ReadString('\n')
			if err != nil {
				panic(k_red(err))
			}
			in_ch <- msg
		}
	}()

	for {
		select {
		case msg := <-in_ch:
			producer.Input() <- &sarama.ProducerMessage{
				Topic: kf.Topics[0],
				Value: sarama.StringEncoder(msg),
			}

		case err := <-producer.Errors():
			fmt.Println(k_red(err))

		}
	}
}
