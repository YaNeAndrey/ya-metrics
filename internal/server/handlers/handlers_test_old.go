package handlers
/*
import (
	"net/http"
	"testing"
	"net/http/httptest"

	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleUpdateMetrics(t *testing.T) {
	type args struct {
		r  *http.Request
		ms *storage.MemStorage
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "First test",
			args: args{
				r: httptest.NewRequest(http.MethodPost, "/update/counter/NewMetric/123", nil),
				ms: storage.NewMemStorage(),
			},
			want: http.StatusOK,
		},
		{
			name: "Second test",
			args: args{
				r: httptest.NewRequest(http.MethodPost, "/update/counter1/NewMetric/value", nil),
				ms: storage.NewMemStorage(),
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//request := httptest.NewRequest(http.MethodPost, tt.request, nil)
            w := httptest.NewRecorder()
            
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if http.MethodPost == r.Method {
					HandlePostUpdateMetrics(w, tt.args.r, tt.args.ms)
				}else {
					w.WriteHeader(http.StatusMethodNotAllowed)
				}
			})

            h(w, tt.args.r)
			result := w.Result()
			assert.Equal(t, tt.want, result.StatusCode)

			err := result.Body.Close()
            require.NoError(t, err)
		})
	}
}

func Test_updateMetrics(t *testing.T) {
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
			name: "First test. Add gauge metric",
			args: args{
				metricType: "gauge",
				metricName: "SomeMetric",
				metricValueStr: "123",
				ms: storage.NewMemStorage(),
				},
			want: http.StatusOK,
		},		
		{
			name: "Second test. Add counter metric",
			args: args{
				metricType: "counter",
				metricName: "SomeMetric",
				metricValueStr: "123",
				ms: storage.NewMemStorage(),
				},
			want: http.StatusOK,
			},		
			{
				name: "Third test. Incorrect metric type ",
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
			assert.Equal(t,updateMetrics(tt.args.metricType, tt.args.metricName, tt.args.metricValueStr, tt.args.ms),tt.want)
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
			name: "First test",
			args: args{
				metricName: "SomeMetric",
				metricValueStr: "123",
				ms: storage.NewMemStorage(),
			},
			want: http.StatusOK,
		},
		{
			name: "Second test. Incorrect value",
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
			name: "First test",
			args: args{
				metricName: "SomeMetric",
				metricValueStr: "123",
				ms: storage.NewMemStorage(),
			},
			want: http.StatusOK,
		},
		{
			name: "Second test. Incorrect value",
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
*/