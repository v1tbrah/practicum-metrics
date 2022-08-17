package memory

import (
	"encoding/json"
	"errors"

	"io"
	"log"
	"os"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/model"
	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

type MemStorage struct {
	Data      *model.Data
	storeFile string
}

// New returns new memory storage.
func New(storeFile string) *MemStorage {
	return &MemStorage{Data: model.NewData(), storeFile: storeFile}
}

func (m *MemStorage) GetMetric(ID string) (metric.Metrics, bool, error) {
	thisMetric, ok := m.Data.Metrics[ID]
	return thisMetric, ok, nil
}

func (m *MemStorage) SetMetric(ID string, thisMetric metric.Metrics) error {
	m.Data.Metrics[ID] = thisMetric
	return nil
}

func (m *MemStorage) SetListMetrics(listMetrics []metric.Metrics) error {
	for _, currMetric := range listMetrics {
		m.Data.Metrics[currMetric.ID] = currMetric
	}
	return nil
}
func (m *MemStorage) GetData() (*model.Data, error) {
	return m.Data, nil
}

func (m *MemStorage) RestoreData() error {
	file, err := os.Open(m.storeFile)
	if err != nil {
		return err
	}
	newMetrics := model.NewData()
	if err = json.NewDecoder(file).Decode(newMetrics); err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
	}
	m.Data = newMetrics
	return nil
}

func (m *MemStorage) StoreData() error {
	if m.storeFile == "" {
		return errors.New("file name is empty")
	}
	file, err := os.Create(m.storeFile)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	dataMetrics, err := m.GetData()
	if err != nil {
		log.Println(err)
		return err
	}
	dataMetricsJSON, err := json.Marshal(dataMetrics)
	if err != nil {
		log.Println(err)
		return err
	}
	if _, err = file.Write(dataMetricsJSON); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
