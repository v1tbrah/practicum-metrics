package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

type pathInfo struct {
	typeM string
	nameM string
	valM  string
}

func newInfoUpdateURL(urlPath string) *pathInfo {
	log.Debug().Str("urlPath", urlPath).Msg("api.newInfoUpdateURL started")
	defer log.Debug().Msg("api.newInfoUpdateURL ended")

	newInfoM := pathInfo{}
	arrInfoM := strings.Split(strings.TrimPrefix(urlPath, "/update/"), "/")
	lenArrInfoM := len(arrInfoM)
	if lenArrInfoM > 0 {
		newInfoM.typeM = arrInfoM[0]
	}
	if lenArrInfoM > 1 {
		newInfoM.nameM = arrInfoM[1]
	}
	if lenArrInfoM > 2 {
		newInfoM.valM = arrInfoM[2]
	}
	return &newInfoM
}

func newInfoGetValueURL(urlPath string) *pathInfo {
	log.Debug().Str("urlPath", urlPath).Msg("api.newInfoGetValueURL started")
	defer log.Debug().Msg("api.newInfoGetValueURL ended")

	newInfoM := pathInfo{}
	arrInfoM := strings.Split(strings.TrimPrefix(urlPath, "/value/"), "/")
	lenArrInfoM := len(arrInfoM)
	if lenArrInfoM > 0 {
		newInfoM.typeM = arrInfoM[0]
	}
	if lenArrInfoM > 1 {
		newInfoM.nameM = arrInfoM[1]
	}
	return &newInfoM
}

func fillMetricFromPathParams(metric *metric.Metrics, handlerType, path string) (int, error) {
	log.Debug().
		Str("metric", metric.String()).
		Str("handlerType", handlerType).
		Str("path", path).
		Msg("api.fillMetricFromPathParams started")
	defer log.Debug().Msg("api.fillMetricFromPathParams ended")

	var pInf *pathInfo
	if handlerType == "update" {
		pInf = newInfoUpdateURL(path)
	} else if handlerType == "value" {
		pInf = newInfoGetValueURL(path)
	}
	if pInf.typeM == "" {
		return http.StatusNotFound, ErrMetricTypeNotSpecified
	}
	metric.MType = pInf.typeM
	if !metric.TypeIsValid() {
		return http.StatusNotImplemented, ErrMetricTypeNotImplemented
	}
	if pInf.nameM == "" {
		return http.StatusNotFound, ErrMetricNameNotSpecified
	}

	metric.ID = pInf.nameM

	if handlerType == "update" {
		if httpStatusCode, err := fillMetricValueFromPathInfo(metric, pInf); err != nil {
			return httpStatusCode, err
		}
	}

	return 0, nil
}

func fillMetricValueFromPathInfo(metric *metric.Metrics, pInf *pathInfo) (int, error) {
	log.Debug().
		Str("metric", metric.String()).
		Str("pathInfo", fmt.Sprint(pInf)).
		Msg("api.fillMetricValueFromPathInfo started")
	defer log.Debug().Msg("api.fillMetricValueFromPathInfo ended")

	if pInf.valM == "" {
		return http.StatusNotFound, ErrMetricValueNotSpecified
	}

	if metric.MType == "gauge" {
		value, err := strconv.ParseFloat(pInf.valM, 64)
		if err != nil {
			return http.StatusBadRequest, err
		}
		metric.Value = &value
	} else if metric.MType == "counter" {
		value, err := strconv.Atoi(pInf.valM)
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
