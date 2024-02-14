package handlers

import (
	"net/http"
	"testing"
	"net/http/httptest"
	"io"
	"log"
	"strings"
	"context"

	"github.com/YaNeAndrey/ya-metrics/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/go-chi/chi/v5"
)

func TestHandleGetRoot(t *testing.T) {
	memStorage := storage.NewMemStorage()

	memStorage.UpdateGaugeMetric("firstGauge", 123.25)
	memStorage.UpdateGaugeMetric("SecondGauge", 2.1)
	memStorage.UpdateCounterMetric("CounterMetric", 444)
	
	result := `
	<table>
		<thead>
			<tr>
				<th>Metric Name</th>
				<th>Metric Value</th>
			</tr>
		</thead>
		<tbody>
			
				<tr>
					<td>CounterMetric</td>
					<td>444</td>
				</tr>
			
				<tr>
					<td>SecondGauge</td>
					<td>2.1</td>
				</tr>
			
				<tr>
					<td>firstGauge</td>
					<td>123.25</td>
				</tr>
			
		</tbody>
	</table>`

	result = strings.ReplaceAll(result,"\n", "")
	result = strings.ReplaceAll(result,"\t", "")
	result = strings.ReplaceAll(result," ", "")

	type args struct {
		req  *http.Request
		ms *storage.MemStorage
	}
	type want struct{
		body string
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "First test. Check root handler",
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/", nil),
				ms: memStorage,
			},
			want: want{
				body: result,
				statusCode: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := chi.NewRouter()

			r.Route("/", func(r chi.Router) {
				r.Get("/",func(rw http.ResponseWriter, r *http.Request) {
					HandleGetRoot(rw,tt.args.req,tt.args.ms)
				})
			})

			r.ServeHTTP(w,tt.args.req)

			result := w.Result()

			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			body, err := io.ReadAll(result.Body)
			if err != nil {
				log.Fatalf("failed to read response body, err:%v", err)
			}
			bodyStr := string(body)
			bodyStr = strings.ReplaceAll(bodyStr,"\n", "")
			bodyStr = strings.ReplaceAll(bodyStr,"\t", "")
			bodyStr = strings.ReplaceAll(bodyStr," ", "")

			assert.Equal(t,tt.want.body,bodyStr)
		})
	}
}

func AddURLParamsForChi(r *http.Request,urlParams map[string]string) *http.Request{
	rctx := chi.NewRouteContext()
	for key,value := range(urlParams) {
		rctx.URLParams.Add(key, value)
	}
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	return r
}

func TestHandleGetMetricValue(t *testing.T) {
	
	memStorage := storage.NewMemStorage()
	memStorage.UpdateGaugeMetric("GaugeMetric", 124.2345)
	memStorage.UpdateCounterMetric("metric", 124)
	
	type args struct {
		req  *http.Request
		ms *storage.MemStorage
	}
	type want struct {
		value string
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
			req: AddURLParamsForChi(httptest.NewRequest(http.MethodGet, "/value/counter/metric", nil),map[string]string {"type" : "counter", "name" : "metric"}),
			ms: memStorage,
			},
			want: want{
				value: "124",
				statusCode: http.StatusOK,
			},
		},		
		{
			name: "Second test. Trying to get metric value without name",
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/value/counter", nil),
				ms: memStorage,
			},
			want: want{
				value: "",
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "Third test. Trying to get metric value with incorrect metric type",
			args: args{
				req: AddURLParamsForChi(httptest.NewRequest(http.MethodGet, "/value/list/GaugeMetric", nil),map[string]string {"type" : "list", "name" : "GaugeMetric"}),
				ms: memStorage,
			},
			want: want{
				value: "",
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
				r.Get("/{type}/{name}",func(rw http.ResponseWriter, r *http.Request) {
					HandleGetMetricValue(rw,tt.args.req,tt.args.ms)
				})
			})
	

			r.ServeHTTP(w,tt.args.req)

			result := w.Result()

			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			bodyStr, err := io.ReadAll(result.Body)
			if err != nil {
				log.Fatalf("failed to read response body, err:%v", err)
			}

			assert.Equal(t,tt.want.value,string(bodyStr))
		})
	}
}

func TestHandlePostUpdateMetricValue(t *testing.T) {
	type args struct {
		req  *http.Request
		ms *storage.MemStorage
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
				ms: storage.NewMemStorage(),
			},
			want: http.StatusOK,
		},
		{
			name: "Second test. Trying to update counter metric with incorrect value",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/update/counter/NewMetric/value", nil),
				ms: storage.NewMemStorage(),
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
					HandlePostUpdateMetricValue(rw,r,tt.args.ms)
				})
			})

			r.ServeHTTP(w,tt.args.req)

			result := w.Result()

			defer result.Body.Close()
			assert.Equal(t, tt.want, result.StatusCode)
		})
	}
}


