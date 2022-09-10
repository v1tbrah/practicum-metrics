// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	metric "github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

// Data is an autogenerated mock type for the Data type
type Data struct {
	mock.Mock
}

// GetData provides a mock function with given fields:
func (_m *Data) GetData() map[string]metric.Metrics {
	ret := _m.Called()

	var r0 map[string]metric.Metrics
	if rf, ok := ret.Get(0).(func() map[string]metric.Metrics); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]metric.Metrics)
		}
	}

	return r0
}

// UpdateAdditional provides a mock function with given fields: keyForUpdateHash
func (_m *Data) UpdateAdditional(keyForUpdateHash string) error {
	ret := _m.Called(keyForUpdateHash)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(keyForUpdateHash)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateBasic provides a mock function with given fields: keyForUpdateHash
func (_m *Data) UpdateBasic(keyForUpdateHash string) {
	_m.Called(keyForUpdateHash)
}
