package main

import (
	"log"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/api"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/memory"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/service"
)

func main() {
	myData := memory.NewMemStorage()
	myService := service.NewService(myData)
	myAPI := api.NewAPI(myService)

	log.Fatalln(myAPI.Run())
}
