package api

//go:generate mockery --all

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type api struct {
	server  *http.Server
	service Service
	cfg     Config
}

// New returns new API.
func New(service Service, cfg Config) *api {
	log.Debug().Msg("api.New started")
	defer log.Debug().Msg("api.New ended")

	newAPI := &api{service: service, cfg: cfg}

	newAPI.server = &http.Server{
		Addr:    newAPI.cfg.ServAddr(),
		Handler: newAPI.newRouter(),
	}

	return newAPI
}

//Run API starts the API.
func (a *api) Run() {
	log.Debug().Msg("api.Run started")
	defer log.Debug().Msg("api.Run ended")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go a.startListener()

	<-exit
	a.service.ShutDown()
	os.Exit(0)
}

func (a *api) startListener() {
	log.Debug().Msg("api.StartListener started")
	defer log.Debug().Msg("api.StartListener ended")

	a.server.ListenAndServe()
	defer a.server.Close()
}

func (a *api) newRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(gzipReadHandle)
	r.Use(gzipWriteHandle)

	r.Get("/", a.handlerGetPage())
	r.Get("/ping", a.handlerPing())
	r.Post("/update/", a.handlerUpdateMetric)
	r.Post("/updates/", a.handlerUpdateListMetrics)
	r.Post("/value/", a.handlerGetMetricValue)
	r.Post("/update/{type}/{metric}/{val}", a.handlerUpdateMetricPathParams())
	r.Get("/value/{type}/{metric}", a.handlerGetMetricValuePathParams())

	return r
}
