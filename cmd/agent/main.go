package main

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/agent"
)

func main() {
	myAgent := agent.NewAgent()

	myAgent.Run()
}
