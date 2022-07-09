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

func UpdateHandler(metric *metric.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, fmt.Sprintf("Method '%s' not allowed", r.Method), http.StatusMethodNotAllowed)
			return
		}
		infoM := strings.Split(strings.TrimPrefix(r.URL.Path, "/update/"), "/")
		infoMLen := len(infoM)
		if infoMLen == 0 {
			http.Error(w, "Have no metric type", http.StatusNotFound)
			return
		}
		if infoMLen == 1 {
			http.Error(w, "Have no metric name", http.StatusNotFound)
			return
		}
		if infoMLen == 2 {
			http.Error(w, "Have no metric value", http.StatusNotFound)
			return
		}
		if handler, ok := supportedHandlers[infoM[0]]; ok {
			handler(metric, infoM[1:], w, r)
			return
		}
		http.Error(w, fmt.Sprintf("Metric type: '%s' unsupported", infoM[0]), http.StatusNotFound)
	}
}

func gaugeHandler(metric *metric.Metrics, infoM []string, w http.ResponseWriter, r *http.Request) {

	strValM := infoM[1]
	if strValM == "" {
		http.Error(w, "have no value", http.StatusNotFound)
		return
	}
	valM, err := strconv.ParseFloat(strValM, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	nameM := infoM[0]
	metric.Set(nameM, valM)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func counterHandler(metric *metric.Metrics, infoM []string, w http.ResponseWriter, r *http.Request) {
	if _, err := strconv.Atoi(infoM[1]); err != nil {
		http.Error(w, fmt.Sprintf("invalid value"), http.StatusBadRequest)
	}
	metric.PollCount++
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}
