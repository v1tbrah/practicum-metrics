package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/model"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/storage"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/storage/memory"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/storage/pg"
	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

var ErrNotDBStorage = errors.New("type of storage is not DB storage")

type Service struct {
	storage storage.Storage
	cfg     Config
}

// New returns new Service.
func New(storage storage.Storage, cfg Config) (*Service, error) {
	log.Debug().
		Str("storage", fmt.Sprint(storage)).
		Str("cfg", cfg.String()).
		Msg("service.New started")
	defer log.Debug().Msg("service.New ended")

	service := &Service{storage: storage, cfg: cfg}

	if cfg.StorageType() == config.StorageTypeMemory {
		inMemStorage, ok := service.storage.(*memory.Memory)
		if !ok {
			return nil, errors.New("type of storage is not *memory.MemStorage with StorageType == StorageTypeMemory")
		}
		if service.cfg.Restore() && service.cfg.StoreFile() != "" {
			if err := inMemStorage.RestoreData(); err != nil {
				return nil, fmt.Errorf("err restoring data: %w", err)
			}
			log.Info().Msg("data restored")
		}

		if intervalIsSet := service.cfg.StoreInterval() != time.Second*0; intervalIsSet {
			go service.writeMetricsToFileWithInterval(context.Background())
		}
	}

	return service, nil
}

func (s *Service) GetMetric(ctx context.Context, ID string) (metric.Metrics, bool, error) {
	return s.storage.GetMetric(ctx, ID)
}

func (s *Service) SetMetric(ctx context.Context, metricData metric.Metrics) error {
	return s.storage.SetMetric(ctx, metricData)
}

func (s *Service) SetListMetrics(ctx context.Context, listMetrics []metric.Metrics) error {
	return s.storage.SetListMetrics(ctx, listMetrics)
}

func (s *Service) GetData(ctx context.Context) (model.Data, error) {
	return s.storage.GetData(ctx)
}

func (s *Service) PingDBStorage(ctx context.Context) error {
	dbStorage, ok := s.storage.(*pg.Pg)
	if !ok {
		return errors.New("type of storage is not DB storage")
	}
	return dbStorage.Ping(ctx)
}

func (s *Service) ShutDown() {
	log.Debug().Msg("service.ShutDown started")
	defer log.Debug().Msg("service.ShutDown ended")

	if s.cfg.StorageType() == config.StorageTypeMemory {
		memStorage, ok := s.storage.(*memory.Memory)
		if !ok {
			log.Error().Err(errors.New("unknown storage type")).Msg("shutting down the service")
			return
		}
		if err := memStorage.StoreData(context.Background()); err != nil {
			log.Error().Err(err).
				Str("storeFile", s.cfg.StoreFile()).
				Msg("unable to store data in file")
			return
		}
		log.Info().Msg(fmt.Sprintf("data saved to file: %s", s.cfg.StoreFile()))
	} else if s.cfg.StorageType() == config.StorageTypeDB {
		DBStorage, ok := s.storage.(*pg.Pg)
		if !ok {
			log.Error().Err(errors.New("unknown storage type")).Msg("shutting down the service")
			return
		}
		DBStorage.ClosePoolConn()
	}
}

func (s *Service) writeMetricsToFileWithInterval(ctx context.Context) {
	log.Debug().
		Str("storeFile", s.cfg.StoreFile()).
		Dur("storeInterval", s.cfg.StoreInterval()).
		Msg("service.writeMetricsToFileWithInterval started")
	defer log.Debug().Msg("service.writeMetricsToFileWithInterval ended")

	if fileNameIsEmpty := s.cfg.StoreFile() == ""; fileNameIsEmpty {
		return
	}
	if intervalIsSet := s.cfg.StoreInterval() != time.Second*0; !intervalIsSet {
		return
	}
	ticker := time.NewTicker(s.cfg.StoreInterval())
	inMemStorage, _ := s.storage.(*memory.Memory)
	for {
		<-ticker.C
		if err := inMemStorage.StoreData(ctx); err != nil {
			log.Error().Err(err).
				Str("storeFile", s.cfg.StoreFile()).
				Msg("unable to store data in file")
		} else {
			log.Info().Msg(fmt.Sprintf("data saved to file: %s", s.cfg.StoreFile()))
		}
	}
}
