package agent

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
	"strings"
	"time"
)

type options struct {
	SrvAddr        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
}

func newDefaultOptions() *options {
	opt := &options{
		SrvAddr:        "127.0.0.1:8080",
		PollInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
	}
	return opt
}

func (o *options) parseFromOsArgs() {
	if strings.Contains(os.Args[0], ".test") {
		return
	}
	defaultOptions := newDefaultOptions()

	flag.StringVar(&o.SrvAddr, "a", defaultOptions.SrvAddr, "api server address")
	flag.DurationVar(&o.PollInterval, "p", defaultOptions.PollInterval, "interval for updating metrics from runtime")
	flag.DurationVar(&o.ReportInterval, "r", defaultOptions.ReportInterval, "interval for report metrics to server")

	if !flag.Parsed() {
		flag.Parse()
	}
}

func (o *options) parseFromEnv() {
	err := env.Parse(o)
	if err != nil {
		log.Println(err)
	}
}
