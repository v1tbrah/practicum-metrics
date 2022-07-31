package service

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/memory"
)

type Service struct {
	MemStorage *memory.Storage
}

// NewService returns new Service.
func NewService(memStorage *memory.Storage) *Service {
	return &Service{MemStorage: memStorage}
}

func (s *Service) SaveMetricsToFile(fileName string) error {
	if fileName == "" {
		return errors.New("file name is empty")
	}
	file, err := os.Create(fileName)
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

func (s *Service) WriteMetricsToFileWithInterval(fileName string, interval time.Duration) {
	if fileNameIsEmpty := fileName == ""; fileNameIsEmpty {
		return
	}
	if intervalIsSet := interval != time.Second*0; !intervalIsSet {
		return
	}
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		if err := s.SaveMetricsToFile(fileName); err != nil {
			log.Println(err)
		} else {
			log.Println("Metrics saved to file:", fileName)
		}
	}
}

func (s *Service) RestoreMetricsFromFile(fileName string) error {
	file, err := os.Open(fileName)
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
