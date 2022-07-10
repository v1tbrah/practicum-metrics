package handler

import (
	"github.com/stretchr/testify/require"
	"github.com/v1tbrah/metricsAndAlerting/internal/metric"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestUpdateHandler(t *testing.T) {
	type args struct {
		request *http.Request
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
			name: "Test Update Not allowed method",
			args: args{
				request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Path: "/update/gauge/Alloc/1.0"},
					Header: map[string][]string{"Content-Type": []string{"text/plain"}}},
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusMethodNotAllowed,
			},
		},
		{
			name: "Test Update Gauge OK",
			args: args{
				request: &http.Request{
					Method: http.MethodPost,
					URL:    &url.URL{Path: "/update/gauge/Alloc/1.0"},
					Header: map[string][]string{"Content-Type": []string{"text/plain"}}},
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
					Header: map[string][]string{"Content-Type": []string{"text/plain"}}},
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
					Header: map[string][]string{"Content-Type": []string{"text/plain"}}},
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
					Header: map[string][]string{"Content-Type": []string{"text/plain"}}},
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
					Header: map[string][]string{"Content-Type": []string{"text/plain"}}},
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
					Header: map[string][]string{"Content-Type": []string{"text/plain"}}},
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
					Header: map[string][]string{"Content-Type": []string{"text/plain"}}},
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
					Header: map[string][]string{"Content-Type": []string{"text/plain"}}},
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
					Header: map[string][]string{"Content-Type": []string{"text/plain"}}},
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
					Header: map[string][]string{"Content-Type": []string{"text/plain"}}},
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
			h := http.HandlerFunc(UpdateHandler(&metric.Metrics{}))
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
	type args struct {
		request *http.Request
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
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name: "Test Value /value/unknown/unknown - Not Found (invalid type)",
			args: args{
				request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Path: "/value/unknown/unknown"},
					Header: map[string][]string{"Content-Type": []string{"text/plain"}}},
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name: "Test Value /value/gauge/unknown - Not Found (invalid name)",
			args: args{
				request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Path: "/value/gauge/unknown"},
					Header: map[string][]string{"Content-Type": []string{"text/plain"}}},
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
			h := http.HandlerFunc(GetValueHandler(&metric.Metrics{}))
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
