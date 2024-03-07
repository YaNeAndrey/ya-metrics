package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/YaNeAndrey/ya-metrics/internal/storage/storagejson"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YaNeAndrey/ya-metrics/internal/constants"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func AddURLParamsForChi(r *http.Request, urlParams map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for key, value := range urlParams {
		rctx.URLParams.Add(key, value)
	}
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	return r
}

// -
func TestHandleGetRoot(t *testing.T) {
	type args struct {
		w  http.ResponseWriter
		r  *http.Request
		st *storage.StorageRepo
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleGetReadMetrics(tt.args.w, tt.args.r, tt.args.st)
		})
	}
}

// ++++
func TestHandleGetMetricValue(t *testing.T) {
	floatValue := float64(124.2345)
	intValue := int64(124)

	testMetrics := []storage.Metrics{
		{
			ID:    "GaugeMetric",
			MType: constants.GaugeMetricType,
			Value: &floatValue,
		},
		{
			ID:    "metric",
			MType: constants.CounterMetricType,
			Delta: &intValue,
		},
	}
	testStorage := storage.StorageRepo(storagejson.NewMemStorageJSON(testMetrics))

	type args struct {
		req *http.Request
		st  *storage.StorageRepo
	}
	type want struct {
		value      string
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "First test. Get metric value by name and type",
			args: args{
				req: AddURLParamsForChi(httptest.NewRequest(http.MethodGet, "/value/counter/metric", nil), map[string]string{"type": "counter", "name": "metric"}),
				st:  &testStorage,
			},
			want: want{
				value:      "124",
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Second test. Trying to get metric value without name",
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/value/counter", nil),
				st:  &testStorage,
			},
			want: want{
				value:      "",
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "Third test. Trying to get metric value with incorrect metric type",
			args: args{
				req: AddURLParamsForChi(httptest.NewRequest(http.MethodGet, "/value/list/GaugeMetric", nil), map[string]string{"type": "list", "name": "GaugeMetric"}),
				st:  &testStorage,
			},
			want: want{
				value:      "",
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.NotFound(func(rw http.ResponseWriter, r *http.Request) {
				rw.WriteHeader(http.StatusNotFound)
			})
			r.Route("/value", func(r chi.Router) {
				r.Get("/{type}/{name}", func(rw http.ResponseWriter, r *http.Request) {
					HandleGetReadOneMetric(rw, tt.args.req, tt.args.st)
				})
			})

			r.ServeHTTP(w, tt.args.req)

			result := w.Result()

			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			bodyStr, err := io.ReadAll(result.Body)
			if err != nil {
				log.Fatalf("failed to read response body, err:%v", err)
			}

			assert.Equal(t, tt.want.value, string(bodyStr))
		})
	}
}

// ++++
func TestHandlePostMetricValueJSON(t *testing.T) {
	floatValue := float64(124.2345)
	intValue := int64(124)
	testMetrics := []storage.Metrics{
		{
			ID:    "GaugeMetric",
			MType: constants.GaugeMetricType,
			Value: &floatValue,
		},
		{
			ID:    "CounterMetric",
			MType: constants.CounterMetricType,
			Delta: &intValue,
		},
	}
	testStorage := storage.StorageRepo(storagejson.NewMemStorageJSON(testMetrics))

	existedMetric, _ := json.Marshal(storage.Metrics{ID: "GaugeMetric", MType: constants.GaugeMetricType, Delta: &intValue, Value: &floatValue})
	notExistedMetric, _ := json.Marshal(storage.Metrics{ID: "NotExistedCounterMetric", MType: constants.CounterMetricType})

	type args struct {
		req *http.Request
		st  *storage.StorageRepo
	}
	type want struct {
		value      string
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "First test. Get metric value wth JSON",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/value/", bytes.NewReader(existedMetric)),
				st:  &testStorage,
			},
			want: want{
				value:      "124",
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Second test. Trying to get metric value without name",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/value/", bytes.NewReader(notExistedMetric)),
				st:  &testStorage,
			},
			want: want{
				value:      "",
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "Third test. Trying to get metric value without JSON in bode",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/value/", bytes.NewReader(nil)),
				st:  &testStorage,
			},
			want: want{
				value:      "",
				statusCode: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.NotFound(func(rw http.ResponseWriter, r *http.Request) {
				rw.WriteHeader(http.StatusNotFound)
			})
			r.Route("/value", func(r chi.Router) {
				r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
					HandlePostMetricValueJSON(rw, tt.args.req, tt.args.st)
				})
			})

			tt.args.req.Header.Add("Content-Type", "application/json")
			r.ServeHTTP(w, tt.args.req)

			result := w.Result()

			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			_, err := io.ReadAll(result.Body)
			if err != nil {
				log.Fatalf("failed to read response body, err:%v", err)
			}

			//assert.Equal(t, tt.want.value, string(bodyStr))
		})
	}
}

// ++++
func TestHandlePostUpdateMetricValue(t *testing.T) {
	testStorage := storage.StorageRepo(storagejson.NewMemStorageJSON([]storage.Metrics{}))

	type args struct {
		req *http.Request
		st  *storage.StorageRepo
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "First test. Update gauge metric",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/update/gauge/NewMetric/123", nil),
				st:  &testStorage,
			},
			want: http.StatusOK,
		},
		{
			name: "Second test. Trying to update counter metric with incorrect value",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/update/counter/NewMetric/value", nil),
				st:  &testStorage,
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := chi.NewRouter()

			r.Route("/update", func(r chi.Router) {
				r.Post("/{type}/{name}/{value}", func(rw http.ResponseWriter, r *http.Request) {
					HandlePostUpdateOneMetric(rw, r, tt.args.st)
				})
			})

			r.ServeHTTP(w, tt.args.req)

			result := w.Result()

			defer result.Body.Close()
			assert.Equal(t, tt.want, result.StatusCode)
		})
	}
}

// ++++
func TestHandlePostUpdateMetricValueJSON(t *testing.T) {
	testStorage := storage.StorageRepo(storagejson.NewMemStorageJSON([]storage.Metrics{}))
	floatValue := float64(123.45)
	correctMetric, _ := json.Marshal(storage.Metrics{ID: "gaugeMetric", MType: constants.GaugeMetricType, Value: &floatValue})
	incorrectMetric, _ := json.Marshal(storage.Metrics{ID: "countermetric", MType: constants.CounterMetricType, Value: &floatValue})
	type args struct {
		req *http.Request
		st  *storage.StorageRepo
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "First test. Update gauge metric with JSON",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/update/", bytes.NewReader(correctMetric)),
				st:  &testStorage,
			},
			want: http.StatusOK,
		},
		{
			name: "Second test. Trying to update counter metric without Delta",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/update/", bytes.NewReader(incorrectMetric)),
				st:  &testStorage,
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := chi.NewRouter()

			r.Route("/update", func(r chi.Router) {
				r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
					HandlePostUpdateOneMetricJSON(rw, r, tt.args.st)
				})
			})

			tt.args.req.Header.Add("Content-Type", "application/json")

			r.ServeHTTP(w, tt.args.req)

			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want, result.StatusCode)
		})
	}
}

func Test_updateMetric(t *testing.T) {
	type args struct {
		metricType     string
		metricName     string
		metricValueStr string
		st             *storage.StorageRepo
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := updateMetric(tt.args.metricType, tt.args.metricName, tt.args.metricValueStr, tt.args.st); got != tt.want {
				t.Errorf("updateMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}
