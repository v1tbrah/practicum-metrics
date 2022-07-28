package service

import (
	"encoding/json"
	"errors"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/memory"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/metric"
	"io"
	"log"
	"os"
	"time"
)

type Service struct {
	MemStorage *memory.MemStorage
}

// Creates a Service.
func NewService(memStorage *memory.MemStorage) *Service {
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

	dataMetrics, err := json.Marshal(s.MemStorage.Metrics)
	if err != nil {
		log.Println(err)
		return err
	}
	if _, err = file.Write([]byte(dataMetrics)); err != nil {
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
	storeInterval := time.NewTicker(interval)
	for {
		<-storeInterval.C
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
	newMetrics := metric.NewAllMetrics()
	if err = json.NewDecoder(file).Decode(newMetrics); err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
	}
	s.MemStorage.Metrics = newMetrics
	return nil
}
