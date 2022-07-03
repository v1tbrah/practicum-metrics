package main

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/serve/handler"
	"net/http"
)

const (
	servAddr = "127.0.0.1"
	servPort = "8080"
)

func installHandlers(s *http.Server) {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/Alloc/", handler.Alloc)
	mux.HandleFunc("/update/BuckHashSys/", handler.BuckHashSys)
	mux.HandleFunc("/update/Frees/", handler.Frees)
	mux.HandleFunc("/update/GCCPUFraction/", handler.GCCPUFraction)
	mux.HandleFunc("/update/GCSys/", handler.GCSys)
	mux.HandleFunc("/update/HeapAlloc/", handler.HeapAlloc)
	mux.HandleFunc("/update/HeapIdle/", handler.HeapIdle)
	mux.HandleFunc("/update/HeapInuse/", handler.HeapInuse)
	mux.HandleFunc("/update/HeapObjects/", handler.HeapObjects)
	mux.HandleFunc("/update/HeapReleased/", handler.HeapReleased)
	mux.HandleFunc("/update/HeapSys/", handler.HeapSys)
	mux.HandleFunc("/update/LastGC/", handler.LastGC)
	mux.HandleFunc("/update/Lookups/", handler.Lookups)
	mux.HandleFunc("/update/MCacheInuse/", handler.MCacheInuse)
	mux.HandleFunc("/update/MCacheSys/", handler.MCacheSys)
	mux.HandleFunc("/update/MSpanInuse/", handler.MSpanInuse)
	mux.HandleFunc("/update/MSpanSys/", handler.MSpanSys)
	mux.HandleFunc("/update/Mallocs/", handler.Mallocs)
	mux.HandleFunc("/update/NextGC/", handler.NextGC)
	mux.HandleFunc("/update/NumForcedGC/", handler.NumForcedGC)
	mux.HandleFunc("/update/NumGC/", handler.NumGC)
	mux.HandleFunc("/update/OtherSys/", handler.OtherSys)
	mux.HandleFunc("/update/PauseTotalNs/", handler.PauseTotalNs)
	mux.HandleFunc("/update/StackInuse/", handler.StackInuse)
	mux.HandleFunc("/update/StackSys/", handler.StackSys)
	mux.HandleFunc("/update/Sys/", handler.Sys)
	mux.HandleFunc("/update/TotalAlloc/", handler.TotalAlloc)
	mux.HandleFunc("/update/PollCount/", handler.PollCount)
	mux.HandleFunc("/update/RandomValue/", handler.RandomValue)
	s.Handler = mux
}

func main() {

	go func() {
		s := &http.Server{
			Addr: servAddr + ":" + servPort,
		}
		installHandlers(s)
		s.ListenAndServe()
	}()

	for {
	}

}
