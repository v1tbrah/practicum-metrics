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
		"Alloc":         NewAlloc(),
		"BuckHashSys":   NewBuckHashSys(),
		"Frees":         NewFrees(),
		"GCCPUFraction": NewGCCPUFraction(),
		"GCSys":         NewGCSys(),
		"HeapAlloc":     NewHeapAlloc(),
		"HeapIdle":      NewHeapIdle(),
		"HeapInuse":     NewHeapInuse(),
		"HeapObjects":   NewHeapObjects(),
		"HeapReleased":  NewHeapReleased(),
		"HeapSys":       NewHeapSys(),
		"LastGC":        NewLastGC(),
		"Lookups":       NewLookups(),
		"MCacheInuse":   NewMCacheInuse(),
		"MCacheSys":     NewMCacheSys(),
		"MSpanInuse":    NewMSpanInuse(),
		"MSpanSys":      NewMSpanSys(),
		"Mallocs":       NewMallocs(),
		"NextGC":        NewNextGC(),
		"NumForcedGC":   NewNumForcedGC(),
		"NumGC":         NewNumGC(),
		"OtherSys":      NewOtherSys(),
		"PauseTotalNs":  NewPauseTotalNs(),
		"StackInuse":    NewStackInuse(),
		"StackSys":      NewStackSys(),
		"Sys":           NewSys(),
		"TotalAlloc":    NewTotalAlloc(),
		"PollCount":     NewPollCount(),
		"RandomValue":   NewRandomValue(),
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

func NewAlloc() Metrics {
	var value float64
	return Metrics{
		ID:    "Alloc",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewBuckHashSys() Metrics {
	var value float64
	return Metrics{
		ID:    "BuckHashSys",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewFrees() Metrics {
	var value float64
	return Metrics{
		ID:    "Frees",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewGCCPUFraction() Metrics {
	var value float64
	return Metrics{
		ID:    "GCCPUFraction",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewGCSys() Metrics {
	var value float64
	return Metrics{
		ID:    "GCSys",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewHeapAlloc() Metrics {
	var value float64
	return Metrics{
		ID:    "HeapAlloc",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewHeapIdle() Metrics {
	var value float64
	return Metrics{
		ID:    "HeapIdle",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewHeapInuse() Metrics {
	var value float64
	return Metrics{
		ID:    "HeapInuse",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewHeapObjects() Metrics {
	var value float64
	return Metrics{
		ID:    "HeapObjects",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewHeapReleased() Metrics {
	var value float64
	return Metrics{
		ID:    "HeapReleased",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewHeapSys() Metrics {
	var value float64
	return Metrics{
		ID:    "HeapSys",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewLastGC() Metrics {
	var value float64
	return Metrics{
		ID:    "LastGC",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewLookups() Metrics {
	var value float64
	return Metrics{
		ID:    "Lookups",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewMCacheInuse() Metrics {
	var value float64
	return Metrics{
		ID:    "MCacheInuse",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewMCacheSys() Metrics {
	var value float64
	return Metrics{
		ID:    "MCacheSys",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewMSpanInuse() Metrics {
	var value float64
	return Metrics{
		ID:    "MSpanInuse",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewMSpanSys() Metrics {
	var value float64
	return Metrics{
		ID:    "MSpanSys",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewMallocs() Metrics {
	var value float64
	return Metrics{
		ID:    "Mallocs",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewNextGC() Metrics {
	var value float64
	return Metrics{
		ID:    "NextGC",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewNumForcedGC() Metrics {
	var value float64
	return Metrics{
		ID:    "NumForcedGC",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewNumGC() Metrics {
	var value float64
	return Metrics{
		ID:    "NumGC",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewOtherSys() Metrics {
	var value float64
	return Metrics{
		ID:    "OtherSys",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewPauseTotalNs() Metrics {
	var value float64
	return Metrics{
		ID:    "PauseTotalNs",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewStackInuse() Metrics {
	var value float64
	return Metrics{
		ID:    "StackInuse",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewStackSys() Metrics {
	var value float64
	return Metrics{
		ID:    "StackSys",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewSys() Metrics {
	var value float64
	return Metrics{
		ID:    "Sys",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewTotalAlloc() Metrics {
	var value float64
	return Metrics{
		ID:    "TotalAlloc",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}

func NewPollCount() Metrics {
	var value int64
	return Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &value,
		Value: nil,
	}
}

func NewRandomValue() Metrics {
	var value float64
	return Metrics{
		ID:    "RandomValue",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
}
