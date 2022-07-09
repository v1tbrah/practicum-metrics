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
			name: "Test StatusOK",
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
			name: "Test StatusMethodNotAllowed",
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
			name: "Test StatusNotFound 1",
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
			name: "Test StatusNotFound 2",
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
			name: "Test StatusNotFound 3",
			args: args{
				request: &http.Request{
					Method: http.MethodPost,
					URL:    &url.URL{Path: "/update/gauge/Alloc/"},
					Header: map[string][]string{"Content-Type": []string{"text/plain"}}},
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name: "Test StatusBadRequest",
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
			name: "TestCounterHandlers/invalid_value",
			args: args{
				request: &http.Request{
					Method: http.MethodPost,
					URL:    &url.URL{Path: "/update/counter/testCounter/none"},
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
