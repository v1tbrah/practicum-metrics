package service

import (
	"encoding/json"
	"errors"

	"io"
	"log"
	"os"
	"time"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/memory"
)

type Service struct {
	MemStorage *memory.Storage
	Cfg        *config.Config
}

// NewService returns new Service.
func NewService(memStorage *memory.Storage, cfg *config.Config) *Service {
	service := &Service{MemStorage: memStorage, Cfg: cfg}
	if service.Cfg.Restore {
		if err := service.restoreMetricsFromFile(); err != nil {
			log.Println(err)
		} else {
			log.Println("Metrics restored from file:", service.Cfg.StoreFile)
		}
	}

	if needWriteMetricsToFileWithInterval := service.Cfg.StoreInterval != time.Second*0; needWriteMetricsToFileWithInterval {
		go service.writeMetricsToFileWithInterval()
	}
	return service
}

func (s *Service) SaveMetricsToFile() error {
	if s.Cfg.StoreFile == "" {
		return errors.New("file name is empty")
	}
	file, err := os.Create(s.Cfg.StoreFile)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	dataMetrics, err := json.Marshal(s.MemStorage.Data)
	if err != nil {
		log.Println(err)
		return err
	}
	if _, err = file.Write(dataMetrics); err != nil {
		log.Println(err)
		return err
	}
	return nil
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
		if err := s.SaveMetricsToFile(); err != nil {
			log.Println(err)
		} else {
			log.Println("Metrics saved to file:", s.Cfg.StoreFile)
		}
	}
}

func (s *Service) restoreMetricsFromFile() error {
	file, err := os.Open(s.Cfg.StoreFile)
	if err != nil {
		return err
	}
	newMetrics := repo.NewData()
	if err = json.NewDecoder(file).Decode(newMetrics); err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
	}
	s.MemStorage.Data = newMetrics
	return nil
}
