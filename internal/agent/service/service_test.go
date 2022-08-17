package service

import (
	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/v1tbrah/metricsAndAlerting/internal/agent/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/agent/service/mockapi"
)

func Test_service_reportMetric(t *testing.T) {

	cfg := config.NewCfg()
	myService := NewService(cfg)

	API := mockapi.NewAPI(myService.cfg.ServerAddr, myService.data)
	defer API.Server.Close()

	tests := []struct {
		name           string
		MID            string
		expectedStatus int
		wantErr        bool
	}{
		{
			name:           "TestOK",
			MID:            "Alloc",
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "TestNotOk",
			MID:            "unknown",
			expectedStatus: http.StatusBadRequest,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceData := myService.data.Metrics
			resp, err := myService.reportMetric(serviceData[tt.MID])
			if tt.wantErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			assert.Equal(t, resp.RawResponse.StatusCode, tt.expectedStatus)
		})
	}
}

func Test_service_reportListMetrics(t *testing.T) {

	cfg := config.NewCfg()
	myService := NewService(cfg)

	API := mockapi.NewAPI(myService.cfg.ServerAddr, myService.data)
	defer API.Server.Close()

	gaugeValue := 2.0

	tests := []struct {
		name           string
		Metrics        []metric.Metrics
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "TestOK",
			Metrics: []metric.Metrics{
				{
					ID:    "Alloc",
					MType: "gauge",
					Delta: nil,
					Value: &gaugeValue,
					Hash:  "",
				},
				{
					ID:    "Malloc",
					MType: "gauge",
					Delta: nil,
					Value: &gaugeValue,
					Hash:  "",
				},
			},
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := myService.reportListMetrics(tt.Metrics)
			if tt.wantErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			assert.Equal(t, resp.RawResponse.StatusCode, tt.expectedStatus)
		})
	}
}

func Test_service_getMetric(t *testing.T) {

	cfg := config.NewCfg()
	myService := NewService(cfg)

	API := mockapi.NewAPI(myService.cfg.ServerAddr, myService.data)
	defer API.Server.Close()

	tests := []struct {
		name           string
		MID            string
		MType          string
		expectedStatus int
		wantErr        bool
	}{
		{
			name:           "TestOK",
			MID:            "Alloc",
			MType:          "gauge",
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "Test empty ID",
			MID:            "",
			MType:          "gauge",
			expectedStatus: http.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:           "Test empty type",
			MID:            "Alloc",
			MType:          "",
			expectedStatus: http.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:           "Test invalid type",
			MID:            "Alloc",
			MType:          "unknown",
			expectedStatus: http.StatusBadRequest,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := myService.getMetric(tt.MID, tt.MType)
			if tt.wantErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			assert.Equal(t, resp.RawResponse.StatusCode, tt.expectedStatus)
		})
	}
}
