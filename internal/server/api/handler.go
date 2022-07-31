package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/service"
)

var (
	ErrMetricTypeNotSpecified   = errors.New("metric type not specified")
	ErrMetricTypeNotImplemented = errors.New("metric type not implemented")
	ErrMetricNameNotSpecified   = errors.New("metric name not specified")
	ErrMetricNotFound           = errors.New("metric not found")
	ErrMetricValueNotSpecified  = errors.New("metric value not specified")
)

func (a *api) updateMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricFromRequest := &repo.Metrics{}
	if err, statusCode := fillMetricFromRequestBody(metricFromRequest, r.Body); err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}
	if err, statusCode := checkValidMetricFromRequest(metricFromRequest, "update"); err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	switch metricFromRequest.MType {
	case "gauge":
		a.updateGaugeMetric(metricFromRequest, w, r)
	case "counter":
		a.updateCounterMetric(metricFromRequest, w, r)
	default:
		http.Error(w, ErrMetricTypeNotImplemented.Error(), http.StatusNotImplemented)
		return
	}
}

func (a *api) getMetricValueHandler(w http.ResponseWriter, r *http.Request) {
	metricFromRequest := &repo.Metrics{}
	if err, statusCode := fillMetricFromRequestBody(metricFromRequest, r.Body); err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}
	if err, statusCode := checkValidMetricFromRequest(metricFromRequest, "value"); err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	metricLocal, ok := a.service.MemStorage.Data.Load(metricFromRequest.ID)
	if !ok {
		http.Error(w, ErrMetricNotFound.Error(), http.StatusNotFound)
		return
	}
	resp, _ := json.Marshal(metricLocal)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.Write(resp)
}

func fillMetricFromRequestBody(metric *repo.Metrics, requestBody io.ReadCloser) (error, int) {
	body, err := io.ReadAll(requestBody)
	if err != nil && err != io.EOF {
		return errors.New("err body reading"), http.StatusBadRequest
	}
	if err = json.Unmarshal(body, metric); err != nil {
		return errors.New("invalid json"), http.StatusBadRequest
	}
	return nil, 0
}

func checkValidMetricFromRequest(metric *repo.Metrics, requestType string) (error, int) {
	if metric.MType == "" {
		return ErrMetricTypeNotSpecified, http.StatusNotFound
	}
	if !metric.TypeIsValid() {
		return ErrMetricTypeNotImplemented, http.StatusNotImplemented
	}
	if metric.ID == "" {
		return ErrMetricNameNotSpecified, http.StatusNotFound
	}
	if requestType == "update" {
		if metric.MType == "gauge" && metric.Value == nil {
			if metric.Value == nil {
				return ErrMetricValueNotSpecified, http.StatusNotFound
			}
		} else if metric.MType == "counter" && metric.Delta == nil {
			return ErrMetricValueNotSpecified, http.StatusNotFound
		}
	}
	return nil, 0
}

func (a *api) updateGaugeMetric(newMetric *repo.Metrics, w http.ResponseWriter, r *http.Request) {

	interfaceMForUpd, ok := a.service.MemStorage.Data.Load(newMetric.ID)
	var mForUpd repo.Metrics
	if !ok {
		mForUpd = repo.Metrics{}
		mForUpd.ID = newMetric.ID
		mForUpd.MType = newMetric.MType
		var value float64
		mForUpd.Value = &value
	} else {
		mForUpd = interfaceMForUpd.(repo.Metrics)
	}
	*mForUpd.Value = *newMetric.Value

	a.service.MemStorage.Data.Store(mForUpd.ID, mForUpd)

	w.WriteHeader(http.StatusOK)
}

func (a *api) updateCounterMetric(newMetric *repo.Metrics, w http.ResponseWriter, r *http.Request) {

	interfaceMForUpd, ok := a.service.MemStorage.Data.Load(newMetric.ID)
	var mForUpd repo.Metrics
	if !ok {
		mForUpd = repo.Metrics{}
		mForUpd.ID = newMetric.ID
		mForUpd.MType = newMetric.MType
		var value int64
		mForUpd.Delta = &value
	} else {
		mForUpd = interfaceMForUpd.(repo.Metrics)
	}

	*mForUpd.Delta += *newMetric.Delta

	a.service.MemStorage.Data.Store(mForUpd.ID, mForUpd)

	w.WriteHeader(http.StatusOK)
}

func (a *api) getPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dataForPage := service.NewDataForPage()
		fillMetricsForPage(&dataForPage.Metrics, a.service.MemStorage.Data)
		page, err := dataForPage.Page()
		if err != nil {
			log.Fatalln(err)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(page))
	}
}

func fillMetricsForPage(dataForPage *[]string, metrics *repo.Data) {
	*dataForPage = append(*dataForPage, sortedMetricsForPage(metrics)...)
}

func sortedMetricsForPage(metrics *repo.Data) []string {
	sortedMetrics := []string{}
	metrics.Range(func(key, value any) bool {
		currMetric := value.(repo.Metrics)
		if currMetric.MType == "gauge" {
			sortedMetrics = append(sortedMetrics, currMetric.ID+": "+fmt.Sprintf("%f", *currMetric.Value))
		} else if currMetric.MType == "counter" {
			sortedMetrics = append(sortedMetrics, currMetric.ID+": "+fmt.Sprintf("%v", *currMetric.Delta))
		}
		return true
	})

	sort.Slice(sortedMetrics, func(i, j int) bool { return sortedMetrics[i] < sortedMetrics[j] })
	return sortedMetrics
}
