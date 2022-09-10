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

	flag.StringVar(&c.serverAddr, "a", c.serverAddr, "server address")
	flag.DurationVar(&c.pollInterval, "p", c.pollInterval, "interval for updating metrics")
	flag.DurationVar(&c.reportInterval, "r", c.reportInterval, "interval for report metrics to server")
	flag.StringVar(&c.hashKey, "k", c.hashKey, "secret key for hash calculation")
	flag.StringVar(&c.logLevel, "l", c.logLevel, "log level")

	if !flag.Parsed() {
		flag.Parse()
	}
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
		ServerAddr     string        `env:"ADDRESS"`
		PollInterval   time.Duration `env:"POLL_INTERVAL"`
		ReportInterval time.Duration `env:"REPORT_INTERVAL"`
		HashKey        string        `env:"KEY"`
		LogLevel       string        `env:"LOGLEVEL"`
	}{}

	if err = env.Parse(c); err != nil {
		return err
	}

	if configEnv.ServerAddr != "" {
		c.serverAddr = configEnv.ServerAddr
	}
	if configEnv.PollInterval != time.Second*0 {
		c.pollInterval = configEnv.PollInterval
	}
	if configEnv.ReportInterval != time.Second*0 {
		c.reportInterval = configEnv.ReportInterval
	}
	if configEnv.ReportInterval != time.Second*0 {
		c.reportInterval = configEnv.ReportInterval
	}
	if configEnv.HashKey != "" {
		c.hashKey = configEnv.HashKey
	}
	if configEnv.LogLevel != "" {
		c.logLevel = configEnv.LogLevel
	}

	return nil
}
