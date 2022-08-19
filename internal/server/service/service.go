package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/memory"
)

type Service struct {
	Storage repo.Storage
	Cfg     *config.Config
}

// New returns new Service.
func New(storage repo.Storage, cfg *config.Config) (*Service, error) {
	log.Debug().
		Str("storage", fmt.Sprint(storage)).
		Str("cfg", fmt.Sprint(cfg)).
		Msg("service.New started")
	defer log.Debug().Msg("service.New ended")

	service := &Service{Storage: storage, Cfg: cfg}

	if cfg.StorageType == config.StorageTypeMemory {
		inMemStorage, ok := service.Storage.(*memory.Memory)
		if !ok {
			return nil, errors.New("type of storage is not *memory.MemStorage with StorageType == StorageTypeMemory")
		}
		if service.Cfg.Restore && service.Cfg.StoreFile != "" {
			if err := inMemStorage.RestoreData(); err != nil {
				return nil, fmt.Errorf("err restoring data: %w", err)
			}
			log.Info().Msg("data restored")
		}

		if intervalIsSet := service.Cfg.StoreInterval != time.Second*0; intervalIsSet {
			go service.writeMetricsToFileWithInterval(context.Background())
		}
	}

	return service, nil
}

func (s *Service) writeMetricsToFileWithInterval(ctx context.Context) {
	log.Debug().
		Str("storeFile", s.Cfg.StoreFile).
		Dur("storeInterval", s.Cfg.StoreInterval).
		Msg("service.writeMetricsToFileWithInterval started")
	defer log.Debug().Msg("service.writeMetricsToFileWithInterval ended")

	if fileNameIsEmpty := s.Cfg.StoreFile == ""; fileNameIsEmpty {
		return
	}
	if intervalIsSet := s.Cfg.StoreInterval != time.Second*0; !intervalIsSet {
		return
	}
	ticker := time.NewTicker(s.Cfg.StoreInterval)
	inMemStorage, _ := s.Storage.(*memory.Memory)
	for {
		<-ticker.C
		if err := inMemStorage.StoreData(ctx); err != nil {
			log.Error().
				Err(err).
				Str("storeFile", s.Cfg.StoreFile).
				Msg("unable to store data in file")
		} else {
			log.Info().Msg(fmt.Sprintf("data saved to file: %s", s.Cfg.StoreFile))
		}
	}
}
