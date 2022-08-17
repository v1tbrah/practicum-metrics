package service

import (
	"log"
	"time"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/memory"
)

type Service struct {
	Storage repo.Storage
	Cfg     *config.Config
}

// NewService returns new Service.
func NewService(storage repo.Storage, cfg *config.Config) *Service {
	service := &Service{Storage: storage, Cfg: cfg}

	if cfg.StorageType == config.InMemory {

		inMemStorage, _ := service.Storage.(*memory.MemStorage)
		if service.Cfg.Restore {
			if err := inMemStorage.RestoreData(); err != nil {
				log.Println(err)
			} else {
				log.Println("Metrics restored.")
			}
		}

		if intervalIsSet := service.Cfg.StoreInterval != time.Second*0; intervalIsSet {
			go service.writeMetricsToFileWithInterval()
		}
	}

	return service
}

func (s *Service) writeMetricsToFileWithInterval() {
	if fileNameIsEmpty := s.Cfg.StoreFile == ""; fileNameIsEmpty {
		return
	}
	if intervalIsSet := s.Cfg.StoreInterval != time.Second*0; !intervalIsSet {
		return
	}
	ticker := time.NewTicker(s.Cfg.StoreInterval)
	inMemStorage, _ := s.Storage.(*memory.MemStorage)
	for {
		<-ticker.C
		if err := inMemStorage.StoreData(); err != nil {
			log.Println(err)
		} else {
			log.Println("Metrics saved to file:", s.Cfg.StoreFile)
		}
	}
}
