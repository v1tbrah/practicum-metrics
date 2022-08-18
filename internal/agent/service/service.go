package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/internal/agent/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/agent/memory"
	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

var (
	ErrMetricIDIsEmpty    = errors.New("metric id is empty")
	ErrInvalidMetricType  = errors.New("invalid type of metric")
	ErrListMetricsIsEmpty = errors.New("list metrics is empty")
)

type service struct {
	client *resty.Client
	cfg    *config.Config
	data   *memory.Data
}

// New returns service.
func New(cfg *config.Config) (*service, error) {
	log.Debug().Str("cfg", fmt.Sprint(cfg)).Msg("service.New started")
	defer log.Debug().Msg("service.New ended")

	if cfg == nil {
		return nil, errors.New("config is empty")
	}

	newService := service{
		client: resty.New(),
		cfg:    cfg,
		data:   memory.NewData()}

	return &newService, nil
}

//Run service updating the metrics once per pollInterval and sends them to the server once per reportInterval.
func (s *service) Run() {
	log.Printf("service.Run started")
	defer log.Printf("service.Run ended")

	log.Info().Msg("agent starting")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	rand.Seed(time.Now().UnixNano())
	go s.updateDataWithInterval(s.cfg.PollInterval)
	go s.reportDataWithInterval(s.cfg.ReportInterval)

	<-shutdown
	log.Info().Msg("agent ended")
	os.Exit(0)

}

func (s *service) updateDataWithInterval(interval time.Duration) {
	log.Debug().Dur("interval", interval).Msg("service.updateDataWithInterval started")
	defer log.Debug().Msg("service.updateDataWithInterval ended")

	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		s.data.Update(s.cfg.HashKey)
	}
}

func (s *service) reportDataWithInterval(interval time.Duration) {
	log.Debug().Dur("interval", interval).Msg("service.reportDataWithInterval started")
	defer log.Printf("service.reportDataWithInterval ended")

	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		// TODO должна ли перед этим range быть блокировка на чтение?
		for _, currMetric := range s.data.Metrics {
			if _, err := s.reportMetric(currMetric); err != nil {
				log.Error().
					Err(err).
					Str("MID", currMetric.ID).
					Str("MType", currMetric.MType).
					Str("deltaPtr", fmt.Sprint(currMetric.Delta)).
					Str("valuePtr", fmt.Sprint(currMetric.Value)).
					Msg("unable reporting metric")
			} else {
				log.Info().Str("MID", currMetric.ID).Msg("metric reported")
			}
		}

	}
}

func (s *service) reportMetric(metricForReport metric.Metrics) (*resty.Response, error) {
	log.Debug().
		Str("MID", metricForReport.ID).
		Str("MType", metricForReport.MType).
		Str("deltaPtr", fmt.Sprint(metricForReport.Delta)).
		Str("valuePtr", fmt.Sprint(metricForReport.Value)).
		Msg("service.reportMetric started")
	defer log.Printf("service.reportMetric ended")

	if metricForReport.ID == "" {
		return nil, ErrMetricIDIsEmpty
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
	log.Debug().
		Str("listMetricsPtr", fmt.Sprint(listMetrics)).
		Msg("service.reportListMetrics started")
	defer log.Printf("service.reportListMetrics ended")

	if listMetrics == nil {
		return nil, errors.New("list metrics is nil ptr")
	}
	if len(listMetrics) == 0 {
		return nil, ErrListMetricsIsEmpty
	}
	for i, curr := range listMetrics {
		if curr.ID == "" {
			return nil, fmt.Errorf("metric with index %d of list: %w", i, ErrMetricIDIsEmpty)
		}
		if !curr.TypeIsValid() {
			return nil, fmt.Errorf("metric with index %d of list: %w", i, ErrInvalidMetricType)
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
	log.Debug().
		Str("ID", ID).
		Str("MType", MType).
		Msg("service.getMetric started")
	defer log.Printf("service.getMetric ended")

	if ID == "" {
		return nil, ErrMetricIDIsEmpty
	}
	metricForRequest := metric.NewMetric(ID, MType)
	if !metricForRequest.TypeIsValid() {
		return nil, ErrInvalidMetricType
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
