package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/v1tbrah/metricsAndAlerting/internal/agent/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/agent/memory"
	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

type service struct {
	client *resty.Client
	cfg    *config.Config
	data   *memory.Data
}

// NewService returns service.
func NewService(cfg *config.Config) *service {
	newService := service{
		client: resty.New(),
		cfg:    cfg,
		data:   memory.NewData()}

	return &newService
}

//Run service updating the metrics once per pollInterval and sends them to the server once per reportInterval.
func (s *service) Run() {

	log.Println("Agent started.")
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	rand.Seed(time.Now().UnixNano())
	go s.updateData()
	go s.reportData()

	<-shutdown
	log.Println("Agent exits.")
	os.Exit(0)

}

func (s *service) updateData() {
	ticker := time.NewTicker(s.cfg.PollInterval)
	for {
		<-ticker.C
		s.data.Update(s.cfg.Key)
	}
}

func (s *service) reportData() {
	ticker := time.NewTicker(s.cfg.ReportInterval)
	for {
		<-ticker.C
		s.data.Lock()
		for _, currMetric := range s.data.Metrics {
			if _, err := s.reportMetric(currMetric); err != nil {
				log.Printf("Error report metric. Metric ID: %s. Reason: %s", currMetric.ID, err.Error())
			}
		}
		s.data.Unlock()
		log.Println("All metrics reported.")
	}
}

func (s *service) reportMetric(metricForReport metric.Metrics) (*resty.Response, error) {

	if metricForReport.ID == "" {
		return nil, errors.New("metric ID is empty")
	}
	if !metricForReport.TypeIsValid() {
		return nil, errors.New("invalid type of metric")
	}

	body, err := json.Marshal(metricForReport)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.NewRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(s.cfg.ReportMetricURL)

	return resp, err
}

func (s *service) reportListMetrics(listMetrics []metric.Metrics) (*resty.Response, error) {

	if len(listMetrics) == 0 {
		return nil, errors.New("list metrics is empty")
	}
	for i, curr := range listMetrics {
		if curr.ID == "" {
			return nil, fmt.Errorf("metric %d: id is empty", i)
		}
		if !curr.TypeIsValid() {
			return nil, fmt.Errorf("metric %d: invalid type of metric", i)
		}
	}

	body, err := json.Marshal(listMetrics)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.NewRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(s.cfg.ReportListMetricsURL)

	return resp, err
}

func (s *service) getMetric(ID, MType string) (*resty.Response, error) {

	if ID == "" {
		return nil, errors.New("metric ID is empty")
	}
	metricForRequest := metric.NewMetric(ID, MType)
	if !metricForRequest.TypeIsValid() {
		return nil, errors.New("invalid type of metric")
	}

	body, err := json.Marshal(metricForRequest)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.NewRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(s.cfg.GetMetricURL)

	return resp, err
}
