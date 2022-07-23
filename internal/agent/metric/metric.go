package metric

import (
	"errors"
	"math/rand"
	"runtime"
	"time"
)

type (
	gauge   = float64
	counter = int64
)

var ErrIsNotAMetric = errors.New("it's not a metric")

type Metrics struct {
	Alloc         gauge
	BuckHashSys   gauge
	Frees         gauge
	GCCPUFraction gauge
	GCSys         gauge
	HeapAlloc     gauge
	HeapIdle      gauge
	HeapInuse     gauge
	HeapObjects   gauge
	HeapReleased  gauge
	HeapSys       gauge
	LastGC        gauge
	Lookups       gauge
	MCacheInuse   gauge
	MCacheSys     gauge
	MSpanInuse    gauge
	MSpanSys      gauge
	Mallocs       gauge
	NextGC        gauge
	NumForcedGC   gauge
	NumGC         gauge
	OtherSys      gauge
	PauseTotalNs  gauge
	StackInuse    gauge
	StackSys      gauge
	Sys           gauge
	TotalAlloc    gauge
	PollCount     counter
	RandomValue   gauge
}

// Creates a Metrics.
func New() *Metrics {
	return &Metrics{}
}

// Updates metrics. Gauge metrics are read from runtime.ReadMemStats.
func (m *Metrics) Update() {
	m.updateGauge()
	m.PollCount++
}

func (m *Metrics) updateGauge() {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	m.Alloc = gauge(stats.Alloc)
	m.BuckHashSys = gauge(stats.BuckHashSys)
	m.Frees = gauge(stats.Frees)
	m.GCCPUFraction = stats.GCCPUFraction
	m.GCSys = gauge(stats.GCSys)
	m.HeapAlloc = gauge(stats.HeapAlloc)
	m.HeapIdle = gauge(stats.HeapIdle)
	m.HeapInuse = gauge(stats.HeapInuse)
	m.HeapObjects = gauge(stats.HeapObjects)
	m.HeapReleased = gauge(stats.HeapReleased)
	m.HeapSys = gauge(stats.HeapSys)
	m.LastGC = gauge(stats.LastGC)
	m.Lookups = gauge(stats.Lookups)
	m.MCacheInuse = gauge(stats.MCacheInuse)
	m.MSpanSys = gauge(stats.MSpanSys)
	m.MSpanInuse = gauge(stats.MSpanInuse)
	m.Mallocs = gauge(stats.Mallocs)
	m.NextGC = gauge(stats.NextGC)
	m.NumForcedGC = gauge(stats.NumForcedGC)
	m.NumGC = gauge(stats.NumGC)
	m.OtherSys = gauge(stats.OtherSys)
	m.PauseTotalNs = gauge(stats.PauseTotalNs)
	m.StackInuse = gauge(stats.StackInuse)
	m.StackSys = gauge(stats.StackSys)
	m.Sys = gauge(stats.Sys)

	rand.Seed(time.Now().UnixNano())
	m.RandomValue = rand.Float64()
}
