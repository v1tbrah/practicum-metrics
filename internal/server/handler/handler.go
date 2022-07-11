package handler

import (
	"fmt"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/api/metric"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/page"
	"log"
	"net/http"
	"strconv"
)

var validUpdateHandlers = map[string]func(metric *metric.Metrics, infoM *metric.Info, w http.ResponseWriter, r *http.Request){
	"gauge":   updateGaugeHandler,
	"counter": updateCounterHandler,
}

func UpdateHandler(metrics *metric.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, fmt.Sprintf("Method '%s' not allowed", r.Method), http.StatusMethodNotAllowed)
			return
		}
		infoFromURL := metric.NewInfoFromUpdateURL(r.URL.Path)
		if infoFromURL.TypeM() == "" {
			http.Error(w, "metric type not specified", http.StatusNotFound)
			return
		}
		handler, ok := validUpdateHandlers[infoFromURL.TypeM()]
		if !ok {
			http.Error(w, fmt.Sprintf("metric type: '%s' not implemented", infoFromURL.TypeM()), http.StatusNotImplemented)
			return
		}
		if infoFromURL.NameM() == "" {
			http.Error(w, "metric name not specified", http.StatusNotFound)
			return
		}
		if infoFromURL.ValM() == "" {
			http.Error(w, "metric value not specified", http.StatusNotFound)
			return
		}

		handler(metrics, infoFromURL, w, r)

	}
}

func updateGaugeHandler(metrics *metric.Metrics, infoFromURL *metric.Info, w http.ResponseWriter, r *http.Request) {

	_, err := strconv.ParseFloat(infoFromURL.ValM(), 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	metricsOfType, err := metrics.MetricsOfType(infoFromURL.TypeM())
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}
	metricsOfType[infoFromURL.NameM()] = infoFromURL.ValM()

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func updateCounterHandler(metrics *metric.Metrics, infoFromURL *metric.Info, w http.ResponseWriter, r *http.Request) {

	valM, err := strconv.Atoi(infoFromURL.ValM())
	if err != nil {
		http.Error(w, "invalid value", http.StatusBadRequest)
		return
	}
	metricsOfType, err := metrics.MetricsOfType(infoFromURL.TypeM())
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusNotImplemented)
		return
	}
	currVal := 0
	currStrVal, ok := metricsOfType[infoFromURL.NameM()]
	if ok {
		currVal, _ = strconv.Atoi(currStrVal)
	}
	metricsOfType[infoFromURL.NameM()] = strconv.Itoa(valM + currVal)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func GetValueHandler(metrics *metric.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		infoFromURL := metric.NewInfoFromGetValueURL(r.URL.Path)
		if infoFromURL.TypeM() == "" {
			http.Error(w, "metric type not specified", http.StatusNotFound)
			return
		}
		if infoFromURL.NameM() == "" {
			http.Error(w, "metric name not specified", http.StatusNotFound)
			return
		}

		metricsOfType, err := metrics.MetricsOfType(infoFromURL.TypeM())
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusNotImplemented)
			return
		}

		valM, ok := metricsOfType[infoFromURL.NameM()]
		if !ok {
			http.Error(w, "metric not found", http.StatusNotFound)
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
