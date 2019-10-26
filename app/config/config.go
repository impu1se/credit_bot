package config

import (
	"flag"
	"log"
	"os"
)

type Config struct {
	ApiToken string
	Port     string
	Address  string
	Debug    bool
	Tls      bool
}

func NewConfig() Config {

	var apiToken, port, addr string

	flag.StringVar(&apiToken, "token", "", "Telegram Bot Token")
	flag.StringVar(&port, "port", "80", "Port for server")
	flag.StringVar(&addr, "addr", "localhost", "Address for server")
	debug := flag.Bool("debug", false, "Debug true/false")
	tls := flag.Bool("tls", false, "TLS true/false")
	flag.Parse()

	if apiToken == "" {
		log.Print("-token is required")
		os.Exit(1)
	}

	return Config{
		ApiToken: apiToken,
		Port:     port,
		Address:  addr,
		Debug:    *debug,
		Tls:      *tls,
	}
}
