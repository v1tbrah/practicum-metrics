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

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	currMetric := metric.Metrics{}
	m := make(chan metric.Metrics)

	go func() {
		for {
			updateTime := time.NewTicker(pollInterval)
			<-updateTime.C
			currMetric.Update()
			m <- currMetric
		}
	}()

	go func() {
		for {
			reportTime := time.NewTicker(reportInterval)
			<-reportTime.C
			send.AllMetrics(<-m)
		}
	}()

	<-shutdown
	os.Exit(0)

}
