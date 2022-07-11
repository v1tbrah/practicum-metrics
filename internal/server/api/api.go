package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/api/metric"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/handler"
	"net/http"
)

const (
	addr = "127.0.0.1:8080"
)

type api struct {
	serv    *http.Server
	metrics *metric.Metrics
}

func NewAPI() *api {
	metrics := metric.NewMetrics()
	return &api{
		serv: &http.Server{
			Addr:    addr,
			Handler: NewRouter(metrics)},
		metrics: metrics}
}

func NewRouter(m *metric.Metrics) chi.Router {
	r := chi.NewRouter()
	r.Get("/", handler.GetAllMetricsHTML(m))
	r.Post("/update/{type}/{metric}/{val}", handler.UpdateHandler(m))
	r.Get("/value/{type}/{metric}", handler.GetValueHandler(m))
	return r
}

func (a *api) Run() error {
	return a.serv.ListenAndServe()
}
