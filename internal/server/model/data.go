package model

import (
	"encoding/json"
	"sync"

	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

type Data struct {
	Metrics map[string]metric.Metrics
	sync.Mutex
}

// NewData returns new data.
func NewData() *Data {
	return &Data{Metrics: map[string]metric.Metrics{}}
}

func (d *Data) MarshalJSON() ([]byte, error) {
	jsonMetrics, err := json.Marshal(d.Metrics)
	if err != nil {
		return nil, err
	}
	return jsonMetrics, nil
}

func (d *Data) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &d.Metrics); err != nil {
		return err
	}
	return nil
}
