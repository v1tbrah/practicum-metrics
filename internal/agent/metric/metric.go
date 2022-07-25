package metric

import (
	"math/rand"
	"reflect"
	"runtime"
	"time"
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type AllMetrics map[string]Metrics

// Creates an AllMetrics.
func NewMetrics() *AllMetrics {
	return &AllMetrics{
		"Alloc":         NewMetric("gauge", "Alloc"),
		"BuckHashSys":   NewMetric("gauge", "BuckHashSys"),
		"Frees":         NewMetric("gauge", "Frees"),
		"GCCPUFraction": NewMetric("gauge", "GCCPUFraction"),
		"GCSys":         NewMetric("gauge", "GCSys"),
		"HeapAlloc":     NewMetric("gauge", "HeapAlloc"),
		"HeapIdle":      NewMetric("gauge", "HeapIdle"),
		"HeapInuse":     NewMetric("gauge", "HeapInuse"),
		"HeapObjects":   NewMetric("gauge", "HeapObjects"),
		"HeapReleased":  NewMetric("gauge", "HeapReleased"),
		"HeapSys":       NewMetric("gauge", "HeapSys"),
		"LastGC":        NewMetric("gauge", "LastGC"),
		"Lookups":       NewMetric("gauge", "Lookups"),
		"MCacheInuse":   NewMetric("gauge", "MCacheInuse"),
		"MCacheSys":     NewMetric("gauge", "MCacheSys"),
		"MSpanInuse":    NewMetric("gauge", "MSpanInuse"),
		"MSpanSys":      NewMetric("gauge", "MSpanSys"),
		"Mallocs":       NewMetric("gauge", "Mallocs"),
		"NextGC":        NewMetric("gauge", "NextGC"),
		"NumForcedGC":   NewMetric("gauge", "NumForcedGC"),
		"NumGC":         NewMetric("gauge", "NumGC"),
		"OtherSys":      NewMetric("gauge", "OtherSys"),
		"PauseTotalNs":  NewMetric("gauge", "PauseTotalNs"),
		"StackInuse":    NewMetric("gauge", "StackInuse"),
		"StackSys":      NewMetric("gauge", "StackSys"),
		"Sys":           NewMetric("gauge", "Sys"),
		"TotalAlloc":    NewMetric("gauge", "TotalAlloc"),
		"PollCount":     NewMetric("counter", "PollCount"),
		"RandomValue":   NewMetric("gauge", "RandomValue"),
	}
}

// Updates all metrics. Gauge metrics are read from runtime.ReadMemStats.
func (m *AllMetrics) Update() {
	mForUpd := *m

	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	rStats := reflect.ValueOf(stats)

	for name, metric := range mForUpd {
		if metric.MType != "gauge" {
			continue
		}
		rStatValue := rStats.FieldByName(name)
		if rStatValue.Kind() == reflect.Invalid || !rStatValue.CanInterface() {
			continue
		}
		var vForUpd float64
		if rStatValue.CanUint() {
			statValue, ok := rStatValue.Interface().(uint64)
			if !ok {
				statValue, _ := rStatValue.Interface().(uint32)
				vForUpd = float64(statValue)
			}
			vForUpd = float64(statValue)
		} else if rStatValue.CanFloat() {
			vForUpd = rStatValue.Interface().(float64)
		}
		metric.Value = &vForUpd
		mForUpd[name] = metric
	}

	mRandomValue := mForUpd["RandomValue"]
	rand.Seed(time.Now().UnixNano())
	newRandomValue := rand.Float64()
	*mRandomValue.Value = newRandomValue

	mPollCount := mForUpd["PollCount"]
	*mPollCount.Delta++

}

func NewMetric(MType, ID string) Metrics {
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
