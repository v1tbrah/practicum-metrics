package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

func fillMetricFromPathParams(metric *metric.Metrics, handlerType, path string) (int, error) {
	log.Debug().
		Str("metric", metric.String()).
		Str("handlerType", handlerType).
		Str("path", path).
		Msg("api.fillMetricFromPathParams started")
	log.Debug().Msg("api.fillMetricFromPathParams ended")

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
	log.Debug().
		Str("metric", fmt.Sprint(metric)).
		Str("pathInfo", fmt.Sprint(pathInfo)).
		Msg("api.fillMetricValueFromPathInfo started")
	log.Debug().Msg("api.fillMetricValueFromPathInfo ended")

	if pathInfo.valM == "" {
		return http.StatusNotFound, ErrMetricValueNotSpecified
	}

	if metric.MType == "gauge" {
		value, err := strconv.ParseFloat(pathInfo.valM, 64)
		if err != nil {
			return http.StatusBadRequest, err
		}
		metric.Value = &value
	} else if metric.MType == "counter" {
		value, err := strconv.Atoi(pathInfo.valM)
		if err != nil {
			return http.StatusBadRequest, err
		}
		valueInt64 := int64(value)
		metric.Delta = &valueInt64
	} else {
		return 0, ErrMetricTypeNotImplemented
	}
	return 0, nil
}
