package main

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/server/api"
	"log"
)

func main() {

	myAPI := api.New()

	log.Fatalln(myAPI.Run())

}
