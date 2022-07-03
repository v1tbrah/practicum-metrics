package main

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/agent/send"
	"github.com/v1tbrah/metricsAndAlerting/internal/metric"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {

	m := metric.Metrics{}

	updateTime := time.NewTicker(pollInterval)
	reportTime := time.NewTicker(reportInterval)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		for {
			select {
			case <-updateTime.C:
				m.Update()
			case <-reportTime.C:
				send.AllMetrics(m)
			case <-shutdown:
				os.Exit(0)
			}
		}
	}()

	for {
	}

}
