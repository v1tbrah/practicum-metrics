package metric

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
)

var (
	ErrInvalidType             = errors.New("invalid type of metric")
	ErrKeyForUpdateHashIsEmpty = errors.New("key for update hash is empty")
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

// NewMetric returns new Metrics.
func NewMetric(ID, MType string) Metrics {
	newMetric := Metrics{
		ID:    ID,
		MType: MType,
	}
	if MType == "gauge" {
		var val float64
		newMetric.Value = &val
	} else if MType == "counter" {
		var delta int64
		newMetric.Delta = &delta
	}
	return newMetric
}

// TypeIsValid checks the validity of metrics.
func (m *Metrics) TypeIsValid() bool {
	return m.MType == "gauge" || m.MType == "counter"
}

func (m *Metrics) UpdateHash(keyForUpdate string) error {

	log.Println("ID:", m.ID, "Type:", m.MType, "value:", m.Value, "delta:", m.Delta)
	if keyForUpdate == "" {
		return ErrKeyForUpdateHashIsEmpty
	}
	if !m.TypeIsValid() {
		return ErrInvalidType
	}

	valueForHash := ""
	if m.MType == "gauge" {
		valueForHash = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
	} else if m.MType == "counter" {
		valueForHash = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	}

	h := hmac.New(sha256.New, []byte(keyForUpdate))
	h.Write([]byte(valueForHash))
	m.Hash = fmt.Sprintf("%x", h.Sum(nil))

	return nil
}
