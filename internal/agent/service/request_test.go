package service

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/v1tbrah/metricsAndAlerting/internal/agent/service/mocks"
	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

func Test_service_reportMetric(t *testing.T) {

	mockData := &mocks.Data{}
	mockCfg := func(servAddr string) *mocks.Config {
		testCfg := mocks.Config{}
		testCfg.On("ReportMetricURL").Return(servAddr).Once()
		testCfg.On("String").Return("mock config").Once()
		return &testCfg
	}

	tests := []struct {
		name        string
		mockCfg     func(servAddr string) *mocks.Config
		metric      metric.Metrics
		handler     func(w http.ResponseWriter, r *http.Request)
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "OK",
			mockCfg: mockCfg,
			metric:  metric.NewMetric("Alloc", "gauge"),
			handler: func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if err == nil {
					w.Write(body)
				}
			},
		},
		{
			name:        "empty metric id",
			mockCfg:     mockCfg,
			metric:      metric.NewMetric("", "gauge"),
			handler:     func(w http.ResponseWriter, r *http.Request) {},
			wantErr:     true,
			expectedErr: ErrMetricIDIsEmpty,
		},
		{
			name:        "invalid metric type",
			mockCfg:     mockCfg,
			metric:      metric.NewMetric("Alloc", "invalid"),
			handler:     func(w http.ResponseWriter, r *http.Request) {},
			wantErr:     true,
			expectedErr: ErrInvalidMetricType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockHandler := http.NewServeMux()
			mockHandler.HandleFunc("/", tt.handler)
			mockServ := httptest.NewServer(mockHandler)

			testService, err := New(mockData, tt.mockCfg(mockServ.URL))
			if err != nil {
				log.Fatal().Err(err).Msg("unable to create test service")
			}
			resp, err := testService.reportMetric(tt.metric)

			if tt.wantErr {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, resp.StatusCode(), http.StatusOK)

				reqBody, _ := json.Marshal(tt.metric)
				assert.Equal(t, resp.Body(), reqBody)
			}
			mockServ.Close()
		})
	}
}

func Test_service_reportListMetrics(t *testing.T) {

	mockData := &mocks.Data{}
	mockCfg := func(servAddr string) *mocks.Config {
		testCfg := mocks.Config{}
		testCfg.On("ReportListMetricsURL").Return(servAddr).Once()
		testCfg.On("String").Return("mock config").Once()
		return &testCfg
	}

	var nilListMetrics []metric.Metrics

	tests := []struct {
		name        string
		mockCfg     func(servAddr string) *mocks.Config
		listMetrics []metric.Metrics
		handler     func(w http.ResponseWriter, r *http.Request)
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "OK",
			mockCfg: mockCfg,
			listMetrics: []metric.Metrics{
				metric.NewMetric("Alloc", "gauge"),
				metric.NewMetric("Malloc", "gauge"),
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if err == nil {
					w.Write(body)
				}
			},
		},
		{
			name:    "have metric with empty id",
			mockCfg: mockCfg,
			listMetrics: []metric.Metrics{
				metric.NewMetric("Alloc", "gauge"),
				metric.NewMetric("", "gauge"),
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if err == nil {
					w.Write(body)
				}
			},
			wantErr:     true,
			expectedErr: ErrMetricIDIsEmpty,
		},
		{
			name:    "have invalid metric type",
			mockCfg: mockCfg,
			listMetrics: []metric.Metrics{
				metric.NewMetric("Alloc", "gauge"),
				metric.NewMetric("Malloc", ""),
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if err == nil {
					w.Write(body)
				}
			},
			wantErr:     true,
			expectedErr: ErrInvalidMetricType,
		},
		{
			name:        "nil list metrics",
			mockCfg:     mockCfg,
			listMetrics: nilListMetrics,
			handler: func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if err == nil {
					w.Write(body)
				}
			},
			wantErr:     true,
			expectedErr: ErrListMetricsIsEmpty,
		},
		{
			name:        "empty list metrics",
			mockCfg:     mockCfg,
			listMetrics: []metric.Metrics{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if err == nil {
					w.Write(body)
				}
			},
			wantErr:     true,
			expectedErr: ErrListMetricsIsEmpty,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockHandler := http.NewServeMux()
			mockHandler.HandleFunc("/", tt.handler)
			mockServ := httptest.NewServer(mockHandler)

			testService, err := New(mockData, tt.mockCfg(mockServ.URL))
			if err != nil {
				log.Fatal().Err(err).Msg("unable to create test service")
			}
			resp, err := testService.reportListMetrics(tt.listMetrics)

			if tt.wantErr {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, resp.StatusCode(), http.StatusOK)

				reqBody, _ := json.Marshal(tt.listMetrics)
				assert.Equal(t, resp.Body(), reqBody)
			}
			mockServ.Close()
		})
	}
}

func Test_service_getMetric(t *testing.T) {

	mockData := &mocks.Data{}
	mockCfg := func(servAddr string) *mocks.Config {
		testCfg := mocks.Config{}
		testCfg.On("GetMetricURL").Return(servAddr).Once()
		testCfg.On("String").Return("mock config").Once()
		return &testCfg
	}

	tests := []struct {
		name        string
		mockCfg     func(servAddr string) *mocks.Config
		ID          string
		MType       string
		handler     func(w http.ResponseWriter, r *http.Request)
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "OK",
			mockCfg: mockCfg,
			ID:      "Alloc",
			MType:   "gauge",
			handler: func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if err == nil {
					w.Write(body)
				}
			},
		},
		{
			name:    "empty metric id",
			mockCfg: mockCfg,
			ID:      "",
			handler: func(w http.ResponseWriter, r *http.Request) {
			},
			wantErr:     true,
			expectedErr: ErrMetricIDIsEmpty,
		},
		{
			name:    "have invalid metric type",
			mockCfg: mockCfg,
			ID:      "Alloc",
			MType:   "invalid",
			handler: func(w http.ResponseWriter, r *http.Request) {
			},
			wantErr:     true,
			expectedErr: ErrInvalidMetricType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockHandler := http.NewServeMux()
			mockHandler.HandleFunc("/", tt.handler)
			mockServ := httptest.NewServer(mockHandler)

			testService, err := New(mockData, tt.mockCfg(mockServ.URL))
			if err != nil {
				log.Fatal().Err(err).Msg("unable to create test service")
			}
			resp, err := testService.getMetric(tt.ID, tt.MType)

			if tt.wantErr {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, resp.StatusCode(), http.StatusOK)

				respMetric := metric.Metrics{}
				json.Unmarshal(resp.Body(), &respMetric)
				assert.Equal(t, tt.ID, respMetric.ID)
				assert.Equal(t, tt.MType, respMetric.MType)
			}
			mockServ.Close()
		})
	}
}
