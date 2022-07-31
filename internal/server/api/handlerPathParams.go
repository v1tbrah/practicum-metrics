package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo"
)

func (a *api) updateMetricHandlerPathParams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricFromRequest := &repo.Metrics{}
		if err, statusCode := fillMetricFromPathParams(metricFromRequest, "update", r.URL.Path); err != nil {
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
		metricFromRequest := &repo.Metrics{}
		if err, statusCode := fillMetricFromPathParams(metricFromRequest, "value", r.URL.Path); err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		metricLocalInterface, ok := a.service.MemStorage.Data.Load(metricFromRequest.ID)
		if !ok {
			http.Error(w, ErrMetricNotFound.Error(), http.StatusNotFound)
			return
		}

		metricLocal := metricLocalInterface.(repo.Metrics)
		if metricFromRequest.MType == "gauge" {
			w.Write([]byte(fmt.Sprintf("%f", *metricLocal.Value)))
		} else if metricFromRequest.MType == "counter" {
			w.Write([]byte(fmt.Sprintf("%d", *metricLocal.Delta)))
		}
	}
}

func fillMetricFromPathParams(metric *repo.Metrics, handlerType, path string) (error, int) {
	var pathInfo *pathInfo
	if handlerType == "update" {
		pathInfo = newInfoUpdateURL(path)
	} else if handlerType == "value" {
		pathInfo = newInfoGetValueURL(path)
	}
	if pathInfo.typeM == "" {
		return ErrMetricTypeNotSpecified, http.StatusNotFound
	}
	metric.MType = pathInfo.typeM
	if !metric.TypeIsValid() {
		return ErrMetricTypeNotImplemented, http.StatusNotImplemented
	}
	if pathInfo.nameM == "" {
		return ErrMetricNameNotSpecified, http.StatusNotFound
	}

	metric.ID = pathInfo.nameM

	if handlerType == "update" {
		if err, httpStatusCode := fillMetricValueFromPathInfo(metric, pathInfo); err != nil {
			return err, httpStatusCode
		}
	}

	return nil, 0
}

func fillMetricValueFromPathInfo(metric *repo.Metrics, pathInfo *pathInfo) (error, int) {
	if pathInfo.valM == "" {
		return ErrMetricValueNotSpecified, http.StatusNotFound
	}

	if metric.MType == "gauge" {
		value, err := strconv.ParseFloat(pathInfo.valM, 64)
		if err != nil {
			return err, http.StatusBadRequest
		}
		metric.Value = &value
	} else if metric.MType == "counter" {
		value, err := strconv.Atoi(pathInfo.valM)
		if err != nil {
			return err, http.StatusBadRequest
		}
		valueInt64 := int64(value)
		metric.Delta = &valueInt64
	} else {
		return ErrMetricTypeNotImplemented, 0
	}
	return nil, 0
}
