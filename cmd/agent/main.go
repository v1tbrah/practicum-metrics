package main

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/agent"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {

	myAgent := agent.NewAgent()

	myAgent.Run()

}
