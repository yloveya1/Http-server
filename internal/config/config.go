package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	Port string       `yaml:"port"`
	Nats NatStreaming `yaml:"natStreaming"`
	Db   Db           `yaml:"db"`
}

type NatStreaming struct {
	ClusterId string `yaml:"clusterId"`
	ClientId  string `yaml:"clientId"`
	NatsURL   string `yaml:"natsURL"`
}

type Db struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Dbname   string `yaml:"dbname"`
	Sslmode  string `yaml:"sslmode"`
}

var Conf *Config
var once sync.Once

func GetConf() *Config {
	once.Do(func() {
		Conf = &Config{}
		if err := cleanenv.ReadConfig("config.yml", Conf); err != nil {
			log.Fatalln(err)
		}
	})
	return Conf
}
