package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

func (c *config) parseFromOsArgs() {
	log.Debug().Msg("config.parseFromOsArgs started")
	defer log.Debug().Msg("config.parseFromOsArgs ended")

	flag.StringVar(&c.servAddr, "a", c.servAddr, "api server address")
	flag.DurationVar(&c.storeInterval, "i", c.storeInterval, "interval for writing metrics to a file")
	flag.StringVar(&c.storeFile, "f", c.storeFile, "path to persistent file storage")
	flag.BoolVar(&c.restore, "r", c.restore, "flag for loading metrics from a file at the start of the api")
	flag.StringVar(&c.hashKey, "k", c.hashKey, "secret key for hash calculation")
	flag.StringVar(&c.pgConnString, "d", c.pgConnString, "postgres db conn string")
	flag.StringVar(&c.logLevel, "l", c.logLevel, "log level")

	flag.Parse()
}

func (c *config) parseFromEnv() error {
	log.Debug().Msg("config.parseFromEnv started")
	var err error
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("config.parseFromEnv ended")
		} else {
			log.Debug().Msg("config.parseFromEnv ended")
		}
	}()

	configEnv := struct {
		ServAddr      string        `env:"ADDRESS"`
		PgConnString  string        `env:"DATABASE_DSN"`
		StoreInterval time.Duration `env:"STORE_INTERVAL"`
		StoreFile     string        `env:"STORE_FILE"`
		Restore       bool          `env:"RESTORE"`
		HashKey       string        `env:"KEY"`
		LogLevel      string        `env:"LOGLEVEL"`
	}{}

	if err = env.Parse(&configEnv); err != nil {
		return err
	}

	if configEnv.ServAddr != "" {
		c.servAddr = configEnv.ServAddr
	}
	if configEnv.PgConnString != "" {
		c.pgConnString = configEnv.PgConnString
	}
	if configEnv.StoreInterval != time.Second*0 {
		c.storeInterval = configEnv.StoreInterval
	}
	if configEnv.StoreFile != "" {
		c.storeFile = configEnv.StoreFile
	}
	if configEnv.Restore {
		c.restore = true
	}
	if configEnv.HashKey != "" {
		c.hashKey = configEnv.HashKey
	}
	if configEnv.LogLevel != "" {
		c.logLevel = configEnv.LogLevel
	}

	return nil
}
