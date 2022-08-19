package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

const (
	WithFlag = "withFlag"
	WithEnv  = "withEnv"
)

type Config struct {
	ServerAddr string `env:"ADDRESS"`

	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`

	HashKey string `env:"KEY"`

	ReportMetricURL      string
	ReportListMetricsURL string
	GetMetricURL         string
}

func New(args ...string) (*Config, error) {
	log.Debug().Msg("config.New started")
	defer log.Debug().Msg("config.New ended")

	cfg := &Config{
		ServerAddr:     "127.0.0.1:8080",
		PollInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
	}

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

	cfg.ReportMetricURL = "http://" + cfg.ServerAddr + "/update/"
	cfg.ReportListMetricsURL = "http://" + cfg.ServerAddr + "/updates/"
	cfg.GetMetricURL = "http://" + cfg.ServerAddr + "/update/"

	return cfg, nil
}

func (c *Config) parseFromOsArgs() {
	log.Debug().Msg("config.parseFromOsArgs started")
	defer log.Debug().Msg("config.parseFromOsArgs ended")

	flag.StringVar(&c.ServerAddr, "a", c.ServerAddr, "server address")
	flag.DurationVar(&c.PollInterval, "p", c.PollInterval, "interval for updating metrics")
	flag.DurationVar(&c.ReportInterval, "r", c.ReportInterval, "interval for report metrics to server")
	flag.StringVar(&c.HashKey, "k", c.HashKey, "secret key for hash calculation")

	if !flag.Parsed() {
		flag.Parse()
	}
}

func (c *Config) parseFromEnv() error {
	log.Debug().Msg("config.parseFromEnv started")
	defer log.Debug().Msg("config.parseFromEnv ended")

	return env.Parse(c)
}