func Test_getGaugeMetricValue(t *testing.T) {

	memStorage := storage.NewMemStorage()
	memStorage.UpdateGaugeMetric("SomeMetric", 124)

	type args struct {
		metricName string
		ms         *storage.MemStorage
	}
	type want struct {
		value  string
		statusCode int
	}
	tests := []struct {
		name  string
		args  args
		want  want
	}{
		{
			name: "First test. Get gauge metric value",
			args: args{
				metricName: "SomeMetric",
				ms: memStorage,
			},
			want: want{
				value: "124",
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Second test. Trying to get value of a non-existent metric",
			args: args{
				metricName: "SomeMetric2",
				ms: memStorage,
			},
			want: want{
				value: "",
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, statusCode := getGaugeMetricValue(tt.args.metricName, tt.args.ms)
			assert.Equal(t,value,tt.want.value )
			assert.Equal(t,statusCode,tt.want.statusCode)
		})
	}
}

func Test_getCounterMetricValue(t *testing.T) {
	memStorage := storage.NewMemStorage()
	memStorage.UpdateCounterMetric("SomeMetric", 124)

	
	type args struct {
		metricName string
		ms         *storage.MemStorage
	}
	type want struct {
		value string
		statusCode int
	}
	tests := []struct {
		name  string
		args  args
		want  want
	}{
		{
			name: "First test. Get counter metric value",
			args: args{
				metricName: "SomeMetric",
				ms: memStorage,
			},
			want: want{
				value: "124",
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Second test. Trying to get value of a non-existent metric",
			args: args{
				metricName: "SomeMetric2",
				ms: memStorage,
			},
			want: want{
				value: "",
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, statusCode := getCounterMetricValue(tt.args.metricName, tt.args.ms)
			assert.Equal(t,value,tt.want.value )
			assert.Equal(t,statusCode,tt.want.statusCode)
		})
	}
}

func Test_updateMetric(t *testing.T) {
	type args struct {
		metricType     string
		metricName     string
		metricValueStr string
		ms             *storage.MemStorage
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "First test. Update gauge metric",
			args: args{
				metricType: "gauge",
				metricName: "SomeMetric",
				metricValueStr: "123",
				ms: storage.NewMemStorage(),
				},
			want: http.StatusOK,
		},		
		{
			name: "Second test. Update counter metric",
			args: args{
				metricType: "counter",
				metricName: "SomeMetric",
				metricValueStr: "123",
				ms: storage.NewMemStorage(),
				},
			want: http.StatusOK,
			},		
			{
				name: "Third test. Trying to update metric with unknown type",
				args: args{
					metricType: "IncorrectMetricType",
					metricName: "SomeMetric",
					metricValueStr: "123",
					ms: storage.NewMemStorage(),
				},
				want: http.StatusBadRequest,
			},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t,updateMetric(tt.args.metricType, tt.args.metricName, tt.args.metricValueStr, tt.args.ms),tt.want)
		})
	}
}

func Test_checkDataAndUpdateGauge(t *testing.T) {
	type args struct {
		metricName     string
		metricValueStr string
		ms             *storage.MemStorage
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "First test. Update gauge metric",
			args: args{
				metricName: "SomeMetric",
				metricValueStr: "123",
				ms: storage.NewMemStorage(),
			},
			want: http.StatusOK,
		},
		{
			name: "Second test. Trying to update gauge metric with incorrect value",
			args: args{
				metricName: "SomeMetric",
				metricValueStr: "incorrect",
				ms: storage.NewMemStorage(),
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkDataAndUpdateGauge(tt.args.metricName, tt.args.metricValueStr, tt.args.ms); got != tt.want {
				t.Errorf("checkDataAndUpdateGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkDataAndUpdateCounter(t *testing.T) {
	type args struct {
		metricName     string
		metricValueStr string
		ms             *storage.MemStorage
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "First test. Update counter metric",
			args: args{
				metricName: "SomeMetric",
				metricValueStr: "123",
				ms: storage.NewMemStorage(),
			},
			want: http.StatusOK,
		},
		{
			name: "Second test. Trying to update counter metric with incorrect value",
			args: args{
				metricName: "SomeMetric",
				metricValueStr: "incorrect",
				ms: storage.NewMemStorage(),
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkDataAndUpdateCounter(tt.args.metricName, tt.args.metricValueStr, tt.args.ms); got != tt.want {
				t.Errorf("checkDataAndUpdateCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}
