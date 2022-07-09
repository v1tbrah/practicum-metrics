package handler

import (
	"fmt"
	"github.com/v1tbrah/metricsAndAlerting/internal/metric"
	"net/http"
	"strconv"
	"strings"
)

type infoM struct {
	typeM string
	nameM string
	valM  string
}

func newInfoM(urlPath string) *infoM {
	newInfoM := infoM{}
	arrInfoM := strings.Split(strings.TrimPrefix(urlPath, "/update/"), "/")
	lenArrInfoM := len(arrInfoM)
	if lenArrInfoM > 0 {
		newInfoM.typeM = arrInfoM[0]
	}
	if lenArrInfoM > 1 {
		newInfoM.nameM = arrInfoM[1]
	}
	if lenArrInfoM > 2 {
		newInfoM.valM = arrInfoM[2]
	}
	return &newInfoM
}

var supportedHandlers = map[string]func(metric *metric.Metrics, infoM *infoM, w http.ResponseWriter, r *http.Request){
	"gauge":   gaugeHandler,
	"counter": counterHandler,
}

func UpdateHandler(metric *metric.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, fmt.Sprintf("Method '%s' not allowed", r.Method), http.StatusMethodNotAllowed)
			return
		}
		infoM := newInfoM(r.URL.Path)
		if infoM.typeM == "" {
			http.Error(w, "metric type not specified", http.StatusNotFound)
			return
		}
		handler, ok := supportedHandlers[infoM.typeM]
		if !ok {
			http.Error(w, fmt.Sprintf("metric type: '%s' not implemented", infoM.typeM), http.StatusNotImplemented)
			return
		}
		if infoM.nameM == "" {
			http.Error(w, "metric name not specified", http.StatusNotFound)
			return
		}
		if infoM.valM == "" {
			http.Error(w, "metric value not specified", http.StatusNotFound)
			return
		}

		handler(metric, infoM, w, r)

	}
}

func gaugeHandler(metric *metric.Metrics, infoM *infoM, w http.ResponseWriter, r *http.Request) {

	valM, err := strconv.ParseFloat(infoM.valM, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	metric.Set(infoM.nameM, valM)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func counterHandler(metric *metric.Metrics, infoM *infoM, w http.ResponseWriter, r *http.Request) {
	if _, err := strconv.Atoi(infoM.valM); err != nil {
		http.Error(w, "invalid value", http.StatusBadRequest)
	}
	metric.PollCount++
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}
