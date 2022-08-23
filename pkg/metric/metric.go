package metric

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
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

// TypeIsValid checks the validity of metrics.
func (m *Metrics) String() string {
	deltaStr := "nil"
	if m.Delta != nil {
		deltaStr = strconv.FormatInt(*m.Delta, 10)
	}
	valueStr := "nil"
	if m.Value != nil {
		valueStr = strconv.FormatFloat(*m.Value, 'f', -1, 64)
	}
	return "ID: " + m.ID +
		", MType: " + m.MType +
		", Delta: " + deltaStr +
		", Value: " + valueStr
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

	if keyForUpdate == "" {
		return ErrKeyForUpdateHashIsEmpty
	}
	if !m.TypeIsValid() {
		return ErrInvalidType
	}

	msgForHash := ""
	if m.MType == "gauge" {
		var valueForHash float64
		if m.Value != nil {
			valueForHash = *m.Value
		}
		msgForHash = fmt.Sprintf("%s:gauge:%f", m.ID, valueForHash)
	} else if m.MType == "counter" {
		var deltaForHash int64
		if m.Delta != nil {
			deltaForHash = *m.Delta
		}
		msgForHash = fmt.Sprintf("%s:counter:%d", m.ID, deltaForHash)
	}

	h := hmac.New(sha256.New, []byte(keyForUpdate))
	h.Write([]byte(msgForHash))
	m.Hash = fmt.Sprintf("%x", h.Sum(nil))

	return nil
}
