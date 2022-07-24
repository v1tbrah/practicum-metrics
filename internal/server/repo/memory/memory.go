package memory

import (
	"sync"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/metric"
)

type MemStorage struct {
	Metrics *sync.Map
}

// Creates an MemStorage.
func NewMemStorage() *MemStorage {
	return &MemStorage{Metrics: metric.NewMetrics()}
}
