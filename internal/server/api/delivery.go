package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/pg"
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

func (a *api) getPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Msg("api.getPageHandler started")
		log.Debug().Msg("api.getPageHandler ended")

		dataForPage := service.NewDataForPage()
		dataMetrics, err := a.service.Storage.GetData(r.Context())
		if err != nil {
			log.Error().
				Err(err).
				Str("func", "Storage.GetData").
				Msg("unable to get data from storage")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fillMetricsForPage(&dataForPage.Metrics, dataMetrics)
		page, err := dataForPage.Page()
		if err != nil {
			log.Error().
				Err(err).
				Str("func", "dataForPage.Page").
				Msg("unable to return page")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(page))
	}
}

func (a *api) checkDBConnHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Msg("api.checkDBConnHandler started")
		log.Debug().Msg("api.checkDBConnHandler ended")

		dbStorage, ok := a.service.Storage.(*pg.PgStorage)
		if !ok {
			http.Error(w, "type of storage is not DB storage", http.StatusBadRequest)
			return
		}

		err := dbStorage.Ping(r.Context())
		if err != nil {
			log.Error().
				Err(err).
				Str("func", "pg.Ping").
				Msg("unable to connect to DB")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("OK"))
	}
}

func (a *api) updateMetricHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("api.updateMetricHandler started")
	log.Debug().Msg("api.updateMetricHandler ended")

	metricFromRequest := &metric.Metrics{}
	if statusCode, err := fillMetricFromRequestBody(metricFromRequest, r.Body); err != nil {
		log.Error().
			Err(err).
			Str("func", "fillMetricFromRequestBody").
			Str("metricFromRequest", fmt.Sprint(metricFromRequest)).
			Str("body", fmt.Sprint(r.Body)).
			Msg("unable to fill metric from request body")
		http.Error(w, err.Error(), statusCode)
		return
	}
	if statusCode, err := checkValidMetricFromRequest(metricFromRequest, "update"); err != nil {
		log.Error().
			Err(err).
			Str("func", "checkValidMetricFromRequest").
			Str("metricFromRequest", fmt.Sprint(metricFromRequest)).
			Str("requestType", "update").
			Msg("unable to check valid metric")
		http.Error(w, err.Error(), statusCode)
		return
	}

	metricForUpd, isExists, err := a.service.Storage.GetMetric(r.Context(), metricFromRequest.ID)
	if err != nil {
		log.Error().
			Err(err).
			Str("func", "Storage.GetMetric").
			Str("MID", metricFromRequest.ID).
			Msg("unable get metric from storage")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !isExists {
		metricForUpd = metric.NewMetric(metricFromRequest.ID, metricFromRequest.MType)
	}
	if metricForUpd.MType == "gauge" {
		*metricForUpd.Value = *metricFromRequest.Value
	} else if metricForUpd.MType == "counter" {
		*metricForUpd.Delta += *metricFromRequest.Delta
	}

	if err = a.service.Storage.SetMetric(r.Context(), metricForUpd); err != nil {
		log.Error().
			Err(err).
			Str("func", "Storage.SetMetric").
			Str("metricForUpd", fmt.Sprint(metricForUpd)).
			Msg("unable set metric to storage")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if a.service.Cfg.HashKey != "" {
		if err = metricForUpd.UpdateHash(a.service.Cfg.HashKey); err != nil {
			log.Error().
				Err(err).
				Str("func", "metricForUpd.UpdateHash").
				Msg("unable to update metric hash")
		}
	}

	resp, _ := json.Marshal(metricForUpd)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)

}

func (a *api) updateListMetricsHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("api.updateListMetricsHandler started")
	log.Debug().Msg("api.updateListMetricsHandler ended")

	listMetricsFromRequest := []metric.Metrics{}
	if statusCode, err := fillListMetricsFromRequestBody(&listMetricsFromRequest, r.Body); err != nil {
		log.Error().
			Err(err).
			Str("func", "fillListMetricsFromRequestBody").
			Str("listMetricsFromRequest", fmt.Sprint(listMetricsFromRequest)).
			Str("body", fmt.Sprint(r.Body)).
			Msg("unable to fill list metrics from request body")
		http.Error(w, err.Error(), statusCode)
		return
	}

	listMetricsForUpdate := []metric.Metrics{}
	counterMetrics := map[string]metric.Metrics{}
	for _, metricFromRequest := range listMetricsFromRequest {
		if statusCode, err := checkValidMetricFromRequest(&metricFromRequest, "update"); err != nil {
			log.Error().
				Err(err).
				Str("func", "checkValidMetricFromRequest").
				Str("metricFromRequest", fmt.Sprint(metricFromRequest)).
				Str("requestType", "update").
				Msg("unable to check valid metric")
			http.Error(w, err.Error(), statusCode)
			return
		}
		metricForUpd, isExists, err := a.service.Storage.GetMetric(r.Context(), metricFromRequest.ID)
		if err != nil {
			log.Error().
				Err(err).
				Str("func", "Storage.GetMetric").
				Str("MID", metricFromRequest.ID).
				Msg("unable get metric from storage")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !isExists {
			metricForUpd = metric.NewMetric(metricFromRequest.ID, metricFromRequest.MType)
		}
		if metricForUpd.MType == "gauge" {
			*metricForUpd.Value = *metricFromRequest.Value
		} else if metricForUpd.MType == "counter" {
			counterMetric, ok := counterMetrics[metricForUpd.ID]
			if ok {
				metricForUpd = counterMetric
			}
			*metricForUpd.Delta += *metricFromRequest.Delta
			counterMetrics[metricForUpd.ID] = metricForUpd
		}

		if a.service.Cfg.HashKey != "" {
			if err = metricForUpd.UpdateHash(a.service.Cfg.HashKey); err != nil {
				log.Error().
					Err(err).
					Str("func", "metricForUpd.UpdateHash").
					Msg("unable to update metric hash")
			}
		}
		listMetricsForUpdate = append(listMetricsForUpdate, metricForUpd)
	}

	if err := a.service.Storage.SetListMetrics(r.Context(), listMetricsForUpdate); err != nil {
		log.Error().
			Err(err).
			Str("func", "Storage.SetListMetrics").
			Str("listMetricsForUpdate", fmt.Sprint(listMetricsForUpdate)).
			Msg("unable to update list metrics")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(""))
}

