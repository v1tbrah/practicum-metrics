package api

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/memory"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/pg"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/service"
)

type api struct {
	server  *http.Server
	service *service.Service
}

// NewAPI returns new API.
func NewAPI(service *service.Service) *api {
	newAPI := &api{service: service}

	newAPI.server = &http.Server{
		Addr:    newAPI.service.Cfg.Addr,
		Handler: newAPI.newRouter(),
	}

	return newAPI
}

//Run API starts the API.
func (a *api) Run() {

	log.Println("API started.")
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		log.Println(a.server.ListenAndServe())
		defer a.server.Close()
	}()

	<-exit
	if a.service.Cfg.StorageType == config.InMemory {
		inMemStorage, _ := a.service.Storage.(*memory.MemStorage)
		if err := inMemStorage.StoreData(); err != nil {
			log.Println(err)
		} else {
			log.Println("Data saved in file.")
		}
	} else if a.service.Cfg.StorageType == config.InDB {
		inDBStorage, _ := a.service.Storage.(*pg.PgStorage)
		inDBStorage.CloseConnection()
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
	r.Get("/ping", a.checkDBConnHandler())
	r.Post("/update/", a.updateMetricHandler)
	r.Post("/updates/", a.updateListMetricsHandler)
	r.Post("/value/", a.getMetricValueHandler)
	r.Post("/update/{type}/{metric}/{val}", a.updateMetricHandlerPathParams())
	r.Get("/value/{type}/{metric}", a.getMetricValueHandlerPathParams())

	return r
}
