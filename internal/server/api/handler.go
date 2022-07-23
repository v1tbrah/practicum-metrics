package api

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/metric"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/service"
)

func (a *api) updateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		infoFromURL := newInfoUpdateURL(r.URL.Path)
		switch infoFromURL.typeM {
		case "gauge":
			a.updateGaugeHandler(infoFromURL, w, r)
		case "counter":
			a.updateCounterHandler(infoFromURL, w, r)
		}
	}
}

func (a *api) getValueHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		infoFromURL := newInfoGetValueURL(r.URL.Path)
		metricsOfType, _ := a.service.MemStorage.Metrics.MetricsOfType(infoFromURL.typeM)
		valM, ok := metricsOfType.Load(infoFromURL.nameM)
		if !ok {
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(valM.(string)))

	}
}

func (a *api) getPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dataForPage := service.NewDataForPage()
		fillMetricsForPage(&dataForPage.Metrics, a.service.MemStorage.Metrics)
		allMetrics, err := dataForPage.Page()
		if err != nil {
			log.Fatalln(err)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(allMetrics))
	}
}

func (a *api) updateGaugeHandler(infoFromURL *infoURL, w http.ResponseWriter, r *http.Request) {

	_, err := strconv.ParseFloat(infoFromURL.valM, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	metricsOfType, err := a.service.MemStorage.Metrics.MetricsOfType(infoFromURL.typeM)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}
	metricsOfType.Store(infoFromURL.nameM, infoFromURL.valM)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func (a *api) updateCounterHandler(infoFromURL *infoURL, w http.ResponseWriter, r *http.Request) {

	valM, err := strconv.Atoi(infoFromURL.valM)
	if err != nil {
		http.Error(w, "invalid value", http.StatusBadRequest)
		return
	}

	metricsOfType, err := a.service.MemStorage.Metrics.MetricsOfType(infoFromURL.typeM)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusNotImplemented)
		return
	}

	currVal := 0
	currStrVal, ok := metricsOfType.Load(infoFromURL.nameM)
	if ok {
		currVal, _ = strconv.Atoi(currStrVal.(string))
	}
	metricsOfType.Store(infoFromURL.nameM, strconv.Itoa(valM+currVal))

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func fillMetricsForPage(dataForPage *[]string, metrics *metric.Metrics) {
	if gaugeMetrics, err := metrics.MetricsOfType("gauge"); err == nil {
		*dataForPage = append(*dataForPage, sortedMetricsForPage(gaugeMetrics)...)
	}
	if counterMetrics, err := metrics.MetricsOfType("counter"); err == nil {
		*dataForPage = append(*dataForPage, sortedMetricsForPage(counterMetrics)...)
	}
}

func sortedMetricsForPage(metrics *sync.Map) []string {

	sortedMetrics := []string{}
	i := 0
	metrics.Range(func(key, value any) bool {
		sortedMetrics = append(sortedMetrics, key.(string)+": "+value.(string))
		i++
		return true
	})

	sort.Slice(sortedMetrics, func(i, j int) bool { return sortedMetrics[i] < sortedMetrics[j] })
	return sortedMetrics
}
