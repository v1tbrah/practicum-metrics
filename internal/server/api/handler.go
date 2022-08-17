package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/model"
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

func (a *api) updateListMetricsHandler(w http.ResponseWriter, r *http.Request) {

	listMetricsFromRequest := []metric.Metrics{}
	if statusCode, err := fillListMetricsFromRequestBody(&listMetricsFromRequest, r.Body); err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	listMetricsForUpdate := []metric.Metrics{}
	counterMetrics := map[string]metric.Metrics{}
	for _, metricFromRequest := range listMetricsFromRequest {
		if statusCode, err := a.checkValidMetricFromRequest(&metricFromRequest, "update"); err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}
		metricForUpd, isExists, err := a.service.Storage.GetMetric(metricFromRequest.ID)
		if err != nil {
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

		if a.service.Cfg.Key != "" {
			if err = metricForUpd.UpdateHash(a.service.Cfg.Key); err != nil {
				log.Println(err)
			}
		}
		listMetricsForUpdate = append(listMetricsForUpdate, metricForUpd)
	}

	if err := a.service.Storage.SetListMetrics(listMetricsForUpdate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(""))
}

func (a *api) getMetricValueHandler(w http.ResponseWriter, r *http.Request) {
	metricFromRequest := &metric.Metrics{}
	if statusCode, err := fillMetricFromRequestBody(metricFromRequest, r.Body); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), statusCode)
		return
	}
	if statusCode, err := a.checkValidMetricFromRequest(metricFromRequest, "value"); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), statusCode)
		return
	}

	metricForResponse, ok, err := a.service.Storage.GetMetric(metricFromRequest.ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if !ok {
		log.Println(ErrMetricNotFound.Error())
		http.Error(w, ErrMetricNotFound.Error(), http.StatusNotFound)
		return
	}

	if a.service.Cfg.Key != "" {
		if err = metricForResponse.UpdateHash(a.service.Cfg.Key); err != nil {
			log.Println(err)
		}
	}

	resp, _ := json.Marshal(metricForResponse)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (a *api) checkDBConnHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbPool, err := pgxpool.Connect(context.Background(), a.service.Cfg.PgConnString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Write([]byte(""))
			defer dbPool.Close()
		}
	}
}

func fillMetricFromRequestBody(metric *metric.Metrics, requestBody io.ReadCloser) (int, error) {
	body, err := io.ReadAll(requestBody)
	if err != nil && err != io.EOF {
		return http.StatusBadRequest, err
	}
	if err = json.Unmarshal(body, metric); err != nil {
		log.Println("DEBUG #1", err)
		return http.StatusBadRequest, err
	}
	return 0, nil
}

func fillListMetricsFromRequestBody(listMetrics *[]metric.Metrics, requestBody io.ReadCloser) (int, error) {
	body, err := io.ReadAll(requestBody)
	if err != nil && err != io.EOF {
		return http.StatusBadRequest, err
	}
	if err = json.Unmarshal(body, listMetrics); err != nil {
		log.Println("DEBUG #2", err)
		return http.StatusBadRequest, err
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

	return 0, nil
}

func (a *api) updateGaugeMetric(newMetric *metric.Metrics, w http.ResponseWriter, r *http.Request) {

	metricForUpd, ok, err := a.service.Storage.GetMetric(newMetric.ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !ok {
		metricForUpd = metric.NewMetric(newMetric.ID, newMetric.MType)
		var value float64
		metricForUpd.Value = &value
	}
	*metricForUpd.Value = *newMetric.Value

	if a.service.Cfg.Key != "" {
		if err = metricForUpd.UpdateHash(a.service.Cfg.Key); err != nil {
			log.Println(err)
		}
	}

	if err = a.service.Storage.SetMetric(metricForUpd.ID, metricForUpd); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, _ := json.Marshal(metricForUpd)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (a *api) updateCounterMetric(newMetric *metric.Metrics, w http.ResponseWriter, r *http.Request) {

	metricForUpd, ok, err := a.service.Storage.GetMetric(newMetric.ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !ok {
		metricForUpd = metric.NewMetric(newMetric.ID, newMetric.MType)
		var delta int64
		metricForUpd.Delta = &delta
	}
	*metricForUpd.Delta += *newMetric.Delta

	if a.service.Cfg.Key != "" {
		if err = metricForUpd.UpdateHash(a.service.Cfg.Key); err != nil {
			log.Println(err)
		}
	}

	if err = a.service.Storage.SetMetric(metricForUpd.ID, metricForUpd); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, _ := json.Marshal(metricForUpd)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (a *api) getPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dataForPage := service.NewDataForPage()
		dataMetrics, err := a.service.Storage.GetData()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fillMetricsForPage(&dataForPage.Metrics, dataMetrics)
		page, err := dataForPage.Page()
		if err != nil {
			log.Fatalln(err)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(page))
	}
}

func fillMetricsForPage(dataForPage *[]string, metrics *model.Data) {
	*dataForPage = append(*dataForPage, sortedMetricsForPage(metrics)...)
}

func sortedMetricsForPage(metrics *model.Data) []string {
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
