package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/memory"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/pg"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/service"
)

type api struct {
	server  *http.Server
	service *service.Service
}

// New returns new API.
func New(service *service.Service) *api {
	log.Debug().Msg("api.New started")
	log.Debug().Msg("api.New ended")

	newAPI := &api{service: service}

	newAPI.server = &http.Server{
		Addr:    newAPI.service.Cfg.Addr,
		Handler: newAPI.newRouter(),
	}

	return newAPI
}

//Run API starts the API.
func (a *api) Run() {
	log.Debug().Msg("api.Run started")
	log.Debug().Msg("api.Run ended")

	log.Info().Msg("api started")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go a.StartListener()

	<-exit
	if a.service.Cfg.StorageType == config.StorageTypeMemory {
		memStorage, _ := a.service.Storage.(*memory.Memory)
		if err := memStorage.StoreData(context.Background()); err != nil {
			log.Error().
				Err(err).
				Str("storeFile", a.service.Cfg.StoreFile).
				Msg("unable to store data in file")
		} else {
			log.Info().Msg(fmt.Sprintf("data saved to file: %s", a.service.Cfg.StoreFile))
		}
	} else if a.service.Cfg.StorageType == config.StorageTypeDB {
		DBStorage, _ := a.service.Storage.(*pg.Pg)
		DBStorage.ClosePoolConn()
	}
	log.Info().Msg("api ended")
	os.Exit(0)
}

func (a *api) StartListener() {
	log.Debug().Msg("api.StartListener started")
	log.Debug().Msg("api.StartListener ended")

	a.server.ListenAndServe()
	defer a.server.Close()
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
