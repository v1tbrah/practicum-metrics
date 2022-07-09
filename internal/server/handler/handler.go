package handler

import (
	"fmt"
	"github.com/v1tbrah/metricsAndAlerting/internal/metric"
	"net/http"
	"strconv"
	"strings"
)

var supportedHandlers = map[string]func(metric *metric.Metrics, infoM []string, w http.ResponseWriter, r *http.Request){
	"gauge":   gaugeHandler,
	"counter": counterHandler,
}

var supportedGaugeMetrics = map[string]struct{}{
	"Alloc":         struct{}{},
	"BuckHashSys":   struct{}{},
	"Frees":         struct{}{},
	"GCCPUFraction": struct{}{},
	"HeapAlloc":     struct{}{},
	"HeapIdle":      struct{}{},
	"HeapInuse":     struct{}{},
	"HeapObjects":   struct{}{},
	"HeapReleased":  struct{}{},
	"HeapSys":       struct{}{},
	"LastGC":        struct{}{},
	"Lookups":       struct{}{},
	"MCacheInuse":   struct{}{},
	"MCacheSys":     struct{}{},
	"MSpanInuse":    struct{}{},
	"MSpanSys":      struct{}{},
	"Mallocs":       struct{}{},
	"NextGC":        struct{}{},
	"NumForcedGC":   struct{}{},
	"NumGC":         struct{}{},
	"OtherSys":      struct{}{},
	"PauseTotalNs":  struct{}{},
	"StackInuse":    struct{}{},
	"StackSys":      struct{}{},
	"Sys":           struct{}{},
	"TotalAlloc":    struct{}{},
	"RandomValue":   struct{}{},
}

var supportedCounterMetrics = map[string]struct{}{
	"PollCount": struct{}{},
}

func UpdateHandler(metric *metric.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, fmt.Sprintf("Method '%s' not allowed", r.Method), http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, fmt.Sprintf("Content-type '%s' unsupported", r.Header.Get("Content-Type")), http.StatusUnsupportedMediaType)
			return
		}
		infoM := strings.Split(strings.TrimPrefix(r.URL.Path, "/update/"), "/")
		infoMLen := len(infoM)
		if infoMLen == 0 {
			http.Error(w, "Have no metric type", http.StatusBadRequest)
			return
		}
		if infoMLen == 1 {
			http.Error(w, "Have no metric name", http.StatusBadRequest)
			return
		}
		if infoMLen == 2 {
			http.Error(w, "Have no metric value", http.StatusBadRequest)
			return
		}
		if handler, ok := supportedHandlers[infoM[0]]; ok {
			handler(metric, infoM[1:], w, r)
			return
		}
		http.Error(w, fmt.Sprintf("Metric type: '%s' unsupported", infoM[0]), http.StatusBadRequest)
	}
}

func gaugeHandler(metric *metric.Metrics, infoM []string, w http.ResponseWriter, r *http.Request) {
	nameM := infoM[0]
	if _, ok := supportedGaugeMetrics[nameM]; !ok {
		http.Error(w, fmt.Sprintf("Metric name: '%s' unsupported", nameM), http.StatusBadRequest)
		return
	}
	valM, err := strconv.ParseFloat(infoM[1], 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}
	err = metric.Set(nameM, valM)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func counterHandler(metric *metric.Metrics, infoM []string, w http.ResponseWriter, r *http.Request) {
	nameM := infoM[0]
	if _, ok := supportedCounterMetrics[nameM]; !ok {
		http.Error(w, fmt.Sprintf("Metric name: '%s' unsupported", nameM), http.StatusBadRequest)
		return
	}
	metric.PollCount++
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}
