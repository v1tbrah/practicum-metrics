package mockapi

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/v1tbrah/metricsAndAlerting/internal/agent/memory"
)

type MockAPI struct {
	Server *httptest.Server
	Data   memory.Data
}

func NewAPI(addr string, data memory.Data) *MockAPI {
	return &MockAPI{
		Server: newServer(addr),
		Data:   data,
	}
}

func newServer(addr string) *httptest.Server {
	newServer := httptest.NewUnstartedServer(router())
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	newServer.Listener = listener
	newServer.Start()
	return newServer
}

func router() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("application/json"))

	r.Post("/update/", reportMetricHandler)
	r.Post("/value/", getMetricHandler)

	return r
}

func reportMetricHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(body)
}

func getMetricHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(body)
}
