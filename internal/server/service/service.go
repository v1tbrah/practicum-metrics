package service

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/server/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo"
	"log"
	"strings"
	"time"
)

type Service struct {
	Storage repo.Storage
	Cfg     *config.Config
}

// NewService returns new Service.
func NewService(storage repo.Storage, cfg *config.Config) *Service {
	service := &Service{Storage: storage, Cfg: cfg}

	if service.Cfg.Restore {
		if err := service.Storage.RestoreData(); err != nil {
			log.Println(err)
		} else {
			log.Println("Metrics restored")
		}
	}

	if haveDBConnection := strings.TrimSpace(cfg.PgConnString) != ""; !haveDBConnection {
		if needWriteMetricsToFileWithInterval := service.Cfg.StoreInterval != time.Second*0; needWriteMetricsToFileWithInterval {
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
	for {
		<-ticker.C
		if err := s.Storage.StoreData(); err != nil {
			log.Println(err)
		} else {
			log.Println("Metrics saved to file:", s.Cfg.StoreFile)
		}
	}
}
