package memory

import (
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"time"
)

type Data map[string]Metrics

// NewData returns all metrics.
func NewData() *Data {
	return &Data{
		"Alloc":         NewMetric("Alloc", "gauge"),
		"BuckHashSys":   NewMetric("BuckHashSys", "gauge"),
		"Frees":         NewMetric("Frees", "gauge"),
		"GCCPUFraction": NewMetric("GCCPUFraction", "gauge"),
		"GCSys":         NewMetric("GCSys", "gauge"),
		"HeapAlloc":     NewMetric("HeapAlloc", "gauge"),
		"HeapIdle":      NewMetric("HeapIdle", "gauge"),
		"HeapInuse":     NewMetric("HeapInuse", "gauge"),
		"HeapObjects":   NewMetric("HeapObjects", "gauge"),
		"HeapReleased":  NewMetric("HeapReleased", "gauge"),
		"HeapSys":       NewMetric("HeapSys", "gauge"),
		"LastGC":        NewMetric("LastGC", "gauge"),
		"Lookups":       NewMetric("Lookups", "gauge"),
		"MCacheInuse":   NewMetric("MCacheInuse", "gauge"),
		"MCacheSys":     NewMetric("MCacheSys", "gauge"),
		"MSpanInuse":    NewMetric("MSpanInuse", "gauge"),
		"MSpanSys":      NewMetric("MSpanSys", "gauge"),
		"Mallocs":       NewMetric("Mallocs", "gauge"),
		"NextGC":        NewMetric("NextGC", "gauge"),
		"NumForcedGC":   NewMetric("NumForcedGC", "gauge"),
		"NumGC":         NewMetric("NumGC", "gauge"),
		"OtherSys":      NewMetric("OtherSys", "gauge"),
		"PauseTotalNs":  NewMetric("PauseTotalNs", "gauge"),
		"StackInuse":    NewMetric("StackInuse", "gauge"),
		"StackSys":      NewMetric("StackSys", "gauge"),
		"Sys":           NewMetric("Sys", "gauge"),
		"TotalAlloc":    NewMetric("TotalAlloc", "gauge"),
		"PollCount":     NewMetric("PollCount", "counter"),
		"RandomValue":   NewMetric("RandomValue", "gauge"),
	}
}

// Update updates all metrics.
func (m *Data) Update() {
	m.updateGaugeMetrics()
	m.updateCounterMetrics()
}

func (m *Data) updateGaugeMetrics() {
	metricsToUpdate := *m

	runtimeStats := runtime.MemStats{}
	runtime.ReadMemStats(&runtimeStats)

	reflectStats := reflect.ValueOf(runtimeStats)

	for name, metric := range metricsToUpdate {
		if metric.MType != "gauge" {
			continue
		}
		reflectStatField := reflectStats.FieldByName(name)
		if reflectStatField.Kind() == reflect.Invalid || !reflectStatField.CanInterface() {
			continue
		}
		var valueForUpd float64
		statValueInterface := reflectStatField.Interface()
		switch statValueInterface.(type) {
		case float64:
			valueForUpd, _ = statValueInterface.(float64)
		case float32:
			statValue, _ := statValueInterface.(float32)
			valueForUpd = float64(statValue)
		case uint64:
			statValue, _ := statValueInterface.(uint64)
			valueForUpd = float64(statValue)
		case uint32:
			statValue, _ := statValueInterface.(uint32)
			valueForUpd = float64(statValue)
		case uint16:
			statValue, _ := statValueInterface.(uint16)
			valueForUpd = float64(statValue)
		case uint8:
			statValue, _ := statValueInterface.(uint8)
			valueForUpd = float64(statValue)
		default:
			log.Fatalln("unsupported metric type")
		}
		metric.Value = &valueForUpd
		metricsToUpdate[name] = metric
	}
}

func (m *Data) updateCounterMetrics() {
	metricsToUpdate := *m

	RandomValue := metricsToUpdate["RandomValue"]
	rand.Seed(time.Now().UnixNano())
	newRandomValue := rand.Float64()
	*RandomValue.Value = newRandomValue

	PollCount := metricsToUpdate["PollCount"]
	*PollCount.Delta++
}
