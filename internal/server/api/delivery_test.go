package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/memory"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/service"
	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

func TestUpdateHandler(t *testing.T) {

	myCfg, err := config.New()
	if err != nil {
		panic(err)
	}
	myServ, err := service.New(memory.New(""), myCfg)
	if err != nil {
		panic(err)
	}
	testAPI := New(myServ)

	localHost := "http://127.0.0.1:8080"

	type args struct {
		request *http.Request
	}

	type want struct {
		statusCode int
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test UpdateBasic Gauge OK",
			args: args{
				request: httptest.NewRequest(http.MethodPost, localHost+"/update/", updateBody("Alloc", "gauge", 1.22, 0)),
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Test UpdateBasic Counter OK",
			args: args{
				request: httptest.NewRequest(http.MethodPost, localHost+"/update/", updateBody("PollCount", "counter", 0.0, 1)),
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Test UpdateBasic /update/ Not Found",
			args: args{
				request: httptest.NewRequest(http.MethodPost, localHost+"/update/", updateBody("", "", 0.0, 0)),
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "Test UpdateBasic /update/unknown/ Not Implemented",
			args: args{
				request: httptest.NewRequest(http.MethodPost, localHost+"/update/", updateBody("", "unknown", 0.0, 0)),
			},
			want: want{
				statusCode: http.StatusNotImplemented,
			},
		},
		{
			name: "Test UpdateBasic /update/gauge/ Not Found",
			args: args{
				request: httptest.NewRequest(http.MethodPost, localHost+"/update/", updateBody("", "gauge", 0.0, 0)),
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "Test UpdateBasic /update/counter/ Not Found",
			args: args{
				request: httptest.NewRequest(http.MethodPost, localHost+"/update/", updateBody("", "counter", 0.0, 0)),
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "Test UpdateBasic /update/gauge/testNameGauge/",
			args: args{
				request: httptest.NewRequest(http.MethodPost, localHost+"/update/", updateBody("testNameGauge", "gauge", 0.0, 0)),
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Test UpdateBasic /update/counter/testNameCounter/",
			args: args{
				request: httptest.NewRequest(http.MethodPost, localHost+"/update/", updateBody("testNameCounter", "counter", 0.0, 0)),
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := tt.args.request
			w := httptest.NewRecorder()
			h := http.HandlerFunc(testAPI.updateMetricHandler)
			h.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()

			require.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func updateBody(MName, MType string, Value float64, Delta int64) *bytes.Buffer {
	allocValue := Value
	deltaValue := Delta
	metricForBody := metric.Metrics{
		ID:    MName,
		MType: MType,
		Delta: &deltaValue,
		Value: &allocValue,
	}
	body, _ := json.Marshal(&metricForBody)
	return bytes.NewBuffer(body)
}

func TestGetValueHandler(t *testing.T) {

	localHost := "http://127.0.0.1:8080"
	ctx := context.Background()

	myCfg, err := config.New()
	if err != nil {
		panic(err)
	}
	myServ, err := service.New(memory.New(""), myCfg)
	if err != nil {
		panic(err)
	}
	testAPI := New(myServ)

	gaugeValue := 2.22
	testAPIWithAllocMetric := testAPI
	testAPIWithAllocMetric.service.Storage.SetMetric(ctx,
		metric.Metrics{
			ID:    "Alloc",
			MType: "gauge",
			Delta: nil,
			Value: &gaugeValue,
		})

	counterValue := int64(2)
	testAPIWithCounterMetric := testAPI
	testAPIWithCounterMetric.service.Storage.SetMetric(ctx,
		metric.Metrics{
			ID:    "PollCount",
			MType: "counter",
			Delta: &counterValue,
			Value: nil,
		})

	type args struct {
		request *http.Request
		api     *api
	}
	type want struct {
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test Value Gauge OK",
			args: args{
				request: httptest.NewRequest(http.MethodPost, localHost+"/value/", getBody("Alloc", "gauge")),
				api:     testAPIWithAllocMetric,
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Test Value Counter OK",
			args: args{
				request: httptest.NewRequest(http.MethodPost, localHost+"/value/", getBody("PollCount", "counter")),
				api:     testAPIWithCounterMetric,
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Test Value /value/ Not Found",
			args: args{
				request: httptest.NewRequest(http.MethodPost, localHost+"/value/", getBody("", "")),
				api:     testAPI,
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "Test Value /value/unknown/unknown - Not Implemented (invalid type)",
			args: args{
				request: httptest.NewRequest(http.MethodPost, localHost+"/value/", getBody("unknown", "unknown")),
				api:     testAPI,
			},
			want: want{
				statusCode: http.StatusNotImplemented,
			},
		},
		{
			name: "Test Value /value/gauge/unknown - Not Found (invalid name)",
			args: args{
				request: httptest.NewRequest(http.MethodPost, localHost+"/value/", getBody("unknown", "gauge")),
				api:     testAPI,
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "Test Value /value/counter/unknown - Not Found (invalid name)",
			args: args{
				request: httptest.NewRequest(http.MethodPost, localHost+"/value/", getBody("unknown", "counter")),
				api:     testAPI,
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := tt.args.request
			w := httptest.NewRecorder()
			h := http.HandlerFunc(tt.args.api.getMetricValueHandler)
			h.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()

			require.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func getBody(MName, MType string) *bytes.Buffer {
	metricForBody := metric.Metrics{
		ID:    MName,
		MType: MType,
	}
	body, _ := json.Marshal(&metricForBody)
	return bytes.NewBuffer(body)
}
