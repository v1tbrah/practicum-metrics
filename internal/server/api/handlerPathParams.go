package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

func (a *api) updateMetricHandlerPathParams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricFromRequest := &metric.Metrics{}
		if statusCode, err := fillMetricFromPathParams(metricFromRequest, "update", r.URL.Path); err != nil {
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
}

func (a *api) getMetricValueHandlerPathParams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricFromRequest := &metric.Metrics{}
		if statusCode, err := fillMetricFromPathParams(metricFromRequest, "value", r.URL.Path); err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		metricLocal, ok := a.service.MemStorage.Data.Metrics[metricFromRequest.ID]
		if !ok {
			http.Error(w, ErrMetricNotFound.Error(), http.StatusNotFound)
			return
		}

		if metricFromRequest.MType == "gauge" {
			w.Write([]byte(fmt.Sprintf("%v", *metricLocal.Value)))
		} else if metricFromRequest.MType == "counter" {
			w.Write([]byte(fmt.Sprintf("%v", *metricLocal.Delta)))
		}
	}
}

func fillMetricFromPathParams(metric *metric.Metrics, handlerType, path string) (int, error) {
	var pathInfo *pathInfo
	if handlerType == "update" {
		pathInfo = newInfoUpdateURL(path)
	} else if handlerType == "value" {
		pathInfo = newInfoGetValueURL(path)
	}
	if pathInfo.typeM == "" {
		return http.StatusNotFound, ErrMetricTypeNotSpecified
	}
	metric.MType = pathInfo.typeM
	if !metric.TypeIsValid() {
		return http.StatusNotImplemented, ErrMetricTypeNotImplemented
	}
	if pathInfo.nameM == "" {
		return http.StatusNotFound, ErrMetricNameNotSpecified
	}

	metric.ID = pathInfo.nameM

	if handlerType == "update" {
		if httpStatusCode, err := fillMetricValueFromPathInfo(metric, pathInfo); err != nil {
			return httpStatusCode, err
		}
	}

	return 0, nil
}

func fillMetricValueFromPathInfo(metric *metric.Metrics, pathInfo *pathInfo) (int, error) {
	if pathInfo.valM == "" {
		return http.StatusNotFound, ErrMetricValueNotSpecified
	}

	if metric.MType == "gauge" {
		value, err := strconv.ParseFloat(pathInfo.valM, 64)
		if err != nil {
			log.Println("ВОТ ЗДЕСЬ #3 ВЕРНУЛСЯ", err)
			return http.StatusBadRequest, err
		}
		metric.Value = &value
	} else if metric.MType == "counter" {
		value, err := strconv.Atoi(pathInfo.valM)
		if err != nil {
			log.Println("ВОТ ЗДЕСЬ #4 ВЕРНУЛСЯ", err)
			return http.StatusBadRequest, err
		}
		valueInt64 := int64(value)
		metric.Delta = &valueInt64
	} else {
		return 0, ErrMetricTypeNotImplemented
	}
	return 0, nil
}
