// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	metric "github.com/v1tbrah/metricsAndAlerting/pkg/metric"

	model "github.com/v1tbrah/metricsAndAlerting/internal/server/model"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// GetData provides a mock function with given fields: _a0
func (_m *Service) GetData(_a0 context.Context) (model.Data, error) {
	ret := _m.Called(_a0)

	var r0 model.Data
	if rf, ok := ret.Get(0).(func(context.Context) model.Data); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.Data)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMetric provides a mock function with given fields: _a0, _a1
func (_m *Service) GetMetric(_a0 context.Context, _a1 string) (metric.Metrics, bool, error) {
	ret := _m.Called(_a0, _a1)

	var r0 metric.Metrics
	if rf, ok := ret.Get(0).(func(context.Context, string) metric.Metrics); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(metric.Metrics)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(context.Context, string) bool); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(_a0, _a1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// PingDBStorage provides a mock function with given fields: _a0
func (_m *Service) PingDBStorage(_a0 context.Context) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetListMetrics provides a mock function with given fields: _a0, _a1
func (_m *Service) SetListMetrics(_a0 context.Context, _a1 []metric.Metrics) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []metric.Metrics) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetMetric provides a mock function with given fields: _a0, _a1
func (_m *Service) SetMetric(_a0 context.Context, _a1 metric.Metrics) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, metric.Metrics) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ShutDown provides a mock function with given fields:
func (_m *Service) ShutDown() {
	_m.Called()
}
