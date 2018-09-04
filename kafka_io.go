package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bsm/sarama-cluster"

	"github.com/Shopify/sarama"
	"github.com/urfave/cli"
)

type KafkaFlags struct {
	Brokers  []string
	Topic    []string
	Group    string
	User     string
	Password string
}

func KafkaShell(c *cli.Context) (err error) {

	config := KafkaFlags{}
	config.Brokers = strings.Split(kflags.Brokers, ",")
	config.Topic = strings.Split(kflags.Topic, ",")
	config.User = kflags.User
	config.Group = kflags.Group
	config.Password = kflags.Pswd

	sub := c.Args().First()
	switch sub {
	case "status":
		kafkaShellStatus(c, &config)
	case "consume":
		kafkaShellConsume(c, &config)
	case "produce":
		kafkaShellProduce(c, &config)
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

func kafkaShellStatus(c *cli.Context, kf *KafkaFlags) error {
	consumerConfig := sarama.NewConfig()
	consumerConfig.Net.SASL.Enable = true
	consumerConfig.Net.SASL.User = kflags.User
	consumerConfig.Net.SASL.Password = kflags.Pswd

	client, err := sarama.NewClient(kf.Brokers, consumerConfig)
	if err != nil {
		fmt.Println("Unable to create kafka client err:  " + err.Error())
		return err
	}
	// Status Offset
	if len(kf.Topic) > 0 {

		partitions, err := client.Partitions(kf.Topic[0])
		if err != nil {
			fmt.Println("Unable to fetch partition IDs for the topic", err, client, kf.Topic)
			return err
		}
		ofmg, err := sarama.NewOffsetManagerFromClient(kf.Group, client)
		if err != nil {
			panic(err)
		}

		for _, p := range partitions {
			pom, err := ofmg.ManagePartition(kf.Topic[0], p)
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
			fmt.Printf("Topic[%s] Partition[%.2d] Offset : %d, %s \n", kf.Topic, p, offset, offstr)
		}
	} else {
		/*
			topics, err := client.Topics()
			if err != nil {
				fmt.Println("Unable to fetch topics", err, kf.client)
				return err
			}
			for _, v := range topics {
				partitions, err := client.Partitions(v)
				if err != nil {
					fmt.Println("Unable to fetch partition IDs for the topic", err, kf.client, v)
					return err
				}
				for _, p := range partitions {
					offset, err := client.GetOffset(v, p, sarama.OffsetNewest)
					if err != nil {
						panic(err)
					}
					fmt.Printf("Topic[%s] Partition[%.2d] Offset : %d \n", v, p, offset)
				}
			}
		*/
	}
	return nil

}

func kafkaShellConsume(c *cli.Context, kf *KafkaFlags) error {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Net.SASL.Enable = true
	config.Net.SASL.User = kf.User
	config.Net.SASL.Password = kf.Password

	consumer, err := cluster.NewConsumer(kf.Brokers, kf.Group, kf.Topic, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()
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
	go func() {
		for msg := range consumer.Messages() {
			fmt.Printf("%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
			consumer.MarkOffset(msg, "")
			return
		}
	}()

	time.Sleep(100 * time.Second)

	return nil
}

func kafkaShellProduce(c *cli.Context, kf *KafkaFlags) error {
	config := sarama.NewConfig()
	config.Net.SASL.Enable = true
	config.Net.SASL.User = kflags.User
	config.Net.SASL.Password = kflags.Pswd

	producer, err := sarama.NewAsyncProducer(kf.Brokers, config)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	msg := c.Args().Get(1)
	fmt.Println("msg:", msg)

	producer.Input() <- &sarama.ProducerMessage{
		Topic: kf.Topic[0],
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