func (a *api) getMetricValueHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("api.getMetricValueHandler started")
	log.Debug().Msg("api.getMetricValueHandler ended")

	metricFromRequest := &metric.Metrics{}
	if statusCode, err := fillMetricFromRequestBody(metricFromRequest, r.Body); err != nil {
		log.Error().
			Err(err).
			Str("func", "fillMetricFromRequestBody").
			Str("metricFromRequest", fmt.Sprint(metricFromRequest)).
			Str("body", fmt.Sprint(r.Body)).
			Msg("unable to fill metric from request body")
		http.Error(w, err.Error(), statusCode)
		return
	}
	if statusCode, err := checkValidMetricFromRequest(metricFromRequest, "value"); err != nil {
		log.Error().
			Err(err).
			Str("func", "checkValidMetricFromRequest").
			Str("metricFromRequest", fmt.Sprint(metricFromRequest)).
			Str("requestType", "update").
			Msg("unable to check valid metric")
		http.Error(w, err.Error(), statusCode)
		return
	}

	metricForResponse, ok, err := a.service.Storage.GetMetric(r.Context(), metricFromRequest.ID)
	if err != nil {
		log.Error().
			Err(err).
			Str("func", "Storage.GetMetric").
			Str("MID", metricFromRequest.ID).
			Msg("unable get metric from storage")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if !ok {
		log.Info().
			Str("func", "Storage.GetMetric").
			Str("MID", metricFromRequest.ID).
			Msg("not found metric")
		http.Error(w, ErrMetricNotFound.Error(), http.StatusNotFound)
		return
	}

	if a.service.Cfg.HashKey != "" {
		if err = metricForResponse.UpdateHash(a.service.Cfg.HashKey); err != nil {
			log.Error().
				Err(err).
				Str("func", "metricForUpd.UpdateHash").
				Msg("unable to update metric hash")
		}
	}

	resp, _ := json.Marshal(metricForResponse)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (a *api) updateMetricHandlerPathParams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Msg("api.updateMetricHandlerPathParams started")
		log.Debug().Msg("api.updateMetricHandlerPathParams ended")

		metricFromRequest := &metric.Metrics{}
		if statusCode, err := fillMetricFromPathParams(metricFromRequest, "update", r.URL.Path); err != nil {
			log.Error().
				Err(err).
				Str("func", "fillMetricFromPathParams").
				Str("metricFromRequest", fmt.Sprint(metricFromRequest)).
				Str("handlerType", "update").
				Str("urlPath", r.URL.Path).
				Msg("unable to fill metric from path params")
			http.Error(w, err.Error(), statusCode)
			return
		}

		metricForUpd, isExists, err := a.service.Storage.GetMetric(r.Context(), metricFromRequest.ID)
		if err != nil {
			log.Error().
				Err(err).
				Str("func", "Storage.GetMetric").
				Str("MID", metricFromRequest.ID).
				Msg("unable get metric from storage")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !isExists {
			metricForUpd = metric.NewMetric(metricFromRequest.ID, metricFromRequest.MType)
		}
		if metricForUpd.MType == "gauge" {
			*metricForUpd.Value = *metricFromRequest.Value
		} else if metricForUpd.MType == "counter" {
			*metricForUpd.Delta += *metricFromRequest.Delta
		}

		if err = a.service.Storage.SetMetric(r.Context(), metricForUpd); err != nil {
			log.Error().
				Err(err).
				Str("func", "Storage.SetMetric").
				Str("metricForUpd", fmt.Sprint(metricForUpd)).
				Msg("unable set metric to storage")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if a.service.Cfg.HashKey != "" {
			if err = metricForUpd.UpdateHash(a.service.Cfg.HashKey); err != nil {
				log.Error().
					Err(err).
					Str("func", "metricForUpd.UpdateHash").
					Msg("unable to update metric hash")
			}
		}

		resp, _ := json.Marshal(metricForUpd)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func (a *api) getMetricValueHandlerPathParams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Msg("api.getMetricValueHandlerPathParams started")
		log.Debug().Msg("api.getMetricValueHandlerPathParams ended")

		metricFromRequest := &metric.Metrics{}
		if statusCode, err := fillMetricFromPathParams(metricFromRequest, "value", r.URL.Path); err != nil {
			log.Error().
				Err(err).
				Str("func", "fillMetricFromPathParams").
				Str("metricFromRequest", fmt.Sprint(metricFromRequest)).
				Str("handlerType", "update").
				Str("urlPath", r.URL.Path).
				Msg("unable to fill metric from path params")
			http.Error(w, err.Error(), statusCode)
			return
		}

		metricLocal, ok, err := a.service.Storage.GetMetric(r.Context(), metricFromRequest.ID)
		if err != nil {
			log.Error().
				Err(err).
				Str("func", "Storage.GetMetric").
				Str("MID", metricFromRequest.ID).
				Msg("unable get metric from storage")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !ok {
			log.Info().
				Str("func", "Storage.GetMetric").
				Str("MID", metricFromRequest.ID).
				Msg("not found metric")
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
