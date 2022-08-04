package memory

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo"
)

type Storage struct {
	Data *repo.Data
}

// NewStorage returns new memory storage.
func NewStorage() *Storage {
	return &Storage{Data: repo.NewData()}
}
