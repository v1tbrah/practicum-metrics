package memory

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

type Memory struct {
	data map[string]metric.Metrics
	mu   sync.RWMutex
}

// New returns all metrics.
func New() *Memory {
	log.Debug().Msg("memory.New started")
	defer log.Debug().Msg("memory.New ended")

	return &Memory{
		data: map[string]metric.Metrics{
			"Alloc":         metric.NewMetric("Alloc", "gauge"),
			"BuckHashSys":   metric.NewMetric("BuckHashSys", "gauge"),
			"Frees":         metric.NewMetric("Frees", "gauge"),
			"GCCPUFraction": metric.NewMetric("GCCPUFraction", "gauge"),
			"GCSys":         metric.NewMetric("GCSys", "gauge"),
			"HeapAlloc":     metric.NewMetric("HeapAlloc", "gauge"),
			"HeapIdle":      metric.NewMetric("HeapIdle", "gauge"),
			"HeapInuse":     metric.NewMetric("HeapInuse", "gauge"),
			"HeapObjects":   metric.NewMetric("HeapObjects", "gauge"),
			"HeapReleased":  metric.NewMetric("HeapReleased", "gauge"),
			"HeapSys":       metric.NewMetric("HeapSys", "gauge"),
			"LastGC":        metric.NewMetric("LastGC", "gauge"),
			"Lookups":       metric.NewMetric("Lookups", "gauge"),
			"MCacheInuse":   metric.NewMetric("MCacheInuse", "gauge"),
			"MCacheSys":     metric.NewMetric("MCacheSys", "gauge"),
			"MSpanInuse":    metric.NewMetric("MSpanInuse", "gauge"),
			"MSpanSys":      metric.NewMetric("MSpanSys", "gauge"),
			"Mallocs":       metric.NewMetric("Mallocs", "gauge"),
			"NextGC":        metric.NewMetric("NextGC", "gauge"),
			"NumForcedGC":   metric.NewMetric("NumForcedGC", "gauge"),
			"NumGC":         metric.NewMetric("NumGC", "gauge"),
			"OtherSys":      metric.NewMetric("OtherSys", "gauge"),
			"PauseTotalNs":  metric.NewMetric("PauseTotalNs", "gauge"),
			"StackInuse":    metric.NewMetric("StackInuse", "gauge"),
			"StackSys":      metric.NewMetric("StackSys", "gauge"),
			"Sys":           metric.NewMetric("Sys", "gauge"),
			"TotalAlloc":    metric.NewMetric("TotalAlloc", "gauge"),
			"PollCount":     metric.NewMetric("PollCount", "counter"),
			"RandomValue":   metric.NewMetric("RandomValue", "gauge"),
		},
	}
}

// GetData return copy of data.
func (d *Memory) GetData() (map[string]metric.Metrics) {
	log.Debug().Msg("memory.GetData started")
	defer log.Debug().Msg("memory.GetData ended")

	result := map[string]metric.Metrics{}

	d.mu.RLock()
	defer d.mu.RUnlock()

	for _, currMetric := range d.data {
		metricForResult := metric.NewMetric(currMetric.ID, currMetric.MType)
		if currMetric.MType == "gauge" {
			valueForResult := *currMetric.Value
			metricForResult.Value = &valueForResult
		} else if currMetric.MType == "counter" {
			deltaForResult := *currMetric.Delta
			metricForResult.Delta = &deltaForResult
		}
		result[metricForResult.ID] = metricForResult
	}

	return result
}

// Update updates all metrics.
func (d *Memory) Update(keyForUpdateHash string) {
	log.Debug().Msg("memory.Update started")
	defer log.Debug().Msg("memory.Update ended")

	d.mu.Lock()
	defer d.mu.Unlock()

	d.updateGaugeMetrics(keyForUpdateHash)
	d.updateCounterMetrics(keyForUpdateHash)
}

func (d *Memory) updateGaugeMetrics(keyForUpdateHash string) {
	log.Debug().Msg("memory.updateGaugeMetrics started")
	defer log.Debug().Msg("memory.updateGaugeMetrics ended")

	metricsToUpdate := d.data

	runtimeStats := runtime.MemStats{}
	runtime.ReadMemStats(&runtimeStats)

	reflectStats := reflect.ValueOf(runtimeStats)

	for name, currMetric := range metricsToUpdate {
		if currMetric.MType != "gauge" {
			continue
		}
		reflectStatField := reflectStats.FieldByName(name)
		if reflectStatField.Kind() == reflect.Invalid || !reflectStatField.CanInterface() {
			continue
		}
		var valueForUpd float64
		switch statValue := reflectStatField.Interface().(type) {
		case float64:
			valueForUpd = statValue
		case float32:
			valueForUpd = float64(statValue)
		case uint64:
			valueForUpd = float64(statValue)
		case uint32:
			valueForUpd = float64(statValue)
		case uint16:
			valueForUpd = float64(statValue)
		case uint8:
			valueForUpd = float64(statValue)
		default:
			log.Error().
				Err(errors.New("unsupported type of metric")).
				Str("MID", currMetric.ID).
				Str("MType", fmt.Sprint(statValue)).
				Msg("unable to update metric")
		}
		currMetric.Value = &valueForUpd
		metricsToUpdate[name] = currMetric
	}

	randomValue := metricsToUpdate["RandomValue"]
	newRandomValue := rand.Float64()
	*randomValue.Value = newRandomValue
	metricsToUpdate["RandomValue"] = randomValue

	for _, currMetric := range metricsToUpdate {
		if keyForUpdateHash != "" {
			if err := currMetric.UpdateHash(keyForUpdateHash); err != nil {
				log.Error().
					Err(err).
					Str("MID", currMetric.ID).
					Msg("unable to computing hash")
			}
		}
	}

}

func (d *Memory) updateCounterMetrics(keyForUpdateHash string) {
	log.Debug().Msg("memory.updateCounterMetrics started")
	defer log.Debug().Msg("memory.updateCounterMetrics ended")

	metricsToUpdate := d.data

	PollCount := metricsToUpdate["PollCount"]
	*PollCount.Delta++
	if keyForUpdateHash != "" {
		if err := PollCount.UpdateHash(keyForUpdateHash); err != nil {
			log.Error().
				Err(err).
				Str("MID", PollCount.ID).
				Msg("unable to computing hash")
		}
	}
	metricsToUpdate["PollCount"] = PollCount
}
