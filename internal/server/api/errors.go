package api

import "errors"

var (
	ErrMetricTypeNotSpecified   = errors.New("metric type not specified")
	ErrMetricTypeNotImplemented = errors.New("metric type not implemented")
	ErrMetricIsNotExists        = errors.New("metric is not exists")
	ErrMetricNameNotSpecified   = errors.New("metric name not specified")
	ErrMetricValueNotSpecified  = errors.New("metric value not specified")
)
