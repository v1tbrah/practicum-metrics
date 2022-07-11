package handler

import (
	"fmt"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/api/metric"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/page"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var validUpdateHandlers = map[string]func(metric *metric.Metrics, infoM *infoM, w http.ResponseWriter, r *http.Request){
	"gauge":   updateGaugeHandler,
	"counter": updateCounterHandler,
}

type infoM struct {
	typeM string
	nameM string
	valM  string
}

func newInfoFromUpdateURL(urlPath string) *infoM {
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

func newInfoFromGetValueURL(urlPath string) *infoM {
	newInfoM := infoM{}
	arrInfoM := strings.Split(strings.TrimPrefix(urlPath, "/value/"), "/")
	lenArrInfoM := len(arrInfoM)
	if lenArrInfoM > 0 {
		newInfoM.typeM = arrInfoM[0]
	}
	if lenArrInfoM > 1 {
		newInfoM.nameM = arrInfoM[1]
	}
	return &newInfoM
}

func UpdateHandler(metrics *metric.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, fmt.Sprintf("Method '%s' not allowed", r.Method), http.StatusMethodNotAllowed)
			return
		}
		infoFromURL := newInfoFromUpdateURL(r.URL.Path)
		if infoFromURL.typeM == "" {
			http.Error(w, "metric type not specified", http.StatusNotFound)
			return
		}
		handler, ok := validUpdateHandlers[infoFromURL.typeM]
		if !ok {
			http.Error(w, fmt.Sprintf("metric type: '%s' not implemented", infoFromURL.typeM), http.StatusNotImplemented)
			return
		}
		if infoFromURL.nameM == "" {
			http.Error(w, "metric name not specified", http.StatusNotFound)
			return
		}
		if infoFromURL.valM == "" {
			http.Error(w, "metric value not specified", http.StatusNotFound)
			return
		}

		handler(metrics, infoFromURL, w, r)

	}
}

func updateGaugeHandler(metrics *metric.Metrics, infoFromURL *infoM, w http.ResponseWriter, r *http.Request) {

	_, err := strconv.ParseFloat(infoFromURL.valM, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	metricsOfType, err := metrics.MetricsOfType(infoFromURL.typeM)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}
	metricsOfType[infoFromURL.nameM] = infoFromURL.valM

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func updateCounterHandler(metrics *metric.Metrics, infoFromURL *infoM, w http.ResponseWriter, r *http.Request) {

	valM, err := strconv.Atoi(infoFromURL.valM)
	if err != nil {
		http.Error(w, "invalid value", http.StatusBadRequest)
		return
	}
	metricsOfType, err := metrics.MetricsOfType(infoFromURL.typeM)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}
	currVal := 0
	currStrVal, ok := metricsOfType[infoFromURL.nameM]
	if ok {
		currVal, _ = strconv.Atoi(currStrVal)
	}
	metricsOfType[infoFromURL.nameM] = strconv.Itoa(valM + currVal)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func GetValueHandler(metrics *metric.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		infoFromURL := newInfoFromGetValueURL(r.URL.Path)
		if infoFromURL.nameM == "" {
			http.Error(w, "metric name not specified", http.StatusNotFound)
			return
		}

		valM, err := metrics.MetricOfTypeAndName(infoFromURL.typeM, infoFromURL.nameM)
		if err != nil {
			http.Error(w, "metric name not specified", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(valM))

	}
}

func GetAllMetricsHTML(metrics *metric.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dataForHTML := page.NewData()
		fillMetricsForHTML(&dataForHTML.Metrics, metrics)
		allMetrics, err := dataForHTML.CompletedTpl()
		if err != nil {
			log.Fatalln(err)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(allMetrics))

	}
}

func fillMetricsForHTML(mForHTML *[]string, metrics *metric.Metrics) {

	if gaugeMetrics, err := metrics.MetricsOfType("gauge"); err == nil {
		for nameM, valM := range gaugeMetrics {
			*mForHTML = append(*mForHTML, string(nameM+": "+valM))
		}
	}

	if counterMetrics, err := metrics.MetricsOfType("counter"); err == nil {
		for nameM, valM := range counterMetrics {
			*mForHTML = append(*mForHTML, string(nameM+": "+valM))
		}
	}

}
