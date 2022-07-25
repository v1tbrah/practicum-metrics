package metric

import (
	"sync"
)

const ErrInvalidType = "invalid type of metrics"

// Creates a Metrics.
func NewMetrics() *sync.Map {
	return &sync.Map{}
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
