package repo

import (
	"errors"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/model"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/memory"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/pg"
	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
	"strings"
)

var ErrEmptyConfig = errors.New("empty config")

type Storage interface {
	GetData() (*model.Data, error)
	GetMetric(ID string) (metric.Metrics, bool, error)
	SetMetric(ID string, thisMetric metric.Metrics) error
}

func New(cfg *config.Config) (Storage, error) {
	if cfg == nil {
		return nil, ErrEmptyConfig
	}

	switch {
	case strings.TrimSpace(cfg.PgConnString) != "":
		return pg.New(cfg.PgConnString), nil
	default:
		return memory.New(cfg.StoreFile), nil
	}

}
