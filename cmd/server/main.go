package main

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/server/api"
	"log"
)

func main() {

	myApi := api.New()

	log.Fatalln(myApi.Run())

}
