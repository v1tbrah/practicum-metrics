package api

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/service"
)

type api struct {
	server  *http.Server
	service *service.Service
	options *options
}

// NewAPI returns new API.
func NewAPI(service *service.Service) *api {
	newAPI := &api{
		service: service,
		options: newDefaultOptions()}

	newAPI.options.parseFromOsArgs()
	newAPI.options.parseFromEnv()

	newAPI.server = &http.Server{
		Addr:    newAPI.options.Addr,
		Handler: newAPI.newRouter(),
	}

	return newAPI
}

//Run API starts the API.
func (a *api) Run() {
	defer a.server.Close()
	log.Println("API started.")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	if a.options.Restore {
		if err := a.service.RestoreMetricsFromFile(a.options.StoreFile); err != nil {
			log.Println(err)
		} else {
			log.Println("Metrics restored from file:", a.options.StoreFile)
		}
	}

	if needWriteMetricsToFileWithInterval := a.options.StoreInterval != time.Second*0; needWriteMetricsToFileWithInterval {
		go a.service.WriteMetricsToFileWithInterval(a.options.StoreFile, a.options.StoreInterval)
	}

	go func() {
		log.Println(a.server.ListenAndServe())
	}()

	<-exit
	if err := a.service.SaveMetricsToFile(a.options.StoreFile); err != nil {
		log.Println(err)
	} else {
		log.Println("Data saved in file:", a.options.StoreFile)
	}
	log.Println("API exits normally.")
	os.Exit(0)
}

func (a *api) newRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(gzipReadHandle)
	r.Use(gzipWriteHandle)

	r.Get("/", a.getPageHandler())
	r.Post("/update/", a.updateMetricHandler)
	r.Post("/value/", a.getMetricValueHandler)
	r.Post("/update/{type}/{metric}/{val}", a.updateMetricHandlerPathParams())
	r.Get("/value/{type}/{metric}", a.getMetricValueHandlerPathParams())

	return r
}
