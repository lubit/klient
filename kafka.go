package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

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
	fmt.Println(config)

	sub := c.Args().First()
	switch sub {
	case "status":
		kafkaShellStatus(c, config)
	case "consume":
		kafkaShellConsume(c, config)
	case "produce":
		kafkaShellProduce(c, config)
	}

	/*
		config.consumer, err = sarama.NewConsumerFromClient(config.client)
		if err != nil {
			fmt.Println("Unable to create new kafka consumer err:", err, config.client)
			return
		}
	*/
	return

}

func kafkaShellStatus(c *cli.Context, kf *KafkaConfigSection) error {
	fmt.Println(kf)

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
		fmt.Println("Unable to create kafka client err:  " + err.Error())
		return err
	}
	// Status Offset
	if len(kf.Topics) < 0 {
		fmt.Println("topic is needed")
	}

	partitions, err := client.Partitions(kf.Topics[0])
	if err != nil {
		fmt.Println("Unable to fetch partition IDs for the topic", err, client, kf.Topics)
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
		offset, offstr := pom.NextOffset()
		/*
			offset, err := config.client.GetOffset(config.Topic, p, sarama.OffsetNewest)
			if err != nil {
				panic(err)
			}
		*/
		fmt.Printf("Topic[%s] Partition[%.2d] Offset : %d, %s \n", kf.Topics, p, offset, offstr)
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

	consumer, err := cluster.NewConsumer(kf.Brokers, kf.Group, kf.Topics, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			fmt.Printf("Error: %s\n", err.Error())
		}
	}()
	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			fmt.Printf("Rebalanced: %+v\n", ntf)
		}
	}()
	// consume msg
	// consume messages, watch signals
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if !ok {
				fmt.Println("kafka no more msg")
				return nil
			}
			fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
			consumer.MarkOffset(msg, "") // mark message as processed
		case err := <-consumer.Errors():
			fmt.Printf("Error: %s\n", err.Error())
		case ntf := <-consumer.Notifications():
			fmt.Printf("Rebalanced: %+v\n", ntf)
		case <-signals:
			return nil
		}
	}

	return nil
}

func kafkaShellProduce(c *cli.Context, kf *KafkaConfigSection) error {
	config := sarama.NewConfig()
	config.Net.SASL.Enable = true
	config.Net.SASL.User = kf.User
	config.Net.SASL.Password = kf.Pswd

	producer, err := sarama.NewAsyncProducer(kf.Brokers, config)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	msg := c.Args().Get(1)
	fmt.Println("msg:", msg)

	producer.Input() <- &sarama.ProducerMessage{
		Topic: kf.Topics[0],
		Value: sarama.StringEncoder(msg),
	}

	go func() {
		for err := range producer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	time.Sleep(1 * time.Second)

	return nil
}
