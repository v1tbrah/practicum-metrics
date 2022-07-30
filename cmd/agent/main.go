package main

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/agent/service"
)

func main() {

	myAgent := service.NewService()

	myAgent.Run()

}
