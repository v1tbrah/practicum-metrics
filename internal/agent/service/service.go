package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

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
	data   *memory.Memory
}

// New returns service.
func New(cfg *config.Config) (*service, error) {
	log.Debug().Str("cfg", cfg.String()).Msg("service.New started")
	defer log.Debug().Msg("service.New ended")

	if cfg == nil {
		return nil, errors.New("config is empty")
	}

	newService := service{
		client: resty.New(),
		cfg:    cfg,
		data:   memory.New()}

	return &newService, nil
}

//Run service updating the metrics once per pollInterval and sends them to the server once per reportInterval.
func (s *service) Run() {
	log.Printf("service.Run started")
	defer log.Printf("service.Run ended")

	log.Info().Msg("agent starting")
	rand.Seed(time.Now().UnixNano())

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return s.updateBasicDataWithInterval(s.cfg.PollInterval, shutdown)
	})

	g.Go(func() error {
		return s.updateAdditionalDataWithInterval(s.cfg.PollInterval, shutdown)
	})

	g.Go(func() error {
		return s.reportDataWithInterval(s.cfg.ReportInterval, shutdown)
	})

	if err := g.Wait(); err != nil {
		log.Error().Err(err).Msg("agent ended")
		return
	}

	log.Info().Msg("agent ended")
	os.Exit(0)

}

func (s *service) updateBasicDataWithInterval(interval time.Duration, shutdown chan os.Signal) error {
	log.Debug().Dur("interval", interval).Msg("service.updateBasicDataWithInterval started")
	defer log.Debug().Msg("service.updateBasicDataWithInterval ended")

	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			s.data.UpdateBasic(s.cfg.HashKey)
		case _, ok := <-shutdown:
			if ok {
				close(shutdown)
			}
			return nil
		}
	}
}

func (s *service) updateAdditionalDataWithInterval(interval time.Duration, shutdown chan os.Signal) error {
	log.Debug().Dur("interval", interval).Msg("service.updateAdditionalDataWithInterval started")
	defer log.Debug().Msg("service.updateAdditionalDataWithInterval ended")

	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			if err := s.data.UpdateAdditional(s.cfg.HashKey); err != nil {
				return err
			}
		case _, ok := <-shutdown:
			if ok {
				close(shutdown)
			}
			return nil
		}
	}
}

func (s *service) reportDataWithInterval(interval time.Duration, shutdown chan os.Signal) error {
	log.Debug().Dur("interval", interval).Msg("service.reportDataWithInterval started")
	defer log.Debug().Msg("service.reportDataWithInterval ended")

	var urlError *url.Error
	tried := 0
	ticker := time.NewTicker(interval)

loop:
	for {
		select {
		case <-ticker.C:
			for _, currMetric := range s.data.GetData() {
				if _, err := s.reportMetric(currMetric); err != nil {
					if errors.As(err, &urlError) {
						tried++
						log.Error().Err(err).Msg("unable reporting metric")
						if tried == 3 {
							return err
						}
						continue loop
					}
					log.Error().Err(err).Str("metric", currMetric.String()).Msg("unable reporting metric")
				} else {
					log.Info().Str("metric", currMetric.String()).Msg("metric reported")
				}
			}
		case _, ok := <-shutdown:
			if ok {
				close(shutdown)
			}
			return nil
		}
	}
}

func (s *service) reportMetric(metricForReport metric.Metrics) (*resty.Response, error) {
	log.Debug().Str("metric", metricForReport.String()).Msg("service.reportMetric started")
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
