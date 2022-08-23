package memory

import (
	"errors"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
  "github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

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
			"Alloc":           metric.NewMetric("Alloc", "gauge"),
			"BuckHashSys":     metric.NewMetric("BuckHashSys", "gauge"),
			"Frees":           metric.NewMetric("Frees", "gauge"),
			"GCCPUFraction":   metric.NewMetric("GCCPUFraction", "gauge"),
			"GCSys":           metric.NewMetric("GCSys", "gauge"),
			"HeapAlloc":       metric.NewMetric("HeapAlloc", "gauge"),
			"HeapIdle":        metric.NewMetric("HeapIdle", "gauge"),
			"HeapInuse":       metric.NewMetric("HeapInuse", "gauge"),
			"HeapObjects":     metric.NewMetric("HeapObjects", "gauge"),
			"HeapReleased":    metric.NewMetric("HeapReleased", "gauge"),
			"HeapSys":         metric.NewMetric("HeapSys", "gauge"),
			"LastGC":          metric.NewMetric("LastGC", "gauge"),
			"Lookups":         metric.NewMetric("Lookups", "gauge"),
			"MCacheInuse":     metric.NewMetric("MCacheInuse", "gauge"),
			"MCacheSys":       metric.NewMetric("MCacheSys", "gauge"),
			"MSpanInuse":      metric.NewMetric("MSpanInuse", "gauge"),
			"MSpanSys":        metric.NewMetric("MSpanSys", "gauge"),
			"Mallocs":         metric.NewMetric("Mallocs", "gauge"),
			"NextGC":          metric.NewMetric("NextGC", "gauge"),
			"NumForcedGC":     metric.NewMetric("NumForcedGC", "gauge"),
			"NumGC":           metric.NewMetric("NumGC", "gauge"),
			"OtherSys":        metric.NewMetric("OtherSys", "gauge"),
			"PauseTotalNs":    metric.NewMetric("PauseTotalNs", "gauge"),
			"StackInuse":      metric.NewMetric("StackInuse", "gauge"),
			"StackSys":        metric.NewMetric("StackSys", "gauge"),
			"Sys":             metric.NewMetric("Sys", "gauge"),
			"TotalAlloc":      metric.NewMetric("TotalAlloc", "gauge"),
			"PollCount":       metric.NewMetric("PollCount", "counter"),
			"RandomValue":     metric.NewMetric("RandomValue", "gauge"),
			"TotalMemory":     metric.NewMetric("TotalMemory", "gauge"),
			"FreeMemory":      metric.NewMetric("FreeMemory", "gauge"),
			"CPUutilization1": metric.NewMetric("CPUutilization1", "gauge"),
		},
	}
}

// GetData return copy of data.
func (d *Memory) GetData() map[string]metric.Metrics {
	log.Debug().Msg("memory.GetData started")
	defer log.Debug().Msg("memory.GetData ended")

	d.mu.RLock()
	defer d.mu.RUnlock()

	result := make(map[string]metric.Metrics, len(d.data))

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

// UpdateBasic updates basic metrics.
func (d *Memory) UpdateBasic(keyForUpdateHash string) {
	log.Debug().Msg("memory.UpdateBasic started")
	defer log.Debug().Msg("memory.UpdateBasic ended")

	d.mu.Lock()
	defer d.mu.Unlock()

	d.updateBasicGaugeMetrics(keyForUpdateHash)
	d.updateBasicCounterMetrics(keyForUpdateHash)
}

// UpdateBasic updates additional metrics.
func (d *Memory) UpdateAdditional(keyForUpdateHash string) error {
	log.Debug().Msg("memory.UpdateAdditional started")
	defer log.Debug().Msg("memory.UpdateAdditional ended")

	d.mu.Lock()
	defer d.mu.Unlock()

	vms, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	totalMemoryForUpd := d.data["TotalMemory"]
	currTotalMemory := float64(vms.Total)
	totalMemoryForUpd.Value = &currTotalMemory
	d.data["TotalMemory"] = totalMemoryForUpd

	freeMemoryForUpd := d.data["FreeMemory"]
	currFreeMemory := float64(vms.Free)
	freeMemoryForUpd.Value = &currFreeMemory
	d.data["FreeMemory"] = freeMemoryForUpd

	CPUUtilForUpd := d.data["CPUutilization1"]
	currCPUUtil, err := cpu.Percent(time.Second*0, false)
	if err != nil {
		return err
	}
	CPUUtilForUpd.Value = &currCPUUtil[0]
	d.data["CPUutilization1"] = CPUUtilForUpd

	if keyForUpdateHash != "" {
		if err := totalMemoryForUpd.UpdateHash(keyForUpdateHash); err != nil {
			log.Error().Err(err).Str("metric", totalMemoryForUpd.String()).Msg("unable to computing hash")
		}
		if err := freeMemoryForUpd.UpdateHash(keyForUpdateHash); err != nil {
			log.Error().Err(err).Str("metric", freeMemoryForUpd.String()).Msg("unable to computing hash")
		}
		if err := CPUUtilForUpd.UpdateHash(keyForUpdateHash); err != nil {
			log.Error().Err(err).Str("metric", CPUUtilForUpd.String()).Msg("unable to computing hash")
		}
	}

	return err
}

func (d *Memory) updateBasicGaugeMetrics(keyForUpdateHash string) {
	log.Debug().Msg("memory.updateBasicGaugeMetrics started")
	defer log.Debug().Msg("memory.updateBasicGaugeMetrics ended")

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
			log.Error().Err(errors.New("unsupported type of metric")).Msg("unable to update metric")
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
				log.Error().Err(err).Str("metric", currMetric.String()).Msg("unable to computing hash")
			}
		}
	}

}

func (d *Memory) updateBasicCounterMetrics(keyForUpdateHash string) {
	log.Debug().Msg("memory.updateBasicCounterMetrics started")
	defer log.Debug().Msg("memory.updateBasicCounterMetrics ended")

	metricsToUpdate := d.data

	PollCount := metricsToUpdate["PollCount"]
	*PollCount.Delta++
	if keyForUpdateHash != "" {
		if err := PollCount.UpdateHash(keyForUpdateHash); err != nil {
			log.Error().Err(err).Str("metric", PollCount.String()).Msg("unable to computing hash")
		}
	}
	metricsToUpdate["PollCount"] = PollCount
}
