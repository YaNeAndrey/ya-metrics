package utils

import (
	"github.com/YaNeAndrey/ya-metrics/internal/constants"
	"github.com/YaNeAndrey/ya-metrics/internal/server/handlers"
	"github.com/YaNeAndrey/ya-metrics/internal/storage/storagejson"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/YaNeAndrey/ya-metrics/internal/agent/config"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
)

// +++
/*
func Test_collectNewMetrics(t *testing.T) {

	testStorage := storage.StorageRepo(storagejson.NewMemStorageJSON([]storage.Metrics{}))

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
	myContext := context.TODO()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collectNewMetrics(tt.args.st)
			metrics, _ := (*tt.args.st).GetAllMetrics(myContext)
			assert.Equal(t, 29, len(metrics))
		})
	}
}*/

// +++
func Test_sendOneMetricUpdate(t *testing.T) {
	floatValue := float64(124.2345)
	intValue := int64(124)

	testStorage := storage.StorageRepo(storagejson.NewMemStorageJSON([]storage.Metrics{}))
	client := http.Client{}
	conf := config.NewConfig()

	type args struct {
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

	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
			handlers.HandlePostUpdateMultipleMetricsJSON(rw, r, &testStorage)
		})
	})
	server := httptest.NewServer(nil)
	defer server.Close()
	host := strings.Split(server.URL, "/")[2]
	hostBufSlice := strings.Split(host, ":")

	hostname := hostBufSlice[0]
	port, _ := strconv.Atoi(hostBufSlice[1])

	conf.SetSrvAddr(hostname)
	conf.SetSrvPort(port)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sendOneMetricUpdate(conf, tt.args.metric, &client)
			if err != nil {
				log.Fatal(err)
			}
		})
	}
}

func Test_sendAllMetricsInOneRequest(t *testing.T) {
	floatArray := []float64{123.456, 456.789, 12, 0}
	intArray := []int64{1, 22, 23}

	testStorage := storage.StorageRepo(storagejson.NewMemStorageJSON([]storage.Metrics{}))
	client := http.Client{}
	conf := config.NewConfig()

	tests := []struct {
		name        string
		metrics     []storage.Metrics
		wantStorage *storagejson.MemStorageJSON
		wantErr     bool
	}{
		{
			name: "First test. Send correct data",
			metrics: []storage.Metrics{
				{
					ID:    "FirstGauge",
					MType: "gauge",
					Value: &floatArray[0],
				},
				{
					ID:    "FirstCounter",
					MType: "counter",
					Delta: &intArray[0],
				},
			},
			wantErr: false,
			wantStorage: storagejson.NewMemStorageJSON([]storage.Metrics{
				{
					ID:    "FirstGauge",
					MType: "gauge",
					Value: &floatArray[0],
				},
				{
					ID:    "FirstCounter",
					MType: "counter",
					Delta: &intArray[0],
				},
			}),
		},
		{
			name: "Second test. Send correct data with repeating metrics",
			metrics: []storage.Metrics{
				{
					ID:    "SecondGauge",
					MType: "gauge",
					Value: &floatArray[0],
				},
				{
					ID:    "SecondCounter",
					MType: "counter",
					Delta: &intArray[0],
				},
				{
					ID:    "SecondGauge",
					MType: "gauge",
					Value: &floatArray[1],
				},
				{
					ID:    "SecondCounter",
					MType: "counter",
					Delta: &intArray[2],
				},
			},
			wantErr: false,
			wantStorage: storagejson.NewMemStorageJSON([]storage.Metrics{
				{
					ID:    "FirstGauge",
					MType: "gauge",
					Value: &floatArray[0],
				},
				{
					ID:    "FirstCounter",
					MType: "counter",
					Delta: &intArray[0],
				},
				{
					ID:    "SecondGauge",
					MType: "gauge",
					Value: &floatArray[1],
				},
				{
					ID:    "SecondCounter",
					MType: "counter",
					Delta: &intArray[2],
				},
			}),
		},
		{
			name: "Third test. Trying to send incorrect data",
			metrics: []storage.Metrics{
				{
					ID:    "SecondGauge",
					MType: "gauge",
					Value: &floatArray[0],
				},
				{
					ID:    "SecondCounter",
					MType: "counter",
					Value: &floatArray[0],
				},
			},
			wantErr: true,
		},
	}
	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
			handlers.HandlePostUpdateMultipleMetricsJSON(rw, r, &testStorage)
		})
	})
	server := httptest.NewServer(nil)
	defer server.Close()
	host := strings.Split(server.URL, "/")[2]
	hostBufSlice := strings.Split(host, ":")

	hostname := hostBufSlice[0]
	port, _ := strconv.Atoi(hostBufSlice[1])

	conf.SetSrvAddr(hostname)
	conf.SetSrvPort(port)

	log.Println(conf)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sendAllMetricsInOneRequest(conf, tt.metrics, &client)
			if !tt.wantErr && err != nil {
				log.Fatal(err)
			}
		})
	}
}
