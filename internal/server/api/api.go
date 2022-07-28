package api

import (
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type options struct {
	Addr          string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
)


type api struct {
	server  *http.Server
	service *service.Service
	options *options
}

// Creates the API.
func NewAPI(service *service.Service) *api {
	server := &http.Server{}
	options := newDefaultOptions()
	err := env.Parse(options)
	if err != nil {
		log.Println(err)
	}
	server.Addr = options.Addr

	newAPI := &api{
		server:  server,
		service: service,
		options: options}

	server.Handler = newAPI.newRouter()
	return newAPI
}

func newDefaultOptions() *options {
	return &options{Addr: "127.0.0.1:8080"}
}

//The API starts.
func (a *api) Run() {
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
	}
	os.Exit(0)
}

func (a *api) newRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", a.getPageHandler())
	r.Post("/update/", a.updateJSONHandler)
	r.Post("/value/", a.getValueJSONHandler)
	r.Post("/update/{type}/{metric}/{val}", checkTypeAndNameMetric("update", a.updateHandler()))
	r.Get("/value/{type}/{metric}", checkTypeAndNameMetric("value", a.getValueHandler()))

	return r
}
