package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/api"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/memory"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/service"
)

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	myData := memory.NewMemStorage()
	myService := service.NewService(myData)
	myAPI := api.NewAPI(myService)

	myAPI.Run()

}
