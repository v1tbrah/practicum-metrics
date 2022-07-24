package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/metric"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/service"
)

var (
	ErrMetricTypeNotSpecified   = errors.New("metric type not specified")
	ErrMetricTypeNotImplemented = errors.New("metric type not implemented")
	ErrMetricNameNotSpecified   = errors.New("metric name not specified")
)

func (a *api) updateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		infoFromURL := newInfoUpdateURL(r.URL.Path)
		metricFromRequest := metric.Metrics{}
		metricFromRequest.MType = infoFromURL.typeM
		metricFromRequest.ID = infoFromURL.nameM

		switch metricFromRequest.MType {
		case "gauge":
			value, err := strconv.ParseFloat(infoFromURL.valM, 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
				return
			}
			metricFromRequest.Value = &value
			a.updateGaugeHandler(&metricFromRequest, w, r)
		case "counter":
			value, err := strconv.Atoi(infoFromURL.valM)
			if err != nil {
				http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
				return
			}
			valueInt64 := int64(value)
			metricFromRequest.Delta = &valueInt64
			a.updateCounterHandler(&metricFromRequest, w, r)
		}
	}
}

func (a *api) getValueHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		infoFromURL := newInfoUpdateURL(r.URL.Path)
		metricFromRequest := metric.Metrics{}
		metricFromRequest.MType = infoFromURL.typeM
		metricFromRequest.ID = infoFromURL.nameM

		metricOnServ, ok := a.service.MemStorage.Metrics.Load(metricFromRequest.ID)
		if !ok {
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}
		resp, _ := json.Marshal(metricOnServ)
		w.Write(resp)
	}
}

func (a *api) updateJSONHandler(w http.ResponseWriter, r *http.Request) {

	metricFromRequest := metric.Metrics{}
	if !tryFillMetricFromRequest(&metricFromRequest, w, r) {
		return
	}

	switch metricFromRequest.MType {
	case "gauge":
		a.updateGaugeHandler(&metricFromRequest, w, r)
	case "counter":
		a.updateCounterHandler(&metricFromRequest, w, r)
	}
}

func (a *api) getValueJSONHandler(w http.ResponseWriter, r *http.Request) {
	metricFromRequest := metric.Metrics{}
	if !tryFillMetricFromRequest(&metricFromRequest, w, r) {
		return
	}

	metricOnServ, ok := a.service.MemStorage.Metrics.Load(metricFromRequest.ID)
	if !ok {
		http.Error(w, "metric not found", http.StatusNotFound)
		return
	}
	resp, _ := json.Marshal(metricOnServ)
	w.Write(resp)
}

func tryFillMetricFromRequest(fillableMetric *metric.Metrics, w http.ResponseWriter, r *http.Request) bool {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	err = json.Unmarshal(body, &fillableMetric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if fillableMetric.MType == "" {
		http.Error(w, fmt.Sprintf("%s", ErrMetricTypeNotSpecified), http.StatusNotFound)
		return false
	}
	if !fillableMetric.TypeIsValid() {
		http.Error(w, fmt.Sprintf("%s", ErrMetricTypeNotImplemented), http.StatusNotImplemented)
		return false
	}
	if fillableMetric.ID == "" {
		http.Error(w, fmt.Sprintf("%s", ErrMetricNameNotSpecified), http.StatusNotFound)
		return false
	}
	return true
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

func (a *api) updateGaugeHandler(newMetric *metric.Metrics, w http.ResponseWriter, r *http.Request) {

	interfaceMForUpd, ok := a.service.MemStorage.Metrics.Load(newMetric.ID)
	var mForUpd *metric.Metrics
	if !ok {
		mForUpd = &metric.Metrics{}
		mForUpd.ID = newMetric.ID
		mForUpd.MType = newMetric.MType
		var value float64
		mForUpd.Value = &value
	} else {
		mForUpd = interfaceMForUpd.(*metric.Metrics)
	}
	*mForUpd.Value = *newMetric.Value

	a.service.MemStorage.Metrics.Store(mForUpd.ID, mForUpd)

	w.WriteHeader(http.StatusOK)
}

func (a *api) updateCounterHandler(newMetric *metric.Metrics, w http.ResponseWriter, r *http.Request) {

	interfaceMForUpd, ok := a.service.MemStorage.Metrics.Load(newMetric.ID)
	var mForUpd *metric.Metrics
	if !ok {
		mForUpd = &metric.Metrics{}
		mForUpd.ID = newMetric.ID
		mForUpd.MType = newMetric.MType
		var value int64
		mForUpd.Delta = &value
	} else {
		mForUpd = interfaceMForUpd.(*metric.Metrics)
	}

	*mForUpd.Delta = *newMetric.Delta

	a.service.MemStorage.Metrics.Store(mForUpd.ID, mForUpd)

	w.WriteHeader(http.StatusOK)
}

func fillMetricsForPage(dataForPage *[]string, metrics *sync.Map) {
	*dataForPage = append(*dataForPage, sortedMetricsForPage(metrics)...)
}

func sortedMetricsForPage(metrics *sync.Map) []string {
	sortedMetrics := []string{}
	metrics.Range(func(key, value any) bool {
		metric := value.(*metric.Metrics)
		if metric.MType == "gauge" {
			sortedMetrics = append(sortedMetrics, metric.ID+": "+fmt.Sprintf("%f", *metric.Value))
		} else if metric.MType == "counter" {
			sortedMetrics = append(sortedMetrics, metric.ID+": "+fmt.Sprintf("%v", *metric.Delta))
		}
		return true
	})

	sort.Slice(sortedMetrics, func(i, j int) bool { return sortedMetrics[i] < sortedMetrics[j] })
	return sortedMetrics
}
