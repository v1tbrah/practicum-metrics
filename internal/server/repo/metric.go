package repo

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

// TypeIsValid checks the validity of metrics.
func (m *Metrics) TypeIsValid() bool {
	return m.MType == "gauge" || m.MType == "counter"
}
