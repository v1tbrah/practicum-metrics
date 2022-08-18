package memory

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/model"
	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

type MemStorage struct {
	Data      *model.Data
	storeFile string
	mu        sync.RWMutex
}

// New returns new memory storage.
func New(storeFile string) *MemStorage {
	log.Debug().Str("storeFile", storeFile).Msg("memory.New started")
	defer log.Debug().Msg("memory.New ended")

	return &MemStorage{Data: model.NewData(), storeFile: storeFile}
}

func (m *MemStorage) GetMetric(ctx context.Context, ID string) (metric.Metrics, bool, error) {
	log.Debug().Str("MID", ID).Msg("memory.GetMetric started")
	defer log.Debug().Msg("memory.GetMetric ended")

	m.mu.RLock()
	defer m.mu.RUnlock()

	thisMetric, ok := m.Data.Metrics[ID]
	return thisMetric, ok, nil
}

func (m *MemStorage) SetMetric(ctx context.Context, thisMetric metric.Metrics) error {
	log.Debug().Str("thisMetric", fmt.Sprint(thisMetric)).Msg("memory.SetMetric started")
	defer log.Debug().Msg("memory.SetMetric ended")

	m.mu.Lock()
	defer m.mu.Unlock()

	m.Data.Metrics[thisMetric.ID] = thisMetric
	return nil
}

func (m *MemStorage) SetListMetrics(ctx context.Context, listMetrics []metric.Metrics) error {
	log.Debug().Str("listMetrics", fmt.Sprint(listMetrics)).Msg("memory.SetListMetrics started")
	defer log.Debug().Msg("memory.SetListMetrics ended")

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, currMetric := range listMetrics {
		m.Data.Metrics[currMetric.ID] = currMetric
	}
	return nil
}

func (m *MemStorage) GetData(ctx context.Context) (*model.Data, error) {
	log.Debug().Msg("memory.GetData started")
	defer log.Debug().Msg("memory.GetData ended")

	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.Data, nil
}

func (m *MemStorage) RestoreData() error {
	log.Debug().Msg("memory.RestoreData started")
	defer log.Debug().Msg("memory.RestoreData ended")

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

func (m *MemStorage) StoreData(ctx context.Context) error {
	log.Debug().Msg("memory.StoreData started")
	defer log.Debug().Msg("memory.StoreData ended")

	if m.storeFile == "" {
		return errors.New("file name is empty")
	}
	file, err := os.Create(m.storeFile)
	if err != nil {
		return err
	}
	defer file.Close()

	dataMetrics, err := m.GetData(ctx)
	if err != nil {
		return err
	}
	dataMetricsJSON, err := json.Marshal(dataMetrics)
	if err != nil {
		return err
	}
	if _, err = file.Write(dataMetricsJSON); err != nil {
		return err
	}
	return nil
}
