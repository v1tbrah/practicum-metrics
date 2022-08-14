package pg

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/server/model"
)

type pgStorage struct {
	Data *model.Data
}

// New returns new postgres storage.
func New() *pgStorage {
	return &pgStorage{Data: model.NewData()}
}

func (p *pgStorage) GetData() *model.Data {
	return p.Data
}

func (p *pgStorage) SetData(data *model.Data) {
	p.Data = data
}
