package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/model"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/memory"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/pg"
	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

var ErrEmptyConfig = errors.New("empty config")

type Storage interface {
	GetData(ctx context.Context) (model.Data, error)
	GetMetric(ctx context.Context, ID string) (metric.Metrics, bool, error)
	SetMetric(ctx context.Context, thisMetric metric.Metrics) error
	SetListMetrics(ctx context.Context, listMetrics []metric.Metrics) error
}

func New(cfg *config.Config) (Storage, error) {
	log.Debug().Str("config", fmt.Sprint(cfg)).Msg("storage.New started")
	defer log.Debug().Msg("storage.New ended")

	if cfg == nil {
		return nil, ErrEmptyConfig
	}

	switch {
	case cfg.StorageType == config.StorageTypeDB:
		return pg.New(cfg.PgConnString)
	default:
		return memory.New(cfg.StoreFile), nil
	}

}
