package repo

import (
	"encoding/json"
	"errors"
	"sync"
)

type Data struct {
	sync.Map
}

// NewData returns new data.
func NewData() *Data {
	return &Data{}
}

func (d *Data) MarshalJSON() ([]byte, error) {
	dataMetrics := make(map[string]Metrics)
	metricsConverted := true
	d.Range(func(key, value interface{}) bool {
		metricName, ok := key.(string)
		if !ok {
			metricsConverted = false
			return false
		}
		metricValue, ok := value.(Metrics)
		if !ok {
			metricsConverted = false
			return false
		}
		dataMetrics[metricName] = metricValue
		return true
	})
	if !metricsConverted {
		errText := "error converting metrics to json"
		return nil, errors.New(errText)
	}
	jsonMetrics, err := json.Marshal(&dataMetrics)
	if err != nil {
		return nil, err
	}
	return jsonMetrics, nil
}

func (d *Data) UnmarshalJSON(data []byte) error {
	tmpMetrics := map[string]Metrics{}
	if err := json.Unmarshal(data, &tmpMetrics); err != nil {
		return err
	}
	for key, value := range tmpMetrics {
		d.Store(key, value)
	}
	return nil
}
