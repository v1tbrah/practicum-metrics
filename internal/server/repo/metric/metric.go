package metric

import (
	"encoding/json"
	"errors"
	"sync"
)

type AllMetrics struct {
	sync.Map
}

// Creates a Metrics.
func NewAllMetrics() *AllMetrics {
	return &AllMetrics{}
}

func (a *AllMetrics) MarshalJSON() ([]byte, error) {
	dataMetrics := make(map[string]Metrics)
	tryConvertMetrics := true
	a.Range(func(key, value interface{}) bool {
		metricName, ok := key.(string)
		if !ok {
			tryConvertMetrics = false
			return false
		}
		metricValue, ok := value.(Metrics)
		if !ok {
			tryConvertMetrics = false
			return false
		}
		dataMetrics[metricName] = metricValue
		return true
	})
	if !tryConvertMetrics {
		errText := "error converting metrics to json"
		return nil, errors.New(errText)
	}
	jsonMetrics, err := json.Marshal(&dataMetrics)
	if err != nil {
		return nil, err
	}
	return jsonMetrics, nil
}

func (a *AllMetrics) UnmarshalJSON(data []byte) error {
	tmpMetrics := map[string]Metrics{}
	if err := json.Unmarshal(data, &tmpMetrics); err != nil {
		return err
	}
	for key, value := range tmpMetrics {
		a.Store(key, value)
	}
	return nil
}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (m *Metrics) TypeIsValid() bool {
	if m.MType == "gauge" || m.MType == "counter" {
		return true
	}
	return false
}
