package main

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/server/api"
	"log"
)

func main() {

	myAPI := api.NewAPI()
	
	log.Fatalln(myAPI.Run())

}
