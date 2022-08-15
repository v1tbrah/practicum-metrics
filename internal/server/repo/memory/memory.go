package memory

import (
	"encoding/json"
	"errors"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/model"
	"io"
	"log"
	"os"
)

type MemStorage struct {
	Data      *model.Data
	storeFile string
}

// New returns new memory storage.
func New(storeFile string) *MemStorage {
	return &MemStorage{Data: model.NewData(), storeFile: storeFile}
}

func (m *MemStorage) GetData() *model.Data {
	return m.Data
}

func (m *MemStorage) SetData(data *model.Data) {
	m.Data = data
}

func (m *MemStorage) RestoreData() error {
	file, err := os.Open(m.storeFile)
	if err != nil {
		return err
	}
	newMetrics := model.NewData()
	if err = json.NewDecoder(file).Decode(newMetrics); err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
	}
	m.SetData(newMetrics)
	return nil
}

func (m *MemStorage) StoreData() error {
	if m.storeFile == "" {
		return errors.New("file name is empty")
	}
	file, err := os.Create(m.storeFile)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	dataMetrics, err := json.Marshal(m.GetData())
	if err != nil {
		log.Println(err)
		return err
	}
	if _, err = file.Write(dataMetrics); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
