package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

func checkValidMetricFromRequest(metric *metric.Metrics, requestType string) (int, error) {
	log.Debug().
		Str("metric", fmt.Sprint(metric)).
		Str("requestType", requestType).
		Msg("api.checkValidMetricFromRequest started")
	log.Debug().Msg("api.checkValidMetricFromRequest ended")

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

	return 0, nil
}

func fillMetricFromRequestBody(metric *metric.Metrics, requestBody io.ReadCloser) (int, error) {
	log.Debug().
		Str("metric", fmt.Sprint(metric)).
		Str("requestBody", fmt.Sprint(requestBody)).
		Msg("api.fillMetricFromRequestBody started")
	log.Debug().Msg("api.fillMetricFromRequestBody ended")

	body, err := io.ReadAll(requestBody)
	if err != nil && err != io.EOF {
		return http.StatusBadRequest, err
	}
	if err = json.Unmarshal(body, metric); err != nil {
		return http.StatusBadRequest, err
	}
	return 0, nil
}

func fillListMetricsFromRequestBody(listMetrics *[]metric.Metrics, requestBody io.ReadCloser) (int, error) {
	log.Debug().
		Str("listMetrics", fmt.Sprint(listMetrics)).
		Str("requestBody", fmt.Sprint(requestBody)).
		Msg("api.fillListMetricsFromRequestBody started")
	log.Debug().Msg("api.fillListMetricsFromRequestBody ended")

	body, err := io.ReadAll(requestBody)
	if err != nil && err != io.EOF {
		return http.StatusBadRequest, err
	}
	if err = json.Unmarshal(body, listMetrics); err != nil {
		return http.StatusBadRequest, err
	}
	return 0, nil
}
