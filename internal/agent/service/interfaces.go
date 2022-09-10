package service

import (
	"time"

	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

type Data interface {
	GetData() map[string]metric.Metrics
	UpdateBasic(keyForUpdateHash string)
	UpdateAdditional(keyForUpdateHash string) error
}

type Config interface {
	String() string
	PollInterval() time.Duration
	ReportInterval() time.Duration
	HashKey() string
	ReportMetricURL() string
	ReportListMetricsURL() string
	GetMetricURL() string
}
