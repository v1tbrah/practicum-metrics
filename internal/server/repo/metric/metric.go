package metric

import (
	"errors"
	"sync"
)

const ErrInvalidType = "invalid type of metric"

type Metrics struct {
	gauge   *sync.Map
	counter *sync.Map
}

// Creates a Metrics.
func NewMetrics() *Metrics {
	return &Metrics{gauge: &sync.Map{}, counter: &sync.Map{}}
}

// Returns metrics by type.
func (m *Metrics) MetricsOfType(typeM string) (*sync.Map, error) {
	switch typeM {
	case "gauge":
		return m.gauge, nil
	case "counter":
		return m.counter, nil
	default:
		return nil, errors.New(ErrInvalidType)
	}
}

// Checks if the metric type exists.
func TypeIsValid(checked string) bool {
	switch checked {
	case "gauge":
		return true
	case "counter":
		return true
	default:
		return false
	}
}
