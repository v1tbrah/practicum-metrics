package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/metric"
)

var (
	ErrMetricTypeNotSpecified   = errors.New("metric type not specified")
	ErrMetricTypeNotImplemented = errors.New("metric type not implemented")
	ErrMetricNameNotSpecified   = errors.New("metric name not specified")
	ErrMetricValueNotSpecified  = errors.New("metric value not specified")
)

func checkTypeAndNameMetric(handler string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		infoFromURL := &infoURL{}
		if handler == "update" {
			infoFromURL = newInfoUpdateURL(r.URL.Path)
		} else if handler == "value" {
			infoFromURL = newInfoGetValueURL(r.URL.Path)
		}

		if infoFromURL.typeM == "" {
			http.Error(w, fmt.Sprintf("%s", ErrMetricTypeNotSpecified), http.StatusNotFound)
			return
		}
		if !metric.TypeIsValid(infoFromURL.typeM) {
			http.Error(w, fmt.Sprintf("%s", ErrMetricTypeNotImplemented), http.StatusNotImplemented)
			return
		}
		if infoFromURL.nameM == "" {
			http.Error(w, fmt.Sprintf("%s", ErrMetricNameNotSpecified), http.StatusNotFound)
			return
		}
		if handler != "value" {
			if infoFromURL.valM == "" {
				http.Error(w, fmt.Sprintf("%s", ErrMetricValueNotSpecified), http.StatusNotFound)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
