package service

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/v1tbrah/metricsAndAlerting/internal/agent/service/mockapi"
)

func Test_service_reportMetric(t *testing.T) {

	myService := NewService()

	API := mockapi.NewAPI(myService.options.ServerAddr, *myService.data)
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
			serviceData := *myService.data
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

func Test_service_getMetric(t *testing.T) {

	myService := NewService()

	API := mockapi.NewAPI(myService.options.ServerAddr, *myService.data)
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
