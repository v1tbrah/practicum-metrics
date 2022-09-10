package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

var (
	ErrMetricIDIsEmpty    = errors.New("metric id is empty")
	ErrInvalidMetricType  = errors.New("invalid type of metric")
	ErrListMetricsIsEmpty = errors.New("list metrics is empty")
)

func (s *service) reportMetric(metricForReport metric.Metrics) (*resty.Response, error) {
	log.Debug().Str("metric", metricForReport.String()).Msg("service.reportMetric started")
	var err error
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("service.reportMetric ended")
		} else {
			log.Debug().Msg("service.reportMetric ended")
		}
	}()

	if metricForReport.ID == "" {
		return nil, ErrMetricIDIsEmpty
	}
	if !metricForReport.TypeIsValid() {
		return nil, ErrInvalidMetricType
	}

	body, err := json.Marshal(metricForReport)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.NewRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(s.cfg.ReportMetricURL())

	return resp, err
}

func (s *service) reportListMetrics(listMetrics []metric.Metrics) (*resty.Response, error) {
	log.Debug().
		Str("listMetrics", fmt.Sprint(listMetrics)).
		Msg("service.reportListMetrics started")
	var err error
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("service.reportListMetrics ended")
		} else {
			log.Debug().Msg("service.reportListMetrics ended")
		}
	}()

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
		Post(s.cfg.ReportListMetricsURL())

	return resp, err
}

func (s *service) getMetric(ID, MType string) (*resty.Response, error) {
	log.Debug().Str("ID", ID).Str("MType", MType).Msg("service.getMetric started")
	var err error
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("service.getMetric ended")
		} else {
			log.Debug().Msg("service.getMetric ended")
		}
	}()

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
		Post(s.cfg.GetMetricURL())

	return resp, err
}
