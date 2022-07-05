package main

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/agent/send"
	"github.com/v1tbrah/metricsAndAlerting/internal/metric"
	"os"
	"os/signal"
	"sync"
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
	var mutex sync.Mutex

	go func(mutex *sync.Mutex) {
		updateTime := time.NewTicker(pollInterval)
		for {
			<-updateTime.C
			mutex.Lock()
			currMetric.Update()
			mutex.Unlock()
		}
	}(&mutex)

	go func(mutex *sync.Mutex) {
		reportTime := time.NewTicker(reportInterval)
		for {
			<-reportTime.C
			mutex.Lock()
			send.AllMetrics(currMetric)
			mutex.Unlock()
		}
	}(&mutex)

	<-shutdown
	os.Exit(0)

}
