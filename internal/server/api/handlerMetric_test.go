package api

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/api/mocks"
	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

func Test_api_handlerUpdateMetric(t *testing.T) {

	testCfg := &mocks.Config{}
	testCfg.On("HashKey").Return("mockHashMetric")

	tests := []struct {
		name         string
		payload      string
		mockService  *mocks.Service
		expectedCode int
	}{
		{
			name:    "gauge OK",
			payload: "{\"id\":\"testName\",\"type\":\"gauge\",\"value\":1.0}",
			mockService: func() *mocks.Service {
				testService := mocks.Service{}

				testService.On("GetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(metric.Metrics{}, true, nil).
					Once()
				testService.On("SetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("metric.Metrics")).
					Return(nil).
					Once()
				return &testService
			}(),
			expectedCode: http.StatusOK,
		},
		{
			name:    "counter OK",
			payload: "{\"id\":\"testName\",\"type\":\"counter\",\"delta\":1}",
			mockService: func() *mocks.Service {
				testService := mocks.Service{}

				testService.On("GetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(metric.Metrics{}, true, nil).
					Once()
				testService.On("SetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("metric.Metrics")).
					Return(nil).
					Once()
				return &testService
			}(),
			expectedCode: http.StatusOK,
		},
		{
			name:    "invalid json body",
			payload: "{\"id\":\"testName\",\"type\":\"counter\",\"invalid}",
			mockService: func() *mocks.Service {
				testService := mocks.Service{}
				testService.On("GetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(metric.Metrics{}, true, nil).
					Once()
				testService.On("SetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("metric.Metrics")).
					Return(nil).
					Once()
				return &testService
			}(),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "empty type metric",
			payload: "{\"id\":\"testName\",\"type\":\"\",\"delta\":1}",
			mockService: func() *mocks.Service {
				testService := mocks.Service{}

				testService.On("GetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(metric.Metrics{}, true, nil).
					Once()
				testService.On("SetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("metric.Metrics")).
					Return(nil).
					Once()
				return &testService
			}(),
			expectedCode: http.StatusNotFound,
		},
		{
			name:    "invalid type metric",
			payload: "{\"id\":\"testName\",\"type\":\"invalid\",\"delta\":1}",
			mockService: func() *mocks.Service {
				testService := mocks.Service{}

				testService.On("GetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(metric.Metrics{}, true, nil).
					Once()
				testService.On("SetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("metric.Metrics")).
					Return(nil).
					Once()
				return &testService
			}(),
			expectedCode: http.StatusNotImplemented,
		},
		{
			name:    "empty name metric",
			payload: "{\"id\":\"\",\"type\":\"gauge\",\"value\":1.0}",
			mockService: func() *mocks.Service {
				testService := mocks.Service{}

				testService.On("GetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(metric.Metrics{}, true, nil).
					Once()
				testService.On("SetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("metric.Metrics")).
					Return(nil).
					Once()
				return &testService
			}(),
			expectedCode: http.StatusNotFound,
		},
		{
			name:    "gauge: empty value",
			payload: "{\"id\":\"\",\"type\":\"gauge\"}",
			mockService: func() *mocks.Service {
				testService := mocks.Service{}

				testService.On("GetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(metric.Metrics{}, true, nil).
					Once()
				testService.On("SetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("metric.Metrics")).
					Return(nil).
					Once()
				return &testService
			}(),
			expectedCode: http.StatusNotFound,
		},
		{
			name:    "counter: empty delta",
			payload: "{\"id\":\"\",\"type\":\"counter\"}",
			mockService: func() *mocks.Service {
				testService := mocks.Service{}

				testService.On("GetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(metric.Metrics{}, true, nil).
					Once()
				testService.On("SetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("metric.Metrics")).
					Return(nil).
					Once()
				return &testService
			}(),
			expectedCode: http.StatusNotFound,
		},
		{
			name:    "error while getting metric",
			payload: "{\"id\":\"testName\",\"type\":\"gauge\",\"value\":1.0}",
			mockService: func() *mocks.Service {
				testService := mocks.Service{}

				testService.On("GetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(metric.Metrics{}, true, errors.New("mock error")).
					Once()
				testService.On("SetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("metric.Metrics")).
					Return(nil).
					Once()
				return &testService
			}(),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "error while setting metric",
			payload: "{\"id\":\"testName\",\"type\":\"gauge\",\"value\":1.0}",
			mockService: func() *mocks.Service {
				testService := mocks.Service{}

				testService.On("GetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(metric.Metrics{}, true, nil).
					Once()
				testService.On("SetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("metric.Metrics")).
					Return(errors.New("mock error")).
					Once()
				return &testService
			}(),
			expectedCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAPI := api{}
			testAPI.service = tt.mockService
			testAPI.cfg = testCfg

			rec := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Post("/updateMetricMockEndpoint", testAPI.handlerUpdateMetric)

			b := &bytes.Buffer{}
			b.WriteString(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/updateMetricMockEndpoint", b)

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestGetValueHandler(t *testing.T) {

	//localHost := "http://127.0.0.1:8080"
	//ctx := context.Background()
	//
	//myCfg, err := config.New()
	//if err != nil {
	//	panic(err)
	//}
	//myServ, err := service.New(memory.New(""), myCfg)
	//if err != nil {
	//	panic(err)
	//}
	//testAPI := New(myServ)
	//
	//gaugeValue := 2.22
	//testAPIWithAllocMetric := testAPI
	//testAPIWithAllocMetric.service.Storage.SetMetric(ctx,
	//	metric.Metrics{
	//		ID:    "Alloc",
	//		MType: "gauge",
	//		Delta: nil,
	//		Value: &gaugeValue,
	//	})
	//
	//counterValue := int64(2)
	//testAPIWithCounterMetric := testAPI
	//testAPIWithCounterMetric.service.Storage.SetMetric(ctx,
	//	metric.Metrics{
	//		ID:    "PollCount",
	//		MType: "counter",
	//		Delta: &counterValue,
	//		Value: nil,
	//	})
	//
	//type args struct {
	//	request *http.Request
	//	api     *api
	//}
	//type want struct {
	//	statusCode int
	//}
	//tests := []struct {
	//	name string
	//	args args
	//	want want
	//}{
	//	{
	//		name: "Test Value Gauge OK",
	//		args: args{
	//			request: httptest.NewRequest(http.MethodPost, localHost+"/value/", getBody("Alloc", "gauge")),
	//			api:     testAPIWithAllocMetric,
	//		},
	//		want: want{
	//			statusCode: http.StatusOK,
	//		},
	//	},
	//	{
	//		name: "Test Value Counter OK",
	//		args: args{
	//			request: httptest.NewRequest(http.MethodPost, localHost+"/value/", getBody("PollCount", "counter")),
	//			api:     testAPIWithCounterMetric,
	//		},
	//		want: want{
	//			statusCode: http.StatusOK,
	//		},
	//	},
	//	{
	//		name: "Test Value /value/ Not Found",
	//		args: args{
	//			request: httptest.NewRequest(http.MethodPost, localHost+"/value/", getBody("", "")),
	//			api:     testAPI,
	//		},
	//		want: want{
	//			statusCode: http.StatusNotFound,
	//		},
	//	},
	//	{
	//		name: "Test Value /value/unknown/unknown - Not Implemented (invalid type)",
	//		args: args{
	//			request: httptest.NewRequest(http.MethodPost, localHost+"/value/", getBody("unknown", "unknown")),
	//			api:     testAPI,
	//		},
	//		want: want{
	//			statusCode: http.StatusNotImplemented,
	//		},
	//	},
	//	{
	//		name: "Test Value /value/gauge/unknown - Not Found (invalid name)",
	//		args: args{
	//			request: httptest.NewRequest(http.MethodPost, localHost+"/value/", getBody("unknown", "gauge")),
	//			api:     testAPI,
	//		},
	//		want: want{
	//			statusCode: http.StatusNotFound,
	//		},
	//	},
	//	{
	//		name: "Test Value /value/counter/unknown - Not Found (invalid name)",
	//		args: args{
	//			request: httptest.NewRequest(http.MethodPost, localHost+"/value/", getBody("unknown", "counter")),
	//			api:     testAPI,
	//		},
	//		want: want{
	//			statusCode: http.StatusNotFound,
	//		},
	//	},
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		request := tt.args.request
	//		w := httptest.NewRecorder()
	//		h := http.HandlerFunc(tt.args.api.getMetricValueHandler)
	//		h.ServeHTTP(w, request)
	//		result := w.Result()
	//		defer result.Body.Close()
	//
	//		require.Equal(t, tt.want.statusCode, result.StatusCode)
	//	})
	//}
}

func Test_api_handlerUpdateListMetrics(t *testing.T) {
	testCfg := &mocks.Config{}
	testCfg.On("HashKey").Return("mockHashMetric")

	tests := []struct {
		name         string
		payload      string
		mockService  *mocks.Service
		expectedCode int
	}{
		{
			name:    "OK",
			payload: "[{\"id\":\"Alloc\",\"type\":\"gauge\",\"value\":1.1},{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":1}]",
			mockService: func() *mocks.Service {
				testService := mocks.Service{}

				testService.On("GetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(metric.Metrics{}, true, nil)
				testService.On("SetListMetrics", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("[]metric.Metrics")).
					Return(nil)
				return &testService
			}(),
			expectedCode: http.StatusOK,
		},
		{
			name:    "invalid json body",
			payload: "[___\"id\":\"Alloc\",\"type\":\"gauge\",\"value\":1.1},{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":1}]",
			mockService: func() *mocks.Service {
				testService := mocks.Service{}

				testService.On("GetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(metric.Metrics{}, true, nil)
				testService.On("SetListMetrics", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("[]metric.Metrics")).
					Return(nil)
				return &testService
			}(),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "empty type metric",
			payload: "[{\"id\":\"Alloc\",\"type\":\"\",\"value\":1.1},{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":1}]",
			mockService: func() *mocks.Service {
				testService := mocks.Service{}

				testService.On("GetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(metric.Metrics{}, true, nil)
				testService.On("SetListMetrics", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("[]metric.Metrics")).
					Return(nil)
				return &testService
			}(),
			expectedCode: http.StatusNotFound,
		},
		{
			name:    "invalid type metric",
			payload: "[{\"id\":\"Alloc\",\"type\":\"INVALID TYPE\",\"value\":1.1},{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":1}]",
			mockService: func() *mocks.Service {
				testService := mocks.Service{}

				testService.On("GetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(metric.Metrics{}, true, nil)
				testService.On("SetListMetrics", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("[]metric.Metrics")).
					Return(nil)
				return &testService
			}(),
			expectedCode: http.StatusNotImplemented,
		},
		{
			name:    "empty name metric",
			payload: "[{\"id\":\"\",\"type\":\"gauge\",\"value\":1.1},{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":1}]",
			mockService: func() *mocks.Service {
				testService := mocks.Service{}

				testService.On("GetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(metric.Metrics{}, true, nil)
				testService.On("SetListMetrics", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("[]metric.Metrics")).
					Return(nil)
				return &testService
			}(),
			expectedCode: http.StatusNotFound,
		},
		{
			name:    "gauge: empty value",
			payload: "[{\"id\":\"Alloc\",\"type\":\"gauge\"},{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":1}]",
			mockService: func() *mocks.Service {
				testService := mocks.Service{}

				testService.On("GetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(metric.Metrics{}, true, nil)
				testService.On("SetListMetrics", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("[]metric.Metrics")).
					Return(nil)
				return &testService
			}(),
			expectedCode: http.StatusNotFound,
		},
		{
			name:    "counter: empty delta",
			payload: "[{\"id\":\"Alloc\",\"type\":\"gauge\",\"value\":1.1},{\"id\":\"PollCount\",\"type\":\"counter\"}]",
			mockService: func() *mocks.Service {
				testService := mocks.Service{}

				testService.On("GetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(metric.Metrics{}, true, nil)
				testService.On("SetListMetrics", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("[]metric.Metrics")).
					Return(nil)
				return &testService
			}(),
			expectedCode: http.StatusNotFound,
		},
		{
			name:    "error while getting metric",
			payload: "[{\"id\":\"Alloc\",\"type\":\"gauge\",\"value\":1.1},{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":1}]",
			mockService: func() *mocks.Service {
				testService := mocks.Service{}

				testService.On("GetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(metric.Metrics{}, true, errors.New("mock error")).
					Once()
				testService.On("SetListMetrics", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("[]metric.Metrics")).
					Return(nil)
				return &testService
			}(),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "error while setting list metrics",
			payload: "[{\"id\":\"Alloc\",\"type\":\"gauge\",\"value\":1.1},{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":1}]",
			mockService: func() *mocks.Service {
				testService := mocks.Service{}

				testService.On("GetMetric", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).
					Return(metric.Metrics{}, true, nil)
				testService.On("SetListMetrics", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("[]metric.Metrics")).
					Return(errors.New("mock error"))
				return &testService
			}(),
			expectedCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAPI := api{}
			testAPI.service = tt.mockService
			testAPI.cfg = testCfg

			rec := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Post("/updateListMetricsMockEndpoint", testAPI.handlerUpdateListMetrics)

			b := &bytes.Buffer{}
			b.WriteString(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/updateListMetricsMockEndpoint", b)

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}
