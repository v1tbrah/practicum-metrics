package api

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
	"strings"
	"time"
)

type options struct {
	Addr          string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func newDefaultOptions() *options {
	return &options{
		Addr:          "127.0.0.1:8080",
		StoreInterval: time.Second * 300,
		Restore:       true}
}

func (o *options) parseFromOsArgs() {
	if strings.Contains(os.Args[0], ".test") {
		return
	}
	defaultOptions := newDefaultOptions()

	flag.StringVar(&o.Addr, "a", defaultOptions.Addr, "api server address")
	flag.DurationVar(&o.StoreInterval, "i", defaultOptions.StoreInterval, "interval for writing metrics to a file")
	flag.StringVar(&o.StoreFile, "f", defaultOptions.StoreFile, "path to persistent file storage")
	flag.BoolVar(&o.Restore, "r", defaultOptions.Restore, "flag for loading metrics from a file at the start of the api")

	flag.Parse()
}

func (o *options) parseFromEnv() {
	err := env.Parse(o)
	if err != nil {
		log.Println(err)
	}
}
