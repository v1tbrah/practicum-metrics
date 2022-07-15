package service

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/memory"
)

type Service struct {
	MemStorage *memory.MemStorage
}

// Creates a Service.
func NewService(memStorage *memory.MemStorage) *Service {
	return &Service{MemStorage: memStorage}
}
