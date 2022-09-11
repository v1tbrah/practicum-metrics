package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

func (a *api) handlerUpdateMetric(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("api.handlerUpdateMetric started")
	defer log.Debug().Msg("api.handlerUpdateMetric ended")

	reqMetric := &metric.Metrics{}
	if statusCode, err := fillMetricFromRequestBody(reqMetric, r.Body); err != nil {
		a.error(w, r, statusCode, err)
		return
	}
	if statusCode, err := checkValidMetricFromRequest(reqMetric, "update"); err != nil {
		a.error(w, r, statusCode, err)
		return
	}

	updMetric, isExists, err := a.service.GetMetric(r.Context(), reqMetric.ID)
	if err != nil {
		a.error(w, r, http.StatusBadRequest, err)
		return
	}
	if !isExists {
		updMetric = metric.NewMetric(reqMetric.ID, reqMetric.MType)
	}
	if updMetric.MType == "gauge" {
		*updMetric.Value = *reqMetric.Value
	} else if updMetric.MType == "counter" {
		*updMetric.Delta += *reqMetric.Delta
	}

	if err = a.service.SetMetric(r.Context(), updMetric); err != nil {
		a.error(w, r, http.StatusBadRequest, err)
		return
	}

	if a.cfg.HashKey() != "" {
		if err = updMetric.UpdateHash(a.cfg.HashKey()); err != nil {
			log.Error().Err(err).Msg("unable to update metric hash")
		}
	}

	resp, _ := json.Marshal(updMetric)
	a.respond(w, http.StatusOK, resp)
}

func (a *api) handlerUpdateListMetrics(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("api.handlerUpdateListMetrics started")
	defer log.Debug().Msg("api.handlerUpdateListMetrics ended")

	reqListMetrics := []metric.Metrics{}
	if statusCode, err := fillListMetricsFromRequestBody(&reqListMetrics, r.Body); err != nil {
		a.error(w, r, statusCode, err)
		return
	}

	updatedListMetrics := []metric.Metrics{}
	counterMetrics := map[string]metric.Metrics{}
	for _, reqMetric := range reqListMetrics {
		if statusCode, err := checkValidMetricFromRequest(&reqMetric, "update"); err != nil {
			a.error(w, r, statusCode, err)
			return
		}
		updatedMetric, isExists, err := a.service.GetMetric(r.Context(), reqMetric.ID)
		if err != nil {
			a.error(w, r, http.StatusBadRequest, err)
			return
		}
		if !isExists {
			updatedMetric = metric.NewMetric(reqMetric.ID, reqMetric.MType)
		}
		if updatedMetric.MType == "gauge" {
			*updatedMetric.Value = *reqMetric.Value
		} else if updatedMetric.MType == "counter" {
			counterMetric, ok := counterMetrics[updatedMetric.ID]
			if ok {
				updatedMetric = counterMetric
			}
			*updatedMetric.Delta += *reqMetric.Delta
			counterMetrics[updatedMetric.ID] = updatedMetric
		}

		if a.cfg.HashKey() != "" {
			if err = updatedMetric.UpdateHash(a.cfg.HashKey()); err != nil {
				log.Error().Err(err).Msg("unable to update metric hash")
			}
		}

		updatedListMetrics = append(updatedListMetrics, updatedMetric)
	}

	if err := a.service.SetListMetrics(r.Context(), updatedListMetrics); err != nil {
		a.error(w, r, http.StatusBadRequest, err)
		return
	}

	a.respond(w, http.StatusOK, []byte{})
}

func (a *api) handlerGetMetricValue(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("api.handlerGetMetricValue started")
	defer log.Debug().Msg("api.handlerGetMetricValue ended")

	reqMetric := &metric.Metrics{}
	if statusCode, err := fillMetricFromRequestBody(reqMetric, r.Body); err != nil {
		a.error(w, r, statusCode, err)
		return
	}
	if statusCode, err := checkValidMetricFromRequest(reqMetric, "value"); err != nil {
		a.error(w, r, statusCode, err)
		return
	}

	respMetric, isExists, err := a.service.GetMetric(r.Context(), reqMetric.ID)
	if err != nil {
		a.error(w, r, http.StatusBadRequest, err)
		return
	} else if !isExists {
		a.error(w, r, http.StatusNotFound, err)
		return
	}

	if a.cfg.HashKey() != "" {
		if err = respMetric.UpdateHash(a.cfg.HashKey()); err != nil {
			log.Error().Err(err).Msg("unable to update metric hash")
		}
	}

	w.Header().Set("Content-Type", "application/json")
	resp, _ := json.Marshal(respMetric)
	a.respond(w, http.StatusOK, resp)
}

func checkValidMetricFromRequest(metric *metric.Metrics, requestType string) (int, error) {
	log.Debug().
		Str("metric", metric.String()).
		Str("requestType", requestType).
		Msg("api.checkValidMetricFromRequest started")
	defer log.Debug().Msg("api.checkValidMetricFromRequest ended")

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
		Str("metric", metric.String()).
		Str("requestBody", fmt.Sprint(requestBody)).
		Msg("api.fillMetricFromRequestBody started")
	defer log.Debug().Msg("api.fillMetricFromRequestBody ended")

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
	defer log.Debug().Msg("api.fillListMetricsFromRequestBody ended")

	body, err := io.ReadAll(requestBody)
	if err != nil && err != io.EOF {
		return http.StatusBadRequest, err
	}
	if err = json.Unmarshal(body, listMetrics); err != nil {
		return http.StatusBadRequest, err
	}
	return 0, nil
}
