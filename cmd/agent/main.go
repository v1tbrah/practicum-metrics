package main

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/agent/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/agent/service"
)

func main() {

	myCfg := config.NewCfg(config.WithFlag, config.WithEnv)
	myAgent := service.NewService(myCfg)

	myAgent.Run()

}
