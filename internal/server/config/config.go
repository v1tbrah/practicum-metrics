package config

import (
	"flag"

	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	WithFlag = "withFlag"
	WithEnv  = "withEnv"
)

type Config struct {
	Addr          string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	Key           string        `env:"KEY"`
}

func NewCfg(args ...string) *Config {
	cfg := &Config{
		Addr:          "127.0.0.1:8080",
		StoreInterval: time.Second * 300,
		Restore:       true}

	for _, arg := range args {
		if arg == WithFlag {
			cfg.parseFromOsArgs()
		}
		if arg == WithEnv {
			cfg.parseFromEnv()
		}
	}

	return cfg
}

func (c *Config) parseFromOsArgs() {

	flag.StringVar(&c.Addr, "a", c.Addr, "api server address")
	flag.DurationVar(&c.StoreInterval, "i", c.StoreInterval, "interval for writing metrics to a file")
	flag.StringVar(&c.StoreFile, "f", c.StoreFile, "path to persistent file storage")
	flag.BoolVar(&c.Restore, "r", c.Restore, "flag for loading metrics from a file at the start of the api")
	flag.StringVar(&c.Key, "k", c.Key, "secret key for hash calculation")

	flag.Parse()
}

func (c *Config) parseFromEnv() {
	err := env.Parse(c)
	if err != nil {
		log.Println(err)
	}
}
