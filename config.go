package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var KlientToml string = ".klient.toml"
var KlientConfig Config

type Config struct {
	Mysql      CommonConfigSection `toml:"mysql"`
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

	f, _ := os.OpenFile(dst, os.O_CREATE|os.O_RDONLY, os.ModePerm|os.ModeTemporary)
	defer f.Close()
	_, err := toml.DecodeReader(f, &KlientConfig)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("%s: %+v \n", KlientToml, KlientConfig)
	}
	return
}

func SaveConfig() error {

	f, err := os.OpenFile(KlientToml, os.O_RDWR|os.O_CREATE, os.ModePerm|os.ModeTemporary)
	if err != nil {
		panic(err)
	}

	if err = toml.NewEncoder(f).Encode(KlientConfig); err != nil {
		panic(err)
	}

	return err
}
