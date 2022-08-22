package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	WithFlag = "withFlag"
	WithEnv  = "withEnv"
)

const (
	StorageTypeMemory = iota
	StorageTypeDB
)

type Config struct {
	Addr string `env:"ADDRESS"`

	StorageType int

	PgConnString string `env:"DATABASE_DSN"`

	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`

	HashKey string `env:"KEY"`
}

func New(args ...string) (*Config, error) {
	cfg := &Config{
		Addr:          "127.0.0.1:8080",
		StoreInterval: time.Second * 300,
		Restore:       true}

	for _, arg := range args {
		switch arg {
		case WithFlag:
			cfg.parseFromOsArgs()
		case WithEnv:
			if err := cfg.parseFromEnv(); err != nil {
				return nil, err
			}
		}
	}

	if haveDBConnection := cfg.PgConnString != ""; haveDBConnection {
		cfg.StorageType = StorageTypeDB
	} else {
		cfg.StorageType = StorageTypeMemory
	}

	return cfg, nil
}

func (c *Config) parseFromOsArgs() {

	flag.StringVar(&c.Addr, "a", c.Addr, "api server address")
	flag.DurationVar(&c.StoreInterval, "i", c.StoreInterval, "interval for writing metrics to a file")
	flag.StringVar(&c.StoreFile, "f", c.StoreFile, "path to persistent file storage")
	flag.BoolVar(&c.Restore, "r", c.Restore, "flag for loading metrics from a file at the start of the api")
	flag.StringVar(&c.HashKey, "k", c.HashKey, "secret key for hash calculation")
	flag.StringVar(&c.PgConnString, "d", c.PgConnString, "postgres db conn string")

	flag.Parse()
}

func (c *Config) parseFromEnv() error {
	return env.Parse(c)
}
