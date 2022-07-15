package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo/memory"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/service"
)

func TestUpdateHandler(t *testing.T) {

	testApi := NewAPI(service.NewService(memory.NewMemStorage()))
	defaultHeader := map[string][]string{"Content-Type": []string{"text/plain"}}

	type args struct {
		request *http.Request
		api     *api
	}

	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test Update Gauge OK",
			args: args{
				request: &http.Request{
					Method: http.MethodPost,
					URL:    &url.URL{Path: "/update/gauge/Alloc/1.0"},
					Header: defaultHeader},
				api: testApi,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
			},
		},
		{
			name: "Test Update Counter OK",
			args: args{
				request: &http.Request{
					Method: http.MethodPost,
					URL:    &url.URL{Path: "/update/counter/PollCount/1"},
					Header: defaultHeader},
				api: testApi,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
			},
		},
		{
			name: "Test Update /update/ Not Found",
			args: args{
				request: &http.Request{
					Method: http.MethodPost,
					URL:    &url.URL{Path: "/update/"},
					Header: defaultHeader},
				api: testApi,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name: "Test Update /update/unknown/ Not Implemented",
			args: args{
				request: &http.Request{
					Method: http.MethodPost,
					URL:    &url.URL{Path: "/update/unknown/testCounter/100"},
					Header: defaultHeader},
				api: testApi,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotImplemented,
			},
		},
		{
			name: "Test Update /update/gauge/ Not Found",
			args: args{
				request: &http.Request{
					Method: http.MethodPost,
					URL:    &url.URL{Path: "/update/gauge/"},
					Header: defaultHeader},
				api: testApi,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name: "Test Update /update/counter/ Not Found",
			args: args{
				request: &http.Request{
					Method: http.MethodPost,
					URL:    &url.URL{Path: "/update/counter/"},
					Header: defaultHeader},
				api: testApi,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name: "Test Update /update/gauge/testNameM/",
			args: args{
				request: &http.Request{
					Method: http.MethodPost,
					URL:    &url.URL{Path: "/update/gauge/testNameM/"},
					Header: defaultHeader},
				api: NewAPI(service.NewService(memory.NewMemStorage())),
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name: "Test Update /update/counter/testNameM/",
			args: args{
				request: &http.Request{
					Method: http.MethodPost,
					URL:    &url.URL{Path: "/update/counter/testNameM/"},
					Header: defaultHeader},
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name: "Test Update /update/gauge/Alloc/- Bad Request (invalid value)",
			args: args{
				request: &http.Request{
					Method: http.MethodPost,
					URL:    &url.URL{Path: "/update/gauge/Alloc/-"},
					Header: defaultHeader},
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name: "Test Update /update/counter/PollCount/1.0 Bad Request (invalid value)",
			args: args{
				request: &http.Request{
					Method: http.MethodPost,
					URL:    &url.URL{Path: "/update/counter/PollCount/1.0"},
					Header: defaultHeader},
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := tt.args.request
			w := httptest.NewRecorder()
			h := http.HandlerFunc(checkTypeAndNameMetric("update", tt.args.api.updateHandler()))
			h.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()

			require.Equal(t, tt.want.statusCode, result.StatusCode)
			if result.StatusCode == http.StatusOK {
				require.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			}
		})
	}
}

func TestGetValueHandler(t *testing.T) {

	testAPI := NewAPI(service.NewService(memory.NewMemStorage()))

	testAPIWithAllocMetric := NewAPI(service.NewService(memory.NewMemStorage()))
	gaugeMetrics, _ := testAPIWithAllocMetric.service.MemStorage.Metrics.MetricsOfType("gauge")
	gaugeMetrics.Store("Alloc", "2.222")

	testApiWithCounterMetric := NewAPI(service.NewService(memory.NewMemStorage()))
	counterMetrics, _ := testApiWithCounterMetric.service.MemStorage.Metrics.MetricsOfType("counter")
	counterMetrics.Store("PollCount", "7")

	type args struct {
		request *http.Request
		api     *api
	}
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test Value Gauge OK",
			args: args{
				request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Path: "/value/gauge/Alloc"},
					Header: map[string][]string{"Content-Type": []string{"text/plain"}}},
				api: testAPIWithAllocMetric,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
			},
		},
		{
			name: "Test Value Counter OK",
			args: args{
				request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Path: "/value/counter/PollCount"},
					Header: map[string][]string{"Content-Type": []string{"text/plain"}}},
				api: testApiWithCounterMetric,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
			},
		},
		{
			name: "Test Value /value/ Not Found",
			args: args{
				request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Path: "/value/"},
					Header: map[string][]string{"Content-Type": []string{"text/plain"}}},
				api: testAPI,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name: "Test Value /value/unknown/unknown - Not Implemented (invalid type)",
			args: args{
				request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Path: "/value/unknown/unknown"},
					Header: map[string][]string{"Content-Type": []string{"text/plain"}}},
				api: testAPI,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotImplemented,
			},
		},
		{
			name: "Test Value /value/gauge/unknown - Not Found (invalid name)",
			args: args{
				request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Path: "/value/gauge/unknown"},
					Header: map[string][]string{"Content-Type": []string{"text/plain"}}},
				api: testAPI,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name: "Test Value /value/counter/unknown - Not Found (invalid name)",
			args: args{
				request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Path: "/value/counter/unknown"},
					Header: map[string][]string{"Content-Type": []string{"text/plain"}}},
				api: testAPI,
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := tt.args.request
			w := httptest.NewRecorder()
			h := http.HandlerFunc(checkTypeAndNameMetric("value", tt.args.api.getValueHandler()))
			h.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()

			require.Equal(t, tt.want.statusCode, result.StatusCode)
			if result.StatusCode == http.StatusOK {
				require.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			}
		})
	}
}
