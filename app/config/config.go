package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

const appName = "credit_bot"

type Config struct {
	ApiToken  string `required:"true"`
	Port      string `default:"80"`
	Address   string `required:"true"`
	Debug     bool   `default:"true"`
	Tls       bool   `default:"false"`
	RedisHost string `default:"localhost"`
	RedisPort string
	RedisDb   int `default:"0"`
}

func NewConfig() *Config {

	var c Config
	err := envconfig.Process(appName, &c)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &c
}
