package api

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/metric"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/handler"
	"net/http"
)

const (
	addr = "127.0.0.1:8080"
)

type api struct {
	Metrics *metric.Metrics
	serv    *http.Server
}

func New() *api {
	return &api{serv: &http.Server{Addr: addr}, Metrics: metric.New()}
}

func (a *api) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handler.UpdateHandler(a.Metrics))
	a.serv.Handler = mux
	return a.serv.ListenAndServe()
}
