package main

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/server/api"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/memory"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/service"
)

func main() {

	myData := memory.NewStorage()
	myService := service.NewService(myData)
	myAPI := api.NewAPI(myService)

	myAPI.Run()

}
