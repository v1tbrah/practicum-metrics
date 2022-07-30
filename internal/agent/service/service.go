package service

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/v1tbrah/metricsAndAlerting/internal/agent/memory"
)

type service struct {
	client  *resty.Client
	options *options
	data    *memory.Data
}

// NewService returns service with default settings.
func NewService() *service {
	newService := service{
		client:  resty.New(),
		options: newOptions(),
		data:    memory.NewData()}

	newService.options.parseFromOsArgs()
	newService.options.parseFromEnv()

	newService.options.reportMetricURL = "http://" + newService.options.ServerAddr + "/update/"
	newService.options.getMetricURL = "http://" + newService.options.ServerAddr + "/update/"

	return &newService
}

//Run service updating the metrics once per pollInterval and sends them to the server once per reportInterval.
func (s *service) Run() {

	log.Println("Agent started.")
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	var mutex sync.Mutex

	go s.updateData(&mutex)
	go s.reportData(&mutex)

	<-shutdown
	log.Println("Agent exits normally.")
	os.Exit(0)

}

func (s *service) updateData(mutex *sync.Mutex) {
	ticker := time.NewTicker(s.options.PollInterval)
	for {
		<-ticker.C
		mutex.Lock()
		s.data.Update()
		log.Println("All metrics updated.")
		mutex.Unlock()
	}
}

func (s *service) reportData(mutex *sync.Mutex) {
	ticker := time.NewTicker(s.options.ReportInterval)
	for {
		<-ticker.C
		mutex.Lock()
		for _, metric := range *s.data {
			if _, err := s.reportMetric(metric); err != nil {
				log.Printf("Error report metric. Metric ID: %s. Reason: %s", metric.ID, err.Error())
			} else {
				log.Printf("Metric %s reported to server.", metric.ID)
			}
		}
		mutex.Unlock()
	}
}

func (s *service) reportMetric(metric memory.Metrics) (*resty.Response, error) {

	if metric.ID == "" {
		return nil, errors.New("metric ID is empty")
	}
	if !metric.TypeIsValid() {
		return nil, errors.New("invalid type of metric")
	}

	body, err := json.Marshal(metric)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.NewRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(s.options.reportMetricURL)

	return resp, err
}

func (s *service) getMetric(ID, MType string) (*resty.Response, error) {

	if ID == "" {
		return nil, errors.New("metric ID is empty")
	}
	metric := memory.NewMetric(ID, MType)
	if !metric.TypeIsValid() {
		return nil, errors.New("invalid type of metric")
	}

	body, err := json.Marshal(metric)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.NewRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(s.options.getMetricURL)

	return resp, err
}
