package memory

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
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
