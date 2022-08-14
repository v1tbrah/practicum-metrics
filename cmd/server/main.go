package main

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/server/api"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/service"
	"log"
)

func main() {

	myCfg := config.NewCfg(config.WithFlag, config.WithEnv)
	myStorage, err := repo.New(myCfg)
	if err != nil {
		log.Fatalln(err)
	}
	myService := service.NewService(myStorage, myCfg)
	myAPI := api.NewAPI(myService)

	myAPI.Run()

}
