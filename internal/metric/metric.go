package metric

import (
	"errors"
	"fmt"
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

func New() *Metrics {
	return &Metrics{}
}

func (m *Metrics) Set (name string, val gauge)  error {
	switch name {
	case "Alloc":
		m.Alloc = val
	case "BuckHashSys":
		m.BuckHashSys = val
	case "Frees":
		m.Frees = val
	case "GCCPUFraction":
		m.GCCPUFraction = val
	case "GCSys":
		m.GCSys = val
	case "HeapAlloc":
		m.HeapAlloc = val
	case "HeapIdle":
		m.HeapIdle = val
	case "HeapInuse":
		m.HeapInuse = val
	case "HeapObjects":
		m.HeapObjects = val
	case "HeapReleased":
		m.HeapReleased = val
	case "HeapSys":
		m.HeapSys = val
	case "LastGC":
		m.LastGC = val
	case "Lookups":
		m.Lookups = val
	case "MCacheInuse":
		m.MCacheInuse = val
	case "MCacheSys":
		m.MCacheSys = val
	case "MSpanInuse":
		m.MSpanInuse = val
	case "MSpanSys":
		m.MSpanSys = val
	case "Mallocs":
		m.Mallocs = val
	case "NextGC":
		m.NextGC = val
	case "NumForcedGC":
		m.NumForcedGC = val
	case "NumGC":
		m.NumGC = val
	case "OtherSys":
		m.OtherSys = val
	case "PauseTotalNs":
		m.PauseTotalNs = val
	case "StackInuse":
		m.StackInuse = val
	case "StackSys":
		m.StackSys = val
	case "Sys":
		m.Sys = val
	case "TotalAlloc":
		m.TotalAlloc = val
	case "RandomValue":
		m.RandomValue = val
	default:
		return ErrIsNotAMetric
	}
	return nil
}

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

type Info struct {
	typeM string
	nameM string
	valM  string
}

func (i *Info) TypeM() string {
	return i.typeM
}

func (i *Info) NameM() string {
	return i.nameM
}

func (i *Info) ValM() string {
	return i.valM
}

func (m *Metrics) Info(name string) (Info, error) {
	info := Info{nameM: name}
	switch name {
	case "Alloc":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.Alloc)
	case "BuckHashSys":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.BuckHashSys)
	case "Frees":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.Frees)
	case "GCCPUFraction":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.GCCPUFraction)
	case "GCSys":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.GCSys)
	case "HeapAlloc":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.HeapAlloc)
	case "HeapIdle":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.HeapIdle)
	case "HeapInuse":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.HeapInuse)
	case "HeapObjects":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.HeapObjects)
	case "HeapReleased":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.HeapReleased)
	case "HeapSys":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.HeapSys)
	case "LastGC":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.LastGC)
	case "Lookups":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.Lookups)
	case "MCacheInuse":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.MCacheInuse)
	case "MCacheSys":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.MCacheSys)
	case "MSpanInuse":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.MSpanInuse)
	case "MSpanSys":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.MSpanSys)
	case "Mallocs":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.Mallocs)
	case "NextGC":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.NextGC)
	case "NumForcedGC":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.NumForcedGC)
	case "NumGC":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.NumGC)
	case "OtherSys":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.OtherSys)
	case "PauseTotalNs":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.PauseTotalNs)
	case "StackInuse":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.StackInuse)
	case "StackSys":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.StackSys)
	case "Sys":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.Sys)
	case "TotalAlloc":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.TotalAlloc)
	case "RandomValue":
		info.typeM = "gauge"
		info.valM = fmt.Sprintf("%f", m.RandomValue)
	case "PollCount":
		info.typeM = "counter"
		info.valM = fmt.Sprintf("%d", m.PollCount)
	default:
		return info, ErrIsNotAMetric
	}
	return info, nil
}
