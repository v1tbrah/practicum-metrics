package service

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
)

type options struct {
	ServerAddr      string        `env:"ADDRESS"`
	PollInterval    time.Duration `env:"POLL_INTERVAL"`
	ReportInterval  time.Duration `env:"REPORT_INTERVAL"`
	reportMetricURL string
	getMetricURL    string
}

func newOptions() *options {
	defaultAddr := "127.0.0.1:8080"
	opt := &options{
		ServerAddr:      defaultAddr,
		PollInterval:    2 * time.Second,
		ReportInterval:  10 * time.Second,
		reportMetricURL: "http://" + defaultAddr + "/update/",
		getMetricURL:    "http://" + defaultAddr + "/update/",
	}
	return opt
}

func (o *options) parseFromOsArgs() {
	if isTestMode := strings.Contains(os.Args[0], ".test"); isTestMode {
		return
	}

	flag.StringVar(&o.ServerAddr, "a", o.ServerAddr, "server address")
	flag.DurationVar(&o.PollInterval, "p", o.PollInterval, "interval for updating metrics")
	flag.DurationVar(&o.ReportInterval, "r", o.ReportInterval, "interval for report metrics to server")

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
