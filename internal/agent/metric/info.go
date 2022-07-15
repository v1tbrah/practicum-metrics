package metric

import "fmt"

type Info struct {
	typeM string
	nameM string
	valM  string
}

// Returns Info about a metric by metric name.
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

// Return metric type.
func (i *Info) TypeM() string {
	return i.typeM
}

// Return metric name.
func (i *Info) NameM() string {
	return i.nameM
}

// Return metric value.
func (i *Info) ValM() string {
	return i.valM
}
