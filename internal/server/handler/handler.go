package handler

import (
	"fmt"
	"github.com/v1tbrah/metricsAndAlerting/internal/metric"
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

	valM, err := strconv.ParseFloat(infoFromURL.valM, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	metrics.Set(infoFromURL.nameM, valM)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func updateCounterHandler(metrics *metric.Metrics, infoFromURL *infoM, w http.ResponseWriter, r *http.Request) {
	if _, err := strconv.Atoi(infoFromURL.valM); err != nil {
		http.Error(w, "invalid value", http.StatusBadRequest)
	}
	metrics.PollCount++
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

		infoM, err := metrics.Info(infoFromURL.nameM)
		if err != nil {
			http.Error(w, "metric name not specified", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(infoM.ValM()))

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

	if infoM, err := metrics.Info("Alloc"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("BuckHashSys"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("Frees"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("GCCPUFraction"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("GCSys"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("HeapAlloc"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("HeapIdle"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("HeapInuse"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("HeapObjects"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("HeapReleased"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("HeapSys"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("LastGC"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("Lookups"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("MCacheInuse"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("MCacheSys"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("MSpanInuse"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("MSpanSys"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("Mallocs"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("NextGC"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("NumForcedGC"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("NumGC"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("OtherSys"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("PauseTotalNs"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("StackInuse"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("StackSys"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("Sys"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("TotalAlloc"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("PollCount"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}
	if infoM, err := metrics.Info("RandomValue"); err == nil {
		*mForHTML = append(*mForHTML, string(infoM.NameM()+": "+infoM.ValM()))
	}

}
