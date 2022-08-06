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
	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

var (
	ErrMetricTypeNotSpecified   = errors.New("metric type not specified")
	ErrMetricTypeNotImplemented = errors.New("metric type not implemented")
	ErrMetricNameNotSpecified   = errors.New("metric name not specified")
	ErrMetricNotFound           = errors.New("metric not found")
	ErrMetricValueNotSpecified  = errors.New("metric value not specified")
)

func (a *api) updateMetricHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	metricFromRequest := &metric.Metrics{}
	if statusCode, err := fillMetricFromRequestBody(metricFromRequest, r.Body); err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}
	if statusCode, err := a.checkValidMetricFromRequest(metricFromRequest, "update"); err != nil {
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
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	metricFromRequest := &metric.Metrics{}
	if statusCode, err := fillMetricFromRequestBody(metricFromRequest, r.Body); err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}
	if statusCode, err := a.checkValidMetricFromRequest(metricFromRequest, "value"); err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	metricForResponse, ok := a.service.MemStorage.Data.Metrics[metricFromRequest.ID]
	if !ok {
		http.Error(w, ErrMetricNotFound.Error(), http.StatusNotFound)
		return
	}

	if a.service.Cfg.Key != "" {
		if err := metricForResponse.UpdateHash(a.service.Cfg.Key); err != nil {
			log.Println(err)
		}
	}

	resp, _ := json.Marshal(metricForResponse)
	w.Write(resp)
}

func fillMetricFromRequestBody(metric *metric.Metrics, requestBody io.ReadCloser) (int, error) {
	body, err := io.ReadAll(requestBody)
	if err != nil && err != io.EOF {
		return http.StatusBadRequest, errors.New("err body reading")
	}
	if err = json.Unmarshal(body, metric); err != nil {
		return http.StatusBadRequest, errors.New("invalid json")
	}
	return 0, nil
}

func (a *api) checkValidMetricFromRequest(metric *metric.Metrics, requestType string) (int, error) {
	if metric.MType == "" {
		return http.StatusNotFound, ErrMetricTypeNotSpecified
	}
	if !metric.TypeIsValid() {
		return http.StatusNotImplemented, ErrMetricTypeNotImplemented
	}
	if metric.ID == "" {
		return http.StatusNotFound, ErrMetricNameNotSpecified
	}

	if requestType == "update" {
		if metric.MType == "gauge" && metric.Value == nil {
			return http.StatusNotFound, ErrMetricValueNotSpecified
		} else if metric.MType == "counter" && metric.Delta == nil {
			return http.StatusNotFound, ErrMetricValueNotSpecified
		}
	}

	if a.service.Cfg.Key != "" {
		hashFromRequest := metric.Hash
		metric.UpdateHash(a.service.Cfg.Key)
		newHash := metric.Hash
		if hashFromRequest != newHash {
			return http.StatusBadRequest, errors.New("invalid hash")
		}
	}

	return 0, nil
}

func (a *api) updateGaugeMetric(newMetric *metric.Metrics, w http.ResponseWriter, r *http.Request) {

	a.service.MemStorage.Data.Lock()
	defer a.service.MemStorage.Data.Unlock()

	metricForUpd, ok := a.service.MemStorage.Data.Metrics[newMetric.ID]
	if !ok {
		metricForUpd = metric.Metrics{}
		metricForUpd.ID = newMetric.ID
		metricForUpd.MType = newMetric.MType
		var value float64
		metricForUpd.Value = &value
	}
	*metricForUpd.Value = *newMetric.Value

	a.service.MemStorage.Data.Metrics[metricForUpd.ID] = metricForUpd

	if a.service.Cfg.Key != "" {
		if err := metricForUpd.UpdateHash(a.service.Cfg.Key); err != nil {
			log.Println(err)
		}
	}

	resp, _ := json.Marshal(metricForUpd)
	w.Write(resp)
}

func (a *api) updateCounterMetric(newMetric *metric.Metrics, w http.ResponseWriter, r *http.Request) {

	a.service.MemStorage.Data.Lock()
	defer a.service.MemStorage.Data.Unlock()

	metricForUpd, ok := a.service.MemStorage.Data.Metrics[newMetric.ID]
	if !ok {
		metricForUpd = metric.Metrics{}
		metricForUpd.ID = newMetric.ID
		metricForUpd.MType = newMetric.MType
		var value int64
		metricForUpd.Delta = &value
	}

	*metricForUpd.Delta += *newMetric.Delta

	a.service.MemStorage.Data.Metrics[metricForUpd.ID] = metricForUpd

	if a.service.Cfg.Key != "" {
		if err := metricForUpd.UpdateHash(a.service.Cfg.Key); err != nil {
			log.Println(err)
		}
	}

	resp, _ := json.Marshal(metricForUpd)
	w.Write(resp)
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
	for _, currMetric := range metrics.Metrics {
		if currMetric.MType == "gauge" {
			sortedMetrics = append(sortedMetrics, currMetric.ID+": "+fmt.Sprintf("%f", *currMetric.Value))
		} else if currMetric.MType == "counter" {
			sortedMetrics = append(sortedMetrics, currMetric.ID+": "+fmt.Sprintf("%v", *currMetric.Delta))
		}
	}

	sort.Slice(sortedMetrics, func(i, j int) bool { return sortedMetrics[i] < sortedMetrics[j] })
	return sortedMetrics
}
