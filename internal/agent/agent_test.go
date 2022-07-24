package agent

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/v1tbrah/metricsAndAlerting/internal/agent/metric"
)

var testAgent = NewAgent()

var gaugeValue = 2.222
var counterValue = int64(1)

var metricOnServer = metric.Metrics{
	ID:    "Alloc",
	MType: "gauge",
	Delta: &counterValue,
	Value: &gaugeValue,
}

func mockRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.AllowContentType(testAgent.options.contentTypeJSON))

	r.Post("/update/", mockUpdateHandler)
	r.Post("/value/", mockValueHandler)

	return r
}

func mockUpdateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, _ := io.ReadAll(r.Body)
	w.Write([]byte(body))
}

func mockValueHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	metric := metric.Metrics{}
	err := json.Unmarshal(body, &metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metric.Value = metricOnServer.Value
	metric.Delta = metricOnServer.Delta
	resp, err := json.Marshal(metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func Test_agent_sendMetric(t *testing.T) {
	mockServer := httptest.NewUnstartedServer(mockRouter())
	listener, _ := net.Listen("tcp", testAgent.options.srvAddr)
	mockServer.Listener = listener
	mockServer.Start()
	defer mockServer.Close()

	gaugeValueForOk := 2.222
	metricForOk := metric.Metrics{
		ID:    "Alloc",
		MType: "gauge",
		Delta: nil,
		Value: &gaugeValueForOk,
	}

	type input struct {
		metric metric.Metrics
	}

	tests := []struct {
		name   string
		input  input
		output metric.Metrics
	}{
		{
			name: "TestOK",
			input: input{
				metric: metricForOk,
			},
			output: metricForOk,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := testAgent.sendMetricJSON(tt.input.metric)
			require.Nil(t, err)
			assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
			resMetric := metric.Metrics{}
			err = json.Unmarshal([]byte(resp.Body()), &resMetric)
			require.Nil(t, err)
			assert.Equal(t, tt.output, resMetric)
		})
	}
}
