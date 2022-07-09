package main

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/agent/send"
	"github.com/v1tbrah/metricsAndAlerting/internal/metric"
	"log"
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

	currMetric := metric.New()

	go func() {
		updateTime := time.NewTicker(pollInterval)
		for {
			<-updateTime.C
			currMetric.Update()
		}
	}()

	go func() {
		reportTime := time.NewTicker(reportInterval)
		for {
			<-reportTime.C
			if err := send.AllMetrics(currMetric); err != nil {
				log.Fatalln(err)
			}
		}
	}()

	<-shutdown
	os.Exit(0)

}
