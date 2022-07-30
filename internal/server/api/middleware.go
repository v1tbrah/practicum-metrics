package api

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipReadHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gz.Close()
			r.Body = gz
		}

		next.ServeHTTP(w, r)
	})
}

func gzipWriteHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

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
		if handler == "update" {
			if infoFromURL.valM == "" {
				http.Error(w, fmt.Sprintf("%s", ErrMetricValueNotSpecified), http.StatusNotFound)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
