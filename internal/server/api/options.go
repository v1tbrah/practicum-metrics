package api

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
)

type options struct {
	Addr          string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
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

	flag.StringVar(&o.Addr, "a", o.Addr, "api server address")
	flag.DurationVar(&o.StoreInterval, "i", o.StoreInterval, "interval for writing metrics to a file")
	flag.StringVar(&o.StoreFile, "f", o.StoreFile, "path to persistent file storage")
	flag.BoolVar(&o.Restore, "r", o.Restore, "flag for loading metrics from a file at the start of the api")

	flag.Parse()
}

func (o *options) parseFromEnv() {
	err := env.Parse(o)
	if err != nil {
		log.Println(err)
	}
}
