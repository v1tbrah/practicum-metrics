package config

import (
	"time"

	"github.com/rs/zerolog/log"
)

const (
	WithFlag = "withFlag"
	WithEnv  = "withEnv"
)

type config struct {
	serverAddr string `env:"ADDRESS"`

	pollInterval   time.Duration `env:"POLL_INTERVAL"`
	reportInterval time.Duration `env:"REPORT_INTERVAL"`

	hashKey string `env:"KEY"`

	reportMetricURL      string
	reportListMetricsURL string
	getMetricURL         string

	logLevel string `env:"LOGLEVEL"`
}

func New(args ...string) (*config, error) {
	log.Debug().Strs("args", args).Msg("config.New started")
	cfg := &config{}
	var err error
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("config.New ended")
		} else {
			log.Debug().Str("config", cfg.String()).Msg("config.New ended")
		}
	}()

	cfg.serverAddr = "127.0.0.1:8080"
	cfg.pollInterval = 2 * time.Second
	cfg.reportInterval = 10 * time.Second
	cfg.logLevel = "info"

	for _, arg := range args {
		switch arg {
		case WithFlag:
			cfg.parseFromOsArgs()
		case WithEnv:
			if err = cfg.parseFromEnv(); err != nil {
				return nil, err
			}
		}
	}

	cfg.reportMetricURL = "http://" + cfg.serverAddr + "/update/"
	cfg.reportListMetricsURL = "http://" + cfg.serverAddr + "/updates/"
	cfg.getMetricURL = "http://" + cfg.serverAddr + "/update/"

	return cfg, nil
}

func (c *config) PollInterval() time.Duration {
	return c.pollInterval
}

func (c *config) ReportInterval() time.Duration {
	return c.reportInterval
}

func (c *config) HashKey() string {
	return c.hashKey
}

func (c *config) ReportMetricURL() string {
	return c.reportMetricURL
}

func (c *config) ReportListMetricsURL() string {
	return c.reportListMetricsURL
}

func (c *config) GetMetricURL() string {
	return c.getMetricURL
}

func (c *config) LogLevel() string {
	return c.logLevel
}

func (c *config) String() string {
	return "ServerAddr: " + c.serverAddr +
		", PollInterval: " + c.pollInterval.String() +
		", ReportInterval: " + c.reportInterval.String() +
		", ReportMetricURL: " + c.reportMetricURL +
		", ReportListMetricsURL: " + c.reportListMetricsURL +
		", GetMetricURL: " + c.getMetricURL +
		", LogLevel: " + c.logLevel
}
