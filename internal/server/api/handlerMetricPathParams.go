package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

func (a *api) handlerUpdateMetricPathParams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Msg("api.handlerUpdateMetricPathParams started")
		defer log.Debug().Msg("api.handlerUpdateMetricPathParams ended")

		reqMetric := &metric.Metrics{}
		if statusCode, err := fillMetricFromPathParams(reqMetric, "update", r.URL.Path); err != nil {
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
			*updatedMetric.Delta += *reqMetric.Delta
		}

		if err = a.service.SetMetric(r.Context(), updatedMetric); err != nil {
			a.error(w, r, http.StatusBadRequest, err)
			return
		}

		if a.cfg.HashKey() != "" {
			if err = updatedMetric.UpdateHash(a.cfg.HashKey()); err != nil {
				log.Error().Err(err).Msg("unable to update metric hash")
			}
		}

		w.Header().Set("Content-Type", "application/json")
		resp, _ := json.Marshal(updatedMetric)
		a.respond(w, http.StatusOK, resp)
	}
}

func (a *api) handlerGetMetricValuePathParams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Msg("api.handlerGetMetricValuePathParams started")
		defer log.Debug().Msg("api.handlerGetMetricValuePathParams ended")

		reqMetric := &metric.Metrics{}
		if statusCode, err := fillMetricFromPathParams(reqMetric, "value", r.URL.Path); err != nil {
			a.error(w, r, statusCode, err)
			return
		}

		metricLocal, isExists, err := a.service.GetMetric(r.Context(), reqMetric.ID)
		if err != nil {
			a.error(w, r, http.StatusBadRequest, err)
			return
		}
		if !isExists {
			a.error(w, r, http.StatusNotFound, err)
			return
		}

		if reqMetric.MType == "gauge" {
			a.respond(w, http.StatusOK, []byte(fmt.Sprintf("%v", *metricLocal.Value)))
		} else if reqMetric.MType == "counter" {
			a.respond(w, http.StatusOK, []byte(fmt.Sprintf("%v", *metricLocal.Delta)))
		}
	}
}
