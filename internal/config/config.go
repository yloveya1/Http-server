package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	Port string `yaml:"port"`
	Db   struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Dbname   string `yaml:"dbname"`
		Sslmode  string `yaml:"sslmode"`
	}
	NatStreaming struct {
		ClusterId string `yaml:"clusterId"`
		ClientId  string `yaml:"clientId"`
		NatsURL   string `yaml:"natsURL"`
	}
}

var Conf *Config
var once sync.Once

func GetConf() *Config {
	once.Do(func() {
		Conf = &Config{}
		if err := cleanenv.ReadConfig("config.yml", Conf); err != nil {
			a, _ := cleanenv.GetDescription(Conf, nil)
			fmt.Println(a)
		}
	})
	return Conf
}
