package main

import (
	"fmt"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/api"
	"log"
	"time"
)

func main() {

	myAPI := api.NewAPI()

	// DEBUG
	go func() {
		tick := time.NewTicker(time.Second * 10)
		for {
			<-tick.C
			fmt.Println(myAPI.Metrics())
		}
	}()
	// DEBUG

	log.Fatalln(myAPI.Run())

}
