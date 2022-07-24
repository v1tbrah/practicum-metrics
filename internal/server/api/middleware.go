package api

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrMetricValueNotSpecified = errors.New("metric value not specified")
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
		if infoFromURL.typeM != "gauge" && infoFromURL.typeM != "counter" {
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
