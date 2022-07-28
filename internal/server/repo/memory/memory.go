package memory

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/metric"
)

type MemStorage struct {
	Metrics *metric.AllMetrics
}

// Creates an MemStorage.
func NewMemStorage() *MemStorage {
	return &MemStorage{Metrics: metric.NewAllMetrics()}
}
