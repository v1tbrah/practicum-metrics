package api

import (
	"context"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/model"
	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

type Service interface {
	GetMetric(context.Context, string) (metric.Metrics, bool, error)
	SetMetric(context.Context, metric.Metrics) error
	GetData(context.Context) (model.Data, error)
	SetListMetrics(context.Context, []metric.Metrics) error
	PingDBStorage(context.Context) error
	ShutDown()
}

type Config interface {
	ServAddr() string
	HashKey() string
}
