package memory

import (
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"sync"

	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

type Data struct {
	Metrics map[string]metric.Metrics
	sync.Mutex
}

// NewData returns all metrics.
func NewData() *Data {
	return &Data{
		Metrics: map[string]metric.Metrics{
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

// Update updates all metrics.
func (d *Data) Update(keyForUpdateHash string) {
	d.Lock()
	defer d.Unlock()
	d.updateGaugeMetrics(keyForUpdateHash)
	d.updateCounterMetrics(keyForUpdateHash)
	log.Println("All metrics updated.")
}

func (d *Data) updateGaugeMetrics(keyForUpdateHash string) {
	metricsToUpdate := d.Metrics

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
			log.Fatalln("unsupported metric type")
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
				log.Println(err.Error())
			}
		}
	}

}

func (d *Data) updateCounterMetrics(keyForUpdateHash string) {
	metricsToUpdate := d.Metrics

	PollCount := metricsToUpdate["PollCount"]
	*PollCount.Delta++
	if keyForUpdateHash != "" {
		if err := PollCount.UpdateHash(keyForUpdateHash); err != nil {
			log.Println(err.Error())
		}
	}
	metricsToUpdate["PollCount"] = PollCount
}
