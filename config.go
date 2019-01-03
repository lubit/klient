package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var KlientToml string = ".klient.toml"
var KlientConfig Config
var KafkaDefaultGroup = "test-consumer-group"

type Config struct {
	Mysql      CommonConfigSection `toml:"mysql"`
	Pgsql      CommonConfigSection `toml:"pgsql"`
	Redis      CommonConfigSection `toml:"redis"`
	Clickhouse CommonConfigSection `toml:"clickhouse"`
	Kafka      KafkaConfigSection  `toml:"kafka"`
}

type CommonConfigSection struct {
	Host string
	Port string
	User string
	Pswd string
	Db   string
}

type KafkaConfigSection struct {
	Brokers []string
	Topics  []string
	Group   string
	User    string
	Pswd    string
}

func LoadConfig() {
	user, _ := user.Current()
	dst := filepath.Join(user.HomeDir, KlientToml)
	/*
		f, _ := os.OpenFile(dst, os.O_CREATE|os.O_RDONLY, os.ModePerm|os.ModeTemporary)
		defer f.Close()
		_, err := toml.DecodeReader(f, &KlientConfig)
		if err != nil {
			panic(k_red(err))
		} else {
			//DEBUG.Printf("%s: %+v \n", dst, KlientConfig)
		}
	*/
	if _, err := toml.DecodeFile(dst, &KlientConfig); err != nil {
		panic(k_red(err))
	}
	return
}

func SaveConfig() error {

	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(KlientConfig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buf.String())

	user, _ := user.Current()
	dst := filepath.Join(user.HomeDir, KlientToml)

	err = ioutil.WriteFile(dst, buf.Bytes(), 0644)
	if err != nil {
		panic(k_red(err))
	}

	/*
		f, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, os.ModePerm|os.ModeTemporary)
		if err != nil {
			panic(k_red(err))

		}
		defer f.Close()

		fmt.Println(KlientConfig)

		if err = toml.NewEncoder(f).Encode(KlientConfig); err != nil {
			panic(k_red(err))
		}
	*/

	return err
}
