package main

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/server/api"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/memory"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/service"
	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

func main() {

	thisMetric := metric.Metrics{MType: "gauge"}
	thisMetric.UpdateHash("fsdfds")

	myData := memory.NewStorage()
	myCfg := config.NewCfg(config.WithFlag, config.WithEnv)
	myService := service.NewService(myData, myCfg)
	myAPI := api.NewAPI(myService)

	myAPI.Run()

}
