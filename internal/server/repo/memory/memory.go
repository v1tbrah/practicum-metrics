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

type Memory struct {
	data      model.Data
	storeFile string
	mu        sync.RWMutex
}

// New returns new memory storage.
func New(storeFile string) *Memory {
	log.Debug().Str("storeFile", storeFile).Msg("memory.New started")
	defer log.Debug().Msg("memory.New ended")

	return &Memory{data: model.NewData(), storeFile: storeFile}
}

func (m *Memory) GetMetric(ctx context.Context, ID string) (metric.Metrics, bool, error) {
	log.Debug().Str("MID", ID).Msg("memory.GetMetric started")
	defer log.Debug().Msg("memory.GetMetric ended")

	m.mu.RLock()
	defer m.mu.RUnlock()

	currMetric, ok := m.data[ID]
	if !ok {
		return currMetric, ok, nil
	}
	resultMetric := metric.NewMetric(currMetric.ID, currMetric.MType)
	if currMetric.MType == "gauge" {
		valueForResult := *currMetric.Value
		resultMetric.Value = &valueForResult
	} else if currMetric.MType == "counter" {
		deltaForResult := *currMetric.Delta
		resultMetric.Delta = &deltaForResult
	}

	return resultMetric, ok, nil
}

func (m *Memory) SetMetric(ctx context.Context, thisMetric metric.Metrics) error {
	log.Debug().Str("thisMetric", fmt.Sprint(thisMetric)).Msg("memory.SetMetric started")
	defer log.Debug().Msg("memory.SetMetric ended")

	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[thisMetric.ID] = thisMetric
	return nil
}

func (m *Memory) SetListMetrics(ctx context.Context, listMetrics []metric.Metrics) error {
	log.Debug().Str("listMetrics", fmt.Sprint(listMetrics)).Msg("memory.SetListMetrics started")
	defer log.Debug().Msg("memory.SetListMetrics ended")

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, currMetric := range listMetrics {
		m.data[currMetric.ID] = currMetric
	}
	return nil
}

func (m *Memory) GetData(ctx context.Context) (model.Data, error) {
	log.Debug().Msg("memory.GetData started")
	defer log.Debug().Msg("memory.GetData ended")

	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(model.Data, len(m.data))

	for _, currMetric := range m.data {
		metricForResult := metric.NewMetric(currMetric.ID, currMetric.MType)
		if currMetric.MType == "gauge" {
			valueForResult := *currMetric.Value
			metricForResult.Value = &valueForResult
		} else if currMetric.MType == "counter" {
			deltaForResult := *currMetric.Delta
			metricForResult.Delta = &deltaForResult
		}
		result[metricForResult.ID] = metricForResult
	}

	return result, nil
}

func (m *Memory) RestoreData() error {
	log.Debug().Msg("memory.RestoreData started")
	defer log.Debug().Msg("memory.RestoreData ended")

	file, err := os.Open(m.storeFile)
	if err != nil {
		return err
	}
	newMetrics := model.NewData()
	if err = json.NewDecoder(file).Decode(&newMetrics); err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
	}
	m.data = newMetrics
	return nil
}

func (m *Memory) StoreData(ctx context.Context) error {
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
