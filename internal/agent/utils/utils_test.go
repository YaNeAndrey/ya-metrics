package utils

import (
	"github.com/YaNeAndrey/ya-metrics/internal/constants"
	"github.com/YaNeAndrey/ya-metrics/internal/server/handlers"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/YaNeAndrey/ya-metrics/internal/agent/config"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"

	"github.com/stretchr/testify/assert"
)

// +++
func Test_collectNewMetrics(t *testing.T) {

	testStorage := storage.StorageRepo(storage.NewMemStorageJSON([]storage.Metrics{}))

	type args struct {
		st *storage.StorageRepo
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "First test. Collect all metrics",
			args: args{
				st: &testStorage,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collectNewMetrics(tt.args.st)
			assert.Equal(t, 29, len((*tt.args.st).GetAllMetrics()))
		})
	}
}

func Test_sendAllMetricsUpdates(t *testing.T) {
	type args struct {
		st *storage.StorageRepo
		c  *config.Config
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendAllMetricsUpdates(tt.args.st, tt.args.c)
		})
	}
}

// +++
func Test_sendOneMetricUpdate(t *testing.T) {
	floatValue := float64(124.2345)
	intValue := int64(124)

	testStorage := storage.StorageRepo(storage.NewMemStorageJSON([]storage.Metrics{}))

	type args struct {
		c      *config.Config
		metric storage.Metrics
	}
	tests := []struct {
		name    string
		args    args
		st      *storage.StorageRepo
		wantErr bool
	}{
		{
			name: "First test. Send gauge metric",
			args: args{
				c: config.NewConfig(),
				metric: storage.Metrics{
					ID:    "NewGauge",
					MType: constants.GaugeMetricType,
					Value: &floatValue,
				},
			},
			st:      &testStorage,
			wantErr: false,
		},

		{
			name: "Second test. Send counter metric",
			args: args{
				c: config.NewConfig(),
				metric: storage.Metrics{
					ID:    "NewCounter",
					MType: constants.CounterMetricType,
					Delta: &intValue,
				},
			},
			st:      &testStorage,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Route("/update", func(r chi.Router) {
				r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
					handlers.HandlePostUpdateMetricValueJSON(rw, r, tt.st)
				})
			})
			server := httptest.NewServer(nil)
			defer server.Close()

			host := strings.Split(server.URL, "/")[2]
			hostBufSlice := strings.Split(host, ":")

			hostname := hostBufSlice[0]
			port, _ := strconv.Atoi(hostBufSlice[1])

			tt.args.c.SetSrvAddr(hostname)
			tt.args.c.SetSrvPort(port)

			err := sendOneMetricUpdate(tt.args.c, tt.args.metric)
			if err != nil {
				log.Fatal(err)
			}
		})
	}
}
