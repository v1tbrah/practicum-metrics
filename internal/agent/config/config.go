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
	ServerAddr           string        `env:"ADDRESS"`
	PollInterval         time.Duration `env:"POLL_INTERVAL"`
	ReportInterval       time.Duration `env:"REPORT_INTERVAL"`
	Key                  string        `env:"KEY"`
	ReportMetricURL      string
	ReportListMetricsURL string
	GetMetricURL         string
}

func NewCfg(args ...string) *Config {
	cfg := &Config{
		ServerAddr:     "127.0.0.1:8080",
		PollInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
	}

	for _, arg := range args {
		if arg == WithFlag {
			cfg.parseFromOsArgs()
		}
		if arg == WithEnv {
			cfg.parseFromEnv()
		}
	}

	cfg.ReportMetricURL = "http://" + cfg.ServerAddr + "/update/"
	cfg.ReportListMetricsURL = "http://" + cfg.ServerAddr + "/updates/"
	cfg.GetMetricURL = "http://" + cfg.ServerAddr + "/update/"

	return cfg
}

func (c *Config) parseFromOsArgs() {

	flag.StringVar(&c.ServerAddr, "a", c.ServerAddr, "server address")
	flag.DurationVar(&c.PollInterval, "p", c.PollInterval, "interval for updating metrics")
	flag.DurationVar(&c.ReportInterval, "r", c.ReportInterval, "interval for report metrics to server")
	flag.StringVar(&c.Key, "k", c.Key, "secret key for hash calculation")

	if !flag.Parsed() {
		flag.Parse()
	}
}

func (c *Config) parseFromEnv() {
	err := env.Parse(c)
	if err != nil {
		log.Println(err)
	}
}
