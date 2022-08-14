package memory

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/server/model"
)

type MemStorage struct {
	Data *model.Data
}

// New returns new memory storage.
func New() *MemStorage {
	return &MemStorage{Data: model.NewData()}
}

func (m *MemStorage) GetData() *model.Data {
	return m.Data
}

func (m *MemStorage) SetData(data *model.Data) {
	m.Data = data
}
