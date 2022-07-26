package api

import (
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/service"
	"net/http"
)

type options struct {
	Addr string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
}

type api struct {
	server  *http.Server
	service *service.Service
}

// Creates the API.
func NewAPI(service *service.Service) *api {
	server := &http.Server{}
	options := newDefaultOptions()
	env.Parse(options)
	server.Addr = options.Addr

	newAPI := &api{
		server:  server,
		service: service}

	server.Handler = newAPI.newRouter()
	return newAPI
}

func newDefaultOptions() *options {
	return &options{Addr: "127.0.0.1:8080"}
}

//The API starts.
func (a *api) Run() error {
	return a.server.ListenAndServe()
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
