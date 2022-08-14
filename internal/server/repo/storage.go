package repo

import (
	"errors"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/model"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/memory"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/pg"
	"strings"
)

var ErrEmptyConfig = errors.New("empty config")

type Storage interface {
	GetData() *model.Data
	SetData(*model.Data)
}

func New(cfg *config.Config) (Storage, error) {
	if cfg == nil {
		return nil, ErrEmptyConfig
	}

	switch {
	case strings.TrimSpace(cfg.PgConnString) != "":
		return pg.New(), nil
	default:
		return memory.New(), nil
	}

}